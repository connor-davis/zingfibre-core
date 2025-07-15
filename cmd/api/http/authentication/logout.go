package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (r *AuthenticationRouter) LogoutRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("User logged out successfully.").
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

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).WithDescription("Unauthorized.").WithContent(openapi3.Content{
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
			).WithDescription("Internal Server Error.").WithContent(openapi3.Content{
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
			Summary:     "Logout",
			Description: "Endpoint for user logout",
			Tags:        []string{"Authentication"},
			Parameters:  nil,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.PostMethod,
		Path:   "/authentication/logout",
		Handler: func(c *fiber.Ctx) error {
			currentSession, err := r.Sessions.Get(c)

			if err != nil {
				log.Errorf("🔥 Error retrieving session: %s", err.Error())

				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"message": constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				})
			}

			if err := currentSession.Destroy(); err != nil {
				log.Errorf("🔥 Error destroying session: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"message": constants.InternalServerError,
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
