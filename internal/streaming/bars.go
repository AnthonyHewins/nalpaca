package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ Transmitter[stream.Bar, *protoStream.Bar] = (*bars)(nil)

type bars struct {
	baseClient[stream.Bar]
	client *stream.StocksClient
}

func newBars(client *Client, c *StreamTypeConfig, s *stream.StocksClient) *bars {
	return &bars{
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

func (b *bars) subject(w *protoStream.Bar) string {
	return fmt.Sprintf("%s.%s", SubscriptionStockBars, w.Symbol)
}

func (x *bars) toWire(b stream.Bar) (*protoStream.Bar, error) {
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

func (b *bars) List() []string { return b.baseClient.list() }

func (b *bars) Unsubscribe(x ...string) error {
	return b.rmSubscription(b.client.UnsubscribeFromBars, x...)
}

func (b *bars) Subscribe(x ...string) error {
	return b.addSubscription(b.client.SubscribeToBars, b.handler, x...)
}

func (b *bars) handler(x stream.Bar) { wrap(b.Client, b, x) }
