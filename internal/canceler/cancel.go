package canceler

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
)

type Canceler struct {
	counters Counters
	logger   *slog.Logger
	client   nalpaca.Interface
	timeout  time.Duration
}

type Counters struct {
	CancelCount, CancelFail       prometheus.Counter
	CancelAllCount, CancelAllFail prometheus.Counter
}

func New(logger *slog.Logger, client nalpaca.Interface, counters Counters, timeout time.Duration) *Canceler {
	return &Canceler{
		counters: counters,
		logger:   logger,
		client:   client,
		timeout:  timeout,
	}
}

func (c *Canceler) ack(m jetstream.Msg) {
	if err := m.Ack(); err != nil {
		c.logger.Error("failed ack", "err", err)
	}
}

func (c *Canceler) nak(m jetstream.Msg) {
	if err := m.Nak(); err != nil {
		c.logger.Error("failed nak", "err", err)
	}
}

func (c *Canceler) term(m jetstream.Msg, reason string) {
	if err := m.TermWithReason(reason); err != nil {
		c.logger.Error("failed term", "err", err, "reason", reason)
	}
}

func (c *Canceler) EventLoop(m jetstream.Msg) {
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
		c.counters.CancelFail.Inc()
		c.nak(m)
		return
	}

	l.InfoContext(ctx, "successful cancel")
	c.counters.CancelCount.Inc()
	c.ack(m)
}

func (c *Canceler) cancelAll(m jetstream.Msg) {
	c.logger.Warn("received cancel all orders request")
	if err := c.client.CancelAllOrders(); err != nil {
		c.logger.Error("failed canceling all orders", "err", err)
		c.counters.CancelAllFail.Inc()
		c.nak(m)
		return
	}

	c.logger.Info("successfully canceled all orders")
	c.counters.CancelAllCount.Inc()
	c.ack(m)
}
