// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: sales_agents.sql

package zing

import (
	"context"
)

const getSalesAgent = `-- name: GetSalesAgent :one
SELECT
    id, name, code
FROM
    SalesAgents
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetSalesAgent(ctx context.Context, id string) (Salesagent, error) {
	row := q.db.QueryRowContext(ctx, getSalesAgent, id)
	var i Salesagent
	err := row.Scan(&i.ID, &i.Name, &i.Code)
	return i, err
}

const getSalesAgents = `-- name: GetSalesAgents :many
SELECT
    id, name, code
FROM
    SalesAgents
LIMIT
    ?
OFFSET
    ?
`

type GetSalesAgentsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetSalesAgents(ctx context.Context, arg GetSalesAgentsParams) ([]Salesagent, error) {
	rows, err := q.db.QueryContext(ctx, getSalesAgents, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Salesagent
	for rows.Next() {
		var i Salesagent
		if err := rows.Scan(&i.ID, &i.Name, &i.Code); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
