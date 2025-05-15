package conf

import (
	"net/http"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/streaming"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Alpaca struct {
	// DisableOptionStream bool             `env:"DISABLE_OPTIONS_STREAM"`
	// OptionsStream       streaming.Stream `envPrefix:"OPTIONS_STREAM"`

	EnableStockStream bool             `env:"ENABLE_STOCK_STREAM" envDefault:"false"`
	StockStream       streaming.Stream `envPrefix:"STOCK_STREAM"`

	BaseURL   string `env:"ALPACA_URL" envDefault:"https://paper-api.alpaca.markets"`
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
