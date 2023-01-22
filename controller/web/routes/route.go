package routes

import (
	"github.com/mkaiho/go-auth-api/controller/web/handlers"
)

type Route struct {
	method   string
	path     string
	handlers handlers.Handlers
}

func (r *Route) Method() string {
	return r.method
}

func (r *Route) Path() string {
	return r.path
}

func (r *Route) Handlers() handlers.Handlers {
	return r.handlers
}

type Routes []*Route
