package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (r *AuthenticationRouter) RegisterRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("201", &constants.CreatedResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("409", &constants.ConflictResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
			"email":    openapi3.NewStringSchema().WithFormat("email"),
			"password": openapi3.NewStringSchema().WithMinLength(8),
		}),
	)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Register",
			Description: "Endpoint for user registration",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: requestBody,
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/register",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			var registerRequest RegisterRequest

			if err := c.BodyParser(&registerRequest); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			_, err := r.Postgres.GetUserByEmail(c.Context(), registerRequest.Email)

			if err == nil {
				return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
					"error":   constants.ConflictError,
					"details": constants.ConflictErrorDetails,
				})
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if _, err := r.Postgres.CreateUser(c.Context(), postgres.CreateUserParams{
				Email:    registerRequest.Email,
				Password: string(hashedPassword),
			}); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
				"message": constants.Created,
				"details": constants.CreatedDetails,
			})
		},
	}
}
