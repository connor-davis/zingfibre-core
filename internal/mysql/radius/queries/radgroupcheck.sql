-- name: GetRadgroupcheck :one
SELECT
    *
FROM
    radgroupcheck
WHERE
    id = ?
LIMIT
    1;

-- name: GetRadgroupchecks :many
SELECT
    *
FROM
    radgroupcheck
LIMIT
    ?
OFFSET
    ?;