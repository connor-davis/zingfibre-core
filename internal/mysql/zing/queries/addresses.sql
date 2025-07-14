-- name: GetAddress :one
SELECT
    *
FROM
    Addresses
WHERE
    id = ?
LIMIT
    1;

-- name: GetAddresses :many
SELECT
    *
FROM
    Addresses
LIMIT
    ?
OFFSET
    ?;