package dynamicQueries

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/constants"
	"github.com/connor-davis/zingfibre-core/internal/models/schemas"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mateuszkardas/toon-go"
	"github.com/openai/openai-go/v3"
	openaiResponses "github.com/openai/openai-go/v3/responses"
	"github.com/valyala/fasthttp"
)

type GenerateDynamicQueryOutput struct {
	SqlQuery string `json:"sql_query" jsonschema_description:"The SQL query to be executed for the dynamic query."`
}

var GenerateDynamicQueryOutputSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"sql_query": map[string]any{
			"type":                   "string",
			"description":            "The SQL query to be executed for the dynamic query.",
			"jsonschema_description": "The SQL query to be executed for the dynamic query.",
		},
	},
	"required":             []string{"sql_query"},
	"additionalProperties": false,
}

func (r *DynamicQueriesRouter) GenerateDynamicQueryRoute() system.Route {
	responses := openapi3.NewResponses()

	responses.Set("200", &openapi3.ResponseRef{
		Value: openapi3.NewResponse().
			WithJSONSchema(
				schemas.SuccessResponseSchema.Value,
			).
			WithDescription("Dynamic Query generated successfully.").
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
			Summary:     "Generate Dynamic Query",
			Description: "Endpoint to generate a dynamic query by ID",
			Tags:        []string{"Dynamic Queries"},
			Parameters:  parameters,
			RequestBody: nil,
			Responses:   responses,
		},
		Method: system.GetMethod,
		Path:   "/dynamic-queries/{id}/generate",
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

			dynamicQuery, err := r.Postgres.GetDynamicQuery(c.Context(), id)

			if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
				log.Errorf("🔥 Error retrieving dynamic query: %s", err.Error())

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			if err != nil && strings.Contains(err.Error(), "no rows in result set") {
				log.Warnf("⚠️ Dynamic Query with ID %s not found", id)

				return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
					"error":   constants.NotFoundError,
					"details": constants.NotFoundErrorDetails,
				})
			}

			log.Infof("Generating dynamic query for Query ID: %s with prompt: %s", dynamicQuery.ID, dynamicQuery.Prompt)

			encodedInstructions, err := toon.Encode(
				map[string]any{
					"role": "You are a professional TrinoDB SQL query developer.",
					"expectations": []string{
						"You are expected to call `list-catalogs` first to view connected databases.",
						"You are expected to call `list-schemas` second to view connected database schemas.",
						"You are expected to call `list-tables` third to view connected database schema tables.",
						"You are expected to call tables like {catalog}.{schema}.{table} when writing your queries.",
						"You are expected to fulfill the users expectation without error.",
						"You are expected to use double quotes to quote identifiers.",
						"Do not quote catalogs, schemas and tables.",
						"You are expected to only utilise local information stores, table and column definitions that you have found. Do not invent any non-existent columns or tables.",
						"Only provide the SQL query in the `sql_query` field of the output JSON.",
					},
					"sql_template": `WITH data_cte AS (
  SELECT
    -- Select the columns you need based on the user's prompt
  FROM "catalog"."schema"."table" AS cst -- Replace with actual catalog, schema, and table
  -- Join any necessary tables based on the user's prompt
  -- Joins must also be "catalog"."schema"."table" format - replace with actual catalog, schema, and table
  -- Add any necessary WHERE clauses, ORDER BY, LIMIT, etc. based on the user's prompt
),
columns_array AS (
  SELECT ARRAY_AGG(col) AS cols FROM (
    -- Define the columns metadata based on the user's prompt
    -- The following is an example structure; modify as needed
    SELECT MAP(
      ARRAY['name', 'type', 'label'],
      ARRAY['Full Name', 'varchar', 'Full Name']
    ) AS col
    UNION ALL
    SELECT MAP(
      ARRAY['name', 'type', 'label'],
      ARRAY['Street Address', 'varchar', 'Street Address']
    )
    UNION ALL
    SELECT MAP(
      ARRAY['name', 'type', 'label'],
      ARRAY['POP', 'varchar', 'POP']
    )
  )
),
-- Build the data array
data_array AS (
  SELECT ARRAY_AGG(
    MAP(
      -- The following is an example structure; modify as needed
      ARRAY['Full Name', 'Street Address', 'POP'],
      ARRAY[
        CAST("Full Name" AS VARCHAR),
        CAST("Street Address" AS VARCHAR),
        CAST("POP" AS VARCHAR)
      ]
    )
  ) AS data
  FROM data_cte
)
SELECT CAST(
  MAP(
    ARRAY['columns', 'data'],
    ARRAY[
      CAST(cols AS JSON),
      CAST(data AS JSON)
    ]
  ) AS JSON
) AS result
FROM columns_array, data_array;`,
					"sql_rules": []string{
						"Only generate SQL queries that are compatible with TrinoDB.",
						"Ensure that you generate efficient queries.",
						"Only use tables and columns that exist in the database schema.",
						"Follow the `sql_template` as it is. But modify selected columns, joins, filters, ordering and limits based on the users prompt.",
					},
				},
				&toon.EncodeOptions{
					Indent: 1,
				},
			)

			if err != nil {
				log.Errorf("failed to encode instructions with token oriented object notation: %v", err)

				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"error":   constants.InternalServerError,
					"details": constants.InternalServerErrorDetails,
				})
			}

			c.Set("Content-Type", "text/event-stream")
			c.Set("Cache-Control", "no-cache")
			c.Set("Connection", "keep-alive")
			c.Set("Transfer-Encoding", "chunked")

			c.Status(fiber.StatusOK).Response().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				if _, _ = w.WriteString(": connected\n\n"); w.Flush() != nil {
					return
				}

				stream := r.OpenAI.Responses.NewStreaming(context.Background(), openaiResponses.ResponseNewParams{
					Model:        openai.ChatModelGPT5Nano,
					Instructions: openai.String(encodedInstructions),
					Input: openaiResponses.ResponseNewParamsInputUnion{
						OfString: openai.String(dynamicQuery.Prompt),
					},
					Text: openaiResponses.ResponseTextConfigParam{
						Format: openaiResponses.ResponseFormatTextConfigUnionParam{
							OfJSONSchema: &openaiResponses.ResponseFormatTextJSONSchemaConfigParam{
								Name:        "create_dynamic_query_output",
								Schema:      GenerateDynamicQueryOutputSchema,
								Strict:      openai.Bool(true),
								Description: openai.String("The output for create dynamic query"),
							},
						},
					},
					Tools: []openaiResponses.ToolUnionParam{
						{
							OfMcp: &openaiResponses.ToolMcpParam{
								ServerLabel:       "zingfibre_mcp",
								ServerDescription: openai.String("The ZingFibre MCP server that allows AI to interact with parts of the ZingFibre Reports Portal system."),
								ServerURL:         openai.String("https://zing-mcp.connor-davis.dev"),
								RequireApproval: openaiResponses.ToolMcpRequireApprovalUnionParam{
									OfMcpToolApprovalFilter: &openaiResponses.ToolMcpRequireApprovalMcpToolApprovalFilterParam{
										Never: openaiResponses.ToolMcpRequireApprovalMcpToolApprovalFilterNeverParam{
											ToolNames: []string{
												"list-catalogs",
												"list-schemas",
												"list-tables",
											},
										},
									},
								},
							},
						},
					},
				})

				for stream.Next() {
					current := stream.Current()

					switch current.Type {
					case "response.completed":
						{
							outputText := current.AsResponseCompleted().Response.OutputText()

							var output GenerateDynamicQueryOutput

							if err := json.Unmarshal([]byte(outputText), &output); err != nil {
								log.Infof("Error while unmarshaling output: %v. Closing http connection.\n", err)

								break
							}

							if _, err := r.Postgres.UpdateDynamicQuery(context.Background(), postgres.UpdateDynamicQueryParams{
								ID:         dynamicQuery.ID,
								Name:       dynamicQuery.Name,
								Query:      pgtype.Text{String: output.SqlQuery, Valid: true},
								ResponseID: pgtype.Text{String: current.AsResponseCompleted().Response.ID, Valid: true},
								Status:     postgres.DynamicQueryStatusComplete,
								Prompt:     dynamicQuery.Prompt,
							}); err != nil {
								log.Infof("Error while updating dynamic query: %v. Closing http connection.\n", err)

								break
							}

							fmt.Fprintf(w, "event: response_completed\n")
							fmt.Fprintf(w, "data: %s\n\n", outputText)

							err = w.Flush()

							if err != nil {
								log.Infof("Error while flushing: %v. Closing http connection.\n", err)

								continue
							}

							continue
						}
					default:
						{
							part := current.Type
							payload, err := json.Marshal(part)

							if err != nil {
								log.Infof("Error while marshaling payload: %v. Closing http connection.\n", err)

								break
							}

							fmt.Fprintf(w, "event: current_type\n")
							fmt.Fprintf(w, "data: %s\n\n", payload)

							err = w.Flush()

							if err != nil {
								log.Infof("Error while flushing: %v. Closing http connection.\n", err)

								continue
							}

							continue
						}
					}
				}

				if stream.Err() != nil {
					log.Errorf("Error during streaming: %v", stream.Err())

					err = w.Flush()

					if err != nil {
						log.Infof("Error while flushing at the end: %v.\n", err)
					}
				} else {
					log.Info("Dynamic query generation completed, closing connection.")

					fmt.Fprintf(w, "event: done\n")
					fmt.Fprintf(w, "data: \n\n")

					err = w.Flush()

					if err != nil {
						log.Infof("Error while flushing at the end: %v.\n", err)
					}
				}
			}))

			return nil
		},
	}
}
