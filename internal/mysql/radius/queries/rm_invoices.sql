-- name: GetRmInvoice :one
SELECT
    *
FROM
    rm_invoices
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmInvoices :many
SELECT
    *
FROM
    rm_invoices
LIMIT
    ?
OFFSET
    ?;