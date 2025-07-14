-- name: GetRmNewuser :one
SELECT
    *
FROM
    rm_newusers
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmNewusers :many
SELECT
    *
FROM
    rm_newusers
LIMIT
    ?
OFFSET
    ?;