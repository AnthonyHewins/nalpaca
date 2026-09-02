package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ transmitter[stream.Quote, *protoStream.Quote] = (*Quotes)(nil)

type Quotes struct {
	baseClient[stream.Quote]
	client *stream.StocksClient
}

func newQuotes(client *ClientFactory, c *StreamTypeConfig, s *stream.StocksClient) *Quotes {
	return &Quotes{
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

func (q *Quotes) subject(w *protoStream.Quote) string {
	return fmt.Sprintf("%s.%s", stocksQuotesSubject, w.Symbol)
}

func (q *Quotes) toWire(x stream.Quote) *protoStream.Quote {
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
	}
}

func (q *Quotes) Unsubscribe(x ...string) error {
	return q.rmSubscription(q.client.UnsubscribeFromQuotes, x...)
}

func (q *Quotes) Subscribe(x ...string) error {
	return q.addSubscription(q.client.SubscribeToQuotes, q.handler, x...)
}

func (q *Quotes) handler(x stream.Quote) { wrap(q.ClientFactory, q, x) }
