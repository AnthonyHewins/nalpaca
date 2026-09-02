package conf

import (
	"net/http"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/streaming"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Alpaca struct {
	StockStream  streaming.StocksConfig  `envPrefix:"STOCK_STREAM_"`
	NewsStream   streaming.NewsConfig    `envPrefix:"NEWS_STREAM_"`
	OptionStream streaming.OptionsConfig `envPrefix:"OPTIONS_STREAM_"`

	APIBaseURL string `env:"ALPACA_API_URL" envDefault:"https://paper-api.alpaca.markets" desc:"Base URL for the Alpaca trading REST API"`
	// This may appear scary because we default to a prod URL, but that's what works.
	// There is no sandbox; but, you're streaming data, not mutating anything, so it's fine
	StreamingBaseURL string `env:"ALPACA_STREAMING_URL" envDefault:"https://stream.data.alpaca.markets" desc:"Base URL for Alpaca's market data streaming API; defaults to the prod endpoint since there's no sandbox for streaming"`

	APIKey    string `env:"ALPACA_API_KEY,required" desc:"Alpaca API key ID used to authenticate REST and streaming requests" sensitive:"true"`
	APISecret string `env:"ALPACA_API_SECRET,required" desc:"Alpaca API secret key used to authenticate REST and streaming requests" sensitive:"true"`

	OAuth      string        `env:"ALPACA_OAUTH" desc:"OAuth token to use instead of API key/secret auth, if set" sensitive:"true"`
	RetryLimit uint          `env:"ALPACA_RETRY_LIMIT" desc:"Number of times to retry a failed Alpaca API request"`
	RetryDelay time.Duration `env:"ALPACA_RETRY_DELAY" desc:"Delay between retries of a failed Alpaca API request"`
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
