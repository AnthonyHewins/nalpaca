package main

import (
	"context"

	"github.com/AnthonyHewins/nalpaca/internal/optionquotes"
	"github.com/nats-io/nats.go/jetstream"
)

type Quotes struct {
	DisableQuotes bool   `env:"DISABLE_QUOTES" envDefault:"false"`
	QuotesStream  string `env:"QUOTES_STREAM" envDefault:"nalpaca-quotes-stream-v0"`
}

func (a *app) initQuotes(ctx context.Context, js jetstream.JetStream, c *config) error {
	if c.Quotes.DisableQuotes {
		return nil
	}

	a.quotes = optionquotes.NewController(a.Logger, a.Nalpaca, js)
	return nil
}
