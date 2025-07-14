-- name: GetDocument :one
SELECT
    *
FROM
    Documents
WHERE
    id = ?
LIMIT
    1;

-- name: GetDocuments :many
SELECT
    *
FROM
    Documents
LIMIT
    ?
OFFSET
    ?;