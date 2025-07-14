package pois

import (
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UpdatePointOfInterestRequest struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (r *PointsOfInterestRouter) UpdatePointOfInterestRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessObjectResponse)
	responses.Set("400", &constants.BadRequestResponse)
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
			Summary:     "Update Point Of Interest",
			Description: "Endpoint to update an existing point of interest",
			Tags:        []string{"Point Of Interests"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.PutMethod,
		Path:   "/pois/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			var updatePointOfInterestRequest UpdatePointOfInterestRequest

			if err := c.BodyParser(&updatePointOfInterestRequest); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			id := c.Params("id")

			poi, err := r.Postgres.GetPointOfInterest(c.Context(), uuid.MustParse(id))

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

			_, err = r.Postgres.UpdatePointOfInterest(c.Context(), postgres.UpdatePointOfInterestParams{
				ID:   poi.ID,
				Name: updatePointOfInterestRequest.Name,
				Key:  updatePointOfInterestRequest.Key,
			})

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
