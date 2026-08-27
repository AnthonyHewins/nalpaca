package nalpaca

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// These must match what scripts/nats provisions. Account events live on their
// own stream, separate from the data stream: trade updates are things that
// happened to your orders, not public market data.
const (
	TradeUpdaterStream   = "nalpaca-account-stream"
	TradeUpdaterConsumer = "nalpaca-tradeupdate-consumer"
)

// Simple wrapper creating the tradeupdater consumer.
// Creates consumer with the correct config already set
func (c *Client) TradeUpdaterConsumer(ctx context.Context) (jetstream.Consumer, error) {
	return c.nc.Consumer(ctx, TradeUpdaterStream, TradeUpdaterConsumer)
}
