package streaming

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"google.golang.org/protobuf/proto"
)

// published is one captured Publish call.
type published struct {
	subject string
	payload []byte
}

// mockPublisher stands in for jetstream.JetStream. Follows the hand-written
// function-field style of internal/bridge.Mock.
//
// jetstream.JetStream is a large interface (stream/consumer/KV/object-store
// management, account info, ...) and ClientFactory only ever calls Publish
// on it. Embedding the (nil) interface satisfies the rest of the method set
// without stubbing it out by hand; anything beyond Publish will nil-panic if
// a test ever exercises it, which is the point - it should update this mock,
// not silently succeed.
type mockPublisher struct {
	jetstream.JetStream

	mu   sync.Mutex
	msgs []published

	// PublishFn, when set, decides the outcome. Leave nil to always succeed.
	PublishFn func(subject string, payload []byte) (*jetstream.PubAck, error)
}

func (m *mockPublisher) Publish(_ context.Context, subject string, payload []byte, _ ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	m.mu.Lock()
	m.msgs = append(m.msgs, published{subject: subject, payload: payload})
	m.mu.Unlock()

	if m.PublishFn != nil {
		return m.PublishFn(subject, payload)
	}
	return &jetstream.PubAck{}, nil
}

func (m *mockPublisher) calls() []published {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]published(nil), m.msgs...)
}

// only asserts exactly one message was published, and returns it.
func (m *mockPublisher) only(t *testing.T) published {
	t.Helper()

	got := m.calls()
	if len(got) != 1 {
		t.Fatalf("expected exactly 1 published message, got %d", len(got))
	}
	return got[0]
}

// failingPublisher returns a publisher whose Publish always errors.
func failingPublisher() *mockPublisher {
	return &mockPublisher{
		PublishFn: func(string, []byte) (*jetstream.PubAck, error) {
			return nil, fmt.Errorf("nats is down")
		},
	}
}

func testLogger() *slog.Logger { return slog.New(slog.DiscardHandler) }

// unmarshal decodes a captured payload into msg.
func unmarshal(t *testing.T, p published, msg proto.Message) {
	t.Helper()

	if err := proto.Unmarshal(p.payload, msg); err != nil {
		t.Fatalf("failed unmarshalling payload on %s: %v", p.subject, err)
	}
}

func counter(t *testing.T, c prometheus.Counter) float64 {
	t.Helper()
	return testutil.ToFloat64(c)
}

func assertSubject(t *testing.T, got published, want string) {
	t.Helper()

	if got.subject != want {
		t.Errorf("wrong subject: want %q, got %q", want, got.subject)
	}
}

func assertStrings(t *testing.T, field string, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("%s: want %v, got %v", field, want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s[%d]: want %q, got %q", field, i, want[i], got[i])
		}
	}
}
