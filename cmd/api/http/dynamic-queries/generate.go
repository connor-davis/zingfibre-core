package dynamicQueries

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/connor-davis/zingfibre-core/common"
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
	"github.com/openai/openai-go/v3"
	openaiResponses "github.com/openai/openai-go/v3/responses"
	"github.com/valyala/fasthttp"
)

type GenerateDynamicQueryOutput struct {
	SqlQuery string `json:"sql_query" jsonschema_description:"The SQL query to be executed for the dynamic query."`
 ThoughtProcess string `json:"thought_process" jsonschema_description:"Your thought process"`
}

var GenerateDynamicQueryOutputSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"sql_query": map[string]any{
			"type":                   "string",
			"description":            "The SQL query to be executed for the dynamic query.",
			"jsonschema_description": "The SQL query to be executed for the dynamic query.",
		},
  "thought_process": map[string]any{
    "type": "string",
    "description": "Your thought process",
  },
	},
	"required":             []string{"sql_query","thought_process"},
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

			rawSystemPrompt := `**Role**
You are an expert TrinoDB SQL query developer.

**Workflow & Tool Usage**
You have access to tools to explore the database. You must follow this exact sequence before providing your final answer:
1. **Discovery:**
   * Call ~list-catalogs~ first to view connected databases.
   * Call ~list-schemas~ second to view schemas within the relevant catalogs.
   * Call ~list-tables~ third to view tables within the relevant schemas.
   * *Note: If the user's request involves correlating data from different systems, be sure to explore multiple catalogs to find the necessary tables.*
2. **Planning (Chain of Thought):**
   * Before writing any SQL, you must internally map every single data point requested by the user to its specific ~catalog.schema.table~.
   * Identify the explicit JOIN conditions needed to connect all these tables. If a requested table seems disconnected, investigate further using your tools until you find the linking keys.
3. **Testing:**
   * Write your query and test it using the ~test-query~ tool.
   * Iterate and fix any errors until ~test-query~ returns a successful result. You must have no errors before finalizing.
   * Do not modify the SQL query at all after a successful test.

**Syntax & Formatting Rules**
* **Target Dialect:** Only generate SQL queries compatible with TrinoDB.
* **Federated Queries:** You are uniquely capable of querying across multiple databases. Utilize cross-catalog joins when necessary by referencing the distinct catalogs in your FROM and JOIN clauses.
* **Table References:** Format table names as ~catalog.schema.table~ in your queries. **Do not** use double quotes around catalogs, schemas, or tables in the FROM clauses.
* **Column Identifiers:** You must use **double quotes** to quote column identifiers (e.g., ~"Column Name"~).
* **Null Handling:** If a column contains a null value, you must represent it as a hyphen (~-~) in the CSV output.
* **Strict Adherence:** Only utilize local information stores, tables, and column definitions you have discovered via your tools. **Do not invent or hallucinate non-existent columns or tables.**
* **Statefulness:** If you are modifying a previous response based on user feedback, only make the specific changes the user requested.

**Output Format Requirements**
* Your final response must be formatted as a valid JSON object containing exactly two keys: ~thought_process~ and ~sql_query~.
* **~thought_process~:** A brief string where you list every requested field, the exact ~catalog.schema.table~ it comes from, and how you will join the tables together.
* **~sql_query~:** Your final, successfully tested SQL query string.

**Required SQL Template**
You must format your SQL query to aggregate the result into a single string blob containing the CSV header and data rows separated by newlines. Follow this exact template structure, utilizing cross-catalog joins if the data resides in separate databases:

~~~sql
WITH data_cte AS (
  SELECT
    -- Select the columns you need based on the user's prompt
    db1."Column 1",
    db2."Column 2"
  FROM catalog_one.schema_a.table_x AS db1 -- First database source (unquoted)
  JOIN catalog_two.schema_b.table_y AS db2 -- Second database source (unquoted)
    ON db1."Shared_ID" = db2."Shared_ID"
  -- Add any necessary WHERE clauses, ORDER BY, LIMIT, etc.
),
csv_rows AS (
  SELECT
    -- Format the rows as CSV strings.
    -- IMPORTANT: Cast all columns to VARCHAR and coalesce nulls to '-'.
    format('%s,%s',
      COALESCE(CAST("Column 1" AS VARCHAR), '-'),
      COALESCE(CAST("Column 2" AS VARCHAR), '-')
    ) AS row_line
  FROM data_cte
)
SELECT
  -- 1. Manually write the Header string based on selected columns
  'Column 1,Column 2' || chr(10) ||
  -- 2. Aggregate the data rows with a newline character
  ARRAY_JOIN(ARRAY_AGG(row_line), chr(10)) AS csv_output
FROM csv_rows;
~~~
`

			systemPrompt := strings.ReplaceAll(rawSystemPrompt, "~", "`")

			c.Set("Content-Type", "text/event-stream")
			c.Set("Cache-Control", "no-cache")
			c.Set("Connection", "keep-alive")
			c.Set("Transfer-Encoding", "chunked")

			c.Status(fiber.StatusOK).Response().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				if _, _ = w.WriteString(": connected\n\n"); w.Flush() != nil {
					return
				}

				streamParams := openaiResponses.ResponseNewParams{
					Model:        openai.ChatModelGPT5Nano,
					Instructions: openai.String(systemPrompt),
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
								ServerURL:         openai.String(common.EnvString("MCP_BASE_URL", "http://localhost:6173/api/mcp")),
								RequireApproval: openaiResponses.ToolMcpRequireApprovalUnionParam{
									OfMcpToolApprovalFilter: &openaiResponses.ToolMcpRequireApprovalMcpToolApprovalFilterParam{
										Never: openaiResponses.ToolMcpRequireApprovalMcpToolApprovalFilterNeverParam{
											ToolNames: []string{
												"list-catalogs",
												"list-schemas",
												"list-tables",
												"test-query",
											},
										},
									},
								},
							},
						},
					},
				}

				if dynamicQuery.ResponseID.Valid {
					streamParams.PreviousResponseID = openai.String(dynamicQuery.ResponseID.String)
				}

				stream := r.OpenAI.Responses.NewStreaming(context.Background(), streamParams)

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
