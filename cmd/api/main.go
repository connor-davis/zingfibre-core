package main

import (
	"fmt"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/zingfibre-core/cmd/api/http"
	"github.com/connor-davis/zingfibre-core/env"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Zingfibre Reporting API",
		ServerHeader: "Zingfibre-API",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	httpRouter := http.NewHttpRouter()

	openapiSpecification := httpRouter.InitializeOpenAPI()

	api := app.Group("/api")

	api.Get("/api-spec", func(c *fiber.Ctx) error {
		data, err := json.MarshalIndent(openapiSpecification, "", "  ")

		if err != nil {
			return c.Status(500).SendString("Failed to generate OpenAPI JSON")
		}

		return c.Type("application/json").Send(data)
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
