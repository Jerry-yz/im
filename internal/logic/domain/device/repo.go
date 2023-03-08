package device

import (
	"learn-im/logger"
	"learn-im/pkg/gerrors"
)

type DeviceRepo struct {
	*Dao
	*DeviceCache
}

func NewDeviceRepo() *DeviceRepo {
	return &DeviceRepo{
		NewDao(),
		NewDeviceCache(),
	}
}

func (d *DeviceRepo) GetDevice(deviceId int) (*Device, error) {
	return d.Dao.GetDevice(deviceId)
}

func (d *DeviceRepo) Save(device *Device) error {
	if device.Id != 0 {
		if err := d.DeviceCache.DelDeviceCache(device.UserId); err != nil {
			logger.Sugar.Error(err)
			return err
		}
	}
	return d.Dao.Save(device)
}

func (d *DeviceRepo) ListOnlineByUserId(userId int) ([]*Device, error) {
	devices := make([]*Device, 0)
	devices, err := d.DeviceCache.GetDeviceCache(userId)
	if err != nil {
		return devices, gerrors.WarpError(err)
	}
	if len(devices) > 0 {
		return devices, nil
	}
	return d.Dao.ListOnlineByUserId(userId)
}

func (d *DeviceRepo) ListOnlineByConnAddr(connAddr string) ([]*Device, error) {
	return d.Dao.ListOnlineByConnAddr(connAddr)
}

func (d *DeviceRepo) UpdateDeviceStatus(device *Device) error {
	rows, err := d.Dao.UpdateStatus(device.Id, DeViceOffLie, device.ConnAddr)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if rows == 1 && device.Id != 0 {
		if err := d.DeviceCache.DelDeviceCache(device.UserId); err != nil {
			return err
		}
	}
	return nil
}
