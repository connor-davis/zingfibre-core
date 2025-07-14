-- name: GetRmMpesa :one
SELECT
    *
FROM
    rm_mpesa
WHERE
    merchantrequestid = ?
LIMIT
    1;

-- name: GetRmMpesas :many
SELECT
    *
FROM
    rm_mpesa
LIMIT
    ?
OFFSET
    ?;