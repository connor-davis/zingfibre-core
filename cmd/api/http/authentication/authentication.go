package authentication

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthenticationRouter struct {
	Postgres   *postgres.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
}

func NewAuthenticationRouter(postgres *postgres.Queries, middleware *middleware.Middleware, sessions *session.Store) *AuthenticationRouter {
	return &AuthenticationRouter{
		Postgres:   postgres,
		Middleware: middleware,
		Sessions:   sessions,
	}
}

func (r *AuthenticationRouter) RegisterRoutes() []system.Route {
	// Register the authentication routes
	// This includes the check, login, register, enable MFA, verify MFA, and disable MFA routes
	checkRoute := r.CheckAuthenticationRoute()
	loginRoute := r.LoginRoute()
	registerRoute := r.RegisterRoute()
	logoutRoute := r.LogoutRoute()

	// Add the MFA routes
	enableMfaRoute := r.EnableMFARoute()
	verifyMfaRoute := r.VerifyMFARoute()
	disableMfaRoute := r.DisableMFARoute()

	return []system.Route{
		checkRoute,
		loginRoute,
		registerRoute,
		logoutRoute,
		enableMfaRoute,
		verifyMfaRoute,
		disableMfaRoute,
	}
}
