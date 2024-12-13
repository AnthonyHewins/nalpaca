package trader

import (
	"fmt"
	"strings"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/AnthonyHewins/nalpaca/internal/protomap"
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

	l := c.logger.With("trade", trade.String())

	if trade.Symbol == "" {
		l.Error("no symbol passed")
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("no symbol given")
	}

	type prices struct {
		name  string
		price string
	}

	decimals := make([]*decimal.Decimal, 6)
	for i, v := range []prices{
		{name: "qty", price: trade.Qty},
		{name: "notional", price: trade.Notional},
		{name: "limit price", price: trade.LimitPrice},
		{name: "stop price", price: trade.StopPrice},
		{name: "trail price", price: trade.TrailPrice},
		{name: "trail percent", price: trade.TrailPercent},
	} {
		l = l.With("name", v.name, "val", v.price)

		x, err := newDecimal(v.price)
		if err != nil {
			l.Error("failed converting price to decimal", "err", err)
			return alpaca.PlaceOrderRequest{}, err
		}

		if x != nil && x.IsNegative() {
			l.Error("invalid decimal", "name", v.name, "value", v.price)
			return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid %s: %s", v.name, v.price)
		}

		decimals[i] = x
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
		TakeProfit:     takeProfit,
		StopLoss:       stopLoss,
		ExtendedHours:  trade.ExtendedHours,
		Qty:            decimals[0],
		Notional:       decimals[1],
		LimitPrice:     decimals[2],
		StopPrice:      decimals[3],
		TrailPrice:     decimals[4],
		TrailPercent:   decimals[5],
		Side:           protomap.Side(trade.Side),
		Type:           protomap.OrderType(trade.OrderType),
		TimeInForce:    protomap.TIF(trade.Tif),
		OrderClass:     protomap.OrderClass(trade.Class),
		PositionIntent: protomap.Intent(trade.PositionIntent),
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
		if v.value == "" {
			return alpaca.PlaceOrderRequest{}, fmt.Errorf("missing required value %s: got zero value", v.name)
		}
	}

	return o, nil
}

func (c *Controller) toStopLoss(s *tradesvc.StopLoss) (*alpaca.StopLoss, error) {
	if s == nil || s.Limit == "" && s.Stop == "" {
		return nil, nil
	}

	limit, err := newDecimal(s.Limit)
	if err != nil {
		c.logger.Error("invalid limit passed to stoploss", "err", err, "limit", s.Limit)
		return nil, err
	}

	stop, err := newDecimal(s.Stop)
	if err != nil {
		c.logger.Error("stop loss config requires a valid stop", "err", err, "stop", stop)
		return nil, err
	}

	return &alpaca.StopLoss{LimitPrice: limit, StopPrice: stop}, nil
}

func (c *Controller) toTakeProfit(t *tradesvc.TakeProfit) (*alpaca.TakeProfit, error) {
	if t == nil || t.LimitPrice == "" {
		return nil, nil
	}

	limit, err := newDecimal(t.LimitPrice)
	if err != nil {
		c.logger.Error("invalid limit price for take profit", "got", t.LimitPrice, "err", err)
		return nil, err
	}

	return &alpaca.TakeProfit{LimitPrice: limit}, nil
}

func newDecimal(x string) (*decimal.Decimal, error) {
	if x == "" {
		return nil, nil
	}

	y, err := decimal.NewFromString(x)
	if err != nil {
		return nil, err
	}

	return &y, nil
}
