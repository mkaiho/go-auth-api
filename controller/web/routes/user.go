package routes

import (
	"net/http"

	"github.com/mkaiho/go-auth-api/controller/web/handlers"
	"github.com/mkaiho/go-auth-api/controller/web/middlewares"
	"github.com/mkaiho/go-auth-api/usecase/port"
)

func NewUserRoutes(
	txm port.TransactionManager,
	credGateway port.UserCredentialGateway,
	userFind *handlers.UserFindHandler,
	userCreate *handlers.UserCreateHandler,
	userGet *handlers.UserGetHandler,
	userUpdate *handlers.UserUpdateHandler,
) Routes {
	return Routes{
		{
			method:   http.MethodGet,
			path:     "/users",
			handlers: handlers.Handlers{middlewares.CheckAuth(txm, credGateway), userFind.Handle},
		},
		{
			method:   http.MethodPost,
			path:     "/users",
			handlers: handlers.Handlers{userCreate.Handle},
		},
		{
			method:   http.MethodGet,
			path:     "/users/:id",
			handlers: handlers.Handlers{middlewares.CheckAuth(txm, credGateway), userGet.Handle},
		},
		{
			method:   http.MethodPut,
			path:     "/users/:id",
			handlers: handlers.Handlers{middlewares.CheckAuth(txm, credGateway), userUpdate.Handle},
		},
	}
}
