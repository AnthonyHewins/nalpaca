package protomap

import "github.com/shopspring/decimal"

func ToDecimal(s string) (*decimal.Decimal, error) {
	if s == "" {
		return nil, nil
	}

	x, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}

	return &x, nil
}

func ToString(x *decimal.Decimal) string {
	if x == nil {
		return ""
	}

	return x.String()
}
