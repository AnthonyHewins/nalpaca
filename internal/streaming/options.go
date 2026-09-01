package streaming

import (
	"context"
	"errors"
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultOptionQuoteBufSize = 128
	defaultOptionTradeBufSize = 128
)

var _ config = (*OptionsConfig)(nil)

type OptionsConfig struct {
	StreamConfig
	Feed  string           `env:"FEED" envDefault:"opra"`
	Quote StreamTypeConfig `envPrefix:"QUOTE_"`
	Trade StreamTypeConfig `envPrefix:"TRADE_"`
}

func (o *OptionsConfig) setDefaults() {}
func (o *OptionsConfig) validate() (bool, error) {
	if !o.Quote.Enabled && !o.Trade.Enabled {
		return false, nil
	}

	if o.Quote.Enabled && len(o.Quote.Symbols) == 0 {
		return false, errors.New("option quote streaming is enabled, but no symbols were given")
	}

	if o.Trade.Enabled && len(o.Trade.Symbols) == 0 {
		return false, errors.New("option trade streaming is enabled, but no symbols were given")
	}

	switch o.Feed {
	case "opra", "indicative":
	default:
		return false, fmt.Errorf("invalid option feed %s", o.Feed)
	}

	if o.Quote.BufSize == 0 {
		o.Quote.BufSize = defaultOptionQuoteBufSize
	}

	if o.Trade.BufSize == 0 {
		o.Trade.BufSize = defaultOptionTradeBufSize
	}

	return true, nil
}

// OptionSubscriptionManagers is a struct collecting all the things that use the options client
// for messaging. Since the options client has several data types it allows, it has multiple subclients
// that control subscriptions
type OptionSubscriptionManagers struct {
	Conn   *stream.OptionClient
	Quotes *optionQuotes
	Trades *optionTrades
}

func (o *OptionSubscriptionManagers) Terminated() <-chan error { return o.Conn.Terminated() }

func (c *ClientFactory) Options(d *OptionsConfig) (OptionSubscriptionManagers, error) {
	if enabled, err := c.prepare(d); err != nil || !enabled {
		return OptionSubscriptionManagers{}, err
	}

	so := []stream.OptionOption{}
	for _, v := range c.streamOpts(&d.StreamConfig) {
		so = append(so, v)
	}

	m := OptionSubscriptionManagers{}

	if d.Quote.Enabled {
		m.Quotes = newOptionQuotes(c, &d.Quote, nil)
		so = append(so, stream.WithOptionQuotes(m.Quotes.handler, m.Quotes.List()...))
	}

	if d.Trade.Enabled {
		m.Trades = newOptionTrades(c, &d.Trade, nil)
		so = append(so, stream.WithOptionTrades(m.Trades.handler, m.Trades.List()...))
	}

	c.l.Info("creating options stream client", "conf", d)
	m.Conn = stream.NewOptionClient(d.Feed, so...)

	if m.Quotes != nil {
		m.Quotes.client = m.Conn
	}
	if m.Trades != nil {
		m.Trades.client = m.Conn
	}

	return m, nil
}

func (m *OptionSubscriptionManagers) Metrics() []prometheus.Collector {
	if m == nil {
		return nil
	}

	var out []prometheus.Collector
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
func (m *OptionSubscriptionManagers) Stream(ctx context.Context) error {
	if m.Conn == nil {
		return nil
	}

	if err := m.Conn.Connect(ctx); err != nil {
		return err
	}

	return <-m.Conn.Terminated()
}
