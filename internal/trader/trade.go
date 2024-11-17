package trader

import (
	"context"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"go.opentelemetry.io/otel/codes"
)

func (c *Controller) trade(ctx context.Context, o alpaca.PlaceOrderRequest) error {
	ctx, span := c.tracer.Start(ctx, "executing trade "+o.ClientOrderID)
	defer span.End()

	var err error
	defer func() {
		if err == nil {
			span.SetStatus(codes.Ok, "executed")
			return
		}

		span.SetStatus(codes.Error, "failed trade")
		span.RecordError(err)
	}()

	l := c.logger.With("trace-id", span.SpanContext().TraceID(), "config", o)

	order, err := c.alpaca.PlaceOrder(o)
	if err != nil {
		l.ErrorContext(ctx, "failed placing order", "err", err)
		return err
	}

	l.DebugContext(ctx, "order created", "order", order)
	return nil
}
