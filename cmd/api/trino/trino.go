package trino

import (
	"context"
	"database/sql"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Trino interface {
	ListCatalogs(context context.Context, request *mcp.CallToolRequest, params any) (*mcp.CallToolResult, any, error)
	ListSchemas(context context.Context, request *mcp.CallToolRequest, params ListSchemasParams) (*mcp.CallToolResult, any, error)
	ListTables(context context.Context, request *mcp.CallToolRequest, params ListTablesParams) (*mcp.CallToolResult, any, error)
	TestQuery(context context.Context, request *mcp.CallToolRequest, params TestQueryParams) (*mcp.CallToolResult, any, error)
}

type trino struct {
	db *sql.DB
}

func New(db *sql.DB) Trino {
	return &trino{
		db: db,
	}
}
