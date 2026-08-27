package bridge

import (
	"context"
	"errors"
	"testing"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

// Every Mock method must delegate to its corresponding ...Fn field. GetPositions
// used to call itself instead of PositionsFn, which would recurse until the
// stack blew up the first time anything exercised it.
//
// Driving the calls through AlpacaInterface rather than the concrete type is
// what makes this a real test of the seam the rest of the code depends on.
func TestMockDelegatesEveryMethod(t *testing.T) {
	var (
		wantOrder    = &alpaca.Order{ID: "order-1"}
		wantAcct     = &alpaca.Account{ID: "acct-1"}
		wantPosition = []alpaca.Position{{Symbol: "AAPL"}}
		errBoom      = errors.New("boom")

		called = map[string]bool{}
	)

	mock := Mock{
		StreamTradeUpdatesFn: func(context.Context, func(alpaca.TradeUpdate), alpaca.StreamTradeUpdatesRequest) error {
			called["stream"] = true
			return errBoom
		},
		PlaceOrderFn: func(alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
			called["place"] = true
			return wantOrder, nil
		},
		CancelFn: func(string) error {
			called["cancel"] = true
			return errBoom
		},
		CancelAllOrdersFn: func() error {
			called["cancelAll"] = true
			return errBoom
		},
		GetAcctFn: func() (*alpaca.Account, error) {
			called["acct"] = true
			return wantAcct, nil
		},
		PositionsFn: func() ([]alpaca.Position, error) {
			called["positions"] = true
			return wantPosition, nil
		},
	}

	// GetAccount is currently commented out of AlpacaInterface, so it is
	// exercised on the concrete Mock; everything else goes through the interface.
	var iface AlpacaInterface = mock

	if err := iface.StreamTradeUpdates(context.Background(), nil, alpaca.StreamTradeUpdatesRequest{}); !errors.Is(err, errBoom) {
		t.Errorf("StreamTradeUpdates: want errBoom, got %v", err)
	}

	order, err := iface.PlaceOrder(alpaca.PlaceOrderRequest{})
	if err != nil || order != wantOrder {
		t.Errorf("PlaceOrder: want %v/nil, got %v/%v", wantOrder, order, err)
	}

	if err := iface.CancelOrder("x"); !errors.Is(err, errBoom) {
		t.Errorf("CancelOrder: want errBoom, got %v", err)
	}

	if err := iface.CancelAllOrders(); !errors.Is(err, errBoom) {
		t.Errorf("CancelAllOrders: want errBoom, got %v", err)
	}

	acct, err := mock.GetAccount()
	if err != nil || acct != wantAcct {
		t.Errorf("GetAccount: want %v/nil, got %v/%v", wantAcct, acct, err)
	}

	// The regression: this recursed forever before PositionsFn was wired up.
	positions, err := iface.GetPositions()
	if err != nil {
		t.Fatalf("GetPositions: unexpected error %v", err)
	}
	if len(positions) != 1 || positions[0].Symbol != "AAPL" {
		t.Errorf("GetPositions: want [AAPL], got %v", positions)
	}

	for _, name := range []string{"stream", "place", "cancel", "cancelAll", "acct", "positions"} {
		if !called[name] {
			t.Errorf("%s: the mock did not delegate to its Fn field", name)
		}
	}
}
