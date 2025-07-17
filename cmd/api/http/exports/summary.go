package exports

import (
	"encoding/csv"
	"fmt"
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

func (r *ExportsRouter) SummaryRoute() system.Route {
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
			Summary:     "Summary Report Export",
			Description: "Endpoint to retrieve summary report export in CSV format.",
			Tags:        []string{"Exports"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/exports/summary",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			poi := c.Query("poi")
			months := c.Query("months")

			monthsInt, err := strconv.Atoi(months)

			if err != nil {
				log.Errorf("🔥 Error converting months query parameter to integer: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			summaries, err := r.Zing.GetReportExportsSummary(c.Context(), zing.GetReportExportsSummaryParams{
				Poi:    poi,
				Months: monthsInt,
			})

			if err != nil {
				log.Errorf("🔥 Error fetching summaries from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			now := time.Now()

			disposition := fmt.Sprintf(`attachment; filename="summary_report_%s.csv"`, now.Format(time.DateOnly))

			c.Set(fiber.HeaderContentType, "text/csv")
			c.Set(fiber.HeaderContentDisposition, disposition)

			writer := csv.NewWriter(c.Response().BodyWriter())

			header := []string{"Created On", "Item Name", "Radius Username", "Method", "Amount Gross", "Amount Fee", "Amount Net", "Cash Code", "Cash Amount", "Service ID", "Build Name", "Build Type"}

			if err := writer.Write(header); err != nil {
				log.Errorf("🔥 Error writing CSV header: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			for _, summary := range summaries {
				record := []string{
					summary.DateCreated.Time.Format(time.DateTime),
					string(summary.ItemName.([]byte)),
					summary.RadiusUsername.String,
					summary.Method.String,
					string(summary.AmountGross.([]byte)),
					string(summary.AmountFee.([]byte)),
					string(summary.AmountNet.([]byte)),
					string(summary.CashCode.([]byte)),
					string(summary.CashAmount.([]byte)),
					fmt.Sprintf("%d", summary.ServiceID.Int64),
					summary.BuildName.String,
					summary.BuildType.String,
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
