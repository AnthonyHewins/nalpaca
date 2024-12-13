package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/google/uuid"
)

func (c *controller) cancel() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	id := uuid.New().String()
	_, err := c.client.PushTrade(ctx, id, &tradesvc.Trade{
		Symbol:         "AAPL",
		Qty:            "1",
		Side:           tradesvc.Side_SIDE_BUY,
		OrderType:      tradesvc.OrderType_ORDER_TYPE_LIMIT,
		Tif:            tradesvc.TimeInForce_TIME_IN_FORCE_DAY,
		LimitPrice:     "1.0",
		PositionIntent: tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN,
	})

	if err != nil {
		return err
	}

	return nil
}
