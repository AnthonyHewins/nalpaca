package tradeupdater

import (
	"log/slog"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Controller struct {
	logger *slog.Logger
	client *alpaca.Client
}

func NewController(logger *slog.Logger, client *alpaca.Client) *Controller {
	return &Controller{
		logger: logger,
		client: client,
	}
}
