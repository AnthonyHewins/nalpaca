package trader

import (
	"github.com/AnthonyHewins/falpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func toSide(s tradesvc.Side) alpaca.Side {
	switch s {
	case tradesvc.Side_SIDE_BUY:
		return alpaca.Buy
	case tradesvc.Side_SIDE_SELL:
		return alpaca.Sell
	default:
		return ""
	}
}

func toIntent(i tradesvc.PositionIntent) alpaca.PositionIntent {
	switch i {
	case tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_CLOSE:
		return alpaca.BuyToClose
	case tradesvc.PositionIntent_POSITION_INTENT_BUY_TO_OPEN:
		return alpaca.BuyToOpen
	case tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_CLOSE:
		return alpaca.SellToClose
	case tradesvc.PositionIntent_POSITION_INTENT_SELL_TO_OPEN:
		return alpaca.SellToOpen
	default:
		return ""
	}
}

func toTIF(tif tradesvc.TimeInForce) alpaca.TimeInForce {
	switch tif {
	case tradesvc.TimeInForce_TIME_IN_FORCE_CLS:
		return alpaca.CLS
	case tradesvc.TimeInForce_TIME_IN_FORCE_DAY:
		return alpaca.Day
	case tradesvc.TimeInForce_TIME_IN_FORCE_FOK:
		return alpaca.FOK
	case tradesvc.TimeInForce_TIME_IN_FORCE_GTC:
		return alpaca.GTC
	case tradesvc.TimeInForce_TIME_IN_FORCE_GTD:
		return alpaca.GTD
	case tradesvc.TimeInForce_TIME_IN_FORCE_GTX:
		return alpaca.GTX
	case tradesvc.TimeInForce_TIME_IN_FORCE_IOC:
		return alpaca.IOC
	case tradesvc.TimeInForce_TIME_IN_FORCE_OPG:
		return alpaca.OPG
	default:
		return ""
	}
}

func toOrderClass(c tradesvc.OrderClass) alpaca.OrderClass {
	switch c {
	case tradesvc.OrderClass_ORDER_CLASS_BRACKET:
		return alpaca.Bracket
	case tradesvc.OrderClass_ORDER_CLASS_OCO:
		return alpaca.OCO
	case tradesvc.OrderClass_ORDER_CLASS_OTO:
		return alpaca.OTO
	case tradesvc.OrderClass_ORDER_CLASS_SIMPLE:
		return alpaca.Simple
	default:
		return ""
	}
}

func toOrderType(x tradesvc.OrderType) alpaca.OrderType {
	switch x {
	case tradesvc.OrderType_ORDER_TYPE_LIMIT:
		return alpaca.Limit
	case tradesvc.OrderType_ORDER_TYPE_MARKET:
		return alpaca.Market
	case tradesvc.OrderType_ORDER_TYPE_STOP:
		return alpaca.Stop
	case tradesvc.OrderType_ORDER_TYPE_STOP_LIMIT:
		return alpaca.StopLimit
	case tradesvc.OrderType_ORDER_TYPE_TRAILING_STOP:
		return alpaca.TrailingStop
	default:
		return ""
	}
}
