package bridge

import "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"

var _ = AlpacaInterface(&alpaca.Client{})
var _ = AlpacaInterface(Mock{})
