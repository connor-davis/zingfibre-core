package postgres

import (
	"time"

	"github.com/connor-davis/zingfibre-core/env"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberPg "github.com/gofiber/storage/postgres/v2"
)

func (q *Queries) Sessions() *session.Store {
	sessions := session.New(session.Config{
		Storage: fiberPg.New(fiberPg.Config{
			Table:         "sessions",
			ConnectionURI: string(env.POSTGRES_DSN),
		}),
		KeyLookup:         "cookie:zingfibre_sessions",
		CookieDomain:      string(env.COOKIE_DOMAIN),
		CookiePath:        "/",
		CookieSecure:      true,
		CookieSameSite:    "Strict",
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		Expiration:        5 * time.Minute,
	})

	if sessions == nil {
		log.Errorf("🔥 Failed to create session store")
		return nil
	}

	return sessions
}
