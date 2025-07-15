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

	responses.Set("201", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("Point of interest created successfully.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Created,
						"details": constants.CreatedDetails,
					},
					Schema: schemas.SuccessResponseSchema,
				},
			}),
	})

	responses.Set("400", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Bad Request.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.BadRequestError,
						"details": constants.BadRequestErrorDetails,
					},
					Schema: schemas.ErrorResponseSchema,
				},
			}),
	})

	responses.Set("401", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().WithJSONSchema(
			schemas.ErrorResponseSchema.Value,
		).WithDescription("The user is not authenticated.").WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"error":   constants.UnauthorizedError,
					"details": constants.UnauthorizedErrorDetails,
				},
				Schema: schemas.ErrorResponseSchema,
			},
		}),
	})

	responses.Set("409", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Conflict.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.ConflictError,
						"details": constants.ConflictErrorDetails,
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
			WithDescription("Internal Server Error.").
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
