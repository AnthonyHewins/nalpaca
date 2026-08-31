package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ Transmitter[stream.Quote, *protoStream.Quote] = (*quotes)(nil)

type quotes struct {
	baseClient[stream.Quote]
	client *stream.StocksClient
}

func newQuotes(client *Client, c *StreamTypeConfig, s *stream.StocksClient) *quotes {
	return &quotes{
		baseClient: newBaseClient[stream.Quote](
			SubscriptionStockQuotes,
			client,
			c.Symbols,
			int(c.BufSize),
			c.Timeout,
		),
		client: s,
	}
}

func (q *quotes) subject(w *protoStream.Quote) string {
	return fmt.Sprintf("%s.%s", SubscriptionStockQuotes, w.Symbol)
}

func (q *quotes) toWire(x stream.Quote) (*protoStream.Quote, error) {
	return &protoStream.Quote{
		Symbol:      x.Symbol,
		BidExchange: x.BidExchange,
		BidPrice:    x.BidPrice,
		BidSize:     x.BidSize,
		AskExchange: x.AskExchange,
		AskPrice:    x.AskPrice,
		AskSize:     x.AskSize,
		Timestamp:   timestamppb.New(x.Timestamp),
		Conditions:  x.Conditions,
		Tape:        x.Tape,
	}, nil
}

func (q *quotes) List() []string { return q.baseClient.list() }

func (q *quotes) Unsubscribe(x ...string) error {
	return q.rmSubscription(q.client.UnsubscribeFromQuotes, x...)
}

func (q *quotes) Subscribe(x ...string) error {
	return q.addSubscription(q.client.SubscribeToQuotes, q.handler, x...)
}

func (q *quotes) handler(x stream.Quote) { wrap(q.Client, q, x) }
