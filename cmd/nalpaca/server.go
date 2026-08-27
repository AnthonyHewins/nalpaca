package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotEnabled = status.Error(codes.FailedPrecondition, "this particular entity is not enabled. This means that you are missing an ENABLE_* environment variable. Enable that env var"+
	" and the service will be capable of responding correctly")

// The three helpers below collapse what would otherwise be the same handful of
// lines repeated across every RPC: verify the backing stream is running, then
// delegate to it.
//
// Callers pass a method value from a possibly-nil stream (e.g.
// a.stockStream.ListBarSubscriptions). Binding a method value to a nil pointer
// receiver is legal in Go and does not dereference — only calling it would — and
// these helpers never call unless running is true.
func listSubs(running bool, list func() []string) (*stream.ListSubscriptionsResponse, error) {
	if !running {
		return nil, errNotEnabled
	}
	return &stream.ListSubscriptionsResponse{Subscriptions: list()}, nil
}

func addSubs(running bool, add func(...string) error, symbols []string) (*stream.AddSubscriptionsResponse, error) {
	if !running {
		return nil, errNotEnabled
	}
	return &stream.AddSubscriptionsResponse{}, add(symbols...)
}

func removeSubs(running bool, remove func(...string) error, symbols []string) (*stream.RemoveSubscriptionsResponse, error) {
	if !running {
		return nil, errNotEnabled
	}
	return &stream.RemoveSubscriptionsResponse{}, remove(symbols...)
}

//====================================
// Stock bars
//====================================

func (a *app) ListBarSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	return listSubs(a.stockStream != nil, a.stockStream.ListBarSubscriptions)
}

func (a *app) AddBarSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	return addSubs(a.stockStream != nil, a.stockStream.AddBarSubscriptions, req.Symbols)
}

func (a *app) RemoveBarSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	return removeSubs(a.stockStream != nil, a.stockStream.DeleteBarSubscriptions, req.Symbols)
}

//====================================
// Stock quotes
//====================================

func (a *app) ListStockQuoteSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	return listSubs(a.stockStream != nil, a.stockStream.ListQuoteSubscriptions)
}

func (a *app) AddStockQuoteSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	return addSubs(a.stockStream != nil, a.stockStream.AddQuoteSubscriptions, req.Symbols)
}

func (a *app) RemoveStockQuoteSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	return removeSubs(a.stockStream != nil, a.stockStream.DeleteQuoteSubscriptions, req.Symbols)
}

//====================================
// Stock trades
//====================================

func (a *app) ListStockTradeSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	return listSubs(a.stockStream != nil, a.stockStream.ListTradeSubscriptions)
}

func (a *app) AddStockTradeSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	return addSubs(a.stockStream != nil, a.stockStream.AddTradeSubscriptions, req.Symbols)
}

func (a *app) RemoveStockTradeSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	return removeSubs(a.stockStream != nil, a.stockStream.DeleteTradeSubscriptions, req.Symbols)
}

//====================================
// News
//====================================

func (a *app) ListNewsSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	return listSubs(a.newsStream != nil, a.newsStream.ListSubscriptions)
}

func (a *app) AddNewsSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	return addSubs(a.newsStream != nil, a.newsStream.AddSubscriptions, req.Symbols)
}

func (a *app) RemoveNewsSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	return removeSubs(a.newsStream != nil, a.newsStream.DeleteSubscriptions, req.Symbols)
}
