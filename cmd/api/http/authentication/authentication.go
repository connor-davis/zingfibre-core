package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
)

type AuthenticationRouter struct {
	Postgres *postgres.Queries
}

func NewAuthenticationRouter(postgres *postgres.Queries) *AuthenticationRouter {
	return &AuthenticationRouter{
		Postgres: postgres,
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
