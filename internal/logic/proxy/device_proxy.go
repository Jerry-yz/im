package proxy

import (
	"context"
	"learn-im/pkg/protocol/pb"
)

var DevProxy DeviceProxy

type DeviceProxy interface {
	ListOnlineByUserId(ctx context.Context, userId int64) ([]*pb.Device, error)
}
