package falpaca

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/AnthonyHewins/falpaca/internal/conf"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/caarlos0/env/v11"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/sdk/trace"
)

type app struct {
	alpacaClient *alpaca.Client
	logger       *slog.Logger
	nc           *nats.Conn
	health       *conf.HealthServer
	metrics      *http.Server
	tp           *trace.TracerProvider
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
		logger:       b.Logger,
		alpacaClient: b.Alpaca(&c.Alpaca, &http.Client{Timeout: c.HttpClientTimeout}),
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

func (a *app) initStream(ctx context.Context, c *config) {
	js, err := jetstream.New(a.nc)
	if err != nil {
		a.logger.Error("failed connecting to jetstream", "err", err)
		return nil, err
	}

	subject := "falpaca.orders.>"
	if c.StreamPrefix != "" {
		subject = c.StreamPrefix + "." + subject
	}

	consumer, err := js.CreateConsumer(ctx, subject, jetstream.ConsumerConfig{
		Name: "falpaca-order-consumer-v0",
		// Durable:            "",
		Description:   "Falpaca order consumer",
		MaxDeliver:    3,
		BackOff:       c.StreamBackoff,
		FilterSubject: subject,
		// FilterSubjects:     []string{},
	})

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
