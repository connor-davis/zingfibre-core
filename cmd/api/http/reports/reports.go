package reports

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type ReportsRouter struct {
	Zing       *zing.Queries
	Radius     *radius.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewReportsRouter(zing *zing.Queries, radius *radius.Queries, middleware *middleware.Middleware, sessions *session.Store) *ReportsRouter {
	return &ReportsRouter{
		Zing:       zing,
		Radius:     radius,
		Middleware: middleware,
		Sessions:   sessions,
	}
}

func (r *ReportsRouter) RegisterRoutes() []system.Route {
	customersRoute := r.CustomersRoute()
	expiringCustomersRoute := r.ExpiringCustomersRoute()
	rechargesRoute := r.RechargesRoute()
	rechargesSummaryRoute := r.RechargesSummaryRoute()
	summaryRoute := r.SummaryRoute()

	return []system.Route{
		customersRoute,
		expiringCustomersRoute,
		rechargesRoute,
		rechargesSummaryRoute,
		summaryRoute,
	}
}
