package pops

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
)

type POPsRouter struct {
	Postgres   *postgres.Queries
	Zing       *zing.Queries
	Middleware *middleware.Middleware
}

func NewPOPsRouter(postgres *postgres.Queries, zing *zing.Queries, middleware *middleware.Middleware) *POPsRouter {
	return &POPsRouter{
		Postgres:   postgres,
		Zing:       zing,
		Middleware: middleware,
	}
}

func (r *POPsRouter) RegisterRoutes() []system.Route {
	getPOPsRoute := r.GetPOPsRoute()

	return []system.Route{
		getPOPsRoute,
	}
}
