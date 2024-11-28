package main

import (
	"context"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go/jetstream"
)

type app struct {
	*conf.Bootstrapper
	trader *trader.Controller
	order  consumer
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

	if err = a.connectConsumers(ctx, &c); err != nil {
		return nil, err
	}

	a.trader = trader.NewController(
		a.TP.Tracer("trader"),
		a.Logger,
		a.Nalpaca,
		c.ProcessingTimeout,
		uint(c.CacheSize),
	)

	return &a, nil
}

func (a *app) connectConsumers(ctx context.Context, c *config) error {
	js, err := jetstream.New(a.NC)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed connecting to jetstream", "err", err)
		return err
	}

	l := a.Logger.With(
		"stream", c.OrderConsumerStream,
		"consumer", c.OrderConsumerName,
	)

	a.order.ingestor, err = js.Consumer(ctx, c.OrderConsumerStream, c.OrderConsumerName)
	if err != nil {
		l.ErrorContext(ctx, "failed connecting to order consumer", "err", err)
		return err
	}

	l.InfoContext(ctx, "connected NATS consumer")
	return err
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	a.Bootstrapper.Shutdown(ctx)

	if c := a.order.ctx; c != nil {
		c.Drain()
		c.Stop()
	}
}
