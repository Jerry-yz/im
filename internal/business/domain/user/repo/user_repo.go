package repo

import (
	"errors"
	"learn-im/internal/business/domain/user/model"
	"learn-im/pkg/gerrors"

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

func (u *UserRepo) GetUserByNumber(phone string) (*model.User, error) {
	return u.UserDao.GetUserByNumber(phone)
}

func (u *UserRepo) GetUserIds(userIds []int) ([]*model.User, error) {
	return u.UserDao.GetByIds(userIds)
}

func (u *UserRepo) Search(keyword string) ([]*model.User, error) {
	return u.UserDao.SearchUser(keyword)
}

func (u *UserRepo) Save(user *model.User) error {
	if err := u.UserDao.Save(user); err != nil {
		return gerrors.WarpError(err)
	}
	if user.ID != 0 {
		return u.UserCache.Delete(int(user.ID))
	}
	return nil
}
