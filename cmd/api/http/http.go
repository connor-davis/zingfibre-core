package http

import (
	"fmt"
	"regexp"

	"github.com/connor-davis/zingfibre-core/cmd/api/http/authentication"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/pois"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/users"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type HttpRouter struct {
	Routes     []system.Route
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewHttpRouter(postgres *postgres.Queries, middleware *middleware.Middleware, sessions *session.Store) *HttpRouter {
	authentication := authentication.NewAuthenticationRouter(postgres, middleware, sessions)
	authenticationRoutes := authentication.RegisterRoutes()

	users := users.NewUsersRouter(postgres, middleware, sessions)
	usersRoutes := users.RegisterRoutes()

	pois := pois.NewPointOfInterestsRouter(postgres, middleware, sessions)
	poisRoutes := pois.RegisterRoutes()

	routes := []system.Route{}

	routes = append(routes, authenticationRoutes...)
	routes = append(routes, usersRoutes...)
	routes = append(routes, poisRoutes...)

	return &HttpRouter{
		Routes:     routes,
		Postgres:   postgres,
		Middleware: middleware,
		Sessions:   sessions,
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
				URL:         "https://api.example.com",
				Description: "Production",
			},
		},
		Tags: openapi3.Tags{
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
		},
		Paths: paths,
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"User":                  schemas.UserSchema,
				"CreateUser":            schemas.CreateUserSchema,
				"UpdateUser":            schemas.UpdateUserSchema,
				"PointOfInterest":       schemas.PointOfInterestSchema,
				"CreatePointOfInterest": schemas.CreatePointOfInterestSchema,
				"UpdatePointOfInterest": schemas.UpdatePointOfInterestSchema,
				"LoginRequest":          schemas.LoginRequestSchema,
				"SuccessResponse":       schemas.SuccessResponseSchema,
				"ErrorResponse":         schemas.ErrorResponseSchema,
			},
		},
	}
}
