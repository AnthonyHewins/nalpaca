package testcontrol

import (
	"log/slog"

	"github.com/AnthonyHewins/nalpaca/pkg/nalpaca"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Controller struct {
	nc     *nats.Conn
	client *nalpaca.Client
	logger *slog.Logger
}

func NewController(logger *slog.Logger, nc *nats.Conn) (*Controller, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		logger.Error("failed connecting to jetstream", "err", err)
		return nil, err
	}

	return &Controller{
		nc:     nc,
		client: nalpaca.NewClient(js, "nalpaca"),
		logger: logger,
	}, nil
}

func (c *Controller) Shutdown() {
	if c.nc != nil {
		c.nc.Close()
	}
}
