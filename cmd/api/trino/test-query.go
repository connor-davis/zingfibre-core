package trino

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TestQueryParams struct {
	Query string `json:"query"`
}

func (t *trino) TestQuery(context context.Context, request *mcp.CallToolRequest, params TestQueryParams) (*mcp.CallToolResult, any, error) {
	log.Info("Testing query...")

	row := t.db.QueryRow(params.Query)

	log.Infof("Query being tested:\n%s", params.Query)

	var result string

	if err := row.Scan(&result); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("The query failed to execute: %s", err.Error()),
				},
			},
		}, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: "The query executed successfully.",
			},
		},
	}, nil, nil
}
