package dynamicQueries

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateDynamicQueryRequest struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
}

func (r *DynamicQueriesRouter) CreateDynamicQueryRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("201", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("Dynamic Query created successfully.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"message": constants.Success,
						"details": constants.SuccessDetails,
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
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("Unauthorized.").
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
			Summary:     "Create Dynamic Query",
			Description: "Endpoint to create a new dynamic query",
			Tags:        []string{"Dynamic Queries"},
			Parameters:  nil,
			RequestBody: &openapi3.RequestBodyRef{
				Value: openapi3.NewRequestBody().WithJSONSchema(schemas.CreateDynamicQuerySchema.Value),
			},
			Responses: responses,
		},
		Method: system.PostMethod,
		Path:   "/dynamic-queries",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasRole(postgres.RoleTypeAdmin),
		},
		Handler: func(c *fiber.Ctx) error {
			var createDynamicQueryRequest CreateDynamicQueryRequest

			if err := c.BodyParser(&createDynamicQueryRequest); err != nil {
				log.Errorf("🔥 Error parsing request body: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			dynamicQuery, err := r.Postgres.CreateDynamicQuery(
				c.Context(),
				postgres.CreateDynamicQueryParams{
					Name:       createDynamicQueryRequest.Name,
					Query:      pgtype.Text{String: "", Valid: false},
					ResponseID: pgtype.Text{String: "", Valid: false},
					Status:     postgres.DynamicQueryStatusInProgress,
					Prompt:     createDynamicQueryRequest.Prompt,
				},
			)

			if err != nil {
				log.Errorf("🔥 Error creating dynamic query: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    dynamicQuery.ID,
			})
		},
	}
}
