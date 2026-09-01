package streaming

import (
	"context"
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/prometheus/client_golang/prometheus"
)

var _ config = (*StocksConfig)(nil)

const (
	defaultBarBufSize   = 128
	defaultQuoteBufSize = 128
	defaultTradeBufSize = 128
)

type StocksConfig struct {
	// Configuration for the websocket connection to stocks
	StreamConfig
	// What feed to use; this must be of type marketfeed.Feed in alpaca's SDK
	Feed string `env:"FEED" envDefault:"iex"`

	Bar   StreamTypeConfig `envPrefix:"BARS_"`
	Quote StreamTypeConfig `envPrefix:"QUOTES_"`
	Trade StreamTypeConfig `envPrefix:"TRADES_"`
}

// validate implements [config].
func (s *StocksConfig) validate() (bool, error) {
	if !s.Bar.Enabled && !s.Quote.Enabled && !s.Trade.Enabled {
		return false, nil
	}

	switch s.Feed {
	case "sip", "iex", "otc", "delayed_sip":
	default:
		return false, fmt.Errorf("invalid feed %s", s.Feed)
	}

	return true, nil
}

func (s *StocksConfig) setDefaults() {
	if s.Bar.BufSize == 0 {
		s.Bar.BufSize = defaultBarBufSize
	}

	if s.Quote.BufSize == 0 {
		s.Quote.BufSize = defaultQuoteBufSize
	}

	if s.Trade.BufSize == 0 {
		s.Trade.BufSize = defaultTradeBufSize
	}
}

// StockSubscriptionManagers is a struct collecting all the things that use the stocks client
// for messaging. Since the stocks client has several data types it allows, it has 3 subclients
// that control behavior
type StockSubscriptionManagers struct {
	Conn   *stream.StocksClient
	Bars   *Bars
	Quotes *Quotes
	Trades *Trades
}

func (s *StockSubscriptionManagers) Terminated() <-chan error { return s.Conn.Terminated() }

func (c *ClientFactory) Stocks(d *StocksConfig) (StockSubscriptionManagers, error) {
	if enabled, err := c.prepare(d); err != nil || !enabled {
		return StockSubscriptionManagers{}, err
	}

	so := []stream.StockOption{}
	for _, v := range c.streamOpts(&d.StreamConfig) {
		so = append(so, v)
	}

	m := StockSubscriptionManagers{}

	if d.Bar.Enabled {
		m.Bars = newBars(c, &d.Bar, nil)
		so = append(so, stream.WithBars(m.Bars.handler, m.Bars.List()...))
	}

	if d.Quote.Enabled {
		m.Quotes = newQuotes(c, &d.Quote, nil)
		so = append(so, stream.WithQuotes(m.Quotes.handler, m.Quotes.List()...))
	}

	if d.Trade.Enabled {
		m.Trades = newTrades(c, &d.Trade, nil)
		so = append(so, stream.WithTrades(m.Trades.handler, m.Trades.List()...))
	}

	c.l.Info("creating stocks stream client", "conf", d)
	m.Conn = stream.NewStocksClient(d.Feed, so...)

	if m.Bars != nil {
		m.Bars.client = m.Conn
	}
	if m.Quotes != nil {
		m.Quotes.client = m.Conn
	}
	if m.Trades != nil {
		m.Trades.client = m.Conn
	}

	return m, nil
}

func (m *StockSubscriptionManagers) Metrics() []prometheus.Collector {
	if m == nil {
		return nil
	}

	var out []prometheus.Collector
	if m.Bars != nil {
		out = append(out, m.Bars.Metrics()...)
	}
	if m.Quotes != nil {
		out = append(out, m.Quotes.Metrics()...)
	}
	if m.Trades != nil {
		out = append(out, m.Trades.Metrics()...)
	}
	return out
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (m *StockSubscriptionManagers) Stream(ctx context.Context) error {
	if m.Conn == nil {
		return nil
	}

	if err := m.Conn.Connect(ctx); err != nil {
		return err
	}

	return <-m.Conn.Terminated()
}
