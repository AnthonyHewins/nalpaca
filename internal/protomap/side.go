package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFSide(s alpaca.Side) tradesvc.Side {
	switch s {
	case alpaca.Buy:
		return tradesvc.Side_SIDE_BUY
	case alpaca.Sell:
		return tradesvc.Side_SIDE_SELL
	default:
		return tradesvc.Side_SIDE_UNSPECIFIED
	}
}

func Side(s tradesvc.Side) alpaca.Side {
	switch s {
	case tradesvc.Side_SIDE_BUY:
		return alpaca.Buy
	case tradesvc.Side_SIDE_SELL:
		return alpaca.Sell
	default:
		return ""
	}
}
