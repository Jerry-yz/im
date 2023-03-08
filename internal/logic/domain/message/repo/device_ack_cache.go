package repo

import (
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"
	"strconv"
)

const DeviceACKKey = "device_ack:"

type DeviceACKCache struct {
}

func NewDeviceACKCache() *DeviceACKCache {
	return &DeviceACKCache{}
}

func (d *DeviceACKCache) HSet(userId, deviceId, ack int) error {
	_, err := db.RedisClient.HSet(DeviceACKKey+strconv.FormatInt(int64(userId), 10), strconv.FormatInt(int64(deviceId), 10), strconv.FormatInt(int64(ack), 10)).Result()
	if err != nil {
		return gerrors.WarpError(err)
	}
	return nil
}

func (d *DeviceACKCache) HGet(userId int) (map[int64]int64, error) {
	res, err := db.RedisClient.HGetAll(DeviceACKKey + strconv.FormatInt(int64(userId), 10)).Result()
	if err != nil {
		return nil, gerrors.WarpError(err)
	}
	acks := make(map[int64]int64)
	for k, v := range res {
		deviceId, _ := strconv.ParseInt(k, 10, 64)
		ack, _ := strconv.ParseInt(v, 10, 64)
		acks[deviceId] = ack
	}
	return acks, nil
}
