-- name: GetRmColsetlistradius :one
SELECT
    *
FROM
    rm_colsetlistradius
WHERE
    managername = ?
LIMIT
    1;

-- name: GetRmColsetlistradiuses :many
SELECT
    *
FROM
    rm_colsetlistradius
LIMIT
    ?
OFFSET
    ?;