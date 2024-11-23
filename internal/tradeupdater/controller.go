package tradeupdater

import (
	"log/slog"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/nats-io/nats.go"
)

type Controller struct {
	logger      *slog.Logger
	client      nalpaca.Interface
	timeout     time.Duration
	nc          *nats.Conn
	topicPrefix string
}

func NewController(
	logger *slog.Logger,
	client *alpaca.Client,
	timeout time.Duration,
	nc *nats.Conn,
	topicPrefix string,
) *Controller {
	return &Controller{
		logger:  logger,
		client:  client,
		timeout: timeout,
		nc:      nc,
	}
}
