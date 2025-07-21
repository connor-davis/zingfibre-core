package exports

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"
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

			expiringCustomersRadius, err := r.Radius.GetReportsExpiringCustomers(c.Context())

			if err != nil {
				log.Errorf("🔥 Error fetching expiring customers from Radius: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			expiringCustomersZing, err := r.Zing.GetReportExportsExpiringCustomers(c.Context())

			if err != nil {
				log.Errorf("🔥 Error fetching expiring customers from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			linqExpiringCustomersZing := linq.From(expiringCustomersZing)
			linqExpiringCustomersRadius := linq.From(expiringCustomersRadius)

			expiringCustomersQuery := linqExpiringCustomersZing.
				Join(
					linqExpiringCustomersRadius,
					func(i interface{}) interface{} {
						return strings.ToLower(i.(zing.GetReportExportsExpiringCustomersRow).RadiusUsername.String)
					},
					func(o interface{}) interface{} {
						return strings.ToLower(o.(radius.GetReportsExpiringCustomersRow).Username)
					},
					func(i interface{}, o interface{}) interface{} {
						return system.ReportExpiringCustomer{
							FullName:             i.(zing.GetReportExportsExpiringCustomersRow).FullName,
							Email:                i.(zing.GetReportExportsExpiringCustomersRow).Email.String,
							PhoneNumber:          i.(zing.GetReportExportsExpiringCustomersRow).PhoneNumber.String,
							RadiusUsername:       i.(zing.GetReportExportsExpiringCustomersRow).RadiusUsername.String,
							LastPurchaseDuration: i.(zing.GetReportExportsExpiringCustomersRow).LastPurchaseDuration.String,
							LastPurchaseSpeed:    i.(zing.GetReportExportsExpiringCustomersRow).LastPurchaseSpeed.String,
							Expiration:           o.(radius.GetReportsExpiringCustomersRow).Expiration.Time.Format(time.RFC3339),
							Address:              i.(zing.GetReportExportsExpiringCustomersRow).Address.String,
							POP:                  i.(zing.GetReportExportsExpiringCustomersRow).Pop.String,
						}
					},
				).
				Where(func(i interface{}) bool {
					return strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).POP), strings.ToLower(poi))
				}).
				OrderByDescending(func(i interface{}) interface{} {
					return i.(system.ReportExpiringCustomer).Expiration
				}).
				ThenBy(func(i interface{}) interface{} {
					return i.(system.ReportExpiringCustomer).FullName
				}).
				Results()

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

			for _, customer := range expiringCustomersQuery {
				expiringCustomer := customer.(system.ReportExpiringCustomer)

				record := []string{
					expiringCustomer.Expiration,
					expiringCustomer.FullName,
					expiringCustomer.Email,
					expiringCustomer.PhoneNumber,
					expiringCustomer.RadiusUsername,
					expiringCustomer.LastPurchaseDuration,
					expiringCustomer.LastPurchaseSpeed,
					expiringCustomer.Address,
				}

				if err := writer.Write(record); err != nil {
					log.Errorf("🔥 Error writing CSV record: %s", err.Error())

					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
						"error":   constants.InternalServerError,
						"details": constants.InternalServerErrorDetails,
					})
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
