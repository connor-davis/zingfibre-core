package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetRequest struct {
	Code        string `json:"Code"`
	NewPassword string `json:"NewPassword"`
}

func (r *AuthenticationRouter) PasswordResetRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("MFA code verified successfully.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
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

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Unauthorized.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.UnauthorizedError,
						"details": constants.UnauthorizedErrorDetails,
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

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		schemas.PasswordResetSchema.Value,
	)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Password Reset",
			Description: "Endpoint to reset user password",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: requestBody,
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/password-reset",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(postgres.User)

			var resetRequest PasswordResetRequest

			if err := c.BodyParser(&resetRequest); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			if currentUser.MfaSecret.String == "" {
				log.Warnf("⚠️ User %s has no MFA secret set", currentUser.Email)

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			if !totp.Validate(resetRequest.Code, currentUser.MfaSecret.String) {
				log.Warnf("⚠️ Invalid MFA code for user %s", currentUser.Email)

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetRequest.NewPassword), bcrypt.DefaultCost)

			if err != nil {
				log.Errorf("🔥 Error hashing password: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if _, err := r.Postgres.UpdateUser(c.Context(), postgres.UpdateUserParams{
				ID:          currentUser.ID,
				Email:       currentUser.Email,
				Password:    string(hashedPassword),
				MfaSecret:   currentUser.MfaSecret,
				MfaEnabled:  currentUser.MfaEnabled,
				MfaVerified: currentUser.MfaVerified,
				Role:        currentUser.Role,
			}); err != nil {
				log.Errorf("🔥 Error updating user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
