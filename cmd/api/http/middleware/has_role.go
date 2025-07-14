package middleware

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) HasRole(role postgres.RoleType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("user").(postgres.User)

		if currentUser.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   constants.ForbiddenError,
				"details": constants.ForbiddenErrorDetails,
			})
		}

		return c.Next()
	}

}
