package nalpaca

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/nats-io/nats.go/jetstream"
)

type Interface interface {
	PushTrade(ctx context.Context, order *tradesvc.Trade, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
	Cancel(ctx context.Context, orderID string, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}
