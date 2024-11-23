package nalpaca

import (
	"context"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Interface interface {
	StreamTradeUpdates(context.Context, func(alpaca.TradeUpdate), alpaca.StreamTradeUpdatesRequest) error
	PlaceOrder(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
}
