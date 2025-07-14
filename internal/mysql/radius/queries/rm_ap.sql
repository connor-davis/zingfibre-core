-- name: GetRmAp :one
SELECT
    *
FROM
    rm_ap
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmAps :many
SELECT
    *
FROM
    rm_ap
LIMIT
    ?
OFFSET
    ?;