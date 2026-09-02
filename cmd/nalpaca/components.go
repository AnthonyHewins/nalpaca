package main

import (
	"context"
	"fmt"

	"github.com/AnthonyHewins/nalpaca/internal/canceler"
	"github.com/AnthonyHewins/nalpaca/internal/portfolio"
	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/nats-io/nats.go/jetstream"
)

func (a *app) initCanceler(ctx context.Context, js jetstream.JetStream, c *config) error {
	if !c.EnableCancel {
		return nil
	}

	var err error
	a.cancel.ingestor, err = a.consumer(ctx, js, c.ActionStream, c.CancelConsumer)
	if err != nil {
		return err
	}

	a.canceler = canceler.New(a.Logger, a.Nalpaca, cancelCounters, c.ProcessingTimeout)
	return nil
}

func (a *app) initOrders(ctx context.Context, js jetstream.JetStream, c *config) error {
	if !c.EnableOrders {
		return nil
	}

	var err error
	a.order.ingestor, err = a.consumer(ctx, js, c.ActionStream, c.OrderConsumerName)
	if err != nil {
		return err
	}

	a.trader = trader.NewController(
		a.TP.Tracer("trader"),
		a.Logger,
		orderCounters,
		a.Nalpaca,
		c.ProcessingTimeout,
	)

	return nil
}

func (a *app) initTradeUpdater(js jetstream.JetStream, kv jetstream.KeyValue, c *config) (*portfolio.Controller, error) {
	if !c.EnableTradeUpdater {
		return nil, nil
	}

	return portfolio.NewController(
		a.Logger,
		a.Nalpaca,
		c.ProcessingTimeout,
		js,
		kv,
		fmt.Sprintf("%s.account.tradeupdates", c.Prefix),
	), nil
}

func (a *app) initStockStream(js jetstream.JetStream, c *config) error {
	var err error
	if a.stockStream, err = a.streamClient.Stocks(&c.Alpaca.StockStream); err != nil {
		return err
	}

	return a.Metrics.Register(a.stockStream.Metrics()...)
}

func (a *app) initNewsStream(js jetstream.JetStream, c *config) error {
	var err error
	if a.news, err = a.streamClient.News(&c.Alpaca.NewsStream); err != nil {
		return err
	}

	if a.news == nil { // special case since this doesnt require a manager like options/stocks
		return nil
	}

	return a.Metrics.Register(a.news.Metrics()...)
}

func (a *app) initOptionStream(js jetstream.JetStream, c *config) error {
	var err error
	if a.optionStream, err = a.streamClient.Options(&c.Alpaca.OptionStream); err != nil {
		return err
	}

	return a.Metrics.Register(a.optionStream.Metrics()...)
}
