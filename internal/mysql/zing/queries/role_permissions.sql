-- name: GetRolePermission :one
SELECT
    *
FROM
    RolePermissions
WHERE
    id = ?
LIMIT
    1;

-- name: GetRolePermissions :many
SELECT
    *
FROM
    RolePermissions
LIMIT
    ?
OFFSET
    ?;