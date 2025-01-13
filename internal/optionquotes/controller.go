package optionquotes

import (
	"log/slog"

	"github.com/AnthonyHewins/nalpaca/internal/bridge"
	"github.com/nats-io/nats.go/jetstream"
)

type Controller struct {
	logger *slog.Logger
	np     bridge.AlpacaInterface
	js     jetstream.JetStream
}

func NewController(logger *slog.Logger, n bridge.AlpacaInterface, js jetstream.JetStream) *Controller {
	return &Controller{
		logger: logger,
		np:     n,
		js:     js,
	}
}
