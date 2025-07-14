-- name: GetUserslog :one
SELECT
    *
FROM
    userslog
WHERE
    id = ?
LIMIT
    1;

-- name: GetUserslogs :many
SELECT
    *
FROM
    userslog
LIMIT
    ?
OFFSET
    ?;