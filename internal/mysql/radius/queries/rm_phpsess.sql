-- name: GetRmPhpsess :one
SELECT
    *
FROM
    rm_phpsess
WHERE
    managername = ?
LIMIT
    1;

-- name: GetRmPhpsesses :many
SELECT
    *
FROM
    rm_phpsess
LIMIT
    ?
OFFSET
    ?;