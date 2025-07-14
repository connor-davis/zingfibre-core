package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AuthenticationRouter) CheckAuthenticationRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
					"message": openapi3.NewStringSchema().WithDefault(constants.Success),
					"details": openapi3.NewStringSchema().WithDefault(constants.SuccessDetails),
					"data":    schemas.UserSchema.Value,
				}),
			).
			WithDescription(constants.SuccessDetails).
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
						"data":    map[string]any{},
					},
					Schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
						"message": openapi3.NewStringSchema().WithDefault(constants.Success),
						"details": openapi3.NewStringSchema().WithDefault(constants.SuccessDetails),
						"data":    schemas.UserSchema.Value,
					}).NewRef(),
				},
			}),
	})
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
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    c.Locals("user"),
			})
		},
	}
}
