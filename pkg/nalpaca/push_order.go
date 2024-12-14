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

func (c *Client) PushTrade(ctx context.Context, idemKey string, order *tradesvc.Trade, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	if len(idemKey) > 128 {
		return nil, fmt.Errorf("invalid idempotent order ID: %s. Must be under 128 chars", idemKey)
	}

	if order == nil {
		return nil, ErrMissingOrder
	}

	if order.Symbol == "" {
		return nil, ErrMissingSymbol
	}

	buf, err := proto.Marshal(order)
	if err != nil {
		return nil, err
	}

	return c.nc.Publish(ctx, c.prefix+".orders.v0."+idemKey, buf, opts...)
}
