package conf

import (
	"net/http"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Alpaca struct {
	BaseURL   string `env:"ALPACA_URL" envDefault:"https://paper-api.alpaca.markets"`
	APIKey    string `env:"API_KEY,required"`
	APISecret string `env:"API_SECRET,required"`

	OAuth      string        `env:"ALPACA_OAUTH" envDefault:""`
	RetryLimit uint          `env:"ALPACA_RETRY_LIMIT"`
	RetryDelay time.Duration `env:"ALPACA_RETRY_LIMIT"`
}

func (b *Bootstrapper) Alpaca(a *Alpaca, httpClient *http.Client) *alpaca.Client {
	return alpaca.NewClient(alpaca.ClientOpts{
		APIKey:     a.APIKey,
		APISecret:  a.APISecret,
		BaseURL:    a.BaseURL,
		OAuth:      a.OAuth,
		RetryLimit: int(a.RetryLimit),
		RetryDelay: a.RetryDelay,
		HTTPClient: httpClient,
	})
}
