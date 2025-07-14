-- name: GetRmRadacct :one
SELECT
    *
FROM
    rm_radacct
WHERE
    radacctid = ?
LIMIT
    1;

-- name: GetRmRadaccts :many
SELECT
    *
FROM
    rm_radacct
LIMIT
    ?
OFFSET
    ?;