package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/internal/canceler"
	"github.com/nats-io/nats.go/jetstream"
)

type CancelConf struct {
	EnableCancel   bool   `env:"ENABLE_CANCELER" envDefault:"false"`
	CancelStream   string `env:"CANCEL_STREAM" envDefault:"nalpaca-cancel-stream-v0"`
	CancelConsumer string `env:"CANCEL_CONSUMER" envDefault:"nalpaca-cancel-consumer-v0"`
}

func (a *app) initCanceler(ctx context.Context, js jetstream.JetStream, c *config) error {
	if !c.CancelConf.EnableCancel {
		return nil
	}

	var err error
	a.cancel.ingestor, err = a.connect(ctx, js, c.CancelStream, c.CancelConsumer)
	if err != nil {
		return err
	}

	a.canceler = canceler.New(a.Logger, a.Nalpaca, cancelCounters, c.ProcessingTimeout)
	return nil
}
