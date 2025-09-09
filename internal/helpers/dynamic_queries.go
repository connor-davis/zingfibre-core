package helpers

import (
	"fmt"
	"strings"

	"github.com/connor-davis/zingfibre-core/internal/models/system"
)

func DynamicQueryParser(query system.DynamicQuery) string {
	sqlString := ""

	tableAliases := map[string]string{}

	tableAliases[query.Table.Table] = "t1"

	for joinIdx, join := range query.Joins {
		tableAliases[join.ReferenceTable.Table] = fmt.Sprintf(
			"t%d",
			joinIdx+2,
		)
	}

	sqlString += "SELECT\n"

	columns := []string{}

	for _, column := range query.Columns {
		tableAlias := tableAliases[column.Table.Table]

		if query.IsSubQuery {
			if column.Aggregate != nil {
				columns = append(
					columns,
					fmt.Sprintf(
						"\t%s(%s.%s)",
						*column.Aggregate,
						tableAlias,
						column.Column,
					),
				)
			} else {
				columns = append(
					columns,
					fmt.Sprintf(
						"\t%s.%s",
						tableAlias,
						column.Column,
					),
				)
			}
		} else {
			if column.Aggregate != nil {
				columns = append(
					columns,
					fmt.Sprintf(
						"\t%s(%s.%s) AS `%s`",
						*column.Aggregate,
						tableAlias,
						column.Column,
						column.Label,
					),
				)
			} else {
				columns = append(
					columns,
					fmt.Sprintf(
						"\t%s.%s AS `%s`",
						tableAlias,
						column.Column,
						column.Label,
					),
				)
			}
		}
	}

	sqlString += fmt.Sprintf(
		"%s\n",
		strings.Join(columns, ",\n"),
	)

	sqlString += fmt.Sprintf(
		"FROM %s AS t1",
		query.Table.Table,
	)

	for _, join := range query.Joins {
		tableAlias := tableAliases[join.ReferenceTable.Table]

		if join.SubQuery != nil {
			subQuerySQL := DynamicQueryParser(*join.SubQuery)
			subQuerySQL = strings.ReplaceAll(
				subQuerySQL,
				"\n",
				"\n\t",
			)

			sqlString += fmt.Sprintf(
				"\n%s JOIN (\n\t%s\n) AS %s ON t1.%s = %s.%s",
				join.Type,
				subQuerySQL,
				tableAlias,
				join.LocalColumn,
				tableAlias,
				join.ReferenceColumn,
			)

			continue
		}

		switch join.Type {
		case system.InnerJoin:
			sqlString += fmt.Sprintf(
				"\nINNER JOIN %s AS %s ON t1.%s = %s.%s",
				join.ReferenceTable.Table,
				tableAlias,
				join.LocalColumn,
				tableAlias,
				join.ReferenceColumn,
			)
		case system.LeftJoin:
			sqlString += fmt.Sprintf(
				"\nLEFT JOIN %s AS %s ON t1.%s = %s.%s",
				join.ReferenceTable.Table,
				tableAlias,
				join.LocalColumn,
				tableAlias,
				join.ReferenceColumn,
			)
		case system.RightJoin:
			sqlString += fmt.Sprintf(
				"\nRIGHT JOIN %s AS %s ON t1.%s = %s.%s",
				join.ReferenceTable.Table,
				tableAlias,
				join.LocalColumn,
				tableAlias,
				join.ReferenceColumn,
			)
		case system.OuterJoin:
			sqlString += fmt.Sprintf(
				"\nOUTER JOIN %s AS %s ON t1.%s = %s.%s",
				join.ReferenceTable.Table,
				tableAlias,
				join.LocalColumn,
				tableAlias,
				join.ReferenceColumn,
			)
		}
	}

	if len(query.Filters) > 0 {
		sqlString += "\nWHERE"

		for _, filter := range query.Filters {
			tableAlias := tableAliases[filter.Table.Table]

			switch filter.Type {
			case system.StringFilter:
				sqlString += fmt.Sprintf(
					"\n\t%s.%s %s '%s'",
					tableAlias,
					filter.Column,
					filter.Operator,
					filter.Value,
				)
			case system.DateFilter:
				sqlString += fmt.Sprintf(
					"\n\t%s.%s %s '%s'",
					tableAlias,
					filter.Column,
					filter.Operator,
					filter.Value,
				)
			case system.NumberFilter:
				sqlString += fmt.Sprintf(
					"\n\t%s.%s %s %s",
					tableAlias,
					filter.Column,
					filter.Operator,
					filter.Value,
				)
			case system.BooleanFilter:
				sqlString += fmt.Sprintf(
					"\n\t%s.%s %s %s",
					tableAlias,
					filter.Column,
					filter.Operator,
					filter.Value,
				)
			}
		}
	}

	if len(query.Orders) > 0 {
		sqlString += "\nORDER BY"

		orderColumns := []string{}

		for _, order := range query.Orders {
			tableAlias := tableAliases[order.Table.Table]

			if order.Descending {
				orderColumns = append(orderColumns, fmt.Sprintf("\n\t%s.%s %s", tableAlias, order.Column, "DESC"))
			} else {
				orderColumns = append(orderColumns, fmt.Sprintf("\n\t%s.%s %s", tableAlias, order.Column, "ASC"))
			}
		}

		sqlString += strings.Join(orderColumns, ",")
	}

	return sqlString
}
