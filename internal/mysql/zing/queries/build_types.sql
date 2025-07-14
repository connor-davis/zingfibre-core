-- name: GetBuildType :one
SELECT
    *
FROM
    BuildTypes
WHERE
    id = ?
LIMIT
    1;

-- name: GetBuildTypes :many
SELECT
    *
FROM
    BuildTypes
LIMIT
    ?
OFFSET
    ?;