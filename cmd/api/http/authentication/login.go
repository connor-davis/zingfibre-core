package authentication

import (
	"strings"
	"time"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthenticationRouter) LoginRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessObjectResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		schemas.LoginRequestSchema.Value,
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
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			user, err := r.Postgres.GetUserByEmail(c.Context(), loginRequest.Email)

			if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
				log.Errorf("🔥 Error retrieving user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				log.Warnf("⚠️ User with email %s not found", loginRequest.Email)

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))

			if err != nil {
				log.Warnf("⚠️ Invalid password for user %s: %s", loginRequest.Email, err.Error())

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			currentSession, err := r.Postgres.Sessions().Get(c)

			if err != nil {
				log.Errorf("🔥 Error retrieving session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			currentSession.Set("userId", user.ID.String())
			currentSession.SetExpiry(5 * time.Minute)

			if err := currentSession.Save(); err != nil {
				log.Errorf("🔥 Error saving session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			_, err = r.Postgres.UpdateUser(c.Context(), postgres.UpdateUserParams{
				ID:         user.ID,
				Email:      user.Email,
				Password:   user.Password,
				MfaSecret:  user.MfaSecret,
				MfaEnabled: user.MfaEnabled,
				MfaVerified: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
				Role: user.Role,
			})

			if err != nil {
				log.Errorf("🔥 Error updating user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    user,
			})
		},
	}
}
