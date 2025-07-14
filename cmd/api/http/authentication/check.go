package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AuthenticationRouter) CheckAuthenticationRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessResponse)
	responses.Set("401", &constants.UnauthorizedResponse)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Check Authentication",
			Description: "Endpoint to check if the user is authenticated",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/authentication/check",
		Handler: func(c *fiber.Ctx) error {
			// Implement the logic to check if the user is authenticated
			// This is a placeholder implementation
			isAuthenticated := true // Replace with actual authentication logic

			if isAuthenticated {
				return c.Status(fiber.StatusOK).JSON(true)
			}

			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"error":   constants.UnauthorizedError,
				"details": constants.UnauthorizedErrorDetails,
			})
		},
	}
}
