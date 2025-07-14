package middleware

import (
	"log"
	"time"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *Middleware) Authorized() fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentSession, err := m.Postgres.Sessions().Get(c)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   constants.UnauthorizedError,
				"details": constants.UnauthorizedErrorDetails,
			})
		}

		currentUserId, ok := currentSession.Get("userId").(string)

		if !ok || currentUserId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   constants.UnauthorizedError,
				"details": constants.UnauthorizedErrorDetails,
			})
		}

		log.Printf("🔐 Authorized User with ID: %s", currentUserId)

		currentUser, err := m.Postgres.GetUser(c.Context(), uuid.MustParse(currentUserId))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   constants.UnauthorizedError,
				"details": constants.UnauthorizedErrorDetails,
			})
		}

		c.Locals("userId", currentUser.ID.String())
		c.Locals("user", currentUser)

		currentSession.Set("userId", currentUser.ID.String())
		currentSession.SetExpiry(5 * time.Minute)

		if err := currentSession.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   constants.InternalServerError,
				"details": constants.InternalServerErrorDetails,
			})
		}

		return c.Next()
	}
}
