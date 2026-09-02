package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnthonyHewins/nalpaca/gen/go/stream/v0"
	"github.com/AnthonyHewins/nalpaca/internal/streaming"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *app) ListSubscriptions(ctx context.Context, req *stream.ListSubscriptionsRequest) (*stream.ListSubscriptionsResponse, error) {
	x, err := a.parseSubscription(req.Subscription)
	if err != nil {
		return nil, err
	}

	return &stream.ListSubscriptionsResponse{Subscriptions: x.List()}, err
}

func (a *app) AddSubscriptions(ctx context.Context, req *stream.AddSubscriptionsRequest) (*stream.AddSubscriptionsResponse, error) {
	x, err := a.parseSubscription(req.Subscription)
	if err != nil {
		return nil, err
	}
	return &stream.AddSubscriptionsResponse{}, x.Subscribe(req.Symbols...)
}

func (a *app) RemoveSubscriptions(ctx context.Context, req *stream.RemoveSubscriptionsRequest) (*stream.RemoveSubscriptionsResponse, error) {
	x, err := a.parseSubscription(req.Subscription)
	if err != nil {
		return nil, err
	}
	return &stream.RemoveSubscriptionsResponse{}, x.Unsubscribe(req.Symbols...)
}

func (a *app) parseSubscription(s string) (streaming.Subscriber, error) {
	clean := strings.ReplaceAll(
		strings.ToLower(strings.TrimSpace(s)),
		"-",
		"_",
	)

	sub, err := streaming.SubscriptionString(clean)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "subscription %s doesn't exist. Available: %s", s, strings.Join(streaming.SubscriptionStrings(), ", "))
	}

	var x streaming.Subscriber
	switch sub {
	case streaming.SubscriptionStockBars:
		x = a.stockStream.Bars
	case streaming.SubscriptionNews:
		x = a.news
	case streaming.SubscriptionOptionQuotes:
		x = a.optionStream.Quotes
	case streaming.SubscriptionStockQuotes:
		x = a.stockStream.Quotes
	case streaming.SubscriptionOptionTrades:
		x = a.optionStream.Trades
	case streaming.SubscriptionStockTrades:
		x = a.stockStream.Trades
	default:
		a.Logger.Error("code error unhandled case", "for", sub)
		return nil, fmt.Errorf("missing case statement for %s", sub)
	}

	if x == nil {
		a.Logger.Error("subscription type not turned on", "for", sub)
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"%s is not turned on; you must start nalpaca with the proper ENABLE_ flag to use this stream type",
			s,
		)
	}

	return x, nil
}
