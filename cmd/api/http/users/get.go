package users

import (
	"strconv"
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (r *UsersRouter) GetUsersRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessArrayResponse)
	responses.Set("401", &constants.UnauthorizedResponse)
	responses.Set("500", &constants.InternalServerErrorResponse)

	parameters := []*openapi3.ParameterRef{
		{
			Value: &openapi3.Parameter{
				Name:     "page",
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
	}

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Get Users",
			Description: "Endpoint to retrieve a list of users",
			Tags:        []string{"Users"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/users",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			page, err := strconv.Atoi(c.Query("page"))

			if err != nil {
				page = 1
			}

			users, err := r.Postgres.GetUsers(c.Context(), postgres.GetUsersParams{
				Limit:  10, // Default limit
				Offset: (int32(page) - 1) * 10,
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
				"data":    users,
			})
		},
	}
}

func (r *UsersRouter) GetUserRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessObjectResponse)
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
						Type: &openapi3.Types{
							"string",
						},
					},
				},
			},
		},
	}

	return system.Route{
		OpenAPIMetadata: system.OpenAPIMetadata{
			Summary:     "Get User",
			Description: "Endpoint to retrieve a user by ID",
			Tags:        []string{"Users"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/users/:id",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
		},
		Handler: func(c *fiber.Ctx) error {
			id := c.Params("id")

			user, err := r.Postgres.GetUser(c.Context(), uuid.MustParse(id))

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

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"data":    user,
			})
		},
	}
}
