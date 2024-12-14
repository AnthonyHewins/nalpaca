package testcontrol

import (
	"context"
	"log/slog"

	"github.com/AnthonyHewins/nalpaca/pkg/nalpaca"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Controller struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	kv     jetstream.KeyValue
	client *nalpaca.Client
	logger *slog.Logger
}

func NewController(logger *slog.Logger, nc *nats.Conn) (*Controller, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		logger.Error("failed connecting to jetstream", "err", err)
		return nil, err
	}

	kv, err := js.KeyValue(context.Background(), "nalpaca")
	if err != nil {
		logger.Error("failed connecting to kv", "err", err)
		return nil, err
	}

	return &Controller{
		nc:     nc,
		js:     js,
		kv:     kv,
		client: nalpaca.NewClient(js, kv, "nalpaca"),
		logger: logger,
	}, nil
}

func (c *Controller) Shutdown() {
	if c.nc != nil {
		c.nc.Close()
	}
}
