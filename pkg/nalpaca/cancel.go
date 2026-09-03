package nalpaca

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) Cancel(ctx context.Context, orderID string, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	if len(orderID) > 128 {
		return nil, fmt.Errorf("invalid order ID: %s. Must be under 128 chars", orderID)
	}

	return c.nc.Publish(ctx, fmt.Sprintf("%s.orders.cancel", c.prefix), []byte(orderID), opts...)
}
