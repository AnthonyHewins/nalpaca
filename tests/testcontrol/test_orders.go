package testcontrol

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/google/uuid"
)

func (c *Controller) TestOrders(ctx context.Context) error {
	trades := []tradesvc.Trade{
		{
			Symbol:         "AAPL",
			Qty:            "1",
			Side:           tradesvc.Side_SIDE_BUY,
			OrderType:      tradesvc.OrderType_ORDER_TYPE_MARKET,
			Tif:            tradesvc.TimeInForce_TIME_IN_FORCE_DAY,
			PositionIntent: tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN,
		},
	}

	for i := range trades {
		if _, err := c.client.PushTrade(ctx, uuid.New().String(), &trades[i]); err != nil {
			return err
		}
	}

	return nil
}
