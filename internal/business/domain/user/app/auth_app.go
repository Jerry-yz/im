package app

import (
	"context"
	"learn-im/internal/business/domain/user/service"
)

type AuthApp struct {
	*service.AuthService
}

func NewAuthApp() *AuthApp {
	return &AuthApp{
		service.NewAuthService(),
	}
}

func (a *AuthApp) SignIn(ctx context.Context, phone string, deviceId int) error {
	return a.AuthService.SignIn(ctx, phone, deviceId)
}

func (a *AuthApp) Auth(ctx context.Context, userId, deviceId int, token string) error {
	return a.AuthService.Auth(ctx, userId, deviceId, token)
}
