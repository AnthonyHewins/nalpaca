package nalpaca

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/nats-io/nats.go/jetstream"
)

type Mock struct {
	PushOrderFn func(context.Context, string, *tradesvc.Trade) (*jetstream.PubAck, error)
}

func (m Mock) PushTrade(ctx context.Context, idemKey string, order *tradesvc.Trade) (*jetstream.PubAck, error) {
	return m.PushOrderFn(ctx, idemKey, order)
}
