-- name: GetRmAllowednas :one
SELECT
    *
FROM
    rm_allowednases
WHERE
    srvid = ?
LIMIT
    1;

-- name: GetRmAllowednases :many
SELECT
    *
FROM
    rm_allowednases
LIMIT
    ?
OFFSET
    ?;