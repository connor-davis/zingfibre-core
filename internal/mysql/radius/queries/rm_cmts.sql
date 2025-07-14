-- name: GetRmCmt :one
SELECT
    *
FROM
    rm_cmts
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmCmts :many
SELECT
    *
FROM
    rm_cmts
LIMIT
    ?
OFFSET
    ?;