package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=8"`
	Role     postgres.RoleType `json:"role"`
}

func (r *AuthenticationRouter) RegisterRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("201", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("User registered successfully.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Created,
						"details": constants.CreatedDetails,
					},
					Schema: schemas.SuccessResponseSchema,
				},
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Bad Request.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.BadRequestError,
						"details": constants.BadRequestErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	responses.Set("409", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Conflict.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.ConflictError,
						"details": constants.ConflictErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Internal Server Error.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.InternalServerError,
						"details": constants.InternalServerErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Register",
			Description: "Endpoint for user registration",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().WithJSONSchema(
					schemas.CreateUserSchema.Value,
				),
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/register",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasRole(postgres.RoleTypeAdmin),
		},
		Handler: func(c *fiber.Ctx) error {
			var registerRequest RegisterRequest

			if err := c.BodyParser(&registerRequest); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			_, err := r.Postgres.GetUserByEmail(c.Context(), registerRequest.Email)

			if err == nil {
				log.Warnf("⚠️ User with email %s already exists", registerRequest.Email)

				return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
					"error":   constants.ConflictError,
					"details": constants.ConflictErrorDetails,
				})
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

			if err != nil {
				log.Errorf("🔥 Error hashing password: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			role := postgres.RoleTypeUser

			if registerRequest.Role != "" {
				role = registerRequest.Role
			}

			if _, err := r.Postgres.CreateUser(c.Context(), postgres.CreateUserParams{
				Email:    registerRequest.Email,
				Password: string(hashedPassword),
				Role:     role,
			}); err != nil {
				log.Errorf("🔥 Error creating user: %s", err.Error())

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
