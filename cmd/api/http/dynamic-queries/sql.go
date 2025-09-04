package dynamicQueries

import (
	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/helpers"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func (r *DynamicQueriesRouter) SQL() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithDescription("The SQL query string.").
			WithContent(openapi3.Content{
				"text/plain": &openapi3.MediaType{
					Example: "SELECT * FROM users;",
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{
								"string",
							},
						},
					},
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

	responses.Set("404", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.ErrorResponseSchema.Value,
			).
			WithDescription("User not found.").
			WithContent(openapi3.Content{
				"application/json": &openapi3.MediaType{
					Example: map[string]any{
						"error":   constants.NotFoundError,
						"details": constants.NotFoundErrorDetails,
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

	parameters := []*openapi3.ParameterRef{}

	body := &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody(),
	}

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "SQL",
			Description: "Endpoint to retrieve a SQL query string.",
			Tags:        []string{"Dynamic Queries"},
			Parameters:  parameters,
			RequestBody: body,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/dynamic-queries",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff),
		},
		Handler: func(c *fiber.Ctx) error {
			userEmailsDynamicQuery := system.DynamicQuery{
				Database: "zing",
				Table: system.DynamicQueryTable{
					Table:     "user_emails",
					IsPrimary: false,
				},
				Columns: []system.DynamicQueryColumn{
					{
						Column: "email",
						Label:  "Email",
					},
					{
						Column: "user_id",
						Label:  "UserId",
					},
				},
			}

			dynamicQuery := system.DynamicQuery{
				Database: "zing",
				Table: system.DynamicQueryTable{
					Table:     "users",
					IsPrimary: true,
				},
				Columns: []system.DynamicQueryColumn{
					{
						Column: "t1.id",
						Label:  "Id",
					},
					{
						Column: "t2.email",
						Label:  "Email",
					},
				},
				Joins: []system.DynamicQueryJoin{
					{
						Type:          system.InnerJoin,
						LocalDatabase: "zing",
						LocalTable: system.DynamicQueryTable{
							Table:     "users",
							IsPrimary: true,
						},
						LocalColumn:       "id",
						ReferenceDatabase: "zing",
						ReferenceTable: system.DynamicQueryTable{
							Table:     "user_emails",
							IsPrimary: false,
						},
						ReferenceColumn: "user_id",
						SubQuery:        &userEmailsDynamicQuery,
					},
				},
				Filters: []system.DynamicQueryFilter{
					{
						Column:   "t2.email",
						Operator: "LIKE",
						Type:     system.StringFilter,
						Value:    "%@zingfibre.co.za",
					},
				},
				Orders: []system.DynamicQueryOrder{
					{
						Column:     "t2.email",
						Descending: false,
					},
				},
				SubQueries: []system.DynamicQuery{
					userEmailsDynamicQuery,
				},
			}

			return c.Status(fiber.StatusOK).
				SendString(helpers.DynamicQueryParser(dynamicQuery))
		},
	}
}
