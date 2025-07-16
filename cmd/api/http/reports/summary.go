package reports

import (
	"time"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
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
								AmountGross:    "100.00",
								AmountFee:      "5.00",
								AmountNet:      "95.00",
								CashCode:       "ABC123",
								CashAmount:     "95.00",
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

			summaries, err := r.Zing.GetReportsSummary(c.Context(), poi)

			if err != nil {
				log.Errorf("🔥 Error fetching summaries from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			data := []system.ReportSummary{}

			for _, summary := range summaries {
				cashAmount := summary.CashAmount

				if cashAmount == nil {
					cashAmount = new(float64)
				}

				data = append(data, system.ReportSummary{
					DateCreated:    summary.DateCreated.Time.Format(time.RFC3339),
					ItemName:       string(summary.ItemName.([]byte)),
					RadiusUsername: summary.RadiusUsername.String,
					AmountGross:    string(summary.AmountGross.([]byte)),
					AmountFee:      string(summary.AmountFee.([]byte)),
					AmountNet:      string(summary.AmountNet.([]byte)),
					CashCode:       string(summary.CashCode.([]byte)),
					CashAmount:     string(summary.CashAmount.([]byte)),
					ServiceId:      summary.ServiceID.Int64,
					BuildName:      summary.BuildName.String,
					BuildType:      summary.BuildType.String,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    data,
			})
		},
	}
}
