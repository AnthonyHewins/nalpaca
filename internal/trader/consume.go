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

	if err = c.trade(ctx, trade); err != nil {
		if err = m.Nak(); err != nil {
			c.logger.ErrorContext(ctx, "failed nak", "err", err)
		}
		return
	}

	if err = m.Ack(); err != nil {
		c.logger.ErrorContext(ctx, "failed ACK", "err", err)
		return
	}
}

func (c *Controller) unmarshal(m jetstream.Msg) (alpaca.PlaceOrderRequest, error) {
	var trade tradesvc.Trade
	if err := proto.Unmarshal(m.Data(), &trade); err != nil {
		return alpaca.PlaceOrderRequest{}, err
	}

	if trade.Symbol == "" {
		c.logger.Error("no symbol passed")
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("no symbol given")
	}

	var side alpaca.Side
	switch trade.Side {
	case tradesvc.Side_SIDE_BUY:
		side = alpaca.Buy
	case tradesvc.Side_SIDE_SELL:
		side = alpaca.Sell
	default:
		c.logger.Error("invalid trade side", "side", trade.Side)
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid trade side %s", trade.Side)
	}

	var orderType alpaca.OrderType
	switch trade.OrderType {
	case tradesvc.OrderType_ORDER_TYPE_LIMIT:
		orderType = alpaca.Limit
	case tradesvc.OrderType_ORDER_TYPE_MARKET:
		orderType = alpaca.Market
	case tradesvc.OrderType_ORDER_TYPE_STOP:
		orderType = alpaca.Stop
	case tradesvc.OrderType_ORDER_TYPE_STOP_LIMIT:
		orderType = alpaca.StopLimit
	case tradesvc.OrderType_ORDER_TYPE_TRAILING_STOP:
		orderType = alpaca.TrailingStop
	default:
		c.logger.Error("invalid order type", "type", trade.OrderType)
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid order type %s", trade.OrderType)
	}

	var tif alpaca.TimeInForce
	switch trade.Tif {
	case tradesvc.TimeInForce_TIME_IN_FORCE_CLS:
		tif = alpaca.CLS
	case tradesvc.TimeInForce_TIME_IN_FORCE_DAY:
		tif = alpaca.Day
	case tradesvc.TimeInForce_TIME_IN_FORCE_FOK:
		tif = alpaca.FOK
	case tradesvc.TimeInForce_TIME_IN_FORCE_GTC:
		tif = alpaca.GTC
	case tradesvc.TimeInForce_TIME_IN_FORCE_GTD:
		tif = alpaca.GTD
	case tradesvc.TimeInForce_TIME_IN_FORCE_GTX:
		tif = alpaca.GTX
	case tradesvc.TimeInForce_TIME_IN_FORCE_IOC:
		tif = alpaca.IOC
	case tradesvc.TimeInForce_TIME_IN_FORCE_OPG:
		tif = alpaca.OPG
	default:
		c.logger.Error("invalid TIF", "tif", trade.Tif)
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid tif %s", trade.Tif)
	}

	var orderClass alpaca.OrderClass
	switch trade.Class {
	case tradesvc.OrderClass_ORDER_CLASS_BRACKET:
		orderClass = alpaca.Bracket
	case tradesvc.OrderClass_ORDER_CLASS_OCO:
		orderClass = alpaca.OCO
	case tradesvc.OrderClass_ORDER_CLASS_OTO:
		orderClass = alpaca.OTO
	case tradesvc.OrderClass_ORDER_CLASS_SIMPLE:
		orderClass = alpaca.Simple
	default:
		c.logger.Error("invalid order class", "class", trade.Class)
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid order class %s", trade.Class)
	}

	var takeProfit *alpaca.TakeProfit
	if trade.TakeProfit != nil && trade.TakeProfit.LimitPrice > 0 {
		takeProfit = &alpaca.TakeProfit{
			LimitPrice: newDecimal(trade.TakeProfit.LimitPrice),
		}
	}

	var stopLoss *alpaca.StopLoss
	if s := trade.StopLoss; s != nil {
		stopLoss = &alpaca.StopLoss{
			LimitPrice: newDecimal(s.Limit),
			StopPrice:  newDecimal(s.Stop),
		}
	}

	var intent alpaca.PositionIntent
	switch i := trade.PositionIntent; i {
	case tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_CLOSE:
		intent = alpaca.BuyToClose
	case tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN:
		intent = alpaca.BuyToOpen
	case tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_CLOSE:
		intent = alpaca.SellToClose
	case tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_OPEN:
		intent = alpaca.SellToOpen
	default:
		c.logger.Error("invalid position intent", "intent", i)
		return alpaca.PlaceOrderRequest{}, fmt.Errorf("invalid intent %s", i)
	}

	return alpaca.PlaceOrderRequest{
		Symbol:         strings.ToUpper(trade.Symbol),
		Qty:            newDecimal(trade.Qty),
		Notional:       &decimal.Decimal{},
		Side:           side,
		Type:           orderType,
		TimeInForce:    tif,
		LimitPrice:     newDecimal(trade.LimitPrice),
		ExtendedHours:  trade.ExtendedHours,
		StopPrice:      newDecimal(trade.StopPrice),
		ClientOrderID:  trade.ClientOrderId,
		OrderClass:     orderClass,
		TakeProfit:     takeProfit,
		StopLoss:       stopLoss,
		TrailPrice:     newDecimal(trade.TrailPrice),
		TrailPercent:   newDecimal(trade.TrailPercent),
		PositionIntent: intent,
	}, nil
}

func newDecimal(x float64) *decimal.Decimal {
	y := decimal.NewFromFloat(x)
	return &y
}
