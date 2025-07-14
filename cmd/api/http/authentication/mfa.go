package authentication

import (
	"bytes"
	"encoding/base32"
	"image/png"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type VerifyMFARequest struct {
	Code string `json:"code"`
}

type DisableMFARequest struct {
	UserId string `json:"userId"`
}

func (r *AuthenticationRouter) EnableMFARoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("201", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Content: map[string]*openapi3.MediaType{
				"image/png": {
					Schema: openapi3.NewSchema().WithFormat("binary").NewRef(),
				},
			},
		},
	})
	responses.Set("500", &constants.InternalServerErrorResponse)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Enable MFA",
			Description: "Endpoint to enable multi-factor authentication",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/authentication/mfa/enable",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(postgres.User)

			if currentUser.MfaSecret.String == "" {
				mfaSecret, err := totp.Generate(totp.GenerateOpts{
					Issuer:      "ZingFibre MFA",
					AccountName: currentUser.Email,
					Period:      30,
					Digits:      otp.DigitsSix,
					Algorithm:   otp.AlgorithmSHA1,
					SecretSize:  32,
				})

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
						"error":   constants.InternalServerError,
						"details": constants.InternalServerErrorDetails,
					})
				}

				currentUser.MfaSecret = pgtype.Text{String: mfaSecret.Secret(), Valid: true}

				updatedCurrentUser, err := r.Postgres.UpdateUser(c.Context(), postgres.UpdateUserParams{
					ID:          currentUser.ID,
					Email:       currentUser.Email,
					Password:    currentUser.Password,
					MfaSecret:   currentUser.MfaSecret,
					MfaEnabled:  currentUser.MfaEnabled,
					MfaVerified: currentUser.MfaVerified,
				})

				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
						"error":   constants.InternalServerError,
						"details": constants.InternalServerErrorDetails,
					})
				}

				currentUser = updatedCurrentUser
			}

			mfaSecretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(currentUser.MfaSecret.String)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			mfaSecret, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "ZingFibre MFA",
				AccountName: currentUser.Email,
				Period:      30,
				Digits:      otp.DigitsSix,
				Algorithm:   otp.AlgorithmSHA1,
				SecretSize:  32,
				Secret:      mfaSecretBytes,
			})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			var pngBuffer bytes.Buffer

			mfaImage, err := mfaSecret.Image(256, 256)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err := png.Encode(&pngBuffer, mfaImage); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusCreated).Send(pngBuffer.Bytes())
		},
	}
}

func (r *AuthenticationRouter) VerifyMFARoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessObjectResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
			"code": openapi3.NewStringSchema(),
		}),
	)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Verify MFA",
			Description: "Endpoint to verify multi-factor authentication code",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: requestBody,
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/mfa/verify",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			currentUser := c.Locals("user").(postgres.User)

			var verifyRequest VerifyMFARequest

			if err := c.BodyParser(&verifyRequest); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			if currentUser.MfaSecret.String == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			if !totp.Validate(verifyRequest.Code, currentUser.MfaSecret.String) {
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			if _, err := r.Postgres.UpdateUser(c.Context(), postgres.UpdateUserParams{
				ID:        currentUser.ID,
				Email:     currentUser.Email,
				Password:  currentUser.Password,
				MfaSecret: currentUser.MfaSecret,
				MfaEnabled: pgtype.Bool{
					Bool:  true,
					Valid: true,
				},
				MfaVerified: pgtype.Bool{
					Bool:  true,
					Valid: true,
				},
			}); err != nil {
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

func (r *AuthenticationRouter) DisableMFARoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessObjectResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	requestBody := openapi3.NewRequestBody().WithJSONSchema(
		openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
			"userId": openapi3.NewStringSchema(),
		}),
	)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Disable MFA",
			Description: "Endpoint to disable multi-factor authentication",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: requestBody,
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/mfa/disable",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			var disableRequest DisableMFARequest

			if err := c.BodyParser(&disableRequest); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			if disableRequest.UserId == "" {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			currentUser, err := r.Postgres.GetUser(c.Context(), uuid.MustParse(disableRequest.UserId))

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if _, err := r.Postgres.UpdateUser(c.Context(), postgres.UpdateUserParams{
				ID:       currentUser.ID,
				Email:    currentUser.Email,
				Password: currentUser.Password,
				MfaSecret: pgtype.Text{
					String: "",
					Valid:  false,
				},
				MfaEnabled: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
				MfaVerified: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
			}); err != nil {
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
