package helpers

import (
	"fmt"
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/models/system"
)

func DynamicQueryParser(query system.DynamicQuery) string {
	sqlString := ""

	for subQueryIdx, subQuery := range query.SubQueries {
		if subQueryIdx == 0 {
			sqlString += "WITH (\n"
		}

		sqlString += DynamicQueryParser(subQuery)

		if subQueryIdx == len(query.SubQueries)-1 {
			sqlString += fmt.Sprintf("\n) AS %s_%s", subQuery.Database, subQuery.Table.Table)
		}
	}

	sqlString += "SELECT\n"

	columns := []string{}

	for _, column := range query.Columns {
		if column.Aggregate != nil {
			columns = append(columns, fmt.Sprintf("  %s(%s) AS %s", *column.Aggregate, column.Column, column.Label))
		}

		if column.Aggregate == nil {
			columns = append(columns, fmt.Sprintf("  %s AS %s", column.Column, column.Label))
		}
	}

	sqlString += fmt.Sprintf("%s\n", strings.Join(columns, ",\n"))

	sqlString += fmt.Sprintf("FROM %s.%s AS %s_%s", query.Database, query.Table.Table, query.Database, query.Table.Table)

	return sqlString
}
