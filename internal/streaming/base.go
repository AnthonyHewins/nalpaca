package streaming

import (
	"sync"
	"time"
)

//go:generate enumer -type Subscription -trimprefix Subscription -transform snake
type Subscription byte

const (
	SubscriptionOptionTrades Subscription = iota + 1
	SubscriptionOptionQuotes
	SubscriptionStockBars
	SubscriptionStockQuotes
	SubscriptionStockTrades
	SubscriptionOptionsTrade
	SubscriptionNews
)

type StreamTypeConfig struct {
	Enabled bool     `env:"ENABLED"`
	Symbols []string `env:"SYMBOLS"` // use ',' as delimiter
	// Timeout for nalpaca to send the message
	Timeout time.Duration `env:"TIMEOUT" envDefault:"1s"`
	// This is a separate buffer pool than the socket. This is nalpaca's configured buffer size for when
	// proto messages are serialized, and hence it is smaller and more optimized per message type
	BufSize uint32 `env:"BUFFER_SIZE"`
}

// StreamConfig is the configuration used for particular clients, e.g. if you have a client for stocks or news,
// this config is available for you to use
type StreamConfig struct {
	Version string `env:"VERSION"`

	// Feed       string   `env:"FEED_TYPE"`
	// Symbols    []string `env:"SYMBOLS"` // use ',' as delimiter
	Processors uint16 `env:"PROCESSORS" envDefault:"1"`

	// Below are options directly on the alpaca socket. These options are passed directly on to the
	// SDK
	SocketBufSize  uint32        `env:"SOCKET_BUFFER_SIZE" envDefault:"100000"` // default in lib
	ReconnectLimit uint16        `env:"RECONNECT_LIMIT" envDefault:"20"`        // default in lib
	ReconnectDelay time.Duration `env:"RECONNECT_DELAY" envDefault:"150ms"`     // default in lib
}

type baseClient[X any] struct {
	*Client
	symbolList
	metrics
	pool sync.Pool
	t    time.Duration
	comp Subscription
}

func (b *baseClient[X]) bytePool() *sync.Pool       { return &b.pool }
func (b *baseClient[X]) componentMetrics() *metrics { return &b.metrics }
func (b *baseClient[X]) timeout() time.Duration     { return b.t }
func (b *baseClient[X]) component() Subscription    { return b.comp }

func newBaseClient[X any](component Subscription, c *Client, symbols []string, bufSize int, timeout time.Duration) baseClient[X] {
	return baseClient[X]{
		Client:     c,
		symbolList: newSymbolList(symbols...),
		metrics:    newMetrics(component),
		t:          timeout,
		pool:       sync.Pool{New: func() any { return make([]byte, 0, bufSize) }},
	}
}

func (b *baseClient[X]) addSubscription(addFn func(func(X), ...string) error, handler func(X), s ...string) error {
	for i, v := range s {
		s[i] = b.clean(v)
	}

	l := b.l.With("component", b.comp, "symbols", s)

	if err := addFn(handler, s...); err != nil {
		l.Error("failed adding to subscription", "err", err)
		return err
	}

	b.add(s...)
	l.Debug("added to subscription", "newList", b.list())
	return nil
}

func (b *baseClient[X]) rmSubscription(fn func(...string) error, s ...string) error {
	for i, v := range s {
		s[i] = b.clean(v)
	}

	l := b.l.With("component", b.comp, "symbols", s)

	if err := fn(s...); err != nil {
		l.Error("failed removing symbols from subscriptions", "err", err)
		return err
	}

	b.symbolList.del(s...)
	l.Debug("removed symbols from subscriptions")
	return nil
}

func (b *baseClient[X]) listSubscription() []string { return b.list() }
