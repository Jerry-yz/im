 package device

import (
	"errors"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"

	"gorm.io/gorm"
)

type Dao struct {
}

func NewDao() *Dao {
	return &Dao{}
}

func (d *Dao) Save(device *Device) error {
	return db.DB.Save(device).Error
}

func (d *Dao) GetDevice(deviceId int) (*Device, error) {
	device := new(Device)
	if err := db.DB.Where("id", deviceId).First(device).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return device, gerrors.WarpError(err)
	}
	return device, nil
}

func (d *Dao) ListOnlineByUserId(userId int) ([]*Device, error) {
	devices := make([]*Device, 0)
	if err := db.DB.Where("user_id = ? and status = ?", userId, DeviceOnLine).Find(&devices).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return devices, err
	}
	return devices, nil
}

func (d *Dao) ListOnlineByConnAddr(connAddr string) ([]*Device, error) {
	devices := make([]*Device, 0)
	if err := db.DB.Where("status = ? and conn_addr = ?", DeviceOnLine, connAddr).First(devices).Error; err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		return devices, err
	}
	return devices, nil
}

func (d *Dao) UpdateStatus(deviceId, status int, connAddr string) (int, error) {
	db := db.DB.Where("id = ? and conn_addr = ?", deviceId, connAddr).Update("status", status)
	if db.Error != nil {
		return 0, gerrors.WarpError(db.Error)
	}
	return int(db.RowsAffected), nil
}
