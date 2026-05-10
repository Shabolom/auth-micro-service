package di

import (
	"auth-micro-service/internal/service/auth"
	"auth-micro-service/internal/service/user"
)

func (d *DI) GetAuthService() *auth.Service {
	return auth.New(d.GetAuthRepo(), d.GetInMemoryStorage(), d.Config().Secret, d.Logger())
}

func (d *DI) GetUserService() *user.Service {
	return user.New(d.GetUserRepo(), d.Config().Secret, d.Logger())
}
