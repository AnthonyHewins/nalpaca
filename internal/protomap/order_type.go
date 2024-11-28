package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFOrderType(x alpaca.OrderType) tradesvc.OrderType {
	switch x {
	case alpaca.Limit:
		return tradesvc.OrderType_ORDER_TYPE_LIMIT
	case alpaca.Market:
		return tradesvc.OrderType_ORDER_TYPE_MARKET
	case alpaca.Stop:
		return tradesvc.OrderType_ORDER_TYPE_STOP
	case alpaca.StopLimit:
		return tradesvc.OrderType_ORDER_TYPE_STOP_LIMIT
	case alpaca.TrailingStop:
		return tradesvc.OrderType_ORDER_TYPE_TRAILING_STOP
	default:
		return tradesvc.OrderType_ORDER_TYPE_UNSPECIFIED
	}
}

func OrderType(x tradesvc.OrderType) alpaca.OrderType {
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
