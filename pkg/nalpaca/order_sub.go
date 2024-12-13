package nalpaca

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

const (
	TradeUpdaterStream   = "nalpaca-tradeupdater-stream-v0"
	TradeUpdaterConsumer = "nalpaca-tradeupdater-consumer-v0"
)

// Simple wrapper creating the tradeupdater consumer.
// Creates consumer with the correct config already set
func (c *Client) TradeUpdaterConsumer(ctx context.Context) (jetstream.Consumer, error) {
	return c.nc.Consumer(ctx, TradeUpdaterStream, TradeUpdaterConsumer)
}
