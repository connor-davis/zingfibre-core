-- name: GetSalesAgent :one
SELECT
    *
FROM
    SalesAgents
WHERE
    id = ?
LIMIT
    1;

-- name: GetSalesAgents :many
SELECT
    *
FROM
    SalesAgents
LIMIT
    ?
OFFSET
    ?;