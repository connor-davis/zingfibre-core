-- name: GetRmService :one
SELECT
    *
FROM
    rm_services
WHERE
    srvid = ?
LIMIT
    1;

-- name: GetRmServices :many
SELECT
    *
FROM
    rm_services
LIMIT
    ?
OFFSET
    ?;