package main

import (
	"context"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/canceler"
	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/internal/optionquotes"
	"github.com/AnthonyHewins/nalpaca/internal/portfolio"
	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
)

func newCounter(system, name, desc string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: appName,
		Subsystem: system,
		Name:      name,
		Help:      desc,
	})
}

var (
	cancelCounters = canceler.Counters{
		CancelCount:    newCounter("canceler", "cancel_count", "The number of trades that were canceled"),
		CancelFail:     newCounter("canceler", "order_cancel_errs", "Number of errors encountered canceling orders"),
		CancelAllCount: newCounter("canceler", "cancel_all_count", "The number of times a 'cancel all' was executed"),
		CancelAllFail:  newCounter("canceler", "cancel_all_errs", "The number of times a 'cancel all' failed"),
	}

	orderCounters = trader.Counters{
		OrderCreatedCount: newCounter("orders", "orders_created_count", "Number of times an order was created successfully"),
		OrderFailCount:    newCounter("orders", "orders_failed_count", "Number of times an order creation failed"),
	}
)

type app struct {
	*conf.Server

	canceler *canceler.Canceler
	trader   *trader.Controller
	updater  *portfolio.Controller
	quotes   *optionquotes.Controller

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

	a := app{Server: (*conf.Server)(b)}
	defer func() {
		if err != nil {
			a.shutdown()
		}
	}()

	if err = a.connectNATS(ctx, &c); err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *app) connectNATS(ctx context.Context, c *config) error {
	js, err := jetstream.New(a.NC)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed connecting to jetstream", "err", err)
		return err
	}

	kv, err := js.KeyValue(ctx, c.Bucket)
	if err != nil {
		a.Logger.ErrorContext(ctx, "failed connecting to nalpaca KV bucket", "err", err, "bucket", c.Bucket)
		return err
	}

	for _, fn := range []func(context.Context, jetstream.JetStream, *config) error{
		a.initCanceler,
		a.initOrders,
	} {
		if err = fn(ctx, js, c); err != nil {
			return err
		}
	}

	a.initTradeUpdater(ctx, js, kv, c)
	return err
}

func (a *app) connect(ctx context.Context, js jetstream.JetStream, stream, consumer string) (jetstream.Consumer, error) {
	x, err := js.Consumer(ctx, stream, consumer)
	if err != nil {
		a.Logger.ErrorContext(ctx,
			"failed connecting to consumer",
			"err", err,
			"stream", stream,
			"consumer", consumer,
		)
	}

	return x, err
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

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

	a.Server.Shutdown(ctx)
}
