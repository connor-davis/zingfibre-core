-- name: GetPermission :one
SELECT
    *
FROM
    Permissions
WHERE
    id = ?
LIMIT
    1;

-- name: GetPermissions :many
SELECT
    *
FROM
    Permissions
LIMIT
    ?
OFFSET
    ?;