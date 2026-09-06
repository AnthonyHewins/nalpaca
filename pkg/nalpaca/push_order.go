package nalpaca

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

var (
	ErrMissingSymbol = errors.New("missing symbol")
	ErrMissingOrder  = errors.New("no order passed")
)

// PushTrade asks nalpaca to place an order. order.Id is the idempotency
// key: it becomes the client order ID Alpaca sees (trader reads it straight
// off the payload) and is also set as the NATS dedup ID via
// [jetstream.WithMsgID], so a retried publish with the same order.Id won't
// create a second order. The subject is a single literal endpoint rather
// than being keyed by order.Id or symbol — both already travel in the
// payload, so repeating them in the subject would be redundant, and for
// symbol it would collide with real per-symbol trade prints on
// nalpaca.stocks.trades.<TICKER>.
func (c *Client) PushTrade(ctx context.Context, order *tradesvc.Trade, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	if order == nil {
		return nil, ErrMissingOrder
	}

	if len(order.Id) > 128 {
		return nil, fmt.Errorf("invalid idempotent order ID: %s. Must be under 128 chars", order.Id)
	}

	if order.Symbol == "" {
		return nil, ErrMissingSymbol
	}

	buf, err := proto.Marshal(order)
	if err != nil {
		return nil, err
	}

	opts = append(opts, jetstream.WithMsgID(order.Id))
	return c.nc.Publish(ctx, c.prefix+".stocks.orders.create", buf, opts...)
}
