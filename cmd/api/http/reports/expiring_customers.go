package reports

import (
	"database/sql"
	"math"
	"strconv"

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

func (r *ReportsRouter) ExpiringCustomersRoute() system.Route {
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
		{
			Value: &openapi3.Parameter{
				Name:     "page",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"integer",
						},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:     "pageSize",
				In:       "query",
				Required: false,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"integer",
						},
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:     "search",
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
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("The zingfibre expiring customers report").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
						"data": []system.ReportExpiringCustomer{
							{
								FullName:             "Jane Smith",
								Email:                "jane.smith@example.com",
								PhoneNumber:          "987-654-3210",
								RadiusUsername:       "janesmith",
								LastPurchaseDuration: "30 days",
								LastPurchaseSpeed:    "100 Mbps",
								Expiration:           "2023-12-31T23:59:59Z",
								Address:              "123 Main St, Anytown, USA",
							},
						},
					},
					Schema: schemas.SuccessResponseSchema,
				},
			}),
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
			Summary:     "Expiring Customers Report",
			Description: "Endpoint to retrieve expiring customers report",
			Tags:        []string{"Reports"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/reports/expiring-customers",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			poi := c.Query("poi")
			search := c.Query("search")

			page := c.Query("page")
			pageSize := c.Query("pageSize")

			pageInt, err := strconv.Atoi(page)

			if err != nil {
				pageInt = 1
			}

			pageSizeInt, err := strconv.Atoi(pageSize)

			if err != nil {
				pageSizeInt = 10
			}

			expiringCustomersRadius, err := r.Radius.GetReportsExpiringCustomers(c.Context())

			if err != nil {
				log.Errorf("🔥 Error fetching expiring customers from Radius: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			totalExpiringCustomers, err := r.Zing.GetReportsTotalExpiringCustomers(c.Context(), zing.GetReportsTotalExpiringCustomersParams{
				Poi:    poi,
				Search: search,
				RadiusUsernames: func() []sql.NullString {
					usernames := make([]sql.NullString, len(expiringCustomersRadius))
					for i, customer := range expiringCustomersRadius {
						usernames[i] = sql.NullString{
							String: customer.Username,
							Valid:  true,
						}
					}
					return usernames
				}(),
			})

			if err != nil {
				log.Errorf("🔥 Error fetching total expiring customers from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			pages := int32(math.Ceil(float64(totalExpiringCustomers) / float64(pageSizeInt)))

			expiringCustomers, err := r.Zing.GetReportsExpiringCustomers(c.Context(), zing.GetReportsExpiringCustomersParams{
				Poi:    poi,
				Search: search,
				Limit:  int32(pageSizeInt),
				Offset: int32((pageInt - 1) * pageSizeInt),
				RadiusUsernames: func() []sql.NullString {
					usernames := make([]sql.NullString, len(expiringCustomersRadius))
					for i, customer := range expiringCustomersRadius {
						usernames[i] = sql.NullString{
							String: customer.Username,
							Valid:  true,
						}
					}
					return usernames
				}(),
			})

			if err != nil {
				log.Errorf("🔥 Error fetching expiring customers from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			data := []system.ReportExpiringCustomer{}

			for _, expiringCustomer := range expiringCustomers {
				expiringCustomerRadius := &radius.GetReportsExpiringCustomersRow{}

				for _, radiusCustomer := range expiringCustomersRadius {
					if radiusCustomer.Username == expiringCustomer.RadiusUsername.String {
						expiringCustomerRadius = &radiusCustomer
						break
					}
				}

				data = append(data, system.ReportExpiringCustomer{
					FullName:             expiringCustomer.FullName,
					Email:                expiringCustomer.Email.String,
					PhoneNumber:          expiringCustomer.PhoneNumber.String,
					RadiusUsername:       expiringCustomer.RadiusUsername.String,
					LastPurchaseDuration: expiringCustomer.LastPurchaseDuration.String,
					LastPurchaseSpeed:    expiringCustomer.LastPurchaseSpeed.String,
					Expiration:           expiringCustomerRadius.Expiration.Time.Format("2006-01-02T15:04:05Z07:00"),
					Address:              expiringCustomer.Address.String,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    data,
				"pages":   pages,
			})
		},
	}
}
