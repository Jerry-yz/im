package repo

import (
	"encoding/json"
	// "learn-im/internal/logic/domain/group"

	"learn-im/internal/logic/domain/group/entity"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"
	"learn-im/pkg/util"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	// "k8s.io/apimachinery/pkg/util/json"
)

const GroupKey = "group:"

type GroupCache struct {
}

func NewGroupCache() *GroupCache {
	return &GroupCache{}
}

func (g *GroupCache) GetGroupCache(groupId int) (*entity.Group, error) {
	group := new(entity.Group)
	if err := util.NewRedisUtil().Get(GroupKey+strconv.FormatInt(int64(groupId), 10), &group); err != nil && err != redis.Nil {
		return group, gerrors.WarpError(err)
	}
	return group, nil
}

func (g *GroupCache) SetGroupCache(group *entity.Group) error {
	byt, err := json.Marshal(group)
	if err != nil {
		return gerrors.WarpError(err)
	}
	return util.NewRedisUtil().Set(GroupKey+strconv.FormatInt(int64(group.ID), 10), byt, 24*time.Hour)
}

func (g *GroupCache) DelGroupCache(groupId int) error {
	return db.RedisClient.Del(GroupKey + strconv.FormatInt(int64(groupId), 10)).Err()
}
