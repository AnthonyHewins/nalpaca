package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/conf"
	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/sdk/trace"
)

type app struct {
	trader  *trader.Controller
	logger  *slog.Logger
	nc      *nats.Conn
	health  *conf.HealthServer
	metrics *http.Server
	tp      *trace.TracerProvider

	order        consumer
	tradeupdates jetstream.Stream
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

	b, err := conf.NewBootstrapper(&c.Logger)
	if err != nil {
		return nil, err
	}

	a := app{
		logger: b.Logger,
	}

	defer func() {
		if err != nil {
			a.shutdown()
		}
	}()

	a.nc, err = b.NATSConn(&c.NATS)
	if err != nil {
		return nil, err
	}

	if err = a.connectConsumers(ctx, &c); err != nil {
		return nil, err
	}

	a.health = b.Health(&c.Health)
	if err != nil {
		return nil, err
	}

	a.metrics, err = b.PrometheusHTTP(&c.Metrics)
	if err != nil {
		return nil, err
	}

	a.tp, err = b.Tracer(appName, &c.Tracer)
	if err != nil {
		return nil, err
	}

	a.trader = trader.NewController(
		a.tp.Tracer("trader"),
		a.logger,
		b.Alpaca(&c.Alpaca, &http.Client{Timeout: c.HttpClientTimeout}),
		c.ProcessingTimeout,
		uint(c.CacheSize),
	)

	return &a, nil
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if a.metrics != nil {
		a.metrics.Close()
	}

	if a.health != nil {
		a.health.GracefulStop()
	}

	if a.tp != nil {
		a.tp.Shutdown(ctx)
	}

	if c := a.order.ctx; c != nil {
		c.Drain()
		c.Stop()
	}

	if a.nc != nil {
		a.nc.Close()
	}
}
