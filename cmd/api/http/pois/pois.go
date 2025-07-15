package pois

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type PointsOfInterestRouter struct {
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewPointOfInterestsRouter(postgres *postgres.Queries, middleware *middleware.Middleware, sessions *session.Store) *PointsOfInterestRouter {
	return &PointsOfInterestRouter{
		Postgres:   postgres,
		Middleware: middleware,
	}
}

func (r *PointsOfInterestRouter) RegisterRoutes() []system.Route {
	getPointsOfInterestRoute := r.GetPointsOfInterestRoute()
	getPointOfInterestRoute := r.GetPointOfInterestRoute()
	createPointOfInterestRoute := r.CreatePointOfInterestRoute()
	updatePointOfInterestRoute := r.UpdatePointOfInterestRoute()
	deletePointOfInterestRoute := r.DeletePointOfInterestRoute()

	return []system.Route{
		getPointsOfInterestRoute,
		getPointOfInterestRoute,
		createPointOfInterestRoute,
		updatePointOfInterestRoute,
		deletePointOfInterestRoute,
	}
}
