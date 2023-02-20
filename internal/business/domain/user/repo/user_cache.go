package repo

import (
	"encoding/json"
	"errors"
	"learn-im/internal/business/domain/user/model"
	"learn-im/pkg/db"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	userCache  = "user:"
	expireTime = 2 * time.Hour
)

type UserCache struct {
}

func NewUserCache() *UserCache {
	return &UserCache{}
}

func (u *UserCache) Get(userId int) (*model.User, error) {
	byt, err := db.RedisClient.Get(userCache + strconv.FormatInt(int64(userId), 10)).Bytes()
	if err != nil {
		zap.Error(errors.New("get user cache error"))
		return nil, err
	}
	user := new(model.User)
	if err := json.Unmarshal(byt, &user); err != nil {
		zap.Error(errors.New("unmarshal error"))
		return nil, err
	}
	return user, nil
}

func (u *UserCache) Set(userId int, user *model.User) error {
	byt, err := json.Marshal(user)
	if err != nil {
		zap.Error(errors.New("marshal error"))
		return err
	}
	_, err = db.RedisClient.Set(userCache+strconv.FormatInt(int64(userId), 10), byt, expireTime).Result()
	if err != nil {
		zap.Error(errors.New("set user cache error"))
		return err
	}
	return nil
}

func (u *UserCache) Delete(userId int) error {
	_, err := db.RedisClient.Del(authKey + strconv.FormatInt(int64(userId), 10)).Result()
	if err != nil {
		zap.Error(errors.New("delete user cache error"))
		return err
	}
	return nil
}
