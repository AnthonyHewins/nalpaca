package streaming

import (
	"fmt"

	protoStream "github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ transmitter[stream.OptionTrade, *protoStream.OptionTrade] = (*optionTrades)(nil)

type optionTrades struct {
	baseClient[stream.OptionTrade]
	client *stream.OptionClient
}

func newOptionTrades(client *ClientFactory, c *StreamTypeConfig, s *stream.OptionClient) *optionTrades {
	return &optionTrades{
		baseClient: newBaseClient[stream.OptionTrade](
			SubscriptionOptionTrades,
			client,
			c.Symbols,
			int(c.BufSize),
			c.Timeout,
		),
		client: s,
	}
}

func (t *optionTrades) subject(w *protoStream.OptionTrade) string {
	return fmt.Sprintf("%s.%s", optionsTradesSubject, w.Symbol)
}

func (t *optionTrades) toWire(x stream.OptionTrade) *protoStream.OptionTrade {
	return &protoStream.OptionTrade{
		Symbol:    x.Symbol,
		Exchange:  x.Exchange,
		Price:     x.Price,
		Size:      x.Size,
		Timestamp: timestamppb.New(x.Timestamp),
		Condition: x.Condition,
	}
}

func (t *optionTrades) Unsubscribe(x ...string) error {
	return t.rmSubscription(t.client.UnsubscribeFromTrades, x...)
}

func (t *optionTrades) Subscribe(x ...string) error {
	return t.addSubscription(t.client.SubscribeToTrades, t.handler, x...)
}

func (t *optionTrades) handler(x stream.OptionTrade) { wrap(t.ClientFactory, t, x) }
