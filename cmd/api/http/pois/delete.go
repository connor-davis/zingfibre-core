package pois

import (
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (r *PointsOfInterestRouter) DeletePointOfInterestRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("404", &constants.NotFoundResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:     "id",
				In:       "path",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
				},
			},
		},
	}

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Delete Point Of Interest",
			Description: "Endpoint to delete a point of interest by ID",
			Tags:        []string{"Points Of Interest"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.DeleteMethod,
		Path:   "/pois/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			id := c.Params("id")

			_, err := r.Postgres.GetPointOfInterest(c.Context(), uuid.MustParse(id))

			if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"details": constants.NotFoundErrorDetails,
				})
			}

			_, err = r.Postgres.DeletePointOfInterest(c.Context(), uuid.MustParse(id))

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusNoContent).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
