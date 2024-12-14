package testcontrol

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/google/uuid"
)

func (c *Controller) TestCancel(ctx context.Context) error {
	id := uuid.NewString()
	_, err := c.client.PushTrade(ctx, id, &tradesvc.Trade{
		Symbol:         "AAPL",
		Qty:            "1",
		Side:           tradesvc.Side_SIDE_BUY,
		OrderType:      tradesvc.OrderType_ORDER_TYPE_TRAILING_STOP,
		PositionIntent: tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN,
		Tif:            tradesvc.TimeInForce_TIME_IN_FORCE_DAY,
		TrailPercent:   "5",
	})

	if err != nil {
		return err
	}

	if _, err = c.client.Cancel(ctx, id); err != nil {
		return err
	}

	return nil
}
