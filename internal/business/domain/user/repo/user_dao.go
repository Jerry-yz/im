package repo

import (
	"errors"
	"learn-im/internal/business/domain/user/model"
	"learn-im/pkg/db"

	"gorm.io/gorm"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (d *UserDao) Create(user *model.User) error {
	return db.DB.Create(&user).Error
}

func (d *UserDao) Get(userId int) (*model.User, error) {
	user := new(model.User)
	if err := db.DB.Where("id", userId).First(&user).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, err
	}
	return user, nil
}

func (d *UserDao) Save(user *UserDao) error {
	return db.DB.Save(&user).Error
}

func (d *UserDao) GetByIds(userIds []int) ([]*model.User, error) {
	userList := make([]*model.User, len(userIds))
	if err := db.DB.Where("id", userIds).Find(&userList).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return userList, err
	}
	return userList, nil
}

func (d *UserDao) GetUserByNumber(phone string) (*model.User, error) {
	user := new(model.User)
	if err := db.DB.Where("phone", phone).First(&user).Error; err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		return user, err
	}
	return user, nil
}

func (d *UserDao) List() ([]*model.User, error) {
	userList := make([]*model.User, 0)
	if err := db.DB.Find(&userList).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return userList, err
	}
	return userList, nil
}

func (d *UserDao) SearchUser(keyword string) ([]*model.User, error) {
	userList := make([]*model.User, 0)
	keyword = "%" + keyword + "%"
	if err := db.DB.Where("phone like ? or name_nick ?", keyword, keyword).Find(&userList).Error; err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return userList, err
	}
	return userList, nil
}
