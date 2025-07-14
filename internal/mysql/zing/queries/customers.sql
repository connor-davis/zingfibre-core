-- name: GetCustomer :one
SELECT
    *
FROM
    Customers
WHERE
    id = ?
LIMIT
    1;

-- name: GetCustomers :many
SELECT
    *
FROM
    Customers
LIMIT
    ?
OFFSET
    ?;