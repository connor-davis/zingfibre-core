-- name: GetUser :one
SELECT
    *
FROM
    Users
WHERE
    id = ?
LIMIT
    1;

-- name: GetUsers :many
SELECT
    *
FROM
    Users
LIMIT
    ?
OFFSET
    ?;