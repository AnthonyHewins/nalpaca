package canceler

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/bridge"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
)

type Canceler struct {
	counters Counters
	logger   *slog.Logger
	client   bridge.AlpacaInterface
	timeout  time.Duration
}

type Counters struct {
	CancelCount, CancelFail       prometheus.Counter
	CancelAllCount, CancelAllFail prometheus.Counter
}

// Collectors returns the counters for registration with a prometheus registry.
func (c Counters) Collectors() []prometheus.Collector {
	return []prometheus.Collector{c.CancelCount, c.CancelFail, c.CancelAllCount, c.CancelAllFail}
}

func New(logger *slog.Logger, client bridge.AlpacaInterface, counters Counters, timeout time.Duration) *Canceler {
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

	id := string(m.Data())

	if strings.ToUpper(id) == "ALL" {
		c.cancelAll(m)
		return
	}

	if err := c.client.CancelOrder(id); err != nil {
		c.logger.ErrorContext(ctx, "failed canceling order", "err", err)
		c.counters.CancelFail.Inc()
		c.nak(m)
		return
	}

	c.logger.InfoContext(ctx, "successful cancel")
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
