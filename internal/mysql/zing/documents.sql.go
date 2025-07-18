// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: documents.sql

package zing

import (
	"context"
)

const getDocument = `-- name: GetDocument :one
SELECT
    id, base64data, filename, datecreated, deleted
FROM
    Documents
WHERE
    id = ?
LIMIT
    1
`

func (q *Queries) GetDocument(ctx context.Context, id string) (Document, error) {
	row := q.db.QueryRowContext(ctx, getDocument, id)
	var i Document
	err := row.Scan(
		&i.ID,
		&i.Base64data,
		&i.Filename,
		&i.Datecreated,
		&i.Deleted,
	)
	return i, err
}

const getDocuments = `-- name: GetDocuments :many
SELECT
    id, base64data, filename, datecreated, deleted
FROM
    Documents
LIMIT
    ?
OFFSET
    ?
`

type GetDocumentsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetDocuments(ctx context.Context, arg GetDocumentsParams) ([]Document, error) {
	rows, err := q.db.QueryContext(ctx, getDocuments, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Document
	for rows.Next() {
		var i Document
		if err := rows.Scan(
			&i.ID,
			&i.Base64data,
			&i.Filename,
			&i.Datecreated,
			&i.Deleted,
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
