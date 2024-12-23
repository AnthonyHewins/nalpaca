package conf

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type BootstrapConf struct {
	Logger  Logger
	Metrics Metrics
	Health  Health
	Tracer  Tracer
	NATS    NATS
	Alpaca  Alpaca

	HTTPClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"15s"`
}

type Bootstrapper struct {
	Logger  *slog.Logger
	NC      *nats.Conn
	Health  *HealthServer
	Metrics *http.Server
	TP      *trace.TracerProvider
	Nalpaca nalpaca.Interface
}

func (b *BootstrapConf) New(ctx context.Context, appName string, metrics ...prometheus.Collector) (*Bootstrapper, error) {
	logger, err := b.Logger.Slog()
	if err != nil {
		return nil, err
	}

	a := Bootstrapper{Logger: logger}

	defer func() {
		if err != nil {
			a.Shutdown(ctx)
		}
	}()

	a.NC, err = a.NATSConn(&b.NATS)
	if err != nil {
		return nil, err
	}

	a.Health = a.HealthServer(&b.Health)
	if err != nil {
		return nil, err
	}

	a.Metrics, err = a.PrometheusHTTP(&b.Metrics, metrics...)
	if err != nil {
		return nil, err
	}

	a.TP, err = a.Tracer(appName, &b.Tracer)
	if err != nil {
		return nil, err
	}

	a.Nalpaca, err = a.Alpaca(&b.Alpaca, &http.Client{Timeout: b.HTTPClientTimeout})
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (b *Bootstrapper) Shutdown(ctx context.Context) {
	if b.Metrics != nil {
		b.Metrics.Close()
	}

	if b.Health != nil {
		b.Health.GracefulStop()
	}

	if b.TP != nil {
		b.TP.Shutdown(ctx)
	}

	if b.NC != nil {
		b.NC.Close()
	}
}
