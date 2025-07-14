-- name: GetRmColsetlistuser :one
SELECT
    *
FROM
    rm_colsetlistusers
WHERE
    managername = ?
LIMIT
    1;

-- name: GetRmColsetlistusers :many
SELECT
    *
FROM
    rm_colsetlistusers
LIMIT
    ?
OFFSET
    ?;