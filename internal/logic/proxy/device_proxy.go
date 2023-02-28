package proxy

import (
	"context"
	"learn-im/pkg/protocol/pb"
)

var deviceProxy DeviceProxy

type DeviceProxy interface {
	ListDeviceByUserId(ctx context.Context, userId int64) ([]*pb.Device, error)
}
