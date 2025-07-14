-- name: GetSmartOLTTenant :one
SELECT
    *
FROM
    SmartOLTTenants
WHERE
    id = ?
LIMIT
    1;

-- name: GetSmartOLTTenants :many
SELECT
    *
FROM
    SmartOLTTenants
LIMIT
    ?
OFFSET
    ?;