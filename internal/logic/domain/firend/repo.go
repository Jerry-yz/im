package friend

import (
	"errors"
	"learn-im/pkg/db"
	"learn-im/pkg/gerrors"

	"gorm.io/gorm"
)

type FriendRepo struct {
}

func NewFriendRepo() *FriendRepo {
	return &FriendRepo{}
}

func (f *FriendRepo) GetFriend(userId, friendId int) (*Friend, error) {
	friend := new(Friend)
	if err := db.DB.Where("user_id = ? and friend_id = ?", userId, friendId).
		First(&friend).Error; err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		return friend, gerrors.WarpError(err)
	}
	return friend, nil
}

func (f *FriendRepo) Save(friend *Friend) error {
	return db.DB.Save(&friend).Error
}

func (f *FriendRepo) ListFriend(userId, status int) ([]*Friend, error) {
	friends := make([]*Friend, 0)
	if err := db.DB.Where("user_id = ? and status = ?", userId, status).
		Find(&friends).Error; err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		return friends, gerrors.WarpError(err)
	}
	return friends, nil
}
