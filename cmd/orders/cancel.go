package main

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
)

type canceler struct {
	logger *slog.Logger

	cancelCount, cancelFail       prometheus.Counter
	cancelAllCount, cancelAllFail prometheus.Counter

	client  nalpaca.Interface
	timeout time.Duration
}

func (c *canceler) ack(m jetstream.Msg) {
	if err := m.Ack(); err != nil {
		c.logger.Error("failed ack", "err", err)
	}
}

func (c *canceler) nak(m jetstream.Msg) {
	if err := m.Nak(); err != nil {
		c.logger.Error("failed nak", "err", err)
	}
}

func (c *canceler) term(m jetstream.Msg, reason string) {
	if err := m.TermWithReason(reason); err != nil {
		c.logger.Error("failed term", "err", err, "reason", reason)
	}
}

func (c *canceler) eventLoop(m jetstream.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	l := c.logger.With("subject", m.Subject())

	subj := strings.Split(m.Subject(), ".")
	n := len(subj)
	if n <= 1 {
		l.ErrorContext(ctx, "got invalid subject, could not parse order id")
		c.term(m, "invalid subject received")
		return
	}

	id := subj[n-1]
	l = l.With("id", id)

	if len(id) > 128 {
		l.ErrorContext(ctx, "order ID is too long")
		c.term(m, "invalid ID: "+id)
		return
	}

	if strings.ToUpper(id) == "ALL" {
		c.cancelAll(m)
		return
	}

	if err := c.client.CancelOrder(id); err != nil {
		l.ErrorContext(ctx, "failed canceling order", "err", err)
		c.cancelFail.Inc()
		c.nak(m)
		return
	}

	l.InfoContext(ctx, "successful cancel")
	c.cancelCount.Inc()
	c.ack(m)
}

func (c *canceler) cancelAll(m jetstream.Msg) {
	c.logger.Warn("received cancel all orders request")
	if err := c.client.CancelAllOrders(); err != nil {
		c.logger.Error("failed canceling all orders", "err", err)
		c.cancelAllFail.Inc()
		c.nak(m)
		return
	}

	c.logger.Info("successfully canceled all orders")
	c.cancelAllCount.Inc()
	c.ack(m)
}
