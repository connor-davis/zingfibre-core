package reports

import (
	"math"
	"strconv"
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
				Name:     "sort",
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
			sort := c.Query("sort")

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

			expiringCustomersZing, err := r.Zing.GetReportsExpiringCustomers(c.Context())

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
						return strings.ToLower(i.(zing.GetReportsExpiringCustomersRow).RadiusUsername.String)
					},
					func(o interface{}) interface{} {
						return strings.ToLower(o.(radius.GetReportsExpiringCustomersRow).Username)
					},
					func(i interface{}, o interface{}) interface{} {
						return system.ReportExpiringCustomer{
							FullName:             i.(zing.GetReportsExpiringCustomersRow).FullName,
							Email:                i.(zing.GetReportsExpiringCustomersRow).Email.String,
							PhoneNumber:          i.(zing.GetReportsExpiringCustomersRow).PhoneNumber.String,
							RadiusUsername:       i.(zing.GetReportsExpiringCustomersRow).RadiusUsername.String,
							LastPurchaseDuration: i.(zing.GetReportsExpiringCustomersRow).LastPurchaseDuration.String,
							LastPurchaseSpeed:    i.(zing.GetReportsExpiringCustomersRow).LastPurchaseSpeed.String,
							Expiration:           o.(radius.GetReportsExpiringCustomersRow).Expiration.Time.Format(time.RFC3339),
							Address:              i.(zing.GetReportsExpiringCustomersRow).Address.String,
							POP:                  i.(zing.GetReportsExpiringCustomersRow).Pop.String,
						}
					},
				).
				Where(func(i interface{}) bool {
					return strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).POP), strings.ToLower(poi)) &&
						(strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).FullName), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).Email), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).PhoneNumber), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).RadiusUsername), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).LastPurchaseDuration), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).LastPurchaseSpeed), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).Address), strings.ToLower(search)) ||
							strings.Contains(strings.ToLower(i.(system.ReportExpiringCustomer).Expiration), strings.ToLower(search)))
				})

			orderedExpiringCustomers := []system.ReportExpiringCustomer{}

			switch sort {
			case "full_name_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).FullName)
				}).ToSlice(&orderedExpiringCustomers)
			case "full_name_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).FullName)
				}).ToSlice(&orderedExpiringCustomers)
			case "email_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).Email)
				}).ToSlice(&orderedExpiringCustomers)
			case "email_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).Email)
				}).ToSlice(&orderedExpiringCustomers)
			case "phone_number_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).PhoneNumber)
				}).ToSlice(&orderedExpiringCustomers)
			case "phone_number_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).PhoneNumber)
				}).ToSlice(&orderedExpiringCustomers)
			case "radius_username_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).RadiusUsername)
				}).ToSlice(&orderedExpiringCustomers)
			case "radius_username_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).RadiusUsername)
				}).ToSlice(&orderedExpiringCustomers)
			case "last_purchase_duration_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).LastPurchaseDuration)
				}).ToSlice(&orderedExpiringCustomers)
			case "last_purchase_duration_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).LastPurchaseDuration)
				}).ToSlice(&orderedExpiringCustomers)
			case "last_purchase_speed_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).LastPurchaseSpeed)
				}).ToSlice(&orderedExpiringCustomers)
			case "last_purchase_speed_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).LastPurchaseSpeed)
				}).ToSlice(&orderedExpiringCustomers)
			case "address_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).Address)
				}).ToSlice(&orderedExpiringCustomers)
			case "address_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).Address)
				}).ToSlice(&orderedExpiringCustomers)
			case "expiration_asc":
				expiringCustomersQuery.OrderBy(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).Expiration)
				}).ToSlice(&orderedExpiringCustomers)
			case "expiration_desc":
				expiringCustomersQuery.OrderByDescending(func(i interface{}) interface{} {
					return strings.ToLower(i.(system.ReportExpiringCustomer).Expiration)
				}).ToSlice(&orderedExpiringCustomers)
			}

			totalExpiringCustomers := linq.From(orderedExpiringCustomers).
				Count()

			pages := int(math.Ceil(float64(totalExpiringCustomers) / float64(pageSizeInt)))

			data := linq.From(orderedExpiringCustomers).
				Skip((pageInt - 1) * pageSizeInt).
				Take(pageSizeInt).
				Results()

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    data,
				"pages":   pages,
			})
		},
	}
}
