-- name: GetRadreply :one
SELECT
    *
FROM
    radreply
WHERE
    id = ?
LIMIT
    1;

-- name: GetRadreplies :many
SELECT
    *
FROM
    radreply
LIMIT
    ?
OFFSET
    ?;