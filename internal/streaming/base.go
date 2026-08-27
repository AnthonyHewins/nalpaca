package streaming

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"
)

// publisher is the slice of jetstream.JetStream the handlers actually use.
// Narrowing it here is what makes them testable without a live NATS server.
type publisher interface {
	Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

func newMetric(appName, subsystem, name, help string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: appName,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
}

type Metrics struct {
	TotalErr, MarshalErr, PubErr prometheus.Counter
}

// NewMetrics builds the error counters for one stream. subsystem must be unique
// per stream: prometheus rejects two collectors with the same fully-qualified
// name, so sharing one here makes the whole app fail to start once metrics are
// enabled.
func NewMetrics(appName, subsystem string) Metrics {
	return Metrics{
		TotalErr:   newMetric(appName, subsystem, "total_err", "total error count"),
		MarshalErr: newMetric(appName, subsystem, "marshal_err", "marshal error count"),
		PubErr:     newMetric(appName, subsystem, "pub_err", "nats pub err count"),
	}
}

// Collectors returns the counters for registration with a prometheus registry.
func (m Metrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{m.TotalErr, m.MarshalErr, m.PubErr}
}

type symbolList struct {
	mu      sync.RWMutex
	symbols map[string]struct{}
}

func newSymbolList(x ...string) *symbolList {
	m := make(map[string]struct{}, len(x))
	for _, v := range x {
		m[v] = struct{}{}
	}
	return &symbolList{symbols: m}
}

func (s *symbolList) add(x ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range x {
		s.symbols[v] = struct{}{}
	}
}

func (s *symbolList) del(x ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range x {
		delete(s.symbols, v)
	}
}

func (s *symbolList) list() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	x := make([]string, len(s.symbols))
	i := 0
	for v := range s.symbols {
		x[i] = v
		i++
	}
	return x
}

// Default websocket endpoints per stream type. These differ by path, so a single
// shared envDefault on Stream.BaseURL would silently point one stream type at
// another's endpoint. Each constructor applies its own via Stream.baseURL.
const (
	stocksURL  = "https://stream.data.alpaca.markets/v2"
	newsURL    = "https://stream.data.alpaca.markets/v1beta1/news"
	optionsURL = "https://stream.data.alpaca.markets/v1beta1"
)

type Stream struct {
	Feed string `env:"FEED_TYPE" envDefault:"sip"`

	// Symbols is the fallback symbol list for any message type that doesn't
	// specify its own. Use ',' as delimiter.
	Symbols []string `env:"SYMBOLS"`

	// Per-message-type toggles. These default off because quotes and trades are
	// orders of magnitude higher volume than bars: bars arrive roughly once per
	// minute per symbol, quotes can arrive thousands of times per second.
	Bars   bool `env:"BARS" envDefault:"false"`
	Quotes bool `env:"QUOTES" envDefault:"false"`
	Trades bool `env:"TRADES" envDefault:"false"`

	// Per-message-type symbol lists. Empty means "fall back to Symbols".
	BarSymbols   []string `env:"BAR_SYMBOLS"`
	QuoteSymbols []string `env:"QUOTE_SYMBOLS"`
	TradeSymbols []string `env:"TRADE_SYMBOLS"`

	// BaseURL has no envDefault on purpose; see the URL constants above.
	BaseURL        string        `env:"BASE_URL"`
	Processors     uint16        `env:"PROCESSORS" envDefault:"1"`
	Buffer         uint32        `env:"BUFFER_SIZE" envDefault:"100000"`    // default in lib
	ReconnectLimit uint16        `env:"RECONNECT_LIMIT" envDefault:"20"`    // default in lib
	ReconnectDelay time.Duration `env:"RECONNECT_DELAY" envDefault:"150ms"` // default in lib
}

// baseURL returns the configured URL, falling back to fallback when unset.
func (s *Stream) baseURL(fallback string) string {
	if s.BaseURL == "" {
		return fallback
	}
	return s.BaseURL
}

// barSymbols, quoteSymbols and tradeSymbols resolve a message type's symbol list,
// falling back to the shared Symbols list when no type-specific one is set.
func (s *Stream) barSymbols() []string   { return orFallback(s.BarSymbols, s.Symbols) }
func (s *Stream) quoteSymbols() []string { return orFallback(s.QuoteSymbols, s.Symbols) }
func (s *Stream) tradeSymbols() []string { return orFallback(s.TradeSymbols, s.Symbols) }

func orFallback(specific, fallback []string) []string {
	if len(specific) > 0 {
		return specific
	}
	return fallback
}

// streamOpts builds the options shared by every stream type. url is the resolved
// endpoint from Stream.baseURL, not the raw config field, so the log lines name
// the endpoint actually dialed.
func streamOpts(key, secret, url string, logger *slog.Logger, s *Stream) []stream.Option {
	return []stream.Option{
		stream.WithCredentials(key, secret),
		stream.WithProcessors(int(s.Processors)),
		stream.WithBaseURL(url),
		stream.WithBufferSize(int(s.Buffer)), // default value
		stream.WithReconnectSettings(int(s.ReconnectLimit), s.ReconnectDelay),
		stream.WithBufferFillCallback(func(msg []byte) {
			logger.Info("buffer has been filled, processing interrupted", "len(bufferWaiting)", len(msg))
		}),
		stream.WithDisconnectCallback(func() { logger.Warn("stream was disconnected", "url", url) }),
		stream.WithConnectCallback(func() { logger.Info("stream connected", "url", url) }),
		stream.WithLogger(streamLogger{logger.With("alpaca", true)}),
	}
}

// publish marshals msg and pushes it onto subject, accounting for failures on
// the way. raw is the original SDK struct and is only used for log context.
//
// Every stream handler funnels through here so that error accounting can't drift
// between message types.
func publish(logger *slog.Logger, metrics Metrics, js publisher, subject string, raw any, msg proto.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var err error
	defer func() {
		if err != nil {
			metrics.TotalErr.Inc()
		}
	}()

	buf, err := proto.Marshal(msg)
	if err != nil {
		logger.ErrorContext(ctx, "failed marshal", "err", err, "subject", subject, "raw", raw)
		metrics.MarshalErr.Inc()
		return
	}

	if _, err = js.Publish(ctx, subject, buf); err != nil {
		logger.ErrorContext(ctx, "failed publishing", "err", err, "subject", subject, "raw", raw)
		metrics.PubErr.Inc()
	}
}

// listOp is whether a subscription change adds or removes symbols.
type listOp bool

const (
	add listOp = true
	del listOp = false
)

func (o listOp) String() string {
	if o == add {
		return "added"
	}
	return "removed"
}

// resubscribe applies op to list, then re-sends the entire resulting list to
// alpaca. The SDK's subscribe calls are absolute rather than incremental, so the
// full list has to be sent on every change.
//
// On failure the local list is rolled back, otherwise it would drift out of sync
// with what alpaca actually has subscribed.
func resubscribe[T any](
	logger *slog.Logger,
	kind string,
	list *symbolList,
	subscribe func(func(T), ...string) error,
	handler func(T),
	op listOp,
	delta []string,
) error {
	if len(delta) == 0 {
		return nil
	}

	if op == add {
		list.add(delta...)
	} else {
		list.del(delta...)
	}

	l := list.list()
	if err := subscribe(handler, l...); err != nil {
		// Roll back so the local view matches alpaca's.
		if op == add {
			list.del(delta...)
		} else {
			list.add(delta...)
		}

		logger.Error("failed changing subscriptions",
			"err", err,
			"kind", kind,
			"op", op.String(),
			"delta", delta,
		)
		return err
	}

	logger.Info("changed subscriptions", "kind", kind, "op", op.String(), "delta", delta, "final", l)
	return nil
}

type streamLogger struct {
	l *slog.Logger
}

func (l streamLogger) Infof(format string, v ...interface{})  { l.l.Info(fmt.Sprintf(format, v...)) }
func (l streamLogger) Warnf(format string, v ...interface{})  { l.l.Warn(fmt.Sprintf(format, v...)) }
func (l streamLogger) Errorf(format string, v ...interface{}) { l.l.Error(fmt.Sprintf(format, v...)) }
