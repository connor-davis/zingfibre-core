-- name: GetReportsExpiringCustomers :many
SELECT
    username,
    expiration,
    address
FROM
    rm_users
WHERE
    TRIM(LOWER(username)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%');

-- name: GetReportsUnregisteredCustomers :many
SELECT
    username,
    expiration,
    address
FROM rm_users
WHERE
    lastlogoff IS NOT NULL;