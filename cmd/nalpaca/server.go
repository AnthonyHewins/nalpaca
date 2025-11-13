package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotEnabled = status.Error(codes.FailedPrecondition, "this particular entity is not enabled. This means that you are missing an ENABLE_* environment variable. Enable that env var"+
	" and the service will be capable of responding correctly")

func (a *app) ListBarSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	if a.stockStream == nil {
		return nil, errNotEnabled
	}

	return &stream.ListSubscriptionsResponse{Subscriptions: a.stockStream.ListSubscriptions()}, nil
}

func (a *app) AddBarSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	if a.stockStream == nil {
		return nil, errNotEnabled
	}

	return &stream.AddSubscriptionsResponse{}, a.stockStream.AddSubscriptions(req.Symbols...)
}

func (a *app) RemoveBarSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	if a.stockStream == nil {
		return nil, errNotEnabled
	}

	return &stream.RemoveSubscriptionsResponse{}, a.stockStream.DeleteSubscriptions(req.Symbols...)
}

func (a *app) ListNewsSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	if a.newsStream == nil {
		return nil, errNotEnabled
	}

	return &stream.ListSubscriptionsResponse{Subscriptions: a.newsStream.ListSubscriptions()}, nil
}

func (a *app) AddNewsSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	if a.newsStream == nil {
		return nil, errNotEnabled
	}

	return &stream.AddSubscriptionsResponse{}, a.newsStream.AddSubscriptions(req.Symbols...)
}

func (a *app) RemoveNewsSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	if a.newsStream == nil {
		return nil, errNotEnabled
	}

	return &stream.RemoveSubscriptionsResponse{}, a.newsStream.DeleteSubscriptions(req.Symbols...)
}
