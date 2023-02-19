package repo

import (
	"encoding/json"
	"errors"
	"learn-im/internal/business/domain/user/model"
	"learn-im/pkg/db"
	"strconv"

	"go.uber.org/zap"
)

const authKey = "Auth:"

type AuthCache struct {
}

func NewCacheUser() *AuthCache {
	return &AuthCache{}
}

func (authCache *AuthCache) Get(userId, deviceId int) (*model.Device, error) {
	byt, err := db.RedisClient.HGet(authKey+strconv.FormatInt(int64(userId), 10), strconv.FormatInt(int64(deviceId), 10)).Bytes()
	if err != nil {
		zap.Error(errors.New("get user auth cache error"))
		return nil, err
	}
	device := new(model.Device)
	if err := json.Unmarshal(byt, &device); err != nil {
		zap.Error(errors.New("unmarshal error"))
		return nil, err
	}
	return device, nil
}

func (authCache *AuthCache) HSet(userId, deviceId int, device *model.Device) error {
	byt, err := json.Marshal(device)
	if err != nil {
		zap.Error(errors.New("marshal error"))
		return err
	}
	if _, err := db.RedisClient.HSet(authKey+strconv.FormatInt(int64(userId), 10), strconv.FormatInt(int64(deviceId), 10), byt).Result(); err != nil {
		zap.Error(errors.New("set user auth cache error"))
		return err
	}
	return nil
}
