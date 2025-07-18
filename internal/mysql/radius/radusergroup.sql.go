// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: radusergroup.sql

package radius

import (
	"context"
)

const getRadusergroup = `-- name: GetRadusergroup :one
SELECT
    username, groupname, priority
FROM
    radusergroup
WHERE
    username = ?
LIMIT
    1
`

func (q *Queries) GetRadusergroup(ctx context.Context, username string) (Radusergroup, error) {
	row := q.db.QueryRowContext(ctx, getRadusergroup, username)
	var i Radusergroup
	err := row.Scan(&i.Username, &i.Groupname, &i.Priority)
	return i, err
}

const getRadusergroups = `-- name: GetRadusergroups :many
SELECT
    username, groupname, priority
FROM
    radusergroup
LIMIT
    ?
OFFSET
    ?
`

type GetRadusergroupsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetRadusergroups(ctx context.Context, arg GetRadusergroupsParams) ([]Radusergroup, error) {
	rows, err := q.db.QueryContext(ctx, getRadusergroups, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Radusergroup
	for rows.Next() {
		var i Radusergroup
		if err := rows.Scan(&i.Username, &i.Groupname, &i.Priority); err != nil {
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
