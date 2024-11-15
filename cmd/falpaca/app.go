package falpaca

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AnthonyHewins/falpaca/internal/conf"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/sdk/trace"
)

type app struct {
	logger  *slog.Logger
	nc      *nats.Conn
	health  *conf.HealthServer
	metrics *http.Server
	tp      *trace.TracerProvider
}

func newApp() (*app, error) {
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

	return &a, nil
}

func (a *app) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if a.nc != nil {
		a.nc.Close()
	}

	if a.metrics != nil {
		a.metrics.Close()
	}

	if a.health != nil {
		a.health.GracefulStop()
	}

	if a.tp != nil {
		a.tp.Shutdown(ctx)
	}
}
