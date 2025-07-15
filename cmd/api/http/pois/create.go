package pois

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CreatePointOfInterestRequest struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (r *PointsOfInterestRouter) CreatePointOfInterestRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("201", &constants.SuccessResponse)
	responses.Set("400", &constants.BadRequestResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("409", &constants.ConflictResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Create Point Of Interest",
			Description: "Endpoint to create a new point of interest",
			Tags:        []string{"Points Of Interest"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().WithJSONSchema(schemas.CreatePointOfInterestSchema.Value),
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/pois",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff),
		},
		Handler: func(c *fiber.Ctx) error {
			var createPointOfInterestRequest CreatePointOfInterestRequest

			if err := c.BodyParser(&createPointOfInterestRequest); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			_, err := r.Postgres.CreatePointOfInterest(c.Context(), postgres.CreatePointOfInterestParams{
				Name: createPointOfInterestRequest.Name,
				Key:  createPointOfInterestRequest.Key,
			})

			if err != nil {
				log.Errorf("🔥 Error creating point of interest: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
			})
		},
	}
}
