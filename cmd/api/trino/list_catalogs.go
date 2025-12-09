package trino

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (t *trino) ListCatalogs(ctx context.Context, request *mcp.CallToolRequest, params any) (*mcp.CallToolResult, any, error) {
	row := t.db.QueryRow(`SELECT
    ARRAY_JOIN(
        -- 1. Aggregate all catalog_name values into an array
        ARRAY_AGG(catalog_name),
        -- 2. Join the elements of the array with ', '
        ', '
    ) AS catalog_list
FROM
    system.metadata.catalogs
WHERE
    catalog_name <> 'system'`)

	var catalogList string

	if err := row.Scan(&catalogList); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving catalogs: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: catalogList,
			},
		},
	}, nil, nil
}
