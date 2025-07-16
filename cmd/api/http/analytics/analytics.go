package analytics

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AnalyticsRouter struct {
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewAnalyticsRouter(postgres *postgres.Queries, middleware *middleware.Middleware, sessions *session.Store) *AnalyticsRouter {
	return &AnalyticsRouter{
		Postgres:   postgres,
		Middleware: middleware,
		Sessions:   sessions,
	}
}

func (r *AnalyticsRouter) RegisterRoutes() []system.Route {
	return []system.Route{}
}
