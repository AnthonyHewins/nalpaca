package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFTIF(tif alpaca.TimeInForce) tradesvc.TimeInForce {
	switch tif {
	case alpaca.CLS:
		return tradesvc.TimeInForce_TIME_IN_FORCE_CLS
	case alpaca.Day:
		return tradesvc.TimeInForce_TIME_IN_FORCE_DAY
	case alpaca.FOK:
		return tradesvc.TimeInForce_TIME_IN_FORCE_FOK
	case alpaca.GTC:
		return tradesvc.TimeInForce_TIME_IN_FORCE_GTC
	case alpaca.GTD:
		return tradesvc.TimeInForce_TIME_IN_FORCE_GTD
	case alpaca.GTX:
		return tradesvc.TimeInForce_TIME_IN_FORCE_GTX
	case alpaca.IOC:
		return tradesvc.TimeInForce_TIME_IN_FORCE_IOC
	case alpaca.OPG:
		return tradesvc.TimeInForce_TIME_IN_FORCE_OPG
	default:
		return tradesvc.TimeInForce_TIME_IN_FORCE_UNSPECIFIED
	}
}

func TIF(tif tradesvc.TimeInForce) alpaca.TimeInForce {
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
