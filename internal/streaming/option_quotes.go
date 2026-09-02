package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ transmitter[stream.OptionQuote, *protoStream.OptionQuote] = (*optionQuotes)(nil)

type optionQuotes struct {
	baseClient[stream.OptionQuote]
	client *stream.OptionClient
}

func newOptionQuotes(client *ClientFactory, c *StreamTypeConfig, s *stream.OptionClient) *optionQuotes {
	return &optionQuotes{
		baseClient: newBaseClient[stream.OptionQuote](
			SubscriptionOptionQuotes,
			client,
			c.Symbols,
			int(c.BufSize),
			c.Timeout,
		),
		client: s,
	}
}

func (q *optionQuotes) subject(w *protoStream.OptionQuote) string {
	return fmt.Sprintf("%s.%s", optionsQuotesSubject, w.Symbol)
}

func (q *optionQuotes) toWire(x stream.OptionQuote) *protoStream.OptionQuote {
	return &protoStream.OptionQuote{
		Symbol:      x.Symbol,
		BidExchange: x.BidExchange,
		BidPrice:    x.BidPrice,
		BidSize:     x.BidSize,
		AskExchange: x.AskExchange,
		AskPrice:    x.AskPrice,
		AskSize:     x.AskSize,
		Timestamp:   timestamppb.New(x.Timestamp),
		Condition:   x.Condition,
	}
}

func (q *optionQuotes) Unsubscribe(x ...string) error {
	return q.rmSubscription(q.client.UnsubscribeFromQuotes, x...)
}

func (q *optionQuotes) Subscribe(x ...string) error {
	return q.addSubscription(q.client.SubscribeToQuotes, q.handler, x...)
}

func (q *optionQuotes) handler(x stream.OptionQuote) { wrap(q.ClientFactory, q, x) }
