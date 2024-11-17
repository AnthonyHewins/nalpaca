package trader

import (
	"log/slog"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	lru "github.com/hashicorp/golang-lru/v2"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	tracer            trace.Tracer
	logger            *slog.Logger
	processingTimeout time.Duration
	alpaca            *alpaca.Client

	clientIDCache *lru.Cache[string, struct{}]
}

func NewController(
	tracer trace.Tracer,
	logger *slog.Logger,
	client *alpaca.Client,
	processingTimeout time.Duration,
	cache uint,
) *Controller {
	c, _ := lru.New[string, struct{}](int(cache))
	return &Controller{
		tracer:            tracer,
		logger:            logger,
		processingTimeout: processingTimeout,
		alpaca:            client,
		clientIDCache:     c,
	}
}
