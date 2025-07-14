-- name: GetNas :one
SELECT
    *
FROM
    nas
WHERE
    id = ?
LIMIT
    1;

-- name: GetNases :many
SELECT
    *
FROM
    nas
LIMIT
    ?
OFFSET
    ?;