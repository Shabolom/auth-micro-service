package di

import (
	"auth-micro-service/internal/repository/auth"
	"auth-micro-service/internal/repository/user"
)

func (d *DI) GetUserRepo() *user.Storage {
	return user.New(d.GetPgDatabase())
}

func (d *DI) GetAuthRepo() *auth.Storage {
	return auth.New(d.GetPgDatabase())
}
