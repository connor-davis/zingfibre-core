package users

import (
	"math"
	"strconv"
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func (r *UsersRouter) GetUsersRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &constants.SuccessPagingResponse)
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
		{
			Value: &openapi3.Parameter{
				Name:     "search",
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
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff),
		},
		Handler: func(c *fiber.Ctx) error {
			page, err := strconv.Atoi(c.Query("page"))

			if err != nil {
				page = 1
			}

			totalUsers, err := r.Postgres.GetTotalUsers(c.Context())

			if err != nil {
				log.Errorf("🔥 Error retrieving total users: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			users, err := r.Postgres.GetUsers(c.Context(), postgres.GetUsersParams{
				Limit:      10, // Default limit
				Offset:     (int32(page) - 1) * 10,
				SearchTerm: c.Query("search"),
			})

			if err != nil {
				log.Errorf("🔥 Error retrieving users: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			pages := int32(math.Ceil(float64(totalUsers) / 10))

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"message": constants.Success,
				"details": constants.SuccessDetails,
				"pages":   pages,
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
		Path:   "/users/{id}",
		Middlewares: []fiber.Handler{
			r.Middleware.Authorized(),
			r.Middleware.HasAnyRole(postgres.RoleTypeAdmin, postgres.RoleTypeStaff, postgres.RoleTypeUser),
		},
		Handler: func(c *fiber.Ctx) error {
			id, err := uuid.Parse(c.Params("id"))

			if err != nil {
				log.Errorf("🔥 Invalid UUID format: %s", err.Error())

				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"error":   constants.BadRequestError,
					"details": constants.BadRequestErrorDetails,
				})
			}

			user, err := r.Postgres.GetUser(c.Context(), id)

			if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
				log.Errorf("🔥 Error retrieving user: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				log.Warnf("⚠️ User with ID %s not found", id)

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
