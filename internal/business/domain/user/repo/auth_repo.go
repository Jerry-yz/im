package repo

import "learn-im/internal/business/domain/user/model"

type AuthRepo struct {
	*AuthCache
}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		NewAuthCache(),
	}
}

func (a *AuthRepo) Get(deviceId, userId int) (*model.Device, error) {
	return a.AuthCache.Get(userId, deviceId)
}

func (a *AuthRepo) GetAll(userId int) (map[int]*model.Device, error) {
	return a.AuthCache.HGetAll(userId)
}

func (a *AuthRepo) HSet(userId, deviceId int, device *model.Device) error {
	return a.AuthCache.HSet(userId, deviceId, device)
}
