package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func (a *app) connectConsumers(ctx context.Context, c *config) error {
	js, err := jetstream.New(a.nc)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed connecting to jetstream", "err", err)
		return err
	}

	if c.StreamPrefix == "" {
		c.StreamPrefix = "nalpaca"
	}

	subject := fmt.Sprintf("%s.orders.v0.*", c.StreamPrefix)

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:         "nalpaca-order-stream-v0",
		Description:  "Nalpaca order stream",
		Subjects:     []string{subject},
		Retention:    jetstream.WorkQueuePolicy,
		MaxConsumers: 0,
		MaxMsgs:      1000,
		MaxBytes:     1024 * 50,
		Discard:      jetstream.DiscardOld,
		MaxAge:       time.Hour,
		MaxMsgSize:   1024,
		Storage:      jetstream.MemoryStorage,
		Replicas:     0,
		Compression:  jetstream.S2Compression,
	})

	if err != nil {
		a.logger.ErrorContext(ctx, "failed creating stream")
		return err
	}

	cfg := jetstream.ConsumerConfig{
		Name: "falpaca-order-consumer-v0",
		// Durable:            "",
		Description:   "Falpaca order consumer",
		MaxDeliver:    3,
		BackOff:       c.StreamBackoff,
		FilterSubject: fmt.Sprintf("%s.orders.v0.*", c.StreamPrefix),
	}

	a.order.ingestor, err = s.CreateConsumer(ctx, cfg)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed creating order consumer", "err", err, "config", cfg)
		return err
	}

	return nil
}
