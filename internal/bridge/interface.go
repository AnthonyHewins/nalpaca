package bridge

import (
	"context"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type AlpacaInterface interface {
	StreamTradeUpdates(context.Context, func(alpaca.TradeUpdate), alpaca.StreamTradeUpdatesRequest) error
	PlaceOrder(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
	CancelOrder(orderID string) error
	CancelAllOrders() error
	// GetAccount() (*alpaca.Account, error)
	GetPositions() ([]alpaca.Position, error)
}

type OptionStreamInterface interface {
}
