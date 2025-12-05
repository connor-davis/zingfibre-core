package trino

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListTablesParams struct {
	Catalog string `json:"catalog"`
	Schema  string `json:"schema"`
}

func (t *trino) ListTables(ctx context.Context, request *mcp.CallToolRequest, params ListTablesParams) (*mcp.CallToolResult, any, error) {
	row := t.db.QueryRow(fmt.Sprintf(`WITH params AS (
  SELECT '%s' AS schema_name, '%s' AS catalog_name
),
cols AS (
  SELECT
    table_schema,
    table_name,
    ordinal_position,
    "column_name" || ' ' ||
    "data_type" AS column_ddl
  FROM %s.information_schema.columns
  JOIN params p ON columns.table_schema = p.schema_name
),
tables_ddl AS (
  SELECT
    t.table_schema,
    t.table_name,
    CONCAT(
      'CREATE TABLE ',
      p.catalog_name, '.', t.table_schema, '.', t.table_name,
      ' (', CHR(10), '  ',
      ARRAY_JOIN(ARRAY_AGG(c.column_ddl ORDER BY c.ordinal_position), CONCAT(',', CHR(10), '  ')),
      CHR(10), ')'
    ) AS create_statement
  FROM %s.information_schema.tables t
  JOIN params p ON t.table_schema = p.schema_name
  LEFT JOIN cols c
    ON c.table_schema = t.table_schema
    AND c.table_name = t.table_name
  GROUP BY p.catalog_name, t.table_schema, t.table_name
)
SELECT ARRAY_JOIN(ARRAY_AGG(create_statement), CONCAT(CHR(10), CHR(10))) AS full_schema_sql
FROM tables_ddl`, params.Schema, params.Catalog, params.Catalog, params.Catalog))

	var fullSchemaSQL string

	if err := row.Scan(&fullSchemaSQL); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error retrieving tables: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fullSchemaSQL,
			},
		},
	}, nil, nil
}
