package helpers

import (
	"fmt"
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/models/system"
)

func DynamicQueryParser(query system.DynamicQuery) string {
	sqlString := ""

	tableAliases := map[string]string{}

	// Base table alias
	tableAliases[fmt.Sprintf("%s.%s", query.Database, query.Table.Table)] = "t1"

	// Join table aliases
	for joinIdx, join := range query.Joins {
		tableAliases[fmt.Sprintf("%s.%s", join.ReferenceDatabase, join.ReferenceTable.Table)] =
			fmt.Sprintf("t%d", joinIdx+2)
	}

	// SELECT clause
	sqlString += "SELECT\n"

	columns := []string{}
	for _, column := range query.Columns {
		tableAlias := tableAliases[fmt.Sprintf("%s.%s", column.Database, column.Table.Table)]

		if column.Aggregate != nil {
			columns = append(columns,
				fmt.Sprintf("\t%s(%s.%s) AS `%s`", *column.Aggregate, tableAlias, column.Column, column.Label))
		} else {
			columns = append(columns,
				fmt.Sprintf("\t%s.%s AS `%s`", tableAlias, column.Column, column.Label))
		}
	}

	sqlString += fmt.Sprintf("%s\n", strings.Join(columns, ",\n"))

	// FROM clause
	sqlString += fmt.Sprintf("FROM %s.%s AS t1", query.Database, query.Table.Table)

	// JOINs
	for _, join := range query.Joins {
		refKey := fmt.Sprintf("%s.%s", join.ReferenceDatabase, join.ReferenceTable.Table)
		tableAlias := tableAliases[refKey]

		if join.SubQuery != nil {
			subQuerySQL := DynamicQueryParser(*join.SubQuery)
			subQuerySQL = strings.ReplaceAll(subQuerySQL, "\n", "\n\t")

			sqlString += fmt.Sprintf("\n%s JOIN (\n\t%s\n) AS %s ON t1.%s = %s.%s",
				join.Type, subQuerySQL, tableAlias,
				join.LocalColumn, tableAlias, join.ReferenceColumn)
			continue
		}

		switch join.Type {
		case system.InnerJoin:
			sqlString += fmt.Sprintf("\nINNER JOIN %s.%s AS %s ON t1.%s = %s.%s",
				join.ReferenceDatabase, join.ReferenceTable.Table, tableAlias,
				join.LocalColumn, tableAlias, join.ReferenceColumn)
		case system.LeftJoin:
			sqlString += fmt.Sprintf("\nLEFT JOIN %s.%s AS %s ON t1.%s = %s.%s",
				join.ReferenceDatabase, join.ReferenceTable.Table, tableAlias,
				join.LocalColumn, tableAlias, join.ReferenceColumn)
		case system.RightJoin:
			sqlString += fmt.Sprintf("\nRIGHT JOIN %s.%s AS %s ON t1.%s = %s.%s",
				join.ReferenceDatabase, join.ReferenceTable.Table, tableAlias,
				join.LocalColumn, tableAlias, join.ReferenceColumn)
		case system.OuterJoin:
			sqlString += fmt.Sprintf("\nOUTER JOIN %s.%s AS %s ON t1.%s = %s.%s",
				join.ReferenceDatabase, join.ReferenceTable.Table, tableAlias,
				join.LocalColumn, tableAlias, join.ReferenceColumn)
		}
	}

	return sqlString
}
