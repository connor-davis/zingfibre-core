-- name: GetUserPermission :one
SELECT
    *
FROM
    UserPermissions
WHERE
    id = ?
LIMIT
    1;

-- name: GetUserPermissions :many
SELECT
    *
FROM
    UserPermissions
LIMIT
    ?
OFFSET
    ?;