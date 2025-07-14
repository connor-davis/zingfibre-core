-- name: GetRecharge :one
SELECT
    *
FROM
    Recharges
WHERE
    id = ?
LIMIT
    1;

-- name: GetRecharges :many
SELECT
    *
FROM
    Recharges
LIMIT
    ?
OFFSET
    ?;