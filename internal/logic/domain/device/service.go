package device

import (
	"context"
	"learn-im/logger"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
	"learn-im/pkg/rpc"
	"time"

	"go.uber.org/zap"
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
	if device == nil {
		return nil
	}
	device.OnLien(userId, connAddr, clientAddr)
	return d.DeviceRepo.Save(device)
}

func (d *DeviceService) Auth(ctx context.Context, userId, deviceId int, token string) error {
	_, err := rpc.GetBusinessIntClient().Auth(ctx, &pb.AuthReq{
		UserId:   int64(userId),
		DeviceId: int64(deviceId),
		Token:    token,
	})
	if err != nil {
		return gerrors.WarpError(err)
	}
	return nil
}

func (d *DeviceService) ListOnlineByUserId(ctx context.Context, userId int) ([]*pb.Device, error) {
	deviceList, err := d.DeviceRepo.ListOnlineByUserId(userId)
	if err != nil {
		return []*pb.Device{}, gerrors.WarpError(err)
	}
	devices := make([]*pb.Device, len(deviceList))
	for k := range deviceList {
		devices[k] = deviceList[k].ToProto()
	}
	return devices, nil
}

func (d *DeviceService) ServerStop(ctx context.Context, connAddr string) error {
	devices, err := d.DeviceRepo.ListOnlineByConnAddr(connAddr)
	if err != nil {
		return gerrors.WarpError(err)
	}
	for k := range devices {
		if err := d.DeviceRepo.UpdateDeviceStatus(devices[k]); err != nil {
			logger.Logger.Error("update status error", zap.Any("device", devices[k]))
			return gerrors.WarpError(err)
		}
		time.Sleep(2 * time.Millisecond)
	}
	return nil
}
