package service

import (
	"context"
	"learn-im/internal/business/domain/user/model"
	"learn-im/internal/business/domain/user/repo"
	"learn-im/pkg/gerrors"
	"regexp"
	"time"
)

type AuthService struct {
	*repo.UserRepo
	*repo.AuthRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo.NewUserRepo(),
		repo.NewAuthRepo(),
	}
}

func (a *AuthService) SignIn(ctx context.Context, phone string, deviceId int) error {
	if !VerifyPhone(phone) {
		return gerrors.ErrBadRequest
	}
	user, err := a.GetUserByNumber(phone)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if user == nil {
		user.Phone = phone
		return a.UserRepo.Save(user)
	}
	// grpc获取设备信息
	return a.AuthRepo.HSet(int(user.ID), deviceId, &model.Device{})
}

func (a *AuthService) Auth(ctx context.Context, userId, deviceId int, token string) error {
	device, err := a.AuthRepo.Get(deviceId, userId)
	if err != nil {
		return gerrors.WarpError(err)
	}
	if device == nil {
		return gerrors.ErrUnauthorized
	}
	if device.Expire < time.Now().Unix() {
		return gerrors.ErrUnauthorized
	}
	if device.Token != token {
		return gerrors.ErrUnauthorized
	}
	return nil
}

func VerifyPhone(phone string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)

}
