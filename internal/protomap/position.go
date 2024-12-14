package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

func PBFPosition(x *alpaca.Position) *tradesvc.Position {
	return &tradesvc.Position{
		AssetId:                x.AssetID,
		Symbol:                 x.Symbol,
		Exchange:               x.Exchange,
		AssetClass:             PBFAssetClass(x.AssetClass),
		AssetMarginable:        x.AssetMarginable,
		Qty:                    x.Qty.String(),
		QtyAvailable:           x.QtyAvailable.String(),
		AvgEntryPrice:          x.AvgEntryPrice.String(),
		Side:                   x.Side,
		MarketValue:            ToString(x.MarketValue),
		CostBasis:              x.CostBasis.String(),
		UnrealizedPl:           ToString(x.UnrealizedPL),
		UnrealizedPlpc:         ToString(x.UnrealizedPLPC),
		UnrealizedIntradayPl:   ToString(x.UnrealizedIntradayPL),
		UnrealizedIntradayPlpc: ToString(x.UnrealizedIntradayPLPC),
		CurrentPrice:           ToString(x.CurrentPrice),
		LastdayPrice:           ToString(x.LastdayPrice),
		ChangeToday:            ToString(x.ChangeToday),
	}
}
