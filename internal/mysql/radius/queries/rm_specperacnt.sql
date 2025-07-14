-- name: GetRmSpecperacnt :one
SELECT
    *
FROM
    rm_specperacnt
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmSpecperacnts :many
SELECT
    *
FROM
    rm_specperacnt
LIMIT
    ?
OFFSET
    ?;