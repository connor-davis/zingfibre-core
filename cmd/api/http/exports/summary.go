package exports

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
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

			summaries, err := r.Zing.GetReportsSummary(c.Context(), poi)

			if err != nil {
				log.Errorf("🔥 Error fetching summaries from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			now := time.Now()

			disposition := fmt.Sprintf(`attachment; filename="customers_report_%s.csv"`, now.Format(time.DateOnly))

			c.Set(fiber.HeaderContentType, "text/csv")
			c.Set(fiber.HeaderContentDisposition, disposition)

			writer := csv.NewWriter(c.Response().BodyWriter())

			header := []string{"Created On", "Email", "First Name", "Surname", "Item Name", "Amount", "Successful", "Service ID", "Build Name", "Build Type"}

			if err := writer.Write(header); err != nil {
				log.Errorf("🔥 Error writing CSV header: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			for _, summary := range summaries {
				record := []string{
					summary.DateCreated.Time.Format(time.RFC3339),
					string(summary.ItemName.([]byte)),
					summary.RadiusUsername.String,
					summary.AmountGross,
					summary.AmountFee,
					summary.AmountNet,
					summary.CashCode,
					fmt.Sprintf("%.2f", float32(summary.CashAmount)),
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
