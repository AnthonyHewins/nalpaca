package portfolio

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/AnthonyHewins/nalpaca/internal/protomap"
	"google.golang.org/protobuf/proto"
)

const PositionsKey = "positions"

func (c *Controller) UpdatePositionsKV(ctx context.Context) error {
	p, err := c.client.GetPositions()
	if err != nil {
		c.logger.ErrorContext(ctx, "failed getting positions", "err", err)
		return err
	}

	n := len(p)
	if n == 0 {
		c.logger.InfoContext(ctx, "no positions found in acct")
		return nil
	}

	positions := make([]*tradesvc.Position, n)
	for i := range p {
		positions[i] = protomap.PBFPosition(&p[i])
	}

	b, err := proto.Marshal(&tradesvc.Positions{Positions: positions})
	if err != nil {
		c.logger.ErrorContext(ctx, "failed marshal of positions", "err", err, "val", positions)
		return err
	}

	if _, err = c.portfolioKV.Put(ctx, PositionsKey, b); err != nil {
		c.logger.ErrorContext(ctx, "failed writing portfolio to KV store", "err", err, "positions", positions)
		return err
	}

	c.logger.InfoContext(ctx, "successfully wrote initial portfolio status")
	return nil
}
