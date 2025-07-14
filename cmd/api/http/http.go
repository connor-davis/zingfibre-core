package http

import (
	"fmt"

	"github.com/connor-davis/zingfibre-core/cmd/api/http/authentication"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type HttpRouter struct {
	Routes   []system.Route
	Postgres *postgres.Queries
}

func NewHttpRouter(postgres *postgres.Queries) *HttpRouter {
	authentication := authentication.NewAuthenticationRouter(postgres)
	authenticationRoutes := authentication.RegisterRoutes()

	routes := []system.Route{}

	routes = append(routes, authenticationRoutes...)

	return &HttpRouter{
		Routes:   routes,
		Postgres: postgres,
	}
}

func (h *HttpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range h.Routes {
		switch route.Method {
		case system.GetMethod:
			router.Get(route.Path, append(route.Middlewares, route.Handler)...)
		case system.PostMethod:
			router.Post(route.Path, append(route.Middlewares, route.Handler)...)
		case system.PutMethod:
			router.Put(route.Path, append(route.Middlewares, route.Handler)...)
		case system.DeleteMethod:
			router.Delete(route.Path, append(route.Middlewares, route.Handler)...)
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

		paths.Set(path, pathItem)
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
		},
		Paths: paths,
	}
}
