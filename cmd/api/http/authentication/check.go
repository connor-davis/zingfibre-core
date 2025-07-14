package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AuthenticationRouter) CheckAuthenticationRoute() system.Route {
	responses := openapi3.NewResponses()

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Check Authentication",
			Description: "Endpoint to check if the user is authenticated",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method:  system.GetMethod,
		Path:    "/authentication/check",
		Handler: r.CheckAuthenticationHandler,
	}
}

func (r *AuthenticationRouter) CheckAuthenticationHandler(c *fiber.Ctx) error {
	// Implement the logic to check if the user is authenticated
	// This is a placeholder implementation
	isAuthenticated := true // Replace with actual authentication logic

	if isAuthenticated {
		return c.JSON(true)
	}

	return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
}
