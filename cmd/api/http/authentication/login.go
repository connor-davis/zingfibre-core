package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AuthenticationRouter) LoginRoute() system.Route {
	responses := openapi3.NewResponses()

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
			"username": openapi3.NewStringSchema(),
			"password": openapi3.NewStringSchema().WithMinLength(8),
		}),
	)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Login",
			Description: "Endpoint for user login",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: requestBody,
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/login",
		Handler: func(c *fiber.Ctx) error {
			// Implement the login logic here
			// This is a placeholder implementation
			var requestBody map[string]string
			if err := c.BodyParser(&requestBody); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON("Invalid request body")
			}

			username := requestBody["username"]
			password := requestBody["password"]

			// Replace with actual authentication logic
			if username == "admin" && password == "password" {
				return c.JSON("Login successful")
			}

			return c.Status(fiber.StatusUnauthorized).JSON("Invalid credentials")
		},
	}
}
