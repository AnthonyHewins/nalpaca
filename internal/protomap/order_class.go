package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFOrderClass(c alpaca.OrderClass) tradesvc.OrderClass {
	switch c {
	case alpaca.Bracket:
		return tradesvc.OrderClass_ORDER_CLASS_BRACKET
	case alpaca.OCO:
		return tradesvc.OrderClass_ORDER_CLASS_OCO
	case alpaca.OTO:
		return tradesvc.OrderClass_ORDER_CLASS_OTO
	case alpaca.Simple:
		return tradesvc.OrderClass_ORDER_CLASS_SIMPLE
	default:
		return tradesvc.OrderClass_ORDER_CLASS_UNSPECIFIED
	}
}

func OrderClass(c tradesvc.OrderClass) alpaca.OrderClass {
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
