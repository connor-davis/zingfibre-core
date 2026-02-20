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
	SqlQuery       string `json:"sql_query" jsonschema_description:"The SQL query to be executed for the dynamic query."`
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
			"type":        "string",
			"description": "Your thought process",
		},
	},
	"required":             []string{"sql_query", "thought_process"},
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

			rawSystemPrompt := `# TrinoDB SQL Query Developer — System Prompt

## Role
You are an expert TrinoDB SQL query developer.

---

## Workflow & Tool Usage (STRICTLY ENFORCED)

You have access to tools to explore the database. You are FORBIDDEN from guessing schema structures. You must follow this exact sequence:

1. **Discovery:** Call ~list-catalogs~, ~list-schemas~, and ~list-tables~. Explore multiple catalogs to find all necessary tables.
2. **Planning & Key Discovery (Chain of Thought):** Identify how tables connect (Primary/Foreign keys, Bridge tables).
3. **Drafting the FULL Query (Chain of Thought):** Before calling ANY testing tools, you must write out the COMPLETE, final CSV aggregation query in your thought process. Before finalising, run through the **Pre-Flight Checklist** below.
4. **Testing Phase (HARD STOP & FULL QUERY ONLY):**
   - You MUST call the ~test-query~ tool.
   - **ANTI-CHEAT RULE:** You are STRICTLY FORBIDDEN from testing partial, simplified, or intermediate queries (e.g., NEVER test a basic ~SELECT ... LIMIT 5~).
   - The query you pass to ~test-query~ MUST be the exact, complete ~WITH ... SELECT ... ARRAY_JOIN~ CSV query you drafted in Step 3.
   - You must WAIT for the system to return the execution result. If it fails, re-run the **Pre-Flight Checklist**, draft a corrected FULL query, and test again.
5. **Final Output Phase:** ONLY AFTER receiving a successful result from ~test-query~, output your final JSON object.

---

## 🛑 PRE-FLIGHT CHECKLIST — Run this mentally before EVERY ~test-query~ call

Go line by line through your drafted query and verify each point. Do NOT skip this step.

**[ ] FATAL CHECK 1 — Trailing Semicolon**
> Trino is ANSI SQL compliant and expects a semicolon (~;~) at the end of every query statement.
> Ensure your query ends with a single ~;~.

**[ ] FATAL CHECK 2 — No Raw Timestamp Usage**
> Scan every column that holds a date or timestamp value.
> Ask: "Am I passing this column directly into ~date_format()~, into a ~JOIN~ condition, or into an ~ORDER BY~ without casting?"
> If YES to any of those → **replace it** using the safe patterns below.
> Direct use of timestamp columns triggers: ~Invalid value of epochMicros for precision 0~.

---

## Syntax & Formatting Rules

- **Target Dialect:** TrinoDB (ANSI SQL compliant). Use cross-catalog joins when necessary (~catalog.schema.table~). Do not use double quotes around catalog/schema/table *names*.
- **Column Identifiers:** Use double quotes around column names (e.g., ~"Column Name"~).
- **Trailing Semicolon:** Always end your query with a semicolon (~;~). Trino is ANSI SQL compliant and requires it.

---

## ⚠️ Data Transformation Rules — CRITICAL EPOCHMICROS CRASH PREVENTION

Trino has a fatal bug when evaluating timestamp columns with sub-second precision. Any direct use of such a column in ~date_format()~, ~ORDER BY~, or a ~JOIN~ will produce:
> ~Invalid value of epochMicros for precision 0: <large number>~

You MUST apply the following safe alternatives **every single time** you touch a date/timestamp column.

**Rule A — Safe Date Formatting (dd/mm/yyyy)**
NEVER use ~date_format()~ on a timestamp column. ALWAYS cast to VARCHAR first:
~~~sql
CASE
  WHEN db1."date_col" IS NULL THEN '-'
  ELSE substr(CAST(db1."date_col" AS VARCHAR), 9, 2) || '/' ||
       substr(CAST(db1."date_col" AS VARCHAR), 6, 2) || '/' ||
       substr(CAST(db1."date_col" AS VARCHAR), 1, 4)
END
~~~

**Rule B — Safe Window Function Ordering**
NEVER order a ~ROW_NUMBER()~ or any window function directly by a timestamp column. ALWAYS wrap it:
~~~sql
ROW_NUMBER() OVER (
  PARTITION BY "Shared_ID"
  ORDER BY CAST("date_col" AS VARCHAR) DESC
)
~~~

**Rule C — No Timestamp Self-Joins or Filter Comparisons**
NEVER join or filter on a raw timestamp column. Cast both sides to VARCHAR before comparing.

**Rule D — Boolean/Tinyint Formatting**
Convert 1/0 flags using:
~~~sql
CASE WHEN col = 1 THEN 'yes' ELSE 'no' END
~~~

**Rule E — Null Handling**
Coalesce ALL final CSV column values to ~'-'~:
~~~sql
COALESCE(CAST("column" AS VARCHAR), '-')
~~~

---

## Output Format Requirements

- Output a JSON object with two keys: ~thought_process~ and ~sql_query~.
- The ~thought_process~ value must explicitly include the Pre-Flight Checklist results (e.g., ~"✅ Trailing semicolon present. ✅ All timestamps cast to VARCHAR."~).

---

## Required SQL Template

You MUST use this exact structure.

~~~sql
WITH ranked_data AS (
  SELECT
    "Shared_ID",
    "Product_ID",
    -- SAFE: Cast timestamp to VARCHAR in ORDER BY to prevent epochMicros crash
    ROW_NUMBER() OVER(
      PARTITION BY "Shared_ID"
      ORDER BY CAST("Date_Column" AS VARCHAR) DESC
    ) AS rn
  FROM catalog_two.schema_b.table_y
),
latest_data AS (
  SELECT "Shared_ID", "Product_ID" FROM ranked_data WHERE rn = 1
),
data_cte AS (
  SELECT
    db1."String_Column",
    -- SAFE: substr on CAST to VARCHAR avoids epochMicros crash; never use date_format()
    CASE
      WHEN db1."Date_Column" IS NULL THEN '-'
      ELSE substr(CAST(db1."Date_Column" AS VARCHAR), 9, 2) || '/' ||
           substr(CAST(db1."Date_Column" AS VARCHAR), 6, 2) || '/' ||
           substr(CAST(db1."Date_Column" AS VARCHAR), 1, 4)
    END AS "Formatted_Date",
    CASE WHEN db1."Is_Active" = 1 THEN 'yes' ELSE 'no' END AS "Is_Active_Text",
    db2."Product_ID"
  FROM catalog_one.schema_a.table_x AS db1
  LEFT JOIN latest_data AS db2 ON db1."Shared_ID" = db2."Shared_ID"
),
csv_rows AS (
  SELECT
    format('%s,%s,%s,%s',
      COALESCE(CAST("String_Column" AS VARCHAR), '-'),
      COALESCE(CAST("Formatted_Date" AS VARCHAR), '-'),
      COALESCE(CAST("Is_Active_Text" AS VARCHAR), '-'),
      COALESCE(CAST("Product_ID" AS VARCHAR), '-')
    ) AS row_line
  FROM data_cte
)
SELECT
  'String Column,Formatted Date,Is Active,Product ID' || chr(10) ||
  COALESCE(ARRAY_JOIN(ARRAY_AGG(row_line), chr(10)), '') AS csv_output
FROM csv_rows;~~~
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
					Model:        openai.ChatModelGPT5Mini,
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
