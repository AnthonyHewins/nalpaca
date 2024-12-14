package portfolio

import (
	"log/slog"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/nats-io/nats.go/jetstream"
)

type Controller struct {
	logger  *slog.Logger
	client  nalpaca.Interface
	timeout time.Duration
	js      jetstream.JetStream

	portfolioKV jetstream.KeyValue

	topicPrefix string
}

func NewController(
	logger *slog.Logger,
	client nalpaca.Interface,
	timeout time.Duration,
	nc jetstream.JetStream,
	portfolioKV jetstream.KeyValue,
	topicPrefix string,
) *Controller {
	return &Controller{
		logger:      logger,
		client:      client,
		timeout:     timeout,
		js:          nc,
		topicPrefix: topicPrefix,
	}
}
