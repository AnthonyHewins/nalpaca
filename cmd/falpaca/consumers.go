package main

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

func (a *app) connectConsumers(ctx context.Context, c *config) error {
	js, err := jetstream.New(a.nc)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed connecting to jetstream", "err", err)
		return err
	}

	subject := "falpaca.orders.>"
	if c.StreamPrefix != "" {
		subject = c.StreamPrefix + "." + subject
	}

	cfg := jetstream.ConsumerConfig{
		Name: "falpaca-order-consumer-v0",
		// Durable:            "",
		Description:   "Falpaca order consumer",
		MaxDeliver:    3,
		BackOff:       c.StreamBackoff,
		FilterSubject: subject,
		// FilterSubjects:     []string{},
	}
	a.consumers.order, err = js.CreateConsumer(ctx, subject, cfg)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed creating order consumer", "err", err, "config", cfg)
		return err
	}

	return nil
}
