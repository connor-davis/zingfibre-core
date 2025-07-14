-- name: GetRmSpecperbw :one
SELECT
    *
FROM
    rm_specperbw
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmSpecperbws :many
SELECT
    *
FROM
    rm_specperbw
LIMIT
    ?
OFFSET
    ?;