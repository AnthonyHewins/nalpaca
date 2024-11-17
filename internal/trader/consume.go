package trader

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnthonyHewins/falpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

func (c *Controller) term(ctx context.Context, m jetstream.Msg, reason string) {
	c.logger.DebugContext(ctx, "terminating msg", "reason", reason)
	if err := m.TermWithReason(reason); err != nil {
		c.logger.ErrorContext(ctx, "failed termination", "reason", reason, "err", err)
	}
}

func (c *Controller) Consume(m jetstream.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), c.processingTimeout)
	defer cancel()

	trade, err := c.unmarshal(m)
	if err != nil {
		c.term(ctx, m, err.Error())
		return
	}

	if c.clientIDCache.Contains(trade.ClientOrderID) {
		c.logger.DebugContext(ctx, "client ID already seen, quitting early", "id", trade.ClientOrderID)
		return
	}

	if err = c.trade(ctx, trade); err != nil {
		if err = m.Nak(); err != nil {
			c.logger.ErrorContext(ctx, "failed nak", "err", err)
		}
		return
	}

	c.logger.DebugContext(ctx,
		"added to cache",
		"evicted?", c.clientIDCache.Add(trade.ClientOrderID, struct{}{}),
	)

	if err = m.Ack(); err != nil {
		c.logger.ErrorContext(ctx, "failed ACK", "err", err)
		return
	}
}

func (c *Controller) unmarshal(m jetstream.Msg) (alpaca.PlaceOrderRequest, error) {
	trade, err := c.getMsg(m)
	if err != nil {
		return alpaca.PlaceOrderRequest{}, err
	}

	var takeProfit *alpaca.TakeProfit
	if trade.TakeProfit != nil && trade.TakeProfit.LimitPrice > 0 {
		takeProfit = &alpaca.TakeProfit{
			LimitPrice: newDecimal(trade.TakeProfit.LimitPrice),
		}
	}

	var stopLoss *alpaca.StopLoss
	if s := trade.StopLoss; s != nil {
		if s.Limit <= 0 {
			c.logger.Error("stop loss config requires a limit", "limit", s.Limit)
			return alpaca.PlaceOrderRequest{}, fmt.Errorf("stop loss requires a valid limit, got %f", s.Limit)
		}

		if s.Stop <= 0 {
			c.logger.Error("stop loss config requires a stop", "stop", s.Stop)
			return alpaca.PlaceOrderRequest{}, fmt.Errorf("stop loss requires a valid stop, got %f", s.Stop)
		}

		stopLoss = &alpaca.StopLoss{
			LimitPrice: newDecimal(s.Limit),
			StopPrice:  newDecimal(s.Stop),
		}
	}

	return alpaca.PlaceOrderRequest{
		Symbol:         strings.ToUpper(trade.Symbol),
		Qty:            newDecimal(trade.Qty),
		Notional:       newDecimal(trade.Notional),
		Side:           toSide(trade.Side),
		Type:           toOrderType(trade.OrderType),
		TimeInForce:    toTIF(trade.Tif),
		LimitPrice:     newDecimal(trade.LimitPrice),
		ExtendedHours:  trade.ExtendedHours,
		StopPrice:      newDecimal(trade.StopPrice),
		ClientOrderID:  trade.ClientOrderId,
		OrderClass:     toOrderClass(trade.Class),
		TakeProfit:     takeProfit,
		StopLoss:       stopLoss,
		TrailPrice:     newDecimal(trade.TrailPrice),
		TrailPercent:   newDecimal(trade.TrailPercent),
		PositionIntent: toIntent(trade.PositionIntent),
	}, nil
}

func (c *Controller) getMsg(m jetstream.Msg) (*tradesvc.Trade, error) {
	var trade tradesvc.Trade
	if err := proto.Unmarshal(m.Data(), &trade); err != nil {
		return nil, err
	}

	if len(trade.ClientOrderId) > 128 {
		c.logger.Error("client order id was too long", "id", trade.ClientOrderId)
		return nil, fmt.Errorf("client order ID is too long: %s", trade.ClientOrderId)
	}

	if trade.Symbol == "" {
		c.logger.Error("no symbol passed")
		return nil, fmt.Errorf("no symbol given")
	}

	type prices struct {
		name  string
		price float64
	}

	for _, v := range []prices{
		{name: "limit price", price: trade.LimitPrice},
		{name: "notional", price: trade.Notional},
		{name: "qty", price: trade.Qty},
		{name: "stop price", price: trade.StopPrice},
		{name: "trail percent", price: trade.TrailPercent},
		{name: "trail price", price: trade.TrailPrice},
	} {
		if v.price < 0 {
			c.logger.Error("invalid price", "name", v.name, "value", v.price)
			return nil, fmt.Errorf("invalid %s: %f", v.name, v.price)
		}
	}

	if trade.Qty != 0 && trade.Notional != 0 {
		c.logger.Error(
			"quantity and notional can't both be set",
			"notional", trade.Notional,
			"qty", trade.Qty,
		)

		return nil, fmt.Errorf(
			"trade and quantity can't both be set, got qty %f and notional %f",
			trade.Qty,
			trade.Notional,
		)
	}

	return &trade, nil
}

func newDecimal(x float64) *decimal.Decimal {
	if x == 0 {
		return nil
	}

	y := decimal.NewFromFloat(x)
	return &y
}
