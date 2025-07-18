// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: radgroupcheck.sql

package radius

import (
	"context"
)

const getRadgroupcheck = `-- name: GetRadgroupcheck :one
SELECT
    id, groupname, attribute, op, value
FROM
    radgroupcheck
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetRadgroupcheck(ctx context.Context, id uint32) (Radgroupcheck, error) {
	row := q.db.QueryRowContext(ctx, getRadgroupcheck, id)
	var i Radgroupcheck
	err := row.Scan(
		&i.ID,
		&i.Groupname,
		&i.Attribute,
		&i.Op,
		&i.Value,
	)
	return i, err
}

const getRadgroupchecks = `-- name: GetRadgroupchecks :many
SELECT
    id, groupname, attribute, op, value
FROM
    radgroupcheck
LIMIT
    ?
OFFSET
    ?
`

type GetRadgroupchecksParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetRadgroupchecks(ctx context.Context, arg GetRadgroupchecksParams) ([]Radgroupcheck, error) {
	rows, err := q.db.QueryContext(ctx, getRadgroupchecks, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Radgroupcheck
	for rows.Next() {
		var i Radgroupcheck
		if err := rows.Scan(
			&i.ID,
			&i.Groupname,
			&i.Attribute,
			&i.Op,
			&i.Value,
		); err != nil {
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
