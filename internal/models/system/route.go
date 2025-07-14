package system

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type RouteMethod string

const (
	GetMethod    RouteMethod = "GET"
	PostMethod   RouteMethod = "POST"
	PutMethod    RouteMethod = "PUT"
	DeleteMethod RouteMethod = "DELETE"
)

type OpenAPIMetadata struct {
	Summary     string
	Description string
	Tags        []string
	Parameters  []*openapi3.ParameterRef
	RequestBody *openapi3.RequestBodyRef
	Responses   *openapi3.Responses
}

type Route struct {
	OpenAPIMetadata

	Method      RouteMethod
	Path        string
	Middlewares []fiber.Handler
	Handler     func(*fiber.Ctx) error
}
