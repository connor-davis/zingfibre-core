-- name: GetRadcheck :one
SELECT
    *
FROM
    radcheck
WHERE
    id = ?
LIMIT
    1;

-- name: GetRadchecks :many
SELECT
    *
FROM
    radcheck
LIMIT
    ?
OFFSET
    ?;