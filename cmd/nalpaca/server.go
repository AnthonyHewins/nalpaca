package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
)

func (a *app) ListBarSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	return &stream.ListSubscriptionsResponse{Subscriptions: a.stockStream.ListSubscriptions()}, nil
}

func (a *app) AddBarSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	return &stream.AddSubscriptionsResponse{}, a.stockStream.AddSubscriptions(req.Symbols...)
}

func (a *app) RemoveBarSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	return &stream.RemoveSubscriptionsResponse{}, a.stockStream.DeleteSubscriptions(req.Symbols...)
}

func (a *app) ListNewsSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	return &stream.ListSubscriptionsResponse{Subscriptions: a.newsStream.ListSubscriptions()}, nil
}

func (a *app) AddNewsSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	return &stream.AddSubscriptionsResponse{}, a.newsStream.AddSubscriptions(req.Symbols...)
}

func (a *app) RemoveNewsSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	return &stream.RemoveSubscriptionsResponse{}, a.newsStream.DeleteSubscriptions(req.Symbols...)
}
