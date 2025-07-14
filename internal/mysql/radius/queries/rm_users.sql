-- name: GetRmUser :one
SELECT
    *
FROM
    rm_users
WHERE
    username = ?
LIMIT
    1;

-- name: GetRmUsers :many
SELECT
    *
FROM
    rm_users
LIMIT
    ?
OFFSET
    ?;