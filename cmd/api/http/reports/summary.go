package reports

import (
	"math"
	"strconv"
	"time"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (r *ReportsRouter) SummaryRoute() system.Route {
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
				Name:     "months",
				In:       "query",
				Required: true,
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
				Name:     "page",
				In:       "query",
				Required: true,
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
				Required: true,
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
			WithDescription("The zingfibre summary report").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
						"data": []system.ReportSummary{
							{
								DateCreated:    "2023-01-01T12:00:00Z",
								ItemName:       "Intro Package",
								RadiusUsername: "johnny",
								Method:         "Credit Card",
								Amount:         "100.00",
								ServiceId:      1,
								BuildName:      "Zing Build",
								BuildType:      "Standard",
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
			Summary:     "Summary Report",
			Description: "Endpoint to retrieve summary report",
			Tags:        []string{"Reports"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/reports/summary",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			poi := c.Query("poi")
			search := c.Query("search")

			months := c.Query("months")

			monthsInt, err := strconv.Atoi(months)

			if err != nil {
				log.Errorf("🔥 Error parsing months query parameter: %s", err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

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

			totalSummaries, err := r.Zing.GetReportsTotalSummaries(c.Context(), zing.GetReportsTotalSummariesParams{
				Poi:    poi,
				Search: search,
				Months: monthsInt,
			})

			if err != nil {
				log.Errorf("🔥 Error fetching total summaries from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			pages := int32(math.Ceil(float64(totalSummaries) / 10))

			summaries, err := r.Zing.GetReportsSummary(c.Context(), zing.GetReportsSummaryParams{
				Poi:    poi,
				Search: search,
				Months: monthsInt,
				Limit:  int32(pageSizeInt),
				Offset: int32((pageInt - 1) * pageSizeInt),
			})

			if err != nil {
				log.Errorf("🔥 Error fetching summaries from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			data := []system.ReportSummary{}

			for _, summary := range summaries {
				data = append(data, system.ReportSummary{
					DateCreated:    summary.DateCreated.Time.Format(time.RFC3339),
					ItemName:       string(summary.ItemName.([]byte)),
					RadiusUsername: summary.RadiusUsername.String,
					Method:         summary.Method.String,
					Amount:         summary.Amount.String,
					ServiceId:      summary.ServiceID.Int64,
					BuildName:      summary.BuildName.String,
					BuildType:      summary.BuildType.String,
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
