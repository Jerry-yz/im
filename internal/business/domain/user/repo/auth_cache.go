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

func NewAuthCache() *AuthCache {
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

func (authCache *AuthCache) HGetAll(userId int) (map[int]*model.Device, error) {
	data, err := db.RedisClient.HGetAll(authKey + strconv.FormatInt(int64(userId), 10)).Result()
	if err != nil {
		zap.Error(errors.New("获取redis缓存失败"))
		return nil, err
	}
	res := make(map[int]*model.Device, len(data))
	device := new(model.Device)
	for k, v := range data {
		deviceId, err := strconv.Atoi(k)
		if err != nil {
			zap.Error(errors.New("err"))
			return nil, err
		}
		if err := json.Unmarshal([]byte(v), device); err != nil {
			zap.Error(errors.New("unmarshal error"))
			return nil, err
		}
		res[(deviceId)] = device
	}
	return res, nil
}
