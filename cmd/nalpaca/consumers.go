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

	for _, fn := range []func(context.Context, jetstream.JetStream, *config) error{
		a.createOrderStream,
		a.createUpdater,
	} {
		if err = fn(ctx, js, c); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) createOrderStream(ctx context.Context, js jetstream.JetStream, c *config) error {
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
		a.logger.ErrorContext(ctx, "failed creating stream", "err", err)
		return err
	}

	cfg := jetstream.ConsumerConfig{
		Name:          "nalpaca-order-consumer-v0",
		Durable:       "",
		Description:   "nalpaca order consumer",
		MaxDeliver:    3,
		BackOff:       c.StreamBackoff,
		FilterSubject: fmt.Sprintf("%s.orders.v0.*", c.StreamPrefix),
	}

	a.order.ingestor, err = s.CreateConsumer(ctx, cfg)
	if err != nil {
		a.logger.ErrorContext(ctx, "failed creating order consumer", "err", err, "config", cfg)
	}

	return err
}

func (a *app) createUpdater(ctx context.Context, js jetstream.JetStream, c *config) error {
	subject := fmt.Sprintf("%s.tradeupdates.v0.*", c.StreamPrefix)

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:         "nalpaca-tradeupdate-stream-v0",
		Description:  "Nalpaca trade update stream",
		Subjects:     []string{subject},
		Retention:    jetstream.WorkQueuePolicy,
		MaxConsumers: 0,
		MaxMsgs:      1000,
		MaxBytes:     1024 * 50,
		Discard:      jetstream.DiscardOld,
		MaxAge:       time.Hour,
		MaxMsgSize:   1024 * 2,
		Storage:      jetstream.MemoryStorage,
		Replicas:     0,
		Compression:  jetstream.S2Compression,
	})

	if err != nil {
		a.logger.ErrorContext(ctx, "failed creating trade update stream", "err", err)
		return err
	}

	cfg := jetstream.ConsumerConfig{
		Name:          "nalpaca-tradeupdate-consumer-v0",
		Description:   "Trade updates",
		MaxDeliver:    3,
		BackOff:       c.StreamBackoff,
		FilterSubject: fmt.Sprintf("%s.tradeupdate.v0.*", c.StreamPrefix),
	}

	if _, err = s.CreateConsumer(ctx, cfg); err != nil {
		a.logger.ErrorContext(ctx, "failed creating consumer", "err", err, "cfg", cfg)
	}

	return err
}
