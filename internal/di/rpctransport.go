package di

import (
	"auth-micro-service/internal/transport/rpctransport"
	"auth-micro-service/internal/transport/rpctransport/auth/login_handler"
	"auth-micro-service/internal/transport/rpctransport/auth/logout"
	"auth-micro-service/internal/transport/rpctransport/auth/refresh_handler"
	"auth-micro-service/internal/transport/rpctransport/auth/register"
	user_delete "auth-micro-service/internal/transport/rpctransport/user/delete"
	user_get "auth-micro-service/internal/transport/rpctransport/user/get"
	user_list "auth-micro-service/internal/transport/rpctransport/user/list"
	user_update "auth-micro-service/internal/transport/rpctransport/user/update"
)

func (d *DI) GetGRPCUsersHandlers() *rpctransport.UsersHandlers {
	return rpctransport.NewUsersHandlers(
		d.GetUserDeleteHandler(),
		d.GetUserHandler(),
		d.GetListUserHandler(),
		d.GetUpdateUserHandler(),
	)
}

func (d *DI) GetUserDeleteHandler() *user_delete.Handler {
	return user_delete.New(d.GetUserService())
}

func (d *DI) GetUserHandler() *user_get.Handler {
	return user_get.New(d.GetUserService())
}

func (d *DI) GetListUserHandler() *user_list.Handler {
	return user_list.New(d.GetUserService())
}

func (d *DI) GetUpdateUserHandler() *user_update.Handler {
	return user_update.New(d.GetUserService())
}

func (d *DI) GetGRPCAuthHandlers() *rpctransport.AuthHandlers {
	return rpctransport.NewAuthHandlers(
		d.GetLoginHandler(),
		d.GetLogoutHandler(),
		d.GetRefreshHandler(),
		d.GetRegisterHandler(),
	)
}

func (d *DI) GetLoginHandler() *login_handler.Handler {
	return login_handler.New(d.GetAuthService())
}

func (d *DI) GetLogoutHandler() *logout_handeler.Handler {
	return logout_handeler.New(d.GetAuthService())
}

func (d *DI) GetRefreshHandler() *refresh_handler.Handler {
	return refresh_handler.New(d.GetAuthService())
}

func (d *DI) GetRegisterHandler() *register_handler.Handler {
	return register_handler.New(d.GetAuthService())
}
