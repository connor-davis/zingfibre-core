package ai

import (
	"context"

	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mateuszkardas/toon-go"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type CreateDynamicQueryOutput struct {
	SqlQuery string `json:"sql_query" jsonschema_description:"The SQL query to be executed for the dynamic query."`
}

var CreateDynamicQueryOutputSchema = map[string]any{
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

func (ai *ai) CreateDynamicQuery(queryId uuid.UUID, prompt string) {
	log.Infof("Creating dynamic query for Query ID: %s with prompt: %s", queryId.String(), prompt)

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

		return
	}

	log.Infof("Creating dynamic query response...")

	response, err := ai.openai.Responses.New(context.Background(), responses.ResponseNewParams{
		Model:        openai.ChatModelGPT5Nano,
		Instructions: openai.String(encodedInstructions),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "create_dynamic_query_output",
					Schema:      CreateDynamicQueryOutputSchema,
					Strict:      openai.Bool(true),
					Description: openai.String("The output for create dynamic query"),
				},
			},
		},
		Tools: []responses.ToolUnionParam{
			{
				OfMcp: &responses.ToolMcpParam{
					ServerLabel:       "zingfibre_mcp",
					ServerDescription: openai.String("The ZingFibre MCP server that allows AI to interact with parts of the ZingFibre Reports Portal system."),
					ServerURL:         openai.String("https://zing-mcp.connor-davis.dev"),
					RequireApproval: responses.ToolMcpRequireApprovalUnionParam{
						OfMcpToolApprovalFilter: &responses.ToolMcpRequireApprovalMcpToolApprovalFilterParam{
							Never: responses.ToolMcpRequireApprovalMcpToolApprovalFilterNeverParam{
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

	if err != nil {
		log.Errorf("failed to create dynamic query response: %v", err)

		return
	}

	log.Infof("Dynamic query response created with ID: %s", response.ID)

	var CreateDynamicQueryOutput CreateDynamicQueryOutput

	if err := json.Unmarshal([]byte(response.OutputText()), &CreateDynamicQueryOutput); err != nil {
		log.Errorf("failed to unmarshal create dynamic query output: %v", err)

		return
	}

	dynamicQuery, err := ai.postgres.GetDynamicQuery(context.Background(), queryId)

	if err != nil {
		log.Errorf("failed to get dynamic query from database: %v", err)

		return
	}

	if _, err := ai.postgres.UpdateDynamicQuery(context.Background(), postgres.UpdateDynamicQueryParams{
		ID:         dynamicQuery.ID,
		Name:       dynamicQuery.Name,
		Query:      pgtype.Text{String: CreateDynamicQueryOutput.SqlQuery, Valid: true},
		ResponseID: pgtype.Text{String: response.ID, Valid: true},
		Status:     postgres.DynamicQueryStatusComplete,
	}); err != nil {
		log.Errorf("failed to update dynamic query in database: %v", err)

		return
	}
}
