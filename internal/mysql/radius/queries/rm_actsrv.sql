-- name: GetRmActsrv :one
SELECT
    *
FROM
    rm_actsrv
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmActsrvs :many
SELECT
    *
FROM
    rm_actsrv
LIMIT
    ?
OFFSET
    ?;