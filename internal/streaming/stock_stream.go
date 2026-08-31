package streaming

import (
	"context"
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/prometheus/client_golang/prometheus"
)

type StocksConfig struct {
	Feed string `env:"FEED"`
	StreamConfig
	Bar   StreamTypeConfig `envPrefix:"BAR_"`
	Quote StreamTypeConfig `envPrefix:"QUOTE_"`
	Trade StreamTypeConfig `envPrefix:"TRADE_"`
}

type StockSubscriptionManagers struct {
	s *stream.StocksClient

	Bars   *bars
	Quotes *quotes
	Trades *trades
}

func (s *StockSubscriptionManagers) Terminated() <-chan error { return s.s.Terminated() }

func (c *Client) Stocks(d *StocksConfig) (*StockSubscriptionManagers, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	switch d.Feed {
	case "sip", "iex", "otc", "delayed_sip":
	default:
		c.l.Error("invalid feed", "feed", d.Feed)
		return nil, fmt.Errorf("invalid feed %s", d.Feed)
	}

	so := []stream.StockOption{}
	for _, v := range c.streamOpts(&d.StreamConfig) {
		so = append(so, v)
	}

	m := &StockSubscriptionManagers{}

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
	m.s = stream.NewStocksClient(d.Feed, so...)

	if m.Bars != nil {
		m.Bars.client = m.s
	}
	if m.Quotes != nil {
		m.Quotes.client = m.s
	}
	if m.Trades != nil {
		m.Trades.client = m.s
	}

	return m, nil
}

func (m *StockSubscriptionManagers) Metrics() []prometheus.Collector {
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
	if err := m.s.Connect(ctx); err != nil {
		return err
	}

	return <-m.s.Terminated()
}
