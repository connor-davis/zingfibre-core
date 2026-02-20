package trino

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListSchemasParams struct {
	Catelog string `json:"catalog"`
}

func (t *trino) ListSchemas(ctx context.Context, request *mcp.CallToolRequest, params ListSchemasParams) (*mcp.CallToolResult, any, error) {
	log.Info("Listing schemas...")

	row := t.db.QueryRow(fmt.Sprintf(`SELECT
    ARRAY_JOIN(
        ARRAY_AGG(schema_name),
        ', '
    ) AS schema_list
FROM
    %s.information_schema.schemata
WHERE
    schema_name NOT IN ('information_schema', 'system', 'pg_catalog')`, params.Catelog))

	var schemaList string

	if err := row.Scan(&schemaList); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving schemas: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: schemaList,
			},
		},
	}, nil, nil
}
