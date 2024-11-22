package trader

import (
	"fmt"
	"strings"

	"github.com/AnthonyHewins/falpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

func (c *Controller) getMsg(m jetstream.Msg, id string) (alpaca.PlaceOrderRequest, error) {
	var trade tradesvc.Trade
	if err := proto.Unmarshal(m.Data(), &trade); err != nil {
		return alpaca.PlaceOrderRequest{}, err
	}

	if trade.Symbol == "" {
		c.logger.Error("no symbol passed")
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("no symbol given")
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
			return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid %s: %f", v.name, v.price)
		}
	}

	if trade.Qty != 0 && trade.Notional != 0 {
		c.logger.Error(
			"quantity and notional can't both be set",
			"notional", trade.Notional,
			"qty", trade.Qty,
		)

		return alpaca.PlaceOrderRequest{}, fmt.Errorf(
			"trade and quantity can't both be set, got qty %f and notional %f",
			trade.Qty,
			trade.Notional,
		)
	}

	takeProfit, err := c.toTakeProfit(trade.TakeProfit)
	if err != nil {
		return alpaca.PlaceOrderRequest{}, err
	}

	stopLoss, err := c.toStopLoss(trade.StopLoss)
	if err != nil {
		return alpaca.PlaceOrderRequest{}, err
	}

	o := alpaca.PlaceOrderRequest{
		ClientOrderID:  id,
		Symbol:         strings.ToUpper(trade.Symbol),
		Qty:            newDecimal(trade.Qty),
		Notional:       newDecimal(trade.Notional),
		Side:           toSide(trade.Side),
		Type:           toOrderType(trade.OrderType),
		TimeInForce:    toTIF(trade.Tif),
		LimitPrice:     newDecimal(trade.LimitPrice),
		ExtendedHours:  trade.ExtendedHours,
		StopPrice:      newDecimal(trade.StopPrice),
		OrderClass:     toOrderClass(trade.Class),
		TakeProfit:     takeProfit,
		StopLoss:       stopLoss,
		TrailPrice:     newDecimal(trade.TrailPrice),
		TrailPercent:   newDecimal(trade.TrailPercent),
		PositionIntent: toIntent(trade.PositionIntent),
	}

	type requiredEnums struct {
		name, value string
	}

	for _, v := range []requiredEnums{
		{name: "side", value: string(o.Side)},
		{name: "type", value: string(o.Type)},
		{name: "time in force", value: string(o.TimeInForce)},
		{name: "position intent", value: string(o.PositionIntent)},
	} {
		if v.value != "" {
			return alpaca.PlaceOrderRequest{}, fmt.Errorf("missing required value %s: got zero value", v.name)
		}
	}

	return o, nil
}

func (c *Controller) toStopLoss(s *tradesvc.StopLoss) (*alpaca.StopLoss, error) {
	if s == nil {
		return nil, nil
	}
	if s.Limit <= 0 {
		c.logger.Error("stop loss config requires a limit", "limit", s.Limit)
		return nil, fmt.Errorf("stop loss requires a valid limit, got %f", s.Limit)
	}

	if s.Stop <= 0 {
		c.logger.Error("stop loss config requires a stop", "stop", s.Stop)
		return nil, fmt.Errorf("stop loss requires a valid stop, got %f", s.Stop)
	}

	return &alpaca.StopLoss{
		LimitPrice: newDecimal(s.Limit),
		StopPrice:  newDecimal(s.Stop),
	}, nil
}

func (c *Controller) toTakeProfit(t *tradesvc.TakeProfit) (*alpaca.TakeProfit, error) {
	if t == nil {
		return nil, nil
	}

	if t.LimitPrice < 0 {
		c.logger.Error("invalid limit price for take profit", "got", t.LimitPrice)
		return nil, fmt.Errorf("invalid limit price for take profit object: %f", t.LimitPrice)
	}

	if t.LimitPrice == 0 {
		return nil, nil
	}

	return &alpaca.TakeProfit{LimitPrice: newDecimal(t.LimitPrice)}, nil
}

func newDecimal(x float64) *decimal.Decimal {
	if x == 0 {
		return nil
	}

	y := decimal.NewFromFloat(x)
	return &y
}
