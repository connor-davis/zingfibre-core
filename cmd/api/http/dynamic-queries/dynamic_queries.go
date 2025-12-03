package dynamicQueries

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type DynamicQueriesRouter struct {
	Postgres   *postgres.Queries
	Zing       *zing.Queries
	Radius     *radius.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewDynamicQueriesRouter(
	postgres *postgres.Queries,
	zing *zing.Queries,
	radius *radius.Queries,
	middleware *middleware.Middleware,
	sessions *session.Store,
) *DynamicQueriesRouter {
	return &DynamicQueriesRouter{
		Postgres:   postgres,
		Zing:       zing,
		Radius:     radius,
		Middleware: middleware,
		Sessions:   sessions,
	}
}

func (r *DynamicQueriesRouter) RegisterRoutes() []system.Route {
	return []system.Route{
		r.GetDynamicQueriesRoute(),
		r.GetDynamicQueryRoute(),
		r.CreateDynamicQueryRoute(),
		r.UpdateDynamicQueryRoute(),
		r.DeleteDynamicQueryRoute(),
	}
}
