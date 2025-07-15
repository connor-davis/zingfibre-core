package authentication

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (r *AuthenticationRouter) LogoutRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

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
