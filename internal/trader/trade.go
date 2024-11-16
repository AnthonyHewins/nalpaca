package trader

import (
	"context"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func (c *Controller) trade(ctx context.Context, o alpaca.PlaceOrderRequest) error {
	l := c.logger.With("config", o)

	order, err := c.alpaca.PlaceOrder(o)
	if err != nil {
		l.ErrorContext(ctx, "failed placing order", "err", err)
		return err
	}

	l.DebugContext(ctx, "order created", "order", order)
	return nil
}
