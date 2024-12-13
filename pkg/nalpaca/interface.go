package nalpaca

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/nats-io/nats.go/jetstream"
)

type Interface interface {
	PushTrade(ctx context.Context, idemKey string, order *tradesvc.Trade) (*jetstream.PubAck, error)
}
