package trader

import (
	"log/slog"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
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
	alpaca            nalpaca.Interface
}

func NewController(
	tracer trace.Tracer,
	logger *slog.Logger,
	counters Counters,
	client nalpaca.Interface,
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
