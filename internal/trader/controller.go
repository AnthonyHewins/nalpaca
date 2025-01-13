package trader

import (
	"log/slog"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/bridge"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
)

type Counters struct {
	OrderCreatedCount prometheus.Counter
	OrderFailCount    prometheus.Counter
}

type Controller struct {
	Counters
	tracer            trace.Tracer
	logger            *slog.Logger
	processingTimeout time.Duration
	alpaca            bridge.AlpacaInterface
}

func NewController(
	tracer trace.Tracer,
	logger *slog.Logger,
	counters Counters,
	client bridge.AlpacaInterface,
	processingTimeout time.Duration,
) *Controller {
	return &Controller{
		Counters:          counters,
		tracer:            tracer,
		logger:            logger,
		processingTimeout: processingTimeout,
		alpaca:            client,
	}
}
