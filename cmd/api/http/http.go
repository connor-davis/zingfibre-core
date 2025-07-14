package http

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type HttpRouter struct {
	Registry *Registry
}

func NewHttpRouter() *HttpRouter {
	registry := NewRegistry()

	return &HttpRouter{
		Registry: registry,
	}
}

func (h *HttpRouter) InitializeRoutes(router fiber.Router) {
	for _, route := range h.Registry.Routes {
		router.Add(string(route.Method), route.Path, append(route.Middlewares, route.Handler)...)
	}
}

func (h *HttpRouter) InitializeOpenAPI() *openapi3.T {
	paths := openapi3.NewPaths()

	for _, route := range h.Registry.Routes {
		paths.Set(route.Path, &openapi3.PathItem{
			Summary:     route.Summary,
			Description: route.Description,
			Options: &openapi3.Operation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        route.Tags,
				Parameters:  route.Parameters,
				RequestBody: route.RequestBody,
				Responses:   route.Responses,
			},
		})
	}

	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "My API",
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
		Paths: paths,
	}
}
