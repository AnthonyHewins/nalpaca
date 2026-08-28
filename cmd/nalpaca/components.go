package main

import (
	"context"
	"fmt"

	"github.com/AnthonyHewins/nalpaca/internal/canceler"
	"github.com/AnthonyHewins/nalpaca/internal/portfolio"
	"github.com/AnthonyHewins/nalpaca/internal/streaming"
	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
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

// streamMetrics holds the per-stream counters. They're built before
// BootstrapConf.New because that's what stands up the prometheus registry, and a
// collector registered after the fact would never be exported.
//
// Each stream gets its own subsystem: prometheus rejects two collectors sharing a
// fully-qualified name, so reusing one subsystem across streams makes the whole
// app fail to start the moment metrics are enabled.
type streamMetrics struct {
	stocks, news, options streaming.Metrics
	collectors            []prometheus.Collector
}

func newStreamMetrics(c *config) streamMetrics {
	var s streamMetrics

	// Only build counters for streams that are actually on, so /metrics doesn't
	// advertise streams that aren't running.
	if c.EnableStockStream {
		s.stocks = streaming.NewMetrics(appName, "stocks_stream")
		s.collectors = append(s.collectors, s.stocks.Collectors()...)
	}
	if c.EnableNewsStream {
		s.news = streaming.NewMetrics(appName, "news_stream")
		s.collectors = append(s.collectors, s.news.Collectors()...)
	}
	if c.EnableOptionStream {
		s.options = streaming.NewMetrics(appName, "options_stream")
		s.collectors = append(s.collectors, s.options.Collectors()...)
	}

	if c.EnableCancel {
		s.collectors = append(s.collectors, cancelCounters.Collectors()...)
	}
	if c.EnableOrders {
		s.collectors = append(s.collectors, orderCounters.Collectors()...)
	}

	return s
}

func (a *app) initStockStream(js jetstream.JetStream, c *config, m streaming.Metrics) (*streaming.Stocks, error) {
	if !c.EnableStockStream {
		return nil, nil
	}

	return streaming.NewStocks(
		a.Logger,
		m,
		js,
		fmt.Sprintf("%s.data.stocks", c.Prefix),
		c.Alpaca.APIKey,
		c.Alpaca.APISecret,
		&c.Alpaca.StockStream,
	)
}

func (a *app) initNewsStream(js jetstream.JetStream, c *config, m streaming.Metrics) (*streaming.News, error) {
	if !c.EnableNewsStream {
		return nil, nil
	}

	return streaming.NewNews(
		a.Logger,
		m,
		js,
		fmt.Sprintf("%s.data.news", c.Prefix),
		c.Alpaca.APIKey,
		c.Alpaca.APISecret,
		&c.Alpaca.NewsStream,
	)
}

func (a *app) initOptionStream(js jetstream.JetStream, c *config, m streaming.Metrics) (*streaming.Options, error) {
	if !c.EnableOptionStream {
		return nil, nil
	}

	return streaming.NewOptions(
		a.Logger,
		m,
		js,
		fmt.Sprintf("%s.data.options", c.Prefix),
		c.Alpaca.APIKey,
		c.Alpaca.APISecret,
		&c.Alpaca.OptionStream,
	)
}
