-- name: GetRadpostauth :one
SELECT
    *
FROM
    radpostauth
WHERE
    id = ?
LIMIT
    1;

-- name: GetRadpostauths :many
SELECT
    *
FROM
    radpostauth
LIMIT
    ?
OFFSET
    ?;