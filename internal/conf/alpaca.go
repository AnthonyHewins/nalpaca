package conf

import (
	"net/http"
	"strings"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

type Alpaca struct {
	OptionsStream

	BaseURL   string `env:"ALPACA_URL" envDefault:"https://paper-api.alpaca.markets"`
	APIKey    string `env:"ALPACA_API_KEY,required"`
	APISecret string `env:"ALPACA_API_SECRET,required"`

	OAuth      string        `env:"ALPACA_OAUTH" envDefault:""`
	RetryLimit uint          `env:"ALPACA_RETRY_LIMIT"`
	RetryDelay time.Duration `env:"ALPACA_RETRY_LIMIT"`
}

func (b *Bootstrapper) Alpaca(a *Alpaca, httpClient *http.Client) (*alpaca.Client, error) {
	secret := strings.TrimSpace(a.APISecret)
	l := b.Logger.With(
		"apikey", a.APIKey,
		"len(secret)>0 after trimming spaces", len(a.APISecret) > 0,
		"baseURL", a.BaseURL,
		"oAuth", a.OAuth,
		"retryLimit", a.RetryLimit,
		"retryDelay", a.RetryDelay,
	)

	l.Info("created alpaca client")
	return alpaca.NewClient(alpaca.ClientOpts{
		APIKey:     a.APIKey,
		APISecret:  secret,
		BaseURL:    a.BaseURL,
		OAuth:      a.OAuth,
		RetryLimit: int(a.RetryLimit),
		RetryDelay: a.RetryDelay,
		HTTPClient: httpClient,
	}), nil
}

//go:generate enumer -type optionFeed -text
type optionFeed byte

const (
	indicative optionFeed = iota
	opra
)

type OptionsStream struct {
	BaseURL    string     `env:"OPTIONS_STREAM_BASE_URL" envDefault:"wss://stream.data.sandbox.alpaca.markets"`
	Feed       optionFeed `env:"OPTIONS_STREAM_FEED_TYPE" envDefault:"indicative"`
	Processors uint16     `env:"OPTIONS_STREAM_GOROUTINES" envDefault:"1"`
}

func (b *Bootstrapper) OptionsStream(a *Alpaca) *stream.OptionClient {
	o := a.OptionsStream

	return stream.NewOptionClient(marketdata.IEX,
		stream.WithBaseURL(o.BaseURL),
		stream.WithCredentials(a.APIKey, a.APISecret),
		stream.WithProcessors(int(o.Processors)),
	)
}
