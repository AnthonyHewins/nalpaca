package main

import (
	"context"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/internal/tradeupdater"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go/jetstream"
)

type app struct {
	*conf.Bootstrapper
	tradeupdater *tradeupdater.Controller
}

type consumer struct {
	ctx      jetstream.ConsumeContext
	ingestor jetstream.Consumer
}

func newApp(ctx context.Context) (*app, error) {
	var c config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	b, err := c.BootstrapConf.New(ctx, appName)
	if err != nil {
		return nil, err
	}

	a := app{Bootstrapper: b}
	defer func() {
		if err != nil {
			a.shutdown()
		}
	}()

	a.tradeupdater = tradeupdater.NewController(
		a.Logger,
		a.Nalpaca,
		c.ProcessingTimeout,
		a.NC,
		c.StreamPrefix,
	)

	return &a, nil
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	a.Bootstrapper.Shutdown(ctx)
}
