package main

import (
	"context"
	"fmt"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/zingfibre-core/cmd/api/http"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/env"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5"
)

func main() {
	context := context.Background()

	databaseConnection, err := pgx.Connect(context, string(env.POSTGRES_DSN))

	if err != nil {
		log.Infof("🔥 Failed to connect to PostgreSQL: %v", err)
	}

	defer databaseConnection.Close(context)

	log.Info("Connected to PostgreSQL successfully")

	postgres := postgres.New(databaseConnection)

	app := fiber.New(fiber.Config{
		AppName:      "Zingfibre Reporting API",
		ServerHeader: "Zingfibre-API",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${latency} ${method} ${url}\n",
	}))

	middleware := middleware.NewMiddleware(postgres)

	httpRouter := http.NewHttpRouter(postgres, middleware)

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
