-- name: GetRmAllowedmanager :one
SELECT
    *
FROM
    rm_allowedmanagers
WHERE
    srvid = ?
LIMIT
    1;

-- name: GetRmAllowedmanagers :many
SELECT
    *
FROM
    rm_allowedmanagers
LIMIT
    ?
OFFSET
    ?;