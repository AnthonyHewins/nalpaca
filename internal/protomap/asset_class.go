package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v1"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFAssetClass(x alpaca.AssetClass) tradesvc.AssetClass {
	switch x {
	case alpaca.Crypto:
		return tradesvc.AssetClass_ASSET_CLASS_CRYPTO
	case alpaca.USEquity:
		return tradesvc.AssetClass_ASSET_CLASS_US_EQUITY
	default:
		return tradesvc.AssetClass_ASSET_CLASS_UNSPECIFIED
	}
}
