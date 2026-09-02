package conf

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/bridge"
	"github.com/AnthonyHewins/nalpaca/internal/metrics"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type BootstrapConf struct {
	Logger  Logger
	Metrics metrics.PromConfig
	Health  Health
	Tracer  Tracer
	NATS    NATS
	Alpaca  Alpaca

	HTTPClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"15s" desc:"Timeout for the shared HTTP client used to talk to the Alpaca REST API"`
}

type Server struct {
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`
	Logger          *slog.Logger
	NC              *nats.Conn
	Health          *HealthServer
	Metrics         metrics.Prom
	TP              *trace.TracerProvider
	Nalpaca         bridge.AlpacaInterface
	Stocks          *stream.StocksClient
}

type Bootstrapper Server

func (b *BootstrapConf) New(ctx context.Context, appName string, metrics ...prometheus.Collector) (*Bootstrapper, error) {
	logger, err := b.Logger.Slog()
	if err != nil {
		return nil, err
	}

	a := &Bootstrapper{Logger: logger}

	defer func() {
		if err != nil {
			(*Server)(a).Shutdown(ctx)
		}
	}()

	if a.NC, err = a.NATSConn(&b.NATS); err != nil {
		return nil, err
	}

	if a.Metrics, err = a.PrometheusHTTP(&b.Metrics, metrics...); err != nil {
		return nil, err
	}

	if a.TP, err = a.Tracer(appName, &b.Tracer); err != nil {
		return nil, err
	}

	a.Nalpaca, err = a.Alpaca(&b.Alpaca, &http.Client{Timeout: b.HTTPClientTimeout})
	if err != nil {
		return nil, err
	}

	a.Health = a.HealthServer(&b.Health)
	return a, nil
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.Metrics.Server != nil {
		s.Logger.InfoContext(ctx, "shutting down metrics")
		if err := s.Metrics.Server.Close(); err != nil {
			s.Logger.ErrorContext(ctx, "failed shutting metrics down", "err", err)
		}
	}

	if s.TP != nil {
		s.Logger.InfoContext(ctx, "shutting down tracers")
		if err := s.TP.Shutdown(ctx); err != nil {
			s.Logger.ErrorContext(ctx, "failed shutting down tracers", "err", err)
		}
	}

	if s.NC != nil {
		s.Logger.InfoContext(ctx, "closing NATS")
		s.NC.Close()
		s.Logger.Info("closed nats conn")
	}

	if s.Health != nil {
		s.Logger.InfoContext(ctx, "shutting down health")
		s.Health.GracefulStop()
		s.Logger.Info("health server shut down")
	}
}
