package nalpaca

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/AnthonyHewins/nalpaca/internal/portfolio"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

func (c *Client) GetPositions(ctx context.Context) ([]*tradesvc.Position, error) {
	positions, err := c.kv.Get(ctx, portfolio.PositionsKey)
	if err != nil {
		return nil, err
	}

	var p tradesvc.Positions
	if err = proto.Unmarshal(positions.Value(), &p); err != nil {
		return nil, err
	}

	return p.Positions, nil
}

func (c *Client) WatchPositions(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
	return c.kv.Watch(ctx, portfolio.PositionsKey, opts...)
}
