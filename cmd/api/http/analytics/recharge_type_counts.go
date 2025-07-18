package analytics

import (
	"slices"
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
						"data": map[string][]system.RechargeTypeCount{
							"Weekly": {
								{
									Count:  1,
									Period: "01-01-1990",
								},
							},
							"Monthly": {
								{
									Count:  1,
									Period: "01-1990",
								},
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

			data := []map[string]interface{}{}
			uniquePeriods := []string{}
			uniqueTypes := []string{}

			for _, row := range rows {
				if !slices.Contains(uniquePeriods, string(row.RechargePeriod.([]byte))) {
					uniquePeriods = append(uniquePeriods, string(row.RechargePeriod.([]byte)))
					continue
				}
			}

			for _, row := range rows {
				if row.RechargeName.String == "" {
					row.RechargeName.String = "Intro Package"
				}

				if !slices.Contains(uniqueTypes, row.RechargeName.String) {
					uniqueTypes = append(uniqueTypes, row.RechargeName.String)
					continue
				}
			}

			for _, period := range uniquePeriods {
				base := map[string]interface{}{
					"Period": period,
				}

				for _, row := range rows {
					if string(row.RechargePeriod.([]byte)) == period {
						if row.RechargeName.String == "" {
							row.RechargeName.String = "Intro Package"
						}

						if _, exists := base[row.RechargeName.String]; !exists {
							base[row.RechargeName.String] = int(row.RechargeCount)
						} else {
							base[row.RechargeName.String] = int(base[row.RechargeName.String].(int64) + row.RechargeCount)
						}
					}
				}

				data = append(data, base)
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data": &fiber.Map{
					"Items": data,
					"Types": uniqueTypes,
				},
			})
		},
	}
}
