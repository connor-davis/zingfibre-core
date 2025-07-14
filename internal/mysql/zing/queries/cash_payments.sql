-- name: GetCashPayment :one
SELECT
    *
FROM
    CashPayments
WHERE
    id = ?
LIMIT
    1;

-- name: GetCashPayments :many
SELECT
    *
FROM
    CashPayments
LIMIT
    ?
OFFSET
    ?;