package nalpaca

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *Client) Cancel(ctx context.Context, idemKey string, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	if len(idemKey) > 128 {
		return nil, fmt.Errorf("invalid idempotent order ID: %s. Must be under 128 chars", idemKey)
	}

	return c.nc.Publish(ctx, fmt.Sprintf("%s.cancel.v0.%s", c.prefix, idemKey), nil, opts...)
}
