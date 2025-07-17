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

-- name: GetPOPs :many
SELECT
    POP AS pop
FROM
    Addresses
WHERE
    POP IS NOT NULL
GROUP BY
    POP
ORDER BY
    pop ASC