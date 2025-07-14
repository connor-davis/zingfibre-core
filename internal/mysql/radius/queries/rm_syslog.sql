-- name: GetRmSyslog :one
SELECT
    *
FROM
    rm_syslog
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmSyslogs :many
SELECT
    *
FROM
    rm_syslog
LIMIT
    ?
OFFSET
    ?;