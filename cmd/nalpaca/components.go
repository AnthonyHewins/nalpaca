package main

import (
	"context"
	"fmt"

	"github.com/AnthonyHewins/nalpaca/internal/canceler"
	"github.com/AnthonyHewins/nalpaca/internal/portfolio"
	"github.com/AnthonyHewins/nalpaca/internal/streaming"
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
		fmt.Sprintf("%s.data.v0.tradeupdates", c.Prefix),
	), nil
}

func (a *app) initStockStream(js jetstream.JetStream, c *config) (*streaming.Stocks, error) {
	if !c.EnableStockStream {
		return nil, nil
	}

	return streaming.NewStocks(
		a.Logger,
		streaming.NewMetrics(appName),
		js,
		fmt.Sprintf("%s.data.v0.stocks", c.Prefix),
		c.Alpaca.APIKey,
		c.Alpaca.APISecret,
		&c.Alpaca.StockStream,
	)
}

func (a *app) initNewsStream(js jetstream.JetStream, c *config) (*streaming.News, error) {
	if !c.EnableNewsStream {
		return nil, nil
	}

	return streaming.NewNews(
		a.Logger,
		streaming.NewMetrics(appName),
		js,
		fmt.Sprintf("%s.data.v0.news", c.Prefix),
		c.Alpaca.APIKey,
		c.Alpaca.APISecret,
		&c.Alpaca.NewsStream,
	)
}
