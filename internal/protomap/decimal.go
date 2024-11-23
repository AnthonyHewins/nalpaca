package protomap

import "github.com/shopspring/decimal"

func Decimal(s string) (*decimal.Decimal, error) {
	if s == "" {
		return nil, nil
	}

	x, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}

	return &x, nil
}
