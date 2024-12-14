package main

import "github.com/AnthonyHewins/nalpaca/internal/tradeupdater"

type tradeupdaterConf struct {
	DisableTradeUpdater bool   `env:"DISABLE_TRADE_UPDATER" envDefault:"false"`
	TradeUpdaterStream  string `env:"TRADE_UPDATER_STREAM_NAME" envDefault:"nalpaca-tradeupdate-stream-v0"`
}

func (a *app) initTradeUpdater(c *config) {
	if c.DisableTradeUpdater {
		return
	}

	a.updater = tradeupdater.NewController(a.Logger, a.Nalpaca, c.ProcessingTimeout, a.NC, c.StreamPrefix)
}
