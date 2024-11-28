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
		Notional:       o.Notional.String(),
		Qty:            o.Qty.String(),
		FilledQty:      o.FilledQty.String(),
		FilledAvgPrice: o.FilledAvgPrice.String(),
		LimitPrice:     o.LimitPrice.String(),
		StopPrice:      o.StopPrice.String(),
		TrailPrice:     o.TrailPrice.String(),
		TrailPercent:   o.TrailPercent.String(),
		Hwm:            o.HWM.String(),
		ExtHours:       o.ExtendedHours,
		Legs:           legs,
	}
}
