package authentication

import (
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthenticationRouter) LoginRoute() system.Route {
	responses := openapi3.NewResponses()

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
			"email":    openapi3.NewStringSchema(),
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
			var loginRequest LoginRequest

			if err := c.BodyParser(&loginRequest); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON("Invalid request body")
			}

			user, err := r.Postgres.GetUserByEmail(c.Context(), loginRequest.Email)

			if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(user)
		},
	}
}
