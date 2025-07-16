package analytics

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AnalyticsRouter) RechargeTypeCountsRoute() system.Route {
	responses := openapi3.NewResponses()

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:        "count",
				In:          "query",
				Description: "The number of weeks/months to look back for the report",
				Required:    true,
				Schema:      openapi3.NewSchemaRef("integer", nil),
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:        "period",
				In:          "query",
				Description: "The period for the report (weeks/months)",
				Required:    true,
				Schema:      openapi3.NewSchemaRef("string", nil),
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:        "poi",
				In:          "query",
				Description: "Point of interest for the report",
				Required:    true,
				Schema:      openapi3.NewSchemaRef("string", nil),
			},
		},
	}

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("The recharge type counts for the specified date range").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
						"data": []system.RechargeTypeCount{
							{
								Count:  1,
								Type:   "Example",
								Period: "Weekly",
							},
							{
								Count:  1,
								Type:   "Example",
								Period: "Monthly",
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
			Summary:     "Recharge Type Counts",
			Description: "Endpoint to retrieve recharge type counts over a specified date range",
			Tags:        []string{"Analytics"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/analytics/recharge-type-counts",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			count := c.Query("count")
			period := c.Query("period")
			poi := c.Query("poi")

			rows, err := r.Zing.GetRechargeTypeCounts(c.Context(), zing.GetRechargeTypeCountsParams{
				Period: period,
				Count:  count,
				Poi:    poi,
			})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"message": constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			data := []system.RechargeTypeCount{}

			for _, row := range rows {
				data = append(data, system.RechargeTypeCount{
					Count:  int(row.RechargeCount),
					Type:   row.ProductName.String,
					Period: row.Period.(string),
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
