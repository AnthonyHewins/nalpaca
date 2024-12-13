package protomap

import (
	"github.com/AnthonyHewins/nalpaca/gen/go/tradesvc/v0"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/utils/ptr"
)

func PBFOrder(o *alpaca.Order) *tradesvc.Order {
	if o == nil {
		return nil
	}

	legs := make([]*tradesvc.Order, len(o.Legs))
	for i := range o.Legs {
		legs[i] = PBFOrder(&o.Legs[i])
	}

	return &tradesvc.Order{
		Id:             o.ID,
		ClientOrderId:  o.ClientOrderID,
		CreatedAt:      timestamppb.New(o.CreatedAt),
		UpdatedAt:      timestamppb.New(o.UpdatedAt),
		SubmittedAt:    timestamppb.New(o.SubmittedAt),
		FilledAt:       PBFTimestamp(o.FilledAt),
		ExpiredAt:      PBFTimestamp(o.ExpiredAt),
		CanceledAt:     PBFTimestamp(o.CanceledAt),
		FailedAt:       PBFTimestamp(o.FailedAt),
		ReplacedAt:     PBFTimestamp(o.ReplacedAt),
		ReplacedBy:     ptr.Deref(o.ReplacedBy, ""),
		Replaces:       ptr.Deref(o.Replaces, ""),
		AssetId:        o.AssetID,
		Symbol:         o.Symbol,
		AssetClass:     PBFAssetClass(o.AssetClass),
		OrderClass:     PBFOrderClass(o.OrderClass),
		Type:           PBFOrderType(o.Type),
		Side:           PBFSide(o.Side),
		TimeInForce:    PBFTIF(o.TimeInForce),
		Status:         o.Status,
		Notional:       ToString(o.Notional),
		Qty:            ToString(o.Qty),
		FilledQty:      o.FilledQty.String(),
		FilledAvgPrice: ToString(o.FilledAvgPrice),
		LimitPrice:     ToString(o.LimitPrice),
		StopPrice:      ToString(o.StopPrice),
		TrailPrice:     ToString(o.TrailPrice),
		TrailPercent:   ToString(o.TrailPercent),
		Hwm:            ToString(o.HWM),
		ExtHours:       o.ExtendedHours,
		Legs:           legs,
	}
}
