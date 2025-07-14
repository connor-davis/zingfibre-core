-- name: GetRadippool :one
SELECT
    *
FROM
    radippool
WHERE
    id = ?
LIMIT
    1;

-- name: GetRadippools :many
SELECT
    *
FROM
    radippool
LIMIT
    ?
OFFSET
    ?;