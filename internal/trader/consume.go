package trader

import (
	"context"

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

	trade, err := c.getMsg(m)
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
