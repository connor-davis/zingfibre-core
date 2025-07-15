package middleware

import (
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Middleware struct {
	Postgres *postgres.Queries
	Sessions *session.Store
}

func NewMiddleware(postgres *postgres.Queries, sessions *session.Store) *Middleware {
	return &Middleware{
		Postgres: postgres,
		Sessions: sessions,
	}
}
