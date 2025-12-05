package http

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/connor-davis/zingfibre-core/cmd/api/http/analytics"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/authentication"
	dynamicQueries "github.com/connor-davis/zingfibre-core/cmd/api/http/dynamic-queries"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/exports"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/pops"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/reports"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/users"
	"github.com/connor-davis/zingfibre-core/internal/ai"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type HttpRouter struct {
	Routes     []system.Route
	Postgres   *postgres.Queries
	Zing       *zing.Queries
	Radius     *radius.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
	AI         ai.AI
	Trino      *sql.DB
}

func NewHttpRouter(postgres *postgres.Queries, zing *zing.Queries, radius *radius.Queries, middleware *middleware.Middleware, sessions *session.Store, ai ai.AI, trino *sql.DB) *HttpRouter {
	authentication := authentication.NewAuthenticationRouter(postgres, middleware, sessions)
	authenticationRoutes := authentication.RegisterRoutes()

	users := users.NewUsersRouter(postgres, middleware, sessions)
	usersRoutes := users.RegisterRoutes()

	pois := pops.NewPOPsRouter(postgres, zing, middleware)
	poisRoutes := pois.RegisterRoutes()

	analytics := analytics.NewAnalyticsRouter(zing, middleware, sessions)
	analyticsRoutes := analytics.RegisterRoutes()

	reports := reports.NewReportsRouter(zing, radius, middleware, sessions)
	reportsRoutes := reports.RegisterRoutes()

	exports := exports.NewExportsRouter(zing, radius, middleware, sessions)
	exportsRoutes := exports.RegisterRoutes()

	dynamicQueries := dynamicQueries.NewDynamicQueriesRouter(postgres, zing, radius, middleware, sessions, ai, trino)
	dynamicQueriesRoutes := dynamicQueries.RegisterRoutes()

	routes := []system.Route{}

	routes = append(routes, authenticationRoutes...)
	routes = append(routes, usersRoutes...)
	routes = append(routes, poisRoutes...)
	routes = append(routes, analyticsRoutes...)
	routes = append(routes, reportsRoutes...)
	routes = append(routes, exportsRoutes...)
	routes = append(routes, dynamicQueriesRoutes...)

	return &HttpRouter{
		Routes:     routes,
		Postgres:   postgres,
		Middleware: middleware,
		Sessions:   sessions,
		AI:         ai,
		Trino:      trino,
	}
}

func (h *HttpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range h.Routes {
		path := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(route.Path, ":$1")

		switch route.Method {
		case system.GetMethod:
			router.Get(path, append(route.Middlewares, route.Handler)...)
		case system.PostMethod:
			router.Post(path, append(route.Middlewares, route.Handler)...)
		case system.PutMethod:
			router.Put(path, append(route.Middlewares, route.Handler)...)
		case system.DeleteMethod:
			router.Delete(path, append(route.Middlewares, route.Handler)...)
		}
	}
}

func (h *HttpRouter) InitializeOpenAPI() *openapi3.T {
	paths := openapi3.NewPaths()

	for _, route := range h.Routes {
		pathItem := &openapi3.PathItem{}

		switch route.Method {
		case system.GetMethod:
			pathItem.Get = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		case system.PostMethod:
			pathItem.Post = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case system.PutMethod:
			pathItem.Put = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			}
		case system.DeleteMethod:
			pathItem.Delete = &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				Responses:   route.Responses,
			}
		}

		path := fmt.Sprintf("/api%s", route.Path)

		existingPathItem := paths.Find(path)

		if existingPathItem != nil {
			switch route.Method {
			case system.GetMethod:
				existingPathItem.Get = pathItem.Get
			case system.PostMethod:
				existingPathItem.Post = pathItem.Post
			case system.PutMethod:
				existingPathItem.Put = pathItem.Put
			case system.DeleteMethod:
				existingPathItem.Delete = pathItem.Delete
			}
		} else {
			paths.Set(path, pathItem)
		}
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Zingfibre Reporting API",
			Version: "1.0.0",
		},
		Servers: openapi3.Servers{
			{
				URL:         "http://localhost:4000",
				Description: "Development",
			},
			{
				URL:         "https://reports.core.zingfibre.co.za",
				Description: "Production",
			},
		},
		Tags: openapi3.Tags{
			{
				Name:        "Analytics",
				Description: "Analytics related endpoints",
			},
			{
				Name:        "Authentication",
				Description: "Authentication related endpoints",
			},
			{
				Name:        "Users",
				Description: "User related endpoints",
			},
			{
				Name:        "Points Of Interest",
				Description: "Points Of Interest related endpoints",
			},
			{
				Name:        "Reports",
				Description: "Reports related endpoints",
			},
			{
				Name:        "Exports",
				Description: "Exports related endpoints",
			},
		},
		Paths: paths,
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"User":                    schemas.UserSchema,
				"CreateUser":              schemas.CreateUserSchema,
				"UpdateUser":              schemas.UpdateUserSchema,
				"PointOfPresence":         schemas.PointOfPresenceSchema,
				"PointsOfPresence":        schemas.PointsOfPresenceSchema,
				"DynamicQuery":            schemas.DynamicQuerySchema,
				"CreateDynamicQuery":      schemas.CreateDynamicQuerySchema,
				"UpdateDynamicQuery":      schemas.UpdateDynamicQuerySchema,
				"DynamicQueryResult":      schemas.DynamicQueryResultsSchema,
				"LoginRequest":            schemas.LoginRequestSchema,
				"PasswordReset":           schemas.PasswordResetSchema,
				"SuccessResponse":         schemas.SuccessResponseSchema,
				"ErrorResponse":           schemas.ErrorResponseSchema,
				"RechargeTypeCounts":      schemas.RechargeTypeCountsSchema,
				"ReportCustomer":          schemas.ReportCustomerSchema,
				"ReportCustomers":         schemas.ReportCustomersSchema,
				"ReportExpiringCustomer":  schemas.ReportExpiringCustomerSchema,
				"ReportExpiringCustomers": schemas.ReportExpiringCustomersSchema,
				"ReportRecharge":          schemas.ReportRechargeSchema,
				"ReportRecharges":         schemas.ReportRechargesSchema,
				"ReportRechargeSummary":   schemas.ReportRechargeSummarySchema,
				"ReportRechargeSummaries": schemas.ReportRechargeSummariesSchema,
				"ReportSummary":           schemas.ReportSummarySchema,
				"ReportSummaries":         schemas.ReportSummariesSchema,
				"MonthlyStatistics":       schemas.MonthlyStatisticsSchema,
			},
		},
	}
}
