package users

import (
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Role     postgres.RoleType `json:"role"`
}

func (r *UsersRouter) CreateUserRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("201", &constants.SuccessResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("409", &constants.ConflictResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Create User",
			Description: "Endpoint to create a new user",
			Tags:        []string{"Users"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().WithJSONSchema(schemas.CreateUserSchema.Value),
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/users",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasRole(postgres.RoleTypeAdmin),
		},
		Handler: func(c *fiber.Ctx) error {
			var createUserRequest CreateUserRequest

			if err := c.BodyParser(&createUserRequest); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			_, err := r.Postgres.GetUserByEmail(c.Context(), createUserRequest.Email)

			if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
				log.Errorf("🔥 Error checking if user exists: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err == nil {
				log.Warnf("⚠️ User with email %s already exists", createUserRequest.Email)

				return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
					"error":   constants.ConflictError,
					"details": constants.ConflictErrorDetails,
				})
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), bcrypt.DefaultCost)

			if err != nil {
				log.Errorf("🔥 Error hashing password: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			role := postgres.RoleTypeUser

			if createUserRequest.Role != "" {
				role = createUserRequest.Role
			}

			_, err = r.Postgres.CreateUser(c.Context(), postgres.CreateUserParams{
				Email:    createUserRequest.Email,
				Password: string(hashedPassword),
				Role:     role,
			})

			if err != nil {
				log.Errorf("🔥 Error creating user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
