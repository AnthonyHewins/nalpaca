package main

import (
	"context"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	tradesCanceled = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nalpaca",
		Subsystem: appName,
		Name:      "canceled",
		Help:      "The number of trades that were canceled",
	})

	tradeCancelErrs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nalpaca",
		Subsystem: appName,
		Name:      "order_cancel_errs",
		Help:      "Number of errors encountered canceling orders",
	})

	cancelAllCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nalpaca",
		Subsystem: appName,
		Name:      "cancel_all_count",
		Help:      "The number of times a 'cancel all' was executed",
	})

	cancelAllFails = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nalpaca",
		Subsystem: appName,
		Name:      "cancel_all_errs",
		Help:      "The number of times a 'cancel all' failed",
	})
)

type app struct {
	*conf.Bootstrapper

	canceler canceler
	trader   *trader.Controller

	order  consumer
	cancel consumer
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

	a.canceler = canceler{
		logger:         a.Logger,
		cancelCount:    tradesCanceled,
		cancelFail:     tradeCancelErrs,
		cancelAllCount: cancelAllCounter,
		cancelAllFail:  cancelAllFails,
		client:         a.Nalpaca,
		timeout:        c.ProcessingTimeout,
	}

	return &a, nil
}

func (a *app) connectConsumers(ctx context.Context, c *config) error {
	js, err := jetstream.New(a.NC)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed connecting to jetstream", "err", err)
		return err
	}

	a.order.ingestor, err = a.connect(ctx, js, c.OrderConsumerStream, c.OrderConsumerName)
	if err != nil {
		return err
	}

	a.cancel.ingestor, err = a.connect(ctx, js, c.CancelStream, c.CancelConsumer)
	if err != nil {
		return err
	}

	a.Logger.InfoContext(ctx, "connected all NATS consumers")
	return err
}

func (a *app) connect(ctx context.Context, js jetstream.JetStream, stream, consumer string) (jetstream.Consumer, error) {
	x, err := js.Consumer(ctx, stream, consumer)
	if err != nil {
		a.Logger.ErrorContext(ctx,
			"failed connecting to cancel order",
			"err", err,
			"stream", stream,
			"consumer", consumer,
		)
	}

	return x, err
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	a.Bootstrapper.Shutdown(ctx)

	type consumers struct {
		name     string
		consumer jetstream.ConsumeContext
	}

	for _, v := range []consumers{
		{name: "order consumer", consumer: a.order.ctx},
		{name: "cancel consumer", consumer: a.cancel.ctx},
	} {
		if v.consumer == nil {
			continue
		}

		a.Logger.InfoContext(ctx, "shutting down "+v.name)
		v.consumer.Drain()
		v.consumer.Stop()
		a.Logger.InfoContext(ctx, "shut down "+v.name)
	}
}
