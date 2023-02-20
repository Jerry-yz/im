package repo

import (
	"errors"
	"learn-im/internal/business/domain/user/model"

	"go.uber.org/zap"
)

type UserRepo struct {
	*UserDao
	*UserCache
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		NewUserDao(),
		NewUserCache(),
	}
}

func (u *UserRepo) Get(userId int) (user *model.User, err error) {
	// user := new(model.User)
	user, err = u.UserCache.Get(userId)
	if err != nil {
		zap.Error(errors.New("get user cache error"))
		return user, err
	}
	if user.ID != 0 {
		return user, nil
	}
	user, err = u.UserDao.Get(userId)
	if err != nil {
		zap.Error(errors.New("get user error"))
		return user, err
	}
	if user.ID != 0 {
		return nil, u.UserCache.Set(userId, user)
	}
	return
}
