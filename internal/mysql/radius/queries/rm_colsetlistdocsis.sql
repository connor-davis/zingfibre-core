-- name: GetRmColsetlistdocsis :one
SELECT
    *
FROM
    rm_colsetlistdocsis
WHERE
    managername = ?
LIMIT
    1;

-- name: GetRmColsetlistdocsis :many
SELECT
    *
FROM
    rm_colsetlistdocsis
LIMIT
    ?
OFFSET
    ?;