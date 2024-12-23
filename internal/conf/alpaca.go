package conf

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
)

type Alpaca struct {
	Mock bool `env:"MOCK_ALPACA"`

	BaseURL   string `env:"ALPACA_URL" envDefault:"https://paper-api.alpaca.markets"`
	APIKey    string `env:"ALPACA_API_KEY,required"`
	APISecret string `env:"ALPACA_API_SECRET,required"`

	OAuth      string        `env:"ALPACA_OAUTH" envDefault:""`
	RetryLimit uint          `env:"ALPACA_RETRY_LIMIT"`
	RetryDelay time.Duration `env:"ALPACA_RETRY_LIMIT"`
}

func (b *Bootstrapper) Alpaca(a *Alpaca, httpClient *http.Client) (nalpaca.Interface, error) {
	if a.Mock {
		b.Logger.Info("mock set, mocking Alpaca")
		return nalpaca.Mock{
			StreamTradeUpdatesFn: func(ctx context.Context, fn func(alpaca.TradeUpdate), req alpaca.StreamTradeUpdatesRequest) error {
				ticker := time.NewTicker(time.Second * 5)
				for {
					select {
					case <-ticker.C:
						fn(alpaca.TradeUpdate{
							At:          time.Now(),
							Event:       "",
							EventID:     "",
							ExecutionID: "",
							Order:       alpaca.Order{},
							PositionQty: &decimal.Decimal{},
							Price:       &decimal.Decimal{},
							Qty:         &decimal.Decimal{},
							Timestamp:   &time.Time{},
						})
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			},
			PlaceOrderFn: func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
				return &alpaca.Order{}, nil
			},
		}, nil
	}

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
