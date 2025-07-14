-- name: GetErrorLog :one
SELECT
    *
FROM
    ErrorLogs
WHERE
    id = ?
LIMIT
    1;

-- name: GetErrorLogs :many
SELECT
    *
FROM
    ErrorLogs
LIMIT
    ?
OFFSET
    ?;