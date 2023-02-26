package device

import (
	"context"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
	"learn-im/pkg/rpc"
)

type DeviceService struct {
	*DeviceRepo
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		NewDeviceRepo(),
	}
}

//注册设备
func (d *DeviceService) RegisterDevice(ctx context.Context, device *Device) error {
	return d.DeviceRepo.Save(device)
}

func (d *DeviceService) SignDevice(ctx context.Context, userId, deviceId int, token, connAddr, clientAddr string) error {
	_, err := rpc.GetBusinessIntClient().Auth(ctx, &pb.AuthReq{UserId: int64(userId), DeviceId: int64(deviceId), Token: token})
	if err != nil {
		return gerrors.WarpError(err)
	}
	device, err := d.DeviceRepo.GetDevice(deviceId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	device.OnLien(userId, connAddr, clientAddr)
	return d.DeviceRepo.Save(device)
}
