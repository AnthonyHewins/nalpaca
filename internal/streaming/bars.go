package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ transmitter[stream.Bar, *protoStream.Bar] = (*Bars)(nil)

type Bars struct {
	baseClient[stream.Bar]
	client *stream.StocksClient
}

func newBars(client *ClientFactory, c *StreamTypeConfig, s *stream.StocksClient) *Bars {
	return &Bars{
		baseClient: newBaseClient[stream.Bar](
			SubscriptionStockBars,
			client,
			c.Symbols,
			int(c.BufSize),
			c.Timeout,
		),
		client: s,
	}
}

func (b *Bars) subject(w *protoStream.Bar) string {
	return fmt.Sprintf("%s.%s", SubscriptionStockBars, w.Symbol)
}

func (x *Bars) toWire(b stream.Bar) *protoStream.Bar {
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
	}
}

func (b *Bars) Unsubscribe(x ...string) error {
	return b.rmSubscription(b.client.UnsubscribeFromBars, x...)
}
func (b *Bars) Subscribe(x ...string) error {
	return b.addSubscription(b.client.SubscribeToBars, b.handler, x...)
}

func (b *Bars) handler(x stream.Bar) { wrap(b.ClientFactory, b, x) }
