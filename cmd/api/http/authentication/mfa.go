package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type VerifyMFARequest struct {
	Code string `json:"code"`
}

type DisableMFARequest struct {
	UserId string `json:"userId"`
}

func (r *AuthenticationRouter) EnableMFARoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessResponse)
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
		Method: system.PostMethod,
		Path:   "/authentication/mfa/enable",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			// Implement the logic to enable MFA
			// This is a placeholder implementation
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}

func (r *AuthenticationRouter) VerifyMFARoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessResponse)
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
			// Implement the logic to verify MFA
			// This is a placeholder implementation
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}

func (r *AuthenticationRouter) DisableMFARoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessResponse)
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
			// Implement the logic to disable MFA
			// This is a placeholder implementation
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
