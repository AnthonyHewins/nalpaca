package nalpaca

import (
	"context"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Mock struct {
	StreamTradeUpdatesFn func(context.Context, func(alpaca.TradeUpdate), alpaca.StreamTradeUpdatesRequest) error
	PlaceOrderFn         func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
	CancelFn             func(string) error
	CancelAllOrdersFn    func() error
	GetAcctFn            func() (*alpaca.Account, error)
	PositionsFn          func() ([]alpaca.Position, error)
}

func (m Mock) StreamTradeUpdates(ctx context.Context, fn func(alpaca.TradeUpdate), req alpaca.StreamTradeUpdatesRequest) error {
	return m.StreamTradeUpdatesFn(ctx, fn, req)
}

func (m Mock) PlaceOrder(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
	return m.PlaceOrderFn(req)
}

func (m Mock) CancelOrder(id string) error {
	return m.CancelFn(id)
}

func (m Mock) CancelAllOrders() error {
	return m.CancelAllOrdersFn()
}

func (m Mock) GetAccount() (*alpaca.Account, error) {
	return m.GetAcctFn()
}

func (m Mock) GetPositions() ([]alpaca.Position, error) {
	return m.GetPositions()
}
