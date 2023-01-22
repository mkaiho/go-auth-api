package routes

import (
	"net/http"

	"github.com/mkaiho/go-auth-api/controller/web/handlers"
)

func NewUserRoutes(
	userFind *handlers.UserFindHandler,
	userCreate *handlers.UserCreateHandler,
	userGet *handlers.UserGetHandler,
	userUpdate *handlers.UserUpdateHandler,
) Routes {
	return Routes{
		{
			method:   http.MethodGet,
			path:     "/users",
			handlers: handlers.Handlers{userFind.Handle},
		},
		{
			method:   http.MethodPost,
			path:     "/users",
			handlers: handlers.Handlers{userCreate.Handle},
		},
		{
			method:   http.MethodGet,
			path:     "/users/:id",
			handlers: handlers.Handlers{userGet.Handle},
		},
		{
			method:   http.MethodPut,
			path:     "/users/:id",
			handlers: handlers.Handlers{userUpdate.Handle},
		},
	}
}
