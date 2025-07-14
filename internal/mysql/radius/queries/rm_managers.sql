-- name: GetRmManager :one
SELECT
    *
FROM
    rm_managers
WHERE
    managername = ?
LIMIT
    1;

-- name: GetRmManagers :many
SELECT
    *
FROM
    rm_managers
LIMIT
    ?
OFFSET
    ?;