-- name: GetRole :one
SELECT
    *
FROM
    Roles
WHERE
    id = ?
LIMIT
    1;

-- name: GetRoles :many
SELECT
    *
FROM
    Roles
LIMIT
    ?
OFFSET
    ?;