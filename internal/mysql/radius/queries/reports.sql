-- name: GetReportsExpiringCustomers :many
SELECT
    username,
    expiration
FROM
    rm_users
WHERE
    expiration IS NOT NULL;

-- name: GetReportsTotalExpiringCustomers :one
SELECT
    COUNT(*) AS total_expiring_customers
FROM
    rm_users
WHERE
    expiration IS NOT NULL;