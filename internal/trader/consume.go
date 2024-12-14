package trader

import (
	"context"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
)

func (c *Controller) term(ctx context.Context, m jetstream.Msg, reason string) {
	c.logger.ErrorContext(ctx, "terminating msg", "reason", reason)
	if err := m.TermWithReason(reason); err != nil {
		c.logger.ErrorContext(ctx, "failed termination", "reason", reason, "err", err)
	}

	c.Counters.OrderFailCount.Inc()
}

func (c *Controller) nak(ctx context.Context, m jetstream.Msg) {
	if err := m.Nak(); err != nil {
		c.logger.ErrorContext(ctx, "failed nak", "err", err)
	}

	c.Counters.OrderFailCount.Inc()
}

func (c *Controller) ack(ctx context.Context, m jetstream.Msg) {
	if err := m.Ack(); err != nil {
		c.logger.ErrorContext(ctx, "failed ACK", "err", err)
		return
	}

	c.Counters.OrderCreatedCount.Inc()
}

func (c *Controller) Consume(m jetstream.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), c.processingTimeout)
	defer cancel()

	clientOrderID, err := c.getOrderID(m.Subject())
	if err != nil {
		c.term(ctx, m, err.Error())
		return
	}

	trade, err := c.getMsg(m, clientOrderID)
	if err != nil {
		c.term(ctx, m, err.Error())
		return
	}

	if err = c.trade(ctx, trade); err != nil {
		c.nak(ctx, m)
		return
	}

	c.ack(ctx, m)
}

func (c *Controller) getOrderID(subj string) (string, error) {
	l := c.logger.With("subj", subj)

	s := strings.Split(subj, ".")
	n := len(s)
	if n == 0 {
		l.Error("client order ID was invalid", "subj", subj)
		return "", fmt.Errorf("invalid client order id in NATS subject: %s", subj)
	}

	id := s[n-1]
	if len(id) > 128 {
		l.Error("order ID is too big", "id", id)
		return "", fmt.Errorf("max order ID size is 128, got %s", id)
	}

	return id, nil
}
