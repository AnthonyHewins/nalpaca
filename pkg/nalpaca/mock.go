package nalpaca

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/nats-io/nats.go/jetstream"
)

type Mock struct {
	PushOrderFn func(context.Context, *tradesvc.Trade, ...jetstream.PublishOpt) (*jetstream.PubAck, error)
	CancelFn    func(context.Context, string, ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

func (m Mock) PushTrade(ctx context.Context, order *tradesvc.Trade, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return m.PushOrderFn(ctx, order)
}

func (m Mock) Cancel(ctx context.Context, orderID string, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return m.CancelFn(ctx, orderID, opts...)
}
