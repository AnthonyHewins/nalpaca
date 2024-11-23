package nalpaca

import (
	"context"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Mock struct {
	StreamTradeUpdatesFn func(context.Context, func(alpaca.TradeUpdate), alpaca.StreamTradeUpdatesRequest) error
	PlaceOrderFn         func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
}

func (m Mock) StreamTradeUpdates(ctx context.Context, fn func(alpaca.TradeUpdate), req alpaca.StreamTradeUpdatesRequest) error {
	return m.StreamTradeUpdatesFn(ctx, fn, req)
}

func (m Mock) PlaceOrder(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
	return m.PlaceOrderFn(req)
}
