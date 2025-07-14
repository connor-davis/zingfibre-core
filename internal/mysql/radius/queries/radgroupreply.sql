-- name: GetRadgroupreply :one
SELECT
    *
FROM
    radgroupreply
WHERE
    id = ?
LIMIT
    1;

-- name: GetRadgroupreplies :many
SELECT
    *
FROM
    radgroupreply
LIMIT
    ?
OFFSET
    ?;