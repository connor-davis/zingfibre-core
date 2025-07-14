-- name: GetRmIa :one
SELECT
    *
FROM
    rm_ias
WHERE
    iasid = ?
LIMIT
    1;

-- name: GetRmIas :many
SELECT
    *
FROM
    rm_ias
LIMIT
    ?
OFFSET
    ?;