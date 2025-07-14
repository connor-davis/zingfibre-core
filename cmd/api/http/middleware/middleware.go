package middleware

import "github.com/connor-davis/zingfibre-core/internal/postgres"

type Middleware struct {
	Postgres *postgres.Queries
}

func NewMiddleware(postgres *postgres.Queries) *Middleware {
	return &Middleware{
		Postgres: postgres,
	}
}
