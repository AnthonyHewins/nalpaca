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
	SubscriptionNews
)

// Configuration for a particular data type in a client. E.g. the stocks client can use bars and quotes, but the symbols for either,
// you can pick through these config options
type StreamTypeConfig struct {
	Enabled bool     `env:"ENABLED" desc:"Enable this stream; if false, this data type is never subscribed to"`
	Symbols []string `env:"SYMBOLS" desc:"use ',' as delimiter"` // use ',' as delimiter
	// Timeout for nalpaca to send the message
	Timeout time.Duration `env:"TIMEOUT" envDefault:"1s" desc:"Timeout for nalpaca to send the message"`
	// This is a separate buffer pool than the socket. This is nalpaca's configured buffer size for when
	// proto messages are serialized, and hence it is smaller and more optimized per message type.
	// It is not recommended to adjust this because it is already set to be exactly the correct size for most
	// message types
	BufSize uint32 `env:"BUFFER_SIZE_BYTES" desc:"Buffer size, in bytes, nalpaca uses when serializing this message type; not recommended to change"`
}

// StreamConfig is the configuration used for the 4 clients in streaming: stocks, options, news, crypto.
// Since each client can only have 1 connection it's required that there are only ever 4 config structs
// like this one in your application
type StreamConfig struct {
	// URL override. You probably dont want this, because there is no sandbox URL
	URL string `env:"URL" desc:"URL override. You probably dont want this, because there is no sandbox URL"`
	// Configure SDK processor count, which appears to be a misnomer; the SDK uses this for goroutine counts
	// from what I can tell
	Processors uint16 `env:"PROCESSORS" desc:"Number of goroutines the underlying Alpaca SDK client uses to process incoming socket messages"`

	// Below are options directly on the alpaca socket. These options are passed directly on to the
	// SDK
	SocketBufSize  uint32        `env:"SOCKET_BUFFER_SIZE" desc:"Buffer size, in bytes, for the underlying websocket connection to Alpaca"`
	ReconnectLimit uint16        `env:"RECONNECT_LIMIT" envDefault:"20" desc:"Number of times the client retries reconnecting to Alpaca's websocket before giving up"` // default in lib
	ReconnectDelay time.Duration `env:"RECONNECT_DELAY" envDefault:"150ms" desc:"Delay between websocket reconnect attempts"`                                          // default in lib
}

// baseClient covers most functionality for all message transmitters, and keeps the factory object
type baseClient[X any] struct {
	*ClientFactory
	symbolList
	metrics
	pool sync.Pool
	t    time.Duration
	comp Subscription
}

func newBaseClient[X any](component Subscription, c *ClientFactory, symbols []string, bufSize int, timeout time.Duration) baseClient[X] {
	return baseClient[X]{
		ClientFactory: c,
		symbolList:    newSymbolList(symbols...),
		metrics:       newMetrics(component),
		t:             timeout,
		pool:          sync.Pool{New: func() any { return make([]byte, 0, bufSize) }},
		comp:          component,
	}
}

func (b *baseClient[X]) bytePool() *sync.Pool       { return &b.pool }
func (b *baseClient[X]) componentMetrics() *metrics { return &b.metrics }
func (b *baseClient[X]) timeout() time.Duration     { return b.t }
func (b *baseClient[X]) component() Subscription    { return b.comp }

func (b *baseClient[X]) List() []string { return b.symbolList.list() }

func (b *baseClient[X]) addSubscription(addFn func(func(X), ...string) error, handler func(X), s ...string) error {
	for i, v := range s {
		s[i] = b.clean(v)
	}

	l := b.l.With("component", b.comp.String(), "symbols", s)

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
