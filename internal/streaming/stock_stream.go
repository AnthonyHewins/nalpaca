package streaming

import (
	"context"
	"fmt"
	"sync"
	"time"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const stocksComponent = "stocks"

var (
	_ transmitter[stream.Bar, *protoStream.Bar] = (*Stocks)(nil)
)

type Stocks struct {
	client *Client
	metrics
	sync.Pool
	symbolList

	t time.Duration

	s *stream.StocksClient
}

func (s *Stocks) bytePool() *sync.Pool       { return &s.Pool }
func (s *Stocks) component() string          { return stocksComponent }
func (s *Stocks) componentMetrics() *metrics { return &s.metrics }
func (s *Stocks) timeout() time.Duration     { return s.t }
func (s *Stocks) subject(w *protoStream.Bar) string {
	return stocksComponent + fmt.Sprintf(".%s", w.Symbol)
}
func (s *Stocks) toWire(b stream.Bar) (*protoStream.Bar, error) {
	return &protoStream.Bar{
		Symbol:     b.Symbol,
		Open:       b.Open,
		High:       b.High,
		Low:        b.Low,
		Close:      b.Close,
		Volume:     b.Volume,
		Timestamp:  timestamppb.New(b.Timestamp),
		TradeCount: b.TradeCount,
		Vwap:       b.VWAP,
	}, nil
}

func (c *Client) Stocks(d *Stream) (*Stocks, error) {
	if d == nil {
		return nil, fmt.Errorf("missing stream opts")
	}

	switch d.Feed {
	case "sip", "iex", "otc", "delayed_sip":
	default:
		c.l.Error("invalid feed", "feed", d.Feed)
		return nil, fmt.Errorf("invalid feed %s", d.Feed)
	}

	s := &Stocks{
		client:     c,
		metrics:    newMetrics("stocks"),
		symbolList: newSymbolList(d.Symbols...),
		t:          d.Timeout,
		Pool:       sync.Pool{New: func() any { return make([]byte, d.BufSize) }},
	}

	so := []stream.StockOption{}
	for _, v := range c.streamOpts(d) {
		so = append(so, v)
	}

	symbols := s.symbolList.list()
	c.l.Info("creating stocks stream client", "conf", d)
	s.s = stream.NewStocksClient(d.Feed, append(so, stream.WithBars(s.handler, symbols...))...)
	return s, nil
}

func (s *Stocks) handler(b stream.Bar) { wrap(s.client, s, b) }

func (c *Stocks) AddSubscriptions(x ...string) error {
	if len(x) == 0 {
		return nil
	}

	c.add(x...)
	l := c.list()
	if err := c.s.SubscribeToBars(c.handler, l...); err != nil {
		c.client.l.Error("failed adding symbol", "want", x, "have", l, "err", err)
		return err
	}

	c.client.l.Info("added to subscription", "added", x, "had", l)
	return nil
}

func (c *Stocks) ListSubscriptions() []string { return c.list() }

func (c *Stocks) DeleteSubscriptions(x ...string) error {
	if len(x) == 0 {
		return nil
	}

	c.del(x...)
	l := c.list()
	if err := c.s.SubscribeToBars(c.handler, l...); err != nil {
		c.client.l.Error("failed deleting new subscriptions from stocks stream", "err", err, "wanted", x)
		return err
	}

	c.client.l.Info("removed subscriptions from bars", "delta", x, "final", l)
	return nil
}

// Begin consuming data. Cancel context to initiate a shutdown?
// Unsure the underlying implementation, doesnt say in the alpaca docs
func (c *Stocks) Stream(ctx context.Context) error {
	if err := c.s.Connect(ctx); err != nil {
		c.client.l.ErrorContext(ctx, "failed establishing stocks connection", "err", err)
		return err
	}

	if err := <-c.s.Terminated(); err != nil {
		c.client.l.Error("connection terminated with error", "err", err)
		return err
	}

	c.client.l.Warn("connection terminated gracefully")
	return nil
}
