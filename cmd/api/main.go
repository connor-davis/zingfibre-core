package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/connor-davis/zingfibre-core/cmd/api/http"
	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/common"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/connor-davis/zingfibre-core/internal/sessions"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/trinodb/trino-go-client/trino"
)

func main() {
	context := context.Background()

	log.Info("✅ Starting Zingfibre Reporting API...")
	log.Info("🔃 Connecting to Zingfibre databases...")

	zingConnection, err := sql.Open("mysql", common.EnvString("ZING_DSN", ""))

	if err != nil {
		log.Errorf("🔥 Failed to connect to Zing database: %s", err.Error())

		return
	}

	defer zingConnection.Close()

	log.Info("✅ Connected to Zingfibre Zing database successfully")

	radiusConnection, err := sql.Open("mysql", common.EnvString("RADIUS_DSN", ""))

	if err != nil {
		log.Errorf("🔥 Failed to connect to Radius database: %s", err.Error())

		return
	}

	defer radiusConnection.Close()

	log.Info("✅ Connected to Zingfibre Radius database successfully")

	zingQueries := zing.New(zingConnection)
	radiusQueries := radius.New(radiusConnection)

	log.Info("🔃 Connecting to PostgreSQL database...")

	postgresPoolConfig, err := pgxpool.ParseConfig(common.EnvString("POSTGRES_DSN", ""))

	if err != nil {
		log.Infof("🔥 Failed to parse PostgreSQL connection string: %s", err.Error())

		return
	}

	postgresPool, err := pgxpool.NewWithConfig(context, postgresPoolConfig)

	if err != nil {
		log.Infof("🔥 Failed to connect to PostgreSQL: %s", err.Error())
	}

	defer postgresPool.Close()

	log.Info("✅ Connected to PostgreSQL successfully")
	log.Info("🔃 Connecting to TrinoDB database...")

	trinoDsn := "http://user@trino:8080"
	trinoDb, err := sql.Open("trino", trinoDsn)

	if err != nil {
		log.Fatalf("🔥 Failed to connect to Trino database: %s", err.Error())

		return
	}

	defer trinoDb.Close()

	log.Info("✅ Connected to TrinoDB successfully")

	sessions := sessions.NewSessions(postgresPool)
	postgresQueries := postgres.New(postgresPool)

	log.Info("🔃 Creating default admin user.")

	existingAdmin, err := postgresQueries.GetUserByEmail(context, common.EnvString("ADMIN_EMAIL", ""))

	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		log.Infof("🔥 Error checking for existing admin user: %s", err.Error())
		return
	}

	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		log.Info("🔃 No existing admin user found, creating a new one...")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(common.EnvString("ADMIN_PASSWORD", "")), bcrypt.DefaultCost)

		if err != nil {
			log.Errorf("🔥 Error hashing admin password: %s", err.Error())

			return
		}

		if _, err := postgresQueries.CreateUser(context, postgres.CreateUserParams{
			Email:    string(common.EnvString("ADMIN_EMAIL", "")),
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

	middleware := middleware.NewMiddleware(postgresQueries, sessions)

	httpRouter := http.NewHttpRouter(postgresQueries, zingQueries, radiusQueries, middleware, sessions, trinoDb)

	openapiSpecification := httpRouter.InitializeOpenAPI()

	api := app.Group("/api")

	api.Get("/api-spec", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(openapiSpecification)
	})

	api.Get("/api-doc", func(c *fiber.Ctx) error {
		html, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: func() string {
				if common.EnvString("MODE", "development") == "production" {
					return "https://reports.core.zingfibre.co.za/api/api-spec"
				}

				return fmt.Sprintf("http://localhost:%s/api/api-spec", common.EnvString("PORT", ""))
			}(),
			Theme:  scalar.ThemeDefault,
			Layout: scalar.LayoutModern,
			BaseServerURL: func() string {
				if common.EnvString("MODE", "development") == "production" {
					return "https://reports.core.zingfibre.co.za"
				}

				return fmt.Sprintf("http://localhost:%s", common.EnvString("PORT", ""))
			}(),
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Type("html").SendString(html)
	})

	httpRouter.InitializeRoutes(api)

	log.Info("✅ Starting Zingfibre Reporting API on port 6173...")

	if err := app.Listen(fmt.Sprintf(":%s", common.EnvString("PORT", "6173"))); err != nil {
		log.Errorf("🔥 Failed to start server: %v", err)
	}
}
