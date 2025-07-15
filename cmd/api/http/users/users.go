package users

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UsersRouter struct {
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewUsersRouter(postgres *postgres.Queries, middleware *middleware.Middleware, sessions *session.Store) *UsersRouter {
	return &UsersRouter{
		Postgres:   postgres,
		Middleware: middleware,
		Sessions:   sessions,
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
