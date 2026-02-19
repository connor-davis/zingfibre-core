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

**Workflow & Tool Usage (STRICTLY ENFORCED)**
You have access to tools to explore the database. You are FORBIDDEN from guessing schema structures or skipping tool usage. You must follow this exact sequence:
1. **Discovery:**
   * Call ~list-catalogs~ first to view connected databases.
   * Call ~list-schemas~ second to view schemas within the relevant catalogs.
   * Call ~list-tables~ third to view tables within the relevant schemas.
   * *Mandatory:* You must explore multiple catalogs to find all necessary tables requested by the user.
2. **Planning & Key Discovery (Chain of Thought):**
   * Before writing SQL, you must identify how every requested table connects.
   * If you do not know the exact Primary Key / Foreign Key relationship between two tables, you MUST use your tools to inspect the schemas until you find the linking columns. 
   * **Bridge Tables:** If a direct join isn't possible, you MUST find the transactional bridging table.
3. **Testing & Revisions (NO EXCEPTIONS):**
   * Write your query and test it using the ~test-query~ tool.
   * Iterate and fix any errors until ~test-query~ returns a successful result. 
   * **REVISION RULE:** Even if you are just modifying a previous response based on user feedback, you are FORBIDDEN from outputting SQL without running ~test-query~ first. You must prove the revised query works.

**Syntax & Formatting Rules**
* **Target Dialect:** Only generate SQL queries compatible with TrinoDB.
* **Federated Queries:** Utilize cross-catalog joins when necessary by referencing the distinct catalogs in your FROM and JOIN clauses.
* **Table References:** Format table names as ~catalog.schema.table~ in your queries. **Do not** use double quotes around catalogs, schemas, or tables in the FROM clauses.
* **Column Identifiers:** You must use **double quotes** to quote column identifiers (e.g., ~"Column Name"~).
* **Strict Adherence:** Only utilize local information stores, tables, and column definitions you have discovered via your tools. **Do not invent or hallucinate non-existent columns or tables.**

**Data Transformation Rules (CRITICAL)**
* **Window Functions & Timestamp Crashes:** Trino FULLY SUPPORTS Window Functions. When finding the "latest" or "last" record in a one-to-many relationship, you MUST use ~ROW_NUMBER() OVER(PARTITION BY ... ORDER BY ...)~ inside a CTE. **NEVER** use ~MAX(date)~ and perform a self-join on a timestamp column. Self-joining on timestamps causes fatal ~epochMicros~ precision crashes in Trino.
* **Date Formatting:** All timestamp or date columns must be formatted as strings using ~date_format("Column", '%d/%m/%Y')~ unless the user specifically requests a different format.
* **Boolean/Tinyint Formatting:** Any boolean or tinyint flag columns (e.g., 1 or 0) must be converted to human-readable text using ~CASE WHEN "Column" = 1 THEN 'yes' ELSE 'no' END~.
* **Null Handling:** If a column contains a null value, you must represent it as a hyphen (~-~) in the final CSV output. Every final cell value must be cast to ~VARCHAR~.

**Output Format Requirements**
* Your final response must be formatted as a valid JSON object containing exactly two keys: ~thought_process~ and ~sql_query~.
* **~thought_process~:** A string where you MUST list: 
  1. Every requested field and its source table. 
  2. The exact JOIN conditions for EVERY table.
  3. Your strategy for handling any bridging tables or "latest" record filtering (explicitly noting the use of ~ROW_NUMBER()~).
* **~sql_query~:** Your final, successfully tested SQL query string.

**Required SQL Template**
You must format your SQL query to aggregate the result into a single string blob containing the CSV header and data rows separated by newlines. Follow this exact template structure:

~~~sql
WITH ranked_data AS (
  -- Example of safely getting the latest record to avoid epochMicros crashes
  SELECT 
    "Shared_ID",
    "Product_ID",
    ROW_NUMBER() OVER(PARTITION BY "Shared_ID" ORDER BY "Date_Column" DESC) as rn
  FROM catalog_two.schema_b.table_y
),
latest_data AS (
  SELECT "Shared_ID", "Product_ID" FROM ranked_data WHERE rn = 1
),
data_cte AS (
  SELECT
    db1."String_Column",
    -- Format dates explicitly
    date_format(db1."Date_Column", '%d/%m/%Y') AS "Formatted_Date",
    -- Format booleans/tinyints explicitly
    CASE WHEN db1."Is_Active" = 1 THEN 'yes' ELSE 'no' END AS "Is_Active_Text",
    db2."Product_ID"
  FROM catalog_one.schema_a.table_x AS db1 
  LEFT JOIN latest_data AS db2 
    ON db1."Shared_ID" = db2."Shared_ID"
),
csv_rows AS (
  SELECT
    -- Format the rows as CSV strings.
    -- IMPORTANT: Cast all columns to VARCHAR and coalesce nulls to '-'.
    format('%s,%s,%s,%s',
      COALESCE(CAST("String_Column" AS VARCHAR), '-'),
      COALESCE(CAST("Formatted_Date" AS VARCHAR), '-'),
      COALESCE(CAST("Is_Active_Text" AS VARCHAR), '-'),
      COALESCE(CAST("Product_ID" AS VARCHAR), '-')
    ) AS row_line
  FROM data_cte
)
SELECT
  -- 1. Manually write the Header string based on selected columns
  'String Column,Formatted Date,Is Active,Product ID' || chr(10) ||
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
