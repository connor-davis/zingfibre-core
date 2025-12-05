package sessions

import (
	"time"

	"github.com/connor-davis/zingfibre-core/common"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberPg "github.com/gofiber/storage/postgres/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSessions(db *pgxpool.Pool) *session.Store {
	return session.New(session.Config{
		Storage: fiberPg.New(fiberPg.Config{
			Table: "sessions",
			DB:    db,
		}),
		KeyLookup:         "cookie:zingfibre_sessions",
		CookieDomain:      common.EnvString("COOKIE_DOMAIN", "localhost"),
		CookiePath:        "/",
		CookieSecure:      true,
		CookieSameSite:    "Strict",
		CookieSessionOnly: false,
		CookieHTTPOnly:    false,
		Expiration:        24 * time.Hour,
	})
}
