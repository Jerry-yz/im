package repo

import (
	"errors"
	"learn-im/internal/logic/domain/group/entity"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"

	"gorm.io/gorm"
)

type GroupUserRepo struct {
}

func NewGroupUserRepo() *GroupUserRepo {
	return &GroupUserRepo{}
}

func (g *GroupUserRepo) GetGroupsByUserId(userId int) ([]*entity.Group, error) {
	userGroupIds := make([]int, 0)
	if err := db.DB.Where("user_id", userId).Pluck("group_id", &userGroupIds).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, gerrors.WarpError(err)
	}
	groups := make([]*entity.Group, 0)
	if err := db.DB.Where("id", userGroupIds).Find(&groups).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, gerrors.WarpError(err)
	}
	return groups, nil
}

func (g *GroupUserRepo) ListGroupUsers(groupId int) ([]*entity.GroupUser, error) {
	groupUsers := make([]*entity.GroupUser, 0)
	if err := db.DB.Where("group_id", groupId).Find(&groupUsers).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return groupUsers, gerrors.WarpError(err)
	}
	return groupUsers, nil
}

func (g *GroupUserRepo) GetGroupUser(groupId, userId int) (*entity.GroupUser, error) {
	groupUser := new(entity.GroupUser)
	if err := db.DB.Where("group_id = ? and user_id = ?", groupId, userId).First(&groupUser).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, gerrors.WarpError(err)
	}
	return groupUser, nil
}

func (g *GroupUserRepo) SaveGroupUser(groupUser *entity.GroupUser) error {
	return db.DB.Save(groupUser).Error
}

func (g *GroupUserRepo) DeleteGroupUser(groupId, userId int) error {
	return db.DB.Where("group_id = ? and user_id = ?", groupId, userId).Error
}
