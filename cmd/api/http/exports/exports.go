package exports

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type ExportsRouter struct {
	Zing       *zing.Queries
	Radius     *radius.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewExportsRouter(zing *zing.Queries, radius *radius.Queries, middleware *middleware.Middleware, sessions *session.Store) *ExportsRouter {
	return &ExportsRouter{
		Zing:       zing,
		Radius:     radius,
		Middleware: middleware,
		Sessions:   sessions,
	}
}

func (r *ExportsRouter) RegisterRoutes() []system.Route {
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
