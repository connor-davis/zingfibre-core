-- name: GetPaymentRequest :one
SELECT
    *
FROM
    PaymentRequests
WHERE
    id = ?
LIMIT
    1;

-- name: GetPaymentRequests :many
SELECT
    *
FROM
    PaymentRequests
LIMIT
    ?
OFFSET
    ?;