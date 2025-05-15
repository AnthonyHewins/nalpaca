package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/internal/trader"
	"github.com/nats-io/nats.go/jetstream"
)

type OrdersConf struct {
	EnableOrders        bool   `env:"ENABLE_ORDERS" envDefault:"false"`
	OrderConsumerStream string `env:"ORDER_CONSUMER_STREAM" envDefault:"nalpaca-orders-stream-v0"`
	OrderConsumerName   string `env:"ORDER_CONSUMER_NAME" envDefault:"nalpaca-orders-consumer-v0"`
}

func (a *app) initOrders(ctx context.Context, js jetstream.JetStream, c *config) error {
	if !c.EnableOrders {
		return nil
	}

	var err error
	a.order.ingestor, err = a.connect(ctx, js, c.OrderConsumerStream, c.OrderConsumerName)
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
