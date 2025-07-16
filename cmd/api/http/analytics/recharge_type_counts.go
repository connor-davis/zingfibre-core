package analytics

import (
	"strconv"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (r *AnalyticsRouter) RechargeTypeCountsRoute() system.Route {
	responses := openapi3.NewResponses()

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:     "count",
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
				Name:     "period",
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

			countParsed, err := strconv.Atoi(count)

			if err != nil {
				log.Errorf("🔥 Error converting count to int: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"message": constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			rows, err := r.Zing.GetRechargeTypeCounts(c.Context(), zing.GetRechargeTypeCountsParams{
				Period: period,
				Poi:    poi,
				Count:  countParsed,
			})

			if err != nil {
				log.Errorf("🔥 Error retrieving recharge type counts: %s", err.Error())

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
					Period: string(row.Period.([]byte)),
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
