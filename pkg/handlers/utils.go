package handlers

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func TimeToTimestampProto(in *time.Time) *timestamppb.Timestamp {
	if in == nil {
		return nil
	}

	return timestamppb.New(*in)
}
