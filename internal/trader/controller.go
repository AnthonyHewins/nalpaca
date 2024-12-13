package trader

import (
	"log/slog"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	tracer            trace.Tracer
	logger            *slog.Logger
	processingTimeout time.Duration
	alpaca            nalpaca.Interface
}

func NewController(
	tracer trace.Tracer,
	logger *slog.Logger,
	client nalpaca.Interface,
	processingTimeout time.Duration,
	cache uint,
) *Controller {
	return &Controller{
		tracer:            tracer,
		logger:            logger,
		processingTimeout: processingTimeout,
		alpaca:            client,
	}
}
