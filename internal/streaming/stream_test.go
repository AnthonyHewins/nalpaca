package streaming

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/proto"
)

// fakeJS satisfies jetstream.JetStream by embedding the (nil) interface and
// overriding only Publish, which is all wrap() ever calls. Any other method
// being invoked would panic on the nil embed, which is fine for these tests.
type fakeJS struct {
	jetstream.JetStream

	publish func(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

func (f *fakeJS) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return f.publish(ctx, subject, payload, opts...)
}

// fakeTransmitter is a generic, hand-configurable implementation of the
// transmitter interface used to exercise wrap() in isolation from the real
// Stocks/News components.
type fakeTransmitter struct {
	m       metrics
	pool    sync.Pool
	comp    string
	to      time.Duration
	toWireF func(string) (*protoStream.Bar, error)
	subjF   func(*protoStream.Bar) string
}

func (f *fakeTransmitter) toWire(x string) (*protoStream.Bar, error) { return f.toWireF(x) }
func (f *fakeTransmitter) componentMetrics() *metrics                { return &f.m }
func (f *fakeTransmitter) component() string                         { return f.comp }
func (f *fakeTransmitter) subject(w *protoStream.Bar) string         { return f.subjF(w) }
func (f *fakeTransmitter) timeout() time.Duration                    { return f.to }
func (f *fakeTransmitter) bytePool() *sync.Pool                      { return &f.pool }

var _ transmitter[string, *protoStream.Bar] = (*fakeTransmitter)(nil)

func newTestClient(publish func(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)) *Client {
	c := New(
		"test",
		"http://localhost",
		"key",
		"secret",
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		&fakeJS{publish: publish},
		noop.NewTracerProvider().Tracer("test"),
	)
	return &c
}

// Regression test for the pool.New bug: freshly minted pool buffers must be
// zero-length (with spare capacity), not pre-filled with BufSize zero bytes.
// A buffer built via make([]byte, BufSize) would prepend BufSize zero bytes
// ahead of every marshaled message the first time (and periodically
// thereafter, since sync.Pool is cleared on GC) a pool.New buffer is drawn.
func TestWrap_FreshPoolBufferIsNotPrefixedWithZeros(t *testing.T) {
	const bufSize = 128

	var published []byte
	var subj string
	c := newTestClient(func(_ context.Context, subject string, payload []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		subj = subject
		published = append([]byte(nil), payload...)
		return &jetstream.PubAck{}, nil
	})

	tr := &fakeTransmitter{
		m:    newMetrics("test"),
		pool: sync.Pool{New: func() any { return make([]byte, 0, bufSize) }},
		comp: "test",
		to:   time.Second,
		toWireF: func(s string) (*protoStream.Bar, error) {
			return &protoStream.Bar{Symbol: s}, nil
		},
		subjF: func(w *protoStream.Bar) string { return "bars." + w.Symbol },
	}

	wrap(c, tr, "AAPL")

	if subj != "bars.AAPL" {
		t.Fatalf("expected subject bars.AAPL, got %q", subj)
	}
	if published == nil {
		t.Fatal("expected a message to be published")
	}

	want, err := proto.Marshal(&protoStream.Bar{Symbol: "AAPL"})
	if err != nil {
		t.Fatalf("failed marshaling expected message: %v", err)
	}
	if len(published) != len(want) {
		t.Fatalf("published message has wrong length (likely zero-padded): want %d bytes, got %d bytes: %x", len(want), len(published), published)
	}

	var got protoStream.Bar
	if err := proto.Unmarshal(published, &got); err != nil {
		t.Fatalf("published bytes did not unmarshal cleanly (corrupted): %v", err)
	}
	if got.Symbol != "AAPL" {
		t.Fatalf("expected symbol AAPL, got %q", got.Symbol)
	}
}

// The buffer pool must be safely reusable across multiple calls without
// residue from a prior, larger message leaking into a subsequent, smaller
// one.
func TestWrap_PoolIsReusedCleanlyAcrossCalls(t *testing.T) {
	var published [][]byte
	c := newTestClient(func(_ context.Context, _ string, payload []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		published = append(published, append([]byte(nil), payload...))
		return &jetstream.PubAck{}, nil
	})

	tr := &fakeTransmitter{
		m:    newMetrics("test"),
		pool: sync.Pool{New: func() any { return make([]byte, 0, 8) }},
		comp: "test",
		to:   time.Second,
		toWireF: func(s string) (*protoStream.Bar, error) {
			return &protoStream.Bar{Symbol: s}, nil
		},
		subjF: func(w *protoStream.Bar) string { return w.Symbol },
	}

	wrap(c, tr, "GOOGLEISBIGGERTHANOTHERS") // forces the pool buffer to grow
	wrap(c, tr, "A")                        // subsequent, smaller message

	if len(published) != 2 {
		t.Fatalf("expected 2 published messages, got %d", len(published))
	}

	var second protoStream.Bar
	if err := proto.Unmarshal(published[1], &second); err != nil {
		t.Fatalf("second message did not unmarshal cleanly: %v", err)
	}
	if second.Symbol != "A" {
		t.Fatalf("expected second message symbol %q, got %q (likely leaked bytes from prior message)", "A", second.Symbol)
	}
}

func TestWrap_ToWireErrorDoesNotPublish(t *testing.T) {
	published := false
	c := newTestClient(func(_ context.Context, _ string, _ []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		published = true
		return &jetstream.PubAck{}, nil
	})

	tr := &fakeTransmitter{
		m:    newMetrics("test"),
		pool: sync.Pool{New: func() any { return make([]byte, 0, 8) }},
		comp: "test",
		to:   time.Second,
		toWireF: func(s string) (*protoStream.Bar, error) {
			return nil, errors.New("boom")
		},
		subjF: func(w *protoStream.Bar) string { return "subj" },
	}

	wrap(c, tr, "AAPL")

	if published {
		t.Fatal("expected no publish when toWire fails")
	}
	if got := testutil.ToFloat64(tr.m.transformErr); got != 1 {
		t.Fatalf("expected transformErr=1, got %v", got)
	}
	if got := testutil.ToFloat64(tr.m.totalErr); got != 1 {
		t.Fatalf("expected totalErr=1, got %v", got)
	}
	if got := testutil.ToFloat64(tr.m.publishCount); got != 0 {
		t.Fatalf("expected publishCount=0, got %v", got)
	}
}

func TestWrap_PublishErrorIncrementsErrorMetrics(t *testing.T) {
	c := newTestClient(func(_ context.Context, _ string, _ []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		return nil, errors.New("nats down")
	})

	tr := &fakeTransmitter{
		m:    newMetrics("test"),
		pool: sync.Pool{New: func() any { return make([]byte, 0, 8) }},
		comp: "test",
		to:   time.Second,
		toWireF: func(s string) (*protoStream.Bar, error) {
			return &protoStream.Bar{Symbol: s}, nil
		},
		subjF: func(w *protoStream.Bar) string { return "subj" },
	}

	wrap(c, tr, "AAPL")

	if got := testutil.ToFloat64(tr.m.pubErr); got != 1 {
		t.Fatalf("expected pubErr=1, got %v", got)
	}
	if got := testutil.ToFloat64(tr.m.totalErr); got != 1 {
		t.Fatalf("expected totalErr=1, got %v", got)
	}
	if got := testutil.ToFloat64(tr.m.publishCount); got != 0 {
		t.Fatalf("expected publishCount=0, got %v", got)
	}
}

// The following two tests go through the real Client.Stocks/Client.News
// constructors, exercising the exact sync.Pool{New: ...} line each ships
// with (as opposed to the fakeTransmitter tests above, which build their own
// correct pool). They are the regression tests for the bug where
// make([]byte, BufSize) was used instead of make([]byte, 0, BufSize),
// which left every freshly-minted pool buffer prefixed with BufSize zero
// bytes ahead of the real marshaled message.
func TestStocks_FreshPoolBufferProducesUncorruptedMessage(t *testing.T) {
	var published []byte
	c := newTestClient(func(_ context.Context, _ string, payload []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		published = append([]byte(nil), payload...)
		return &jetstream.PubAck{}, nil
	})

	s, err := c.Stocks(&Stream{
		Feed:    "iex",
		Symbols: []string{"AAPL"},
		Timeout: time.Second,
		BufSize: 80, // matches the production default set in cmd/nalpaca/components.go
	})
	if err != nil {
		t.Fatalf("unexpected error creating stocks client: %v", err)
	}

	bar := stream.Bar{Symbol: "AAPL"}
	s.handler(bar)

	wireWant, err := s.toWire(bar)
	if err != nil {
		t.Fatalf("failed converting expected message: %v", err)
	}
	want, err := proto.Marshal(wireWant)
	if err != nil {
		t.Fatalf("failed marshaling expected message: %v", err)
	}
	if len(published) != len(want) {
		t.Fatalf("published bar has wrong length (likely zero-padded): want %d bytes, got %d bytes: %x", len(want), len(published), published)
	}

	var got protoStream.Bar
	if err := proto.Unmarshal(published, &got); err != nil {
		t.Fatalf("published bar did not unmarshal cleanly (corrupted): %v", err)
	}
	if got.Symbol != "AAPL" {
		t.Fatalf("expected symbol AAPL, got %q", got.Symbol)
	}
}

func TestNews_FreshPoolBufferProducesUncorruptedMessage(t *testing.T) {
	var published []byte
	c := newTestClient(func(_ context.Context, _ string, payload []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		published = append([]byte(nil), payload...)
		return &jetstream.PubAck{}, nil
	})

	n, err := c.News(&Stream{
		Symbols: []string{"AAPL"},
		Timeout: time.Second,
		BufSize: 100000, // matches the production default set in cmd/nalpaca/components.go
	})
	if err != nil {
		t.Fatalf("unexpected error creating news client: %v", err)
	}

	newsItem := stream.News{ID: 1, Headline: "hello"}
	n.handler(newsItem)

	wireWant, err := n.toWire(newsItem)
	if err != nil {
		t.Fatalf("failed converting expected message: %v", err)
	}
	want, err := proto.Marshal(wireWant)
	if err != nil {
		t.Fatalf("failed marshaling expected message: %v", err)
	}
	if len(published) != len(want) {
		t.Fatalf("published news has wrong length (likely zero-padded): want %d bytes, got %d bytes: %x", len(want), len(published), published)
	}

	var got protoStream.News
	if err := proto.Unmarshal(published, &got); err != nil {
		t.Fatalf("published news did not unmarshal cleanly (corrupted): %v", err)
	}
	if got.Headline != "hello" {
		t.Fatalf("expected headline %q, got %q", "hello", got.Headline)
	}
}

func TestNews_NilStreamReturnsError(t *testing.T) {
	c := newTestClient(nil)
	if _, err := c.News(nil); err == nil {
		t.Fatal("expected an error passing a nil *Stream, got nil")
	}
}

func TestWrap_SuccessIncrementsCountsAndUsesSubject(t *testing.T) {
	var gotSubj string
	c := newTestClient(func(_ context.Context, subject string, _ []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
		gotSubj = subject
		return &jetstream.PubAck{}, nil
	})

	tr := &fakeTransmitter{
		m:    newMetrics("test"),
		pool: sync.Pool{New: func() any { return make([]byte, 0, 8) }},
		comp: "test",
		to:   time.Second,
		toWireF: func(s string) (*protoStream.Bar, error) {
			return &protoStream.Bar{Symbol: s}, nil
		},
		subjF: func(w *protoStream.Bar) string { return "stocks." + w.Symbol },
	}

	wrap(c, tr, "MSFT")

	if gotSubj != "stocks.MSFT" {
		t.Fatalf("expected subject stocks.MSFT, got %q", gotSubj)
	}
	if got := testutil.ToFloat64(tr.m.receiveCount); got != 1 {
		t.Fatalf("expected receiveCount=1, got %v", got)
	}
	if got := testutil.ToFloat64(tr.m.publishCount); got != 1 {
		t.Fatalf("expected publishCount=1, got %v", got)
	}
	if got := testutil.ToFloat64(tr.m.totalErr); got != 0 {
		t.Fatalf("expected totalErr=0, got %v", got)
	}
}
