package device

import (
	"learn-im/logger"
	"learn-im/pkg/db"
	"learn-im/pkg/util"
	"strconv"
	"time"
)

const (
	DeviceCacheKey    = "user_device:"
	DeviceCacheExpire = 2 * time.Hour
)

type DeviceCache struct {
}

func NewDeviceCache() *DeviceCache {
	return &DeviceCache{}
}

func (d *DeviceCache) SetDeviceCache(userId int, devices []*Device) error {
	return db.RedisClient.Set(DeviceCacheKey+strconv.FormatInt(int64(userId), 10), devices, DeviceCacheExpire).Err()
}

func (d *DeviceCache) GetDeviceCache(userId int) ([]*Device, error) {
	devices := make([]*Device, 0)
	if err := util.NewRedisUtil().Get(DeviceCacheKey+strconv.FormatInt(int64(userId), 10), &devices); err != nil {
		logger.Sugar.Error(err)
		return devices, err
	}
	return devices, nil
}

func (d *DeviceCache) DelDeviceCache(userId int) error {
	return db.RedisClient.Del(DeviceCacheKey + strconv.FormatInt(int64(userId), 10)).Err()
}
