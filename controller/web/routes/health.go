package routes

import (
	"net/http"

	"github.com/mkaiho/go-auth-api/controller/web/handlers"
)

func NewHealthRoutes(
	healthGet *handlers.HealthGetHandler,
) Routes {
	return Routes{
		{
			method:   http.MethodGet,
			path:     "/health",
			handlers: handlers.Handlers{healthGet.Handle},
		},
	}
}
