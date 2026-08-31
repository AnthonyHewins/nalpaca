package conf

import (
	"net/http"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/streaming"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Alpaca struct {
	StockStream  streaming.StreamConfig `envPrefix:"STOCK_STREAM_"`
	NewsStream   streaming.StreamConfig `envPrefix:"NEWS_STREAM_"`
	OptionStream streaming.StreamConfig `envPrefix:"OPTIONS_STREAM_"`

	APIBaseURL string `env:"ALPACA_API_URL" envDefault:"https://paper-api.alpaca.markets"`
	// This may appear scary because we default to a prod URL, but that's what works.
	// There is no sandbox; but, you're streaming data, not mutating anything, so it's fine
	StreamingBaseURL string `env:"ALPACA_STREAMING_URL" envDefault:"https://stream.data.alpaca.markets"`

	APIKey    string `env:"ALPACA_API_KEY,required"`
	APISecret string `env:"ALPACA_API_SECRET,required"`

	OAuth      string        `env:"ALPACA_OAUTH"`
	RetryLimit uint          `env:"ALPACA_RETRY_LIMIT"`
	RetryDelay time.Duration `env:"ALPACA_RETRY_DELAY"`
}

func (b *Bootstrapper) Alpaca(a *Alpaca, httpClient *http.Client) (*alpaca.Client, error) {
	secret := strings.TrimSpace(a.APISecret)
	l := b.Logger.With(
		"apikey", a.APIKey,
		"len(secret)>0 after trimming spaces", len(a.APISecret) > 0,
		"baseURL", a.APIBaseURL,
		"oAuth", a.OAuth,
		"retryLimit", a.RetryLimit,
		"retryDelay", a.RetryDelay,
	)

	l.Info("created alpaca client")
	return alpaca.NewClient(alpaca.ClientOpts{
		APIKey:     a.APIKey,
		APISecret:  secret,
		BaseURL:    a.APIBaseURL,
		OAuth:      a.OAuth,
		RetryLimit: int(a.RetryLimit),
		RetryDelay: a.RetryDelay,
		HTTPClient: httpClient,
	}), nil
}
