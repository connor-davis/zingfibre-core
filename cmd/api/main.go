package main

import (
	"github.com/connor-davis/zingfibre-core/cmd/api/http"
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

	httpRouter.InitializeOpenAPI()
	httpRouter.InitializeRoutes(app)

	log.Info("Starting Zingfibre Reporting API on port 4000...")

	if err := app.Listen(":4000"); err != nil {
		log.Infof("Failed to start server: %v", err)
	}
}
