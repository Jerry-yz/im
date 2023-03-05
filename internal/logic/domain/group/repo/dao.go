package repo

import (
	"errors"
	"learn-im/internal/logic/domain/group/entity"
	"learn-im/logger"
	"learn-im/pkg/db"

	"gorm.io/gorm"
)

type GroupDao struct {
}

func NewGroupDao() *GroupDao {
	return &GroupDao{}
}

func (g *GroupDao) GetGroup(groupId int) (*entity.Group, error) {
	group := new(entity.Group)
	if err := db.DB.Where("id", groupId).First(&group).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		logger.Sugar.Error(err)
		return group, err
	}
	return group, nil
}

func (g *GroupDao) Save(group *entity.Group) error {
	return db.DB.Save(group).Error
}
