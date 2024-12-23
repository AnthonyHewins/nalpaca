package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/internal/portfolio"
	"github.com/nats-io/nats.go/jetstream"
)

type TradeUpdaterConf struct {
	DisableTradeUpdater bool   `env:"DISABLE_TRADE_UPDATER" envDefault:"false"`
	TradeUpdaterStream  string `env:"TRADE_UPDATER_STREAM_NAME" envDefault:"nalpaca-tradeupdater-stream-v0"`
}

func (a *app) initTradeUpdater(ctx context.Context, js jetstream.JetStream, kv jetstream.KeyValue, c *config) error {
	if c.TradeUpdaterConf.DisableTradeUpdater {
		return nil
	}

	a.updater = portfolio.NewController(a.Logger, a.Nalpaca, c.ProcessingTimeout, js, kv, c.StreamPrefix)
	return nil
}
