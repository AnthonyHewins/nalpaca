package tradeupdater

import (
	"context"
	"errors"
	"time"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/AnthonyHewins/nalpaca/internal/protomap"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Controller) EventLoop(ctx context.Context) error {
	for {
		var lastMessage time.Time
		for {
			req := alpaca.StreamTradeUpdatesRequest{}
			if !lastMessage.IsZero() {
				req.Since = lastMessage.Add(time.Nanosecond)
			}

			err := c.client.StreamTradeUpdates(ctx, func(tu alpaca.TradeUpdate) {
				lastMessage = tu.At
				c.handler(tu)
			}, req)

			if err == nil {
				return nil
			}

			if errors.Is(err, context.Canceled) {
				c.logger.Warn("tradeupdater: ctx canceled", "ctx-err", ctx.Err())
				return err
			}

			c.logger.ErrorContext(ctx, "trade updater error", "err", err)
			return err
		}
	}
}

func (c *Controller) handler(u alpaca.TradeUpdate) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	l := c.logger.With("order", u)

	var t *timestamp.Timestamp
	if u.Timestamp != nil {
		t = timestamppb.New(*u.Timestamp)
	}

	buf, err := proto.Marshal(&tradesvc.TradeUpdate{
		At:          timestamppb.New(u.At),
		Event:       u.Event,
		EventId:     u.EventID,
		ExecutionId: u.ExecutionID,
		Order:       protomap.PBFOrder(&u.Order),
		PositionQty: u.PositionQty.String(),
		Price:       u.Price.String(),
		Qty:         u.Qty.String(),
		Timestamp:   t,
	})

	if err != nil {
		l.ErrorContext(ctx, "failed order marshal", "err", err)
		return
	}

	if err = c.nc.Publish(c.topicPrefix+"."+u.Order.Symbol, buf); err != nil {
		l.ErrorContext(ctx, "failed publishing order", "err", err)
	}
}
