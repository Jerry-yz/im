package service

import (
	"context"
	"learn-im/internal/logic/domain/message/repo"
	"learn-im/pkg/gerrors"
)

type DeviceServiceACK struct {
	*repo.DeviceACKCache
}

func NewDeviceACK() *DeviceServiceACK {
	return &DeviceServiceACK{
		repo.NewDeviceACKCache(),
	}
}

func (d *DeviceServiceACK) UpdateACK(ctx context.Context, userId, deviceId, ack int) error {
	if ack < 0 {
		return nil
	}
	return d.DeviceACKCache.HSet(userId, deviceId, ack)
}

func (d *DeviceServiceACK) GetMaxByUserId(ctx context.Context, userId int) (int, error) {
	acks, err := d.DeviceACKCache.HGet(userId)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	max := 0
	for _, ack := range acks {
		if int(ack) > max {
			max = int(ack)
		}
	}
	return max, nil
}
