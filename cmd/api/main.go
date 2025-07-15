package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/zingfibre-core/cmd/api/http"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/env"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	context := context.Background()

	postgresPoolConfig, err := pgxpool.ParseConfig(string(env.POSTGRES_DSN))

	if err != nil {
		log.Infof("🔥 Failed to parse PostgreSQL connection string: %v", err)
		return
	}

	postgresPool, err := pgxpool.NewWithConfig(context, postgresPoolConfig)

	if err != nil {
		log.Infof("🔥 Failed to connect to PostgreSQL: %v", err)
	}

	defer postgresPool.Close()

	log.Info("✅ Connected to PostgreSQL successfully")

	postgresQueries := postgres.New(postgresPool)

	log.Info("🔃 Creating default admin user.")

	existingAdmin, err := postgresQueries.GetUserByEmail(context, string(env.ADMIN_EMAIL))

	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		log.Infof("🔥 Error checking for existing admin user: %v", err)
		return
	}

	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		log.Info("🔃 No existing admin user found, creating a new one...")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(string(env.ADMIN_PASSWORD)), bcrypt.DefaultCost)

		if err != nil {
			log.Errorf("🔥 Error hashing admin password: %s", err.Error())
			return
		}

		if _, err := postgresQueries.CreateUser(context, postgres.CreateUserParams{
			Email:    string(env.ADMIN_EMAIL),
			Password: string(hashedPassword),
			Role:     postgres.RoleTypeAdmin,
		}); err != nil {
			log.Errorf("🔥 Error creating admin user: %s", err.Error())
			return
		}

		log.Info("✅ Admin user created successfully")
	} else {
		log.Infof("✅ Admin user already exists: %s", existingAdmin.Email)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Zingfibre Reporting API",
		ServerHeader: "Zingfibre-API",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,https://reports.core.zingfibre.co.za",
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE",
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${latency} ${method} ${url}\n",
	}))

	middleware := middleware.NewMiddleware(postgresQueries)

	httpRouter := http.NewHttpRouter(postgresQueries, middleware)

	openapiSpecification := httpRouter.InitializeOpenAPI()

	api := app.Group("/api")

	api.Get("/api-spec", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(openapiSpecification)
	})

	api.Get("/api-doc", func(c *fiber.Ctx) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: func() string {
				if env.MODE == "production" {
					return "https://one-staging.thusa.co.za/api/api-spec"
				}

				return fmt.Sprintf("http://localhost:%s/api/api-spec", env.PORT)
			}(),
			Theme:  scalar.ThemeDefault,
			Layout: scalar.LayoutModern,
			BaseServerURL: func() string {
				if env.MODE == "production" {
					return "https://one-staging.thusa.co.za"
				}

				return fmt.Sprintf("http://localhost:%s", env.PORT)
			}(),
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Type("html").SendString(html)
	})

	httpRouter.InitializeRoutes(api)

	log.Info("Starting Zingfibre Reporting API on port 4000...")

	if err := app.Listen(":4000"); err != nil {
		log.Infof("Failed to start server: %v", err)
	}
}
