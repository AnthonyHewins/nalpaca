package portfolio

import (
	"log/slog"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/nats-io/nats.go/jetstream"
)

type Controller struct {
	logger  *slog.Logger
	nalpaca nalpaca.Interface
	kv      jetstream.KeyValue
}

func NewController(logger *slog.Logger, kv jetstream.KeyValue) *Controller {
	return &Controller{
		logger: logger,
		kv:     kv,
	}
}
