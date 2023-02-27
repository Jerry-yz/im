package device

import (
	"context"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/protocol/pb"
)

type DeviceApp struct {
	*DeviceRepo
	*DeviceService
}

func NewDeviceApp() *DeviceApp {
	return &DeviceApp{
		NewDeviceRepo(),
		NewDeviceService(),
	}
}

func (d *DeviceApp) Register(ctx context.Context, req *pb.Device) (int64, error) {
	device := &Device{
		Id:            int(req.DeviceId),
		UserId:        int(req.UserId),
		Type:          int(req.Type),
		Brand:         req.Brand,
		Model:         req.Model,
		SystemVersion: req.SystemVersion,
		SdkVersion:    req.SdkVersion,
		Status:        int(req.Status),
		ConnAddr:      req.ConnAddr,
		ClientAddr:    req.ClientAddr,
	}
	if !device.IsLegal() {
		return 0, gerrors.ErrBadRequest
	}
	if err := d.DeviceRepo.Save(device); err != nil {
		return 0, gerrors.WarpError(err)
	}
	return int64(device.Id), nil
}

func (d *DeviceApp) SignIn(ctx context.Context, userId, deviceId int, token, connAddr, clientAddr string) error {
	return d.DeviceService.SignDevice(ctx, userId, deviceId, token, connAddr, clientAddr)
}

func (d *DeviceApp) Offline(ctx context.Context, deviceId int, clientAddr string) error {
	device, err := d.DeviceRepo.GetDevice(deviceId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if device == nil {
		return nil
	}
	if device.ClientAddr != clientAddr {
		return nil
	}
	device.Status = DeViceOffLie
	return d.DeviceRepo.Save(device)
}

func (d *DeviceApp) ListOnlineByUserId(ctx context.Context, userId int) ([]*pb.Device, error) {
	return d.DeviceService.ListOnlineByUserId(ctx, userId)
}

func (d *DeviceApp) GetDevice(ctx context.Context, deviceId int) (*pb.Device, error) {
	device, err := d.DeviceRepo.GetDevice(deviceId)
	if err != nil {
		return &pb.Device{}, gerrors.WarpError(err)
	}
	return device.ToProto(), nil
}

func (d *DeviceApp) ServerStop(ctx context.Context, connAddr string) error {
	return d.DeviceService.ServerStop(ctx, connAddr)
}
