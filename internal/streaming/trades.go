package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ transmitter[stream.Trade, *protoStream.Trade] = (*Trades)(nil)

type Trades struct {
	baseClient[stream.Trade]
	client *stream.StocksClient
}

func newTrades(client *ClientFactory, c *StreamTypeConfig, s *stream.StocksClient) *Trades {
	return &Trades{
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

func (t *Trades) subject(w *protoStream.Trade) string {
	return fmt.Sprintf("%s.%s", stocksTradesSubject, w.Symbol)
}

func (t *Trades) toWire(x stream.Trade) *protoStream.Trade {
	return &protoStream.Trade{
		Id:         x.ID,
		Symbol:     x.Symbol,
		Exchange:   x.Exchange,
		Price:      x.Price,
		Size:       x.Size,
		Timestamp:  timestamppb.New(x.Timestamp),
		Conditions: x.Conditions,
		Tape:       x.Tape,
	}
}

func (t *Trades) Unsubscribe(x ...string) error {
	return t.rmSubscription(t.client.UnsubscribeFromTrades, x...)
}

func (t *Trades) Subscribe(x ...string) error {
	return t.addSubscription(t.client.SubscribeToTrades, t.handler, x...)
}

func (t *Trades) handler(x stream.Trade) { wrap(t.ClientFactory, t, x) }
