package reports

import (
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

func (r *ReportsRouter) RechargesRoute() system.Route {
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
				Name:     "startDate",
				In:       "query",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"string",
						},
						Format: "date-time",
					},
				},
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:     "endDate",
				In:       "query",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{
							"string",
						},
						Format: "date-time",
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
			WithDescription("The zingfibre recharges report").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
						"data": []system.ReportRecharge{
							{
								DateCreated: "2023-01-01T12:00:00Z",
								Email:       "john.doe@example.com",
								FirstName:   "John",
								Surname:     "Doe",
								ItemName:    "Monthly Subscription",
								Amount:      29.99,
								Successful:  true,
								ServiceId:   123,
								BuildName:   "Zing Fibre",
								BuildType:   "Fibre",
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
			Summary:     "Recharges Report",
			Description: "Endpoint to retrieve recharges report",
			Tags:        []string{"Reports"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/reports/recharges",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			poi := c.Query("poi")
			startDate := c.Query("startDate")
			endDate := c.Query("endDate")

			startDateParsed, err := time.Parse(time.RFC3339, startDate)

			if err != nil {
				log.Errorf("🔥 Error parsing start date: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			endDateParsed, err := time.Parse(time.RFC3339, endDate)

			if err != nil {
				log.Errorf("🔥 Error parsing end date: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			recharges, err := r.Zing.GetReportsRecharges(c.Context(), zing.GetReportsRechargesParams{
				Poi:       poi,
				StartDate: startDateParsed,
				EndDate:   endDateParsed,
			})

			if err != nil {
				log.Errorf("🔥 Error fetching recharges from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			data := []system.ReportRecharge{}

			for _, recharge := range recharges {
				amount, err := strconv.ParseFloat(recharge.Amount.String, 64)

				if err != nil {
					log.Errorf("🔥 Error parsing amount for recharge: %s", err.Error())

					continue
				}

				data = append(data, system.ReportRecharge{
					DateCreated: recharge.DateCreated.Format(time.RFC3339),
					Email:       recharge.Email.String,
					FirstName:   recharge.FirstName.String,
					Surname:     recharge.Surname.String,
					ItemName:    string(recharge.ItemName.([]byte)),
					Amount:      amount,
					Successful:  recharge.Successful,
					ServiceId:   recharge.ServiceID.Int64,
					BuildName:   recharge.BuildName.String,
					BuildType:   recharge.BuildType.String,
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

func (r *ReportsRouter) RechargesSummaryRoute() system.Route {
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
			WithDescription("The zingfibre recharges summary report").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
						"data": []system.ReportRecharge{
							{
								DateCreated: "2023-01-01T12:00:00Z",
								Email:       "john.doe@example.com",
								FirstName:   "John",
								Surname:     "Doe",
								ItemName:    "Monthly Subscription",
								Amount:      29.99,
								Successful:  true,
								ServiceId:   123,
								BuildName:   "Zing Fibre",
								BuildType:   "Fibre",
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
			Summary:     "Recharges Summary Report",
			Description: "Endpoint to retrieve recharges summary report",
			Tags:        []string{"Reports"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/reports/recharges-summary",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			poi := c.Query("poi")

			rechargeSummaries, err := r.Zing.GetReportsRechargesSummary(c.Context(), poi)

			if err != nil {
				log.Errorf("🔥 Error fetching recharges summary from Zing: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			data := []system.ReportRecharge{}

			for _, rechargeSummary := range rechargeSummaries {
				amount, err := strconv.ParseFloat(rechargeSummary.Amount.String, 64)

				if err != nil {
					log.Errorf("🔥 Error parsing amount for recharge: %s", err.Error())

					continue
				}

				data = append(data, system.ReportRecharge{
					DateCreated: rechargeSummary.DateCreated.Format(time.RFC3339),
					Email:       rechargeSummary.Email.String,
					FirstName:   rechargeSummary.FirstName.String,
					Surname:     rechargeSummary.Surname.String,
					ItemName:    string(rechargeSummary.ItemName.([]byte)),
					Amount:      amount,
					Successful:  rechargeSummary.Successful,
					ServiceId:   rechargeSummary.ServiceID.Int64,
					BuildName:   rechargeSummary.BuildName.String,
					BuildType:   rechargeSummary.BuildType.String,
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
