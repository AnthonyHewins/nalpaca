package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ Transmitter[stream.Trade, *protoStream.Trade] = (*trades)(nil)

type trades struct {
	baseClient[stream.Trade]
	client *stream.StocksClient
}

func newTrades(client *Client, c *StreamTypeConfig, s *stream.StocksClient) *trades {
	return &trades{
		baseClient: newBaseClient[stream.Trade](
			SubscriptionStockTrades,
			client,
			c.Symbols,
			int(c.BufSize),
			c.Timeout,
		),
		client: s,
	}
}

func (t *trades) subject(w *protoStream.Trade) string {
	return fmt.Sprintf("%s.%s", SubscriptionStockTrades, w.Symbol)
}

func (t *trades) toWire(x stream.Trade) (*protoStream.Trade, error) {
	return &protoStream.Trade{
		Id:         x.ID,
		Symbol:     x.Symbol,
		Exchange:   x.Exchange,
		Price:      x.Price,
		Size:       x.Size,
		Timestamp:  timestamppb.New(x.Timestamp),
		Conditions: x.Conditions,
		Tape:       x.Tape,
	}, nil
}

func (t *trades) List() []string { return t.baseClient.list() }

func (t *trades) Unsubscribe(x ...string) error {
	return t.rmSubscription(t.client.UnsubscribeFromTrades, x...)
}

func (t *trades) Subscribe(x ...string) error {
	return t.addSubscription(t.client.SubscribeToTrades, t.handler, x...)
}

func (t *trades) handler(x stream.Trade) { wrap(t.Client, t, x) }
