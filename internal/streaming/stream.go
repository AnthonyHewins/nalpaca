package streaming

import (
	"context"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	prefix, key, secret, baseURL string
	js                           jetstream.JetStream
	l                            *slog.Logger
	t                            trace.Tracer

	marshalOpts proto.MarshalOptions
}

func New(prefix, baseUrl, key, secret string, l *slog.Logger, js jetstream.JetStream, t trace.Tracer) Client {
	return Client{
		prefix:  prefix,
		key:     key,
		secret:  secret,
		l:       l,
		t:       t,
		baseURL: baseUrl,
		js:      js,
	}
}

func (c *Client) streamOpts(s *Stream) []stream.Option {
	x, _ := url.JoinPath(c.baseURL, s.Version)
	return []stream.Option{
		stream.WithCredentials(c.key, c.secret),
		stream.WithProcessors(int(s.Processors)),
		stream.WithBaseURL(x),
		stream.WithBufferSize(int(s.SocketBufSize)), // default value
		stream.WithReconnectSettings(int(s.ReconnectLimit), s.ReconnectDelay),
		stream.WithBufferFillCallback(func(msg []byte) {
			c.l.Info("buffer has been filled, processing interrupted", "len(bufferWaiting)", len(msg))
		}),
		stream.WithDisconnectCallback(func() { c.l.Warn("stream was disconnected", "url", c.baseURL) }),
		stream.WithConnectCallback(func() { c.l.Info("stream connected", "url", c.baseURL) }),
		stream.WithLogger(streamLogger{c.l.With("alpaca", true)}),
	}
}

type transmitter[X any, Y proto.Message] interface {
	toWire(X) (Y, error)
	componentMetrics() *metrics
	component() string
	subject(w Y) string
	timeout() time.Duration
	bytePool() *sync.Pool
}

// i want this to be generic when possible with go1.27
func wrap[X any, W proto.Message](c *Client, t transmitter[X, W], x X) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout())
	defer cancel()

	ctx, span := c.t.Start(ctx, t.component())
	defer span.End()

	m := t.componentMetrics()
	m.receiveCount.Inc()

	var err error
	defer func() {
		if err == nil {
			span.SetStatus(codes.Ok, "success")
			return
		}

		m.totalErr.Inc()
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed")
	}()

	l := c.l.With("raw", x)

	var w W
	if w, err = t.toWire(x); err != nil {
		m.transformErr.Inc()
		l.Error("failed converting message", "err", err)
		return
	}

	pool := t.bytePool()
	buf := pool.Get().([]byte)
	defer func() { pool.Put(buf[:0]) }()

	if buf, err = c.marshalOpts.MarshalAppend(buf, w); err != nil {
		m.marshalErr.Inc()
		c.l.Error("failed marshaling proto", "err", err)
		return
	}

	l = l.With("wire", w)

	subj := t.subject(w)
	if _, err = c.js.Publish(ctx, subj, buf); err != nil {
		m.pubErr.Inc()
		l.Error("failed publishing msg", "err", err, "subj", subj)
		return
	}

	c.l.Debug("published msg", "subj", subj, "len(bytes)", len(buf))
	m.publishCount.Inc()
}
