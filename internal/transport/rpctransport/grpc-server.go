package rpctransport

import (
	"auth-micro-service/internal/transport/rpctransport/auth/check"
	"auth-micro-service/internal/transport/rpctransport/auth/login_handler"
	"auth-micro-service/internal/transport/rpctransport/auth/logout"
	"auth-micro-service/internal/transport/rpctransport/auth/refresh_handler"
	"auth-micro-service/internal/transport/rpctransport/auth/register"
	delet "auth-micro-service/internal/transport/rpctransport/user/delete"
	"auth-micro-service/internal/transport/rpctransport/user/get"
	"auth-micro-service/internal/transport/rpctransport/user/list"
	"auth-micro-service/internal/transport/rpctransport/user/update"
)

type (
	GetLoginHandler    = login_handler.Handler
	GetLogoutHandler   = logout_handeler.Handler
	GetRefreshHandler  = refresh_handler.Handler
	GetRegisterHandler = register_handler.Handler

	GetUserDeleteHandler = delet.Handler
	GetUserHandler       = get.Handler
	GetListUserHandler   = list.Handler
	GetUpdateUserHandler = update.Handler
	GetCheckHandler      = check.Handler
)

type Handlers struct {
	*GetLoginHandler
	*GetLogoutHandler
	*GetRefreshHandler
	*GetRegisterHandler
	*GetCheckHandler
	*GetUserDeleteHandler
	*GetUserHandler
	*GetListUserHandler
	*GetUpdateUserHandler
}

func NewHandlers(
	loginHandler *GetLoginHandler,
	logoutHandler *GetLogoutHandler,
	refreshHandler *GetRefreshHandler,
	registerHandler *GetRegisterHandler,
	getCheckHandler *GetCheckHandler,
	getUserDeleteHandler *GetUserDeleteHandler,
	getUserHandler *GetUserHandler,
	getListUserHandler *GetListUserHandler,
	getUpdateUserHandler *GetUpdateUserHandler,
) *Handlers {
	return &Handlers{
		GetLoginHandler:      loginHandler,
		GetLogoutHandler:     logoutHandler,
		GetRefreshHandler:    refreshHandler,
		GetRegisterHandler:   registerHandler,
		GetCheckHandler:      getCheckHandler,
		GetUserDeleteHandler: getUserDeleteHandler,
		GetUserHandler:       getUserHandler,
		GetListUserHandler:   getListUserHandler,
		GetUpdateUserHandler: getUpdateUserHandler,
	}
}
