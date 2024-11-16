package trader

import (
	"log/slog"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	tracer            trace.Tracer
	logger            *slog.Logger
	processingTimeout time.Duration
	alpaca            *alpaca.Client
}

func NewController(tracer trace.Tracer, logger *slog.Logger, client *alpaca.Client, processingTimeout time.Duration) *Controller {
	return &Controller{
		tracer:            tracer,
		logger:            logger,
		processingTimeout: processingTimeout,
		alpaca:            client,
	}
}
