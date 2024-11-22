package tradeupdater

import (
	"context"
	"errors"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func (c *Controller) EventLoop(ctx context.Context) error {
	for {
		var lastMessage time.Time
		for {
			req := alpaca.StreamTradeUpdatesRequest{}
			if !lastMessage.IsZero() {
				req.Since = lastMessage.Add(time.Nanosecond)
			}

			err := c.client.StreamTradeUpdates(ctx, func(tu alpaca.TradeUpdate) {
				lastMessage = tu.At
				c.handler(tu)
			}, req)

			if err == nil {
				return nil
			}

			if errors.Is(err, context.Canceled) {
				c.logger.Warn("tradeupdater: ctx canceled", "ctx-err", ctx.Err())
				return err
			}

			c.logger.ErrorContext(ctx, "trade updater error", "err", err)
			return err
		}
	}
}

func (c *Controller) handler(u alpaca.TradeUpdate) {

}
