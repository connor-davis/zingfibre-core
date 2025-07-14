-- name: GetProduct :one
SELECT
    *
FROM
    Products
WHERE
    id = ?
LIMIT
    1;

-- name: GetProducts :many
SELECT
    *
FROM
    Products
LIMIT
    ?
OFFSET
    ?;