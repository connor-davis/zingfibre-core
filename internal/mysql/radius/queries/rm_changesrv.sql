-- name: GetRmChangesrv :one
SELECT
    *
FROM
    rm_changesrv
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmChangesrvs :many
SELECT
    *
FROM
    rm_changesrv
LIMIT
    ?
OFFSET
    ?;