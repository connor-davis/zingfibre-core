package analytics

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *AnalyticsRouter) RechargeTypeCountsRoute() system.Route {
	responses := openapi3.NewResponses()

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:        "startDate",
				In:          "query",
				Description: "Start date for the report",
				Required:    true,
				Schema:      openapi3.NewSchemaRef("string", nil),
			},
		},
		{
			Value: &openapi3.Parameter{
				Name:        "endDate",
				In:          "query",
				Description: "End date for the report",
				Required:    true,
				Schema:      openapi3.NewSchemaRef("string", nil),
			},
		},
	}

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
		Path:   "/authentication/check",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
