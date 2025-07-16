package analytics

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AnalyticsRouter struct {
	Zing       *zing.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewAnalyticsRouter(zing *zing.Queries, middleware *middleware.Middleware, sessions *session.Store) *AnalyticsRouter {
	return &AnalyticsRouter{
		Zing:       zing,
		Middleware: middleware,
		Sessions:   sessions,
	}
}

func (r *AnalyticsRouter) RegisterRoutes() []system.Route {
	rechargeTypeCountsRoute := r.RechargeTypeCountsRoute()

	return []system.Route{
		rechargeTypeCountsRoute,
	}
}
