package exports

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (r *ExportsRouter) ExpiringCustomersRoute() system.Route {
	responses := openapi3.NewResponses()

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:     "poi",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"string",
						},
					},
				},
			},
		},
	}

	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Content: map[string]*openapi3.MediaType{
				"text/csv": {
					Schema: openapi3.NewSchema().WithFormat("text").NewRef(),
				},
			},
		},
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("The user is not authenticated.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.UnauthorizedError,
						"details": constants.UnauthorizedErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	responses.Set("500", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Internal server error").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.InternalServerError,
						"details": constants.InternalServerErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Expiring Customers Report Export",
			Description: "Endpoint to retrieve expiring customers report export in CSV format.",
			Tags:        []string{"Exports"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/exports/expiring-customers",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			poi := c.Query("poi")

			expiringCustomersRadius, err := r.Radius.GetReportsExpiringCustomers(c.Context(), poi)

			if err != nil {
				log.Errorf("🔥 Error fetching expiring customers from Radius: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			expiringCustomers, err := r.Zing.GetReportExportsExpiringCustomers(c.Context(), zing.GetReportExportsExpiringCustomersParams{
				Expiration: "",
				Address:    "",
				Poi:        poi,
			})

			if err != nil {
				log.Errorf("🔥 Error fetching expiring customers from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			now := time.Now()

			disposition := fmt.Sprintf(`attachment; filename="expiring_customers_report_%s.csv"`, now.Format(time.DateOnly))

			c.Set(fiber.HeaderContentType, "text/csv")
			c.Set(fiber.HeaderContentDisposition, disposition)

			writer := csv.NewWriter(c.Response().BodyWriter())

			header := []string{"Expires On", "Full Name", "Email", "Phone Number", "Radius Username", "Last Purchase Duration", "Last Purchase Speed", "Address"}

			if err := writer.Write(header); err != nil {
				log.Errorf("🔥 Error writing CSV header: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			for _, expiringCustomer := range expiringCustomers {
				expiringCustomerRadius := &radius.GetReportsExpiringCustomersRow{}

				for _, radiusCustomer := range expiringCustomersRadius {
					if radiusCustomer.Username == expiringCustomer.RadiusUsername.String {
						expiringCustomerRadius = &radiusCustomer
						break
					}
				}

				record := []string{
					expiringCustomerRadius.Expiration.Time.Format(time.DateTime),
					expiringCustomer.FullName,
					expiringCustomer.Email.String,
					expiringCustomer.PhoneNumber.String,
					expiringCustomer.RadiusUsername.String,
					expiringCustomer.LastPurchaseDuration.String,
					expiringCustomer.LastPurchaseSpeed.String,
					expiringCustomerRadius.Address,
				}

				if err := writer.Write(record); err != nil {
					log.Errorf("🔥 Error writing CSV record: %s", err.Error())

					continue
				}
			}

			defer writer.Flush()

			if err := writer.Error(); err != nil {
				log.Errorf("🔥 Error flushing CSV writer: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return nil
		},
	}
}
