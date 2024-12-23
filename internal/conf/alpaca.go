package conf

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AnthonyHewins/nalpaca/internal/nalpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
)

type Alpaca struct {
	Mock bool `env:"MOCK_ALPACA"`

	BaseURL   string `env:"ALPACA_URL" envDefault:"https://paper-api.alpaca.markets"`
	APIKey    string `env:"ALPACA_API_KEY"`
	APISecret string `env:"ALPACA_API_SECRET"`

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

	l := b.Logger.With(
		"apikey", a.APIKey,
		"len(secret)>0", len(a.APISecret) > 0,
		"baseURL", a.BaseURL,
		"oAuth", a.OAuth,
		"retryLimit", a.RetryLimit,
		"retryDelay", a.RetryDelay,
	)

	if len(a.APISecret) == 0 {
		l.Error("missing api secret")
		return nil, fmt.Errorf("no API secret for alpaca")
	}

	return alpaca.NewClient(alpaca.ClientOpts{
		APIKey:     a.APIKey,
		APISecret:  a.APISecret,
		BaseURL:    a.BaseURL,
		OAuth:      a.OAuth,
		RetryLimit: int(a.RetryLimit),
		RetryDelay: a.RetryDelay,
		HTTPClient: httpClient,
	}), nil
}
