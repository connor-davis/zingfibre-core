-- name: GetBuild :one
SELECT
    *
FROM
    Builds
WHERE
    id = ?
LIMIT
    1;

-- name: GetBuilds :many
SELECT
    *
FROM
    Builds
LIMIT
    ?
OFFSET
    ?;