-- name: GetRmCard :one
SELECT
    *
FROM
    rm_cards
WHERE
    id = ?
LIMIT
    1;

-- name: GetRmCards :many
SELECT
    *
FROM
    rm_cards
LIMIT
    ?
OFFSET
    ?;