package protomap

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func PBFTimestamp(x *time.Time) *timestamppb.Timestamp {
	if x == nil {
		return nil
	}

	y := timestamppb.New(*x)
	return y
}
