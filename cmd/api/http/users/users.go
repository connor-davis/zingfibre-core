package users

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
)

type UsersRouter struct {
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
}

func NewUsersRouter(postgres *postgres.Queries, middleware *middleware.Middleware) *UsersRouter {
	return &UsersRouter{
		Postgres:   postgres,
		Middleware: middleware,
	}
}

func (r *UsersRouter) RegisterRoutes() []system.Route {
	getUsersRoute := r.GetUsersRoute()
	getUserRoute := r.GetUserRoute()
	createUserRoute := r.CreateUserRoute()
	updateUserRoute := r.UpdateUserRoute()
	deleteUserRoute := r.DeleteUserRoute()

	return []system.Route{
		getUsersRoute,
		getUserRoute,
		createUserRoute,
		updateUserRoute,
		deleteUserRoute,
	}
}
