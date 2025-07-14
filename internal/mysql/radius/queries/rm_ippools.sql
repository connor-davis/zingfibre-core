-- name: GetRmIppool :one
SELECT
    *
FROM
    rm_ippools
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmIppools :many
SELECT
    *
FROM
    rm_ippools
LIMIT
    ?
OFFSET
    ?;