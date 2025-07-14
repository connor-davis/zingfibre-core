-- name: GetCustomerNote :one
SELECT
    *
FROM
    CustomerNotes
WHERE
    id = ?
LIMIT
    1;

-- name: GetCustomerNotes :many
SELECT
    *
FROM
    CustomerNotes
LIMIT
    ?
OFFSET
    ?;