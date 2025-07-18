// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: rm_phpsess.sql

package radius

import (
	"context"
)

const getRmPhpsess = `-- name: GetRmPhpsess :one
SELECT
    managername, ip, sessid, lastact, closed
FROM
    rm_phpsess
WHERE
    managername = ?
LIMIT
    1
`

func (q *Queries) GetRmPhpsess(ctx context.Context, managername string) (RmPhpsess, error) {
	row := q.db.QueryRowContext(ctx, getRmPhpsess, managername)
	var i RmPhpsess
	err := row.Scan(
		&i.Managername,
		&i.Ip,
		&i.Sessid,
		&i.Lastact,
		&i.Closed,
	)
	return i, err
}

const getRmPhpsesses = `-- name: GetRmPhpsesses :many
SELECT
    managername, ip, sessid, lastact, closed
FROM
    rm_phpsess
LIMIT
    ?
OFFSET
    ?
`

type GetRmPhpsessesParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetRmPhpsesses(ctx context.Context, arg GetRmPhpsessesParams) ([]RmPhpsess, error) {
	rows, err := q.db.QueryContext(ctx, getRmPhpsesses, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RmPhpsess
	for rows.Next() {
		var i RmPhpsess
		if err := rows.Scan(
			&i.Managername,
			&i.Ip,
			&i.Sessid,
			&i.Lastact,
			&i.Closed,
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
