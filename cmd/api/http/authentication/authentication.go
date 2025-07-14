package authentication

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
)

type AuthenticationRouter struct {
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
}

func NewAuthenticationRouter(postgres *postgres.Queries, middleware *middleware.Middleware) *AuthenticationRouter {
	return &AuthenticationRouter{
		Postgres:   postgres,
		Middleware: middleware,
	}
}
func (r *AuthenticationRouter) RegisterRoutes() []system.Route {
	checkRoute := r.CheckAuthenticationRoute()
	loginRoute := r.LoginRoute()
	registerRoute := r.RegisterRoute()

	return []system.Route{
		checkRoute,
		loginRoute,
		registerRoute,
	}
}
