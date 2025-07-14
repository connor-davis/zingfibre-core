-- name: GetRadacct :one
SELECT
    *
FROM
    radacct
WHERE
    radacctid = ?
LIMIT
    1;

-- name: GetRadaccts :many
SELECT
    *
FROM
    radacct
LIMIT
    ?
OFFSET
    ?;