package util

import (
	"encoding/json"
	"learn-im/logger"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
)

type RedisUtil struct {
	Client *redis.Client
}

func NewRedisUtil() *RedisUtil {
	return &RedisUtil{}
}

func (r *RedisUtil) Get(key string, value interface{}) error {
	byt, err := db.RedisClient.Get(key).Bytes()
	if err != nil {
		return gerrors.WarpError(err)
	}
	if err := jsoniter.Unmarshal(byt, &value); err != nil {
		return gerrors.WarpError(err)
	}
	return nil
}

func (r *RedisUtil) Set(key string, value interface{}, duration time.Duration) error {
	byt, err := json.Marshal(value)
	if err != nil {
		logger.Sugar.Error(err)
		return gerrors.WarpError(err)
	}
	return db.RedisClient.Set(key, byt, duration).Err()
}
