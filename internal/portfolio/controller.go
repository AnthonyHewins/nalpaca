package portfolio

import (
	"log/slog"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/bridge"
	"github.com/nats-io/nats.go/jetstream"
)

type Controller struct {
	logger  *slog.Logger
	client  bridge.AlpacaInterface
	timeout time.Duration
	js      jetstream.JetStream
	prefix  string

	portfolioKV jetstream.KeyValue
}

func NewController(
	logger *slog.Logger,
	client bridge.AlpacaInterface,
	timeout time.Duration,
	nc jetstream.JetStream,
	portfolioKV jetstream.KeyValue,
	prefix string,
) *Controller {
	return &Controller{
		logger:      logger,
		client:      client,
		timeout:     timeout,
		js:          nc,
		portfolioKV: portfolioKV,
		prefix:      prefix,
	}
}
