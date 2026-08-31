package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errNotEnabled = status.Error(codes.FailedPrecondition, "this particular entity is not enabled. This means that you are missing an ENABLE_* environment variable. Enable that env var"+
		" and the service will be capable of responding correctly")
)

func (a *app) ListNewsSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	x, err := a.parseSubscription(req.Subscription)
	if err != nil {
		return nil, err
	}

	return listSubs(a.newsStream != nil, a.newsStream.ListSubscriptions)
}

func (a *app) AddNewsSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	x, err := a.parseSubscription(req.Subscription)
	if err != nil {
		return nil, err
	}

	return addSubs(a.newsStream != nil, a.newsStream.AddSubscriptions, req.Symbols)
}

func (a *app) RemoveNewsSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	x, err := a.parseSubscription(req.Subscription)
	if err != nil {
		return nil, err
	}

	return removeSubs(a.newsStream != nil, a.newsStream.DeleteSubscriptions, req.Symbols)
}
