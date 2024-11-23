package nalpaca

import "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"

var _ = Interface(&alpaca.Client{})
var _ = Interface(Mock{})
