package streaming

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type ClientFactory struct {
	prefix, key, secret, baseURL string
	js                           jetstream.JetStream
	l                            *slog.Logger
	t                            trace.Tracer

	marshalOpts proto.MarshalOptions
}

// Create new streaming client. This is the base client that will create subscriptions to alpaca:
// Stock quotes, trades, bars
// Option quotes, trades
// News
func New(prefix, baseUrl, key, secret string, l *slog.Logger, js jetstream.JetStream, t trace.Tracer) ClientFactory {
	return ClientFactory{
		prefix:  prefix,
		key:     key,
		secret:  secret,
		l:       l,
		t:       t,
		baseURL: baseUrl,
		js:      js,
	}
}

func (c *ClientFactory) streamOpts(s *StreamConfig) []stream.Option {
	opts := []stream.Option{
		stream.WithCredentials(c.key, c.secret),
		stream.WithReconnectSettings(int(s.ReconnectLimit), s.ReconnectDelay),
		stream.WithBufferFillCallback(func(msg []byte) {
			c.l.Info("buffer has been filled, processing interrupted", "len(bufferWaiting)", len(msg))
		}),
		stream.WithDisconnectCallback(func() { c.l.Warn("stream was disconnected") }),
		stream.WithConnectCallback(func() { c.l.Info("stream connected") }),
		stream.WithLogger(streamLogger{c.l.With("alpaca", true)}),
	}

	if s.Processors > 0 {
		opts = append(opts, stream.WithProcessors(int(s.Processors)))
	}

	if s.SocketBufSize > 0 {
		opts = append(opts, stream.WithBufferSize(int(s.SocketBufSize)))
	}

	if s.URL != "" {
		opts = append(opts, stream.WithBaseURL(s.URL))
	}

	return opts
}

type Subscriber interface {
	Subscribe(...string) error
	Unsubscribe(...string) error
	List() []string
}

type transmitter[X any, Y proto.Message] interface {
	toWire(X) Y
	componentMetrics() *metrics
	component() Subscription
	subject(w Y) string
	timeout() time.Duration
	bytePool() *sync.Pool
}

type config interface {
	validate() (enabled bool, err error)
	setDefaults()
}

var errMissingOptions = errors.New("options missing for this stream")

func (c *ClientFactory) prepare(x config) (bool, error) {
	if x == nil {
		c.l.Error("options missing")
		return false, errMissingOptions
	}

	enabled, err := x.validate()
	if err != nil {
		c.l.Error("configuration failed validation", "err", err)
		return false, err
	}

	if enabled {
		x.setDefaults()
	}

	return enabled, nil
}

// i want this to be generic when possible with go1.27
func wrap[X any, W proto.Message](c *ClientFactory, t transmitter[X, W], x X) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout())
	defer cancel()

	ctx, span := c.t.Start(ctx, t.component().String())
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

	w := t.toWire(x)

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

type streamLogger struct {
	l *slog.Logger
}

func (l streamLogger) Infof(format string, v ...interface{})  { l.l.Info(fmt.Sprintf(format, v...)) }
func (l streamLogger) Warnf(format string, v ...interface{})  { l.l.Warn(fmt.Sprintf(format, v...)) }
func (l streamLogger) Errorf(format string, v ...interface{}) { l.l.Error(fmt.Sprintf(format, v...)) }
