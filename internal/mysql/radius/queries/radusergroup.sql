-- name: GetRadusergroup :one
SELECT
    *
FROM
    radusergroup
WHERE
    username = ?
LIMIT
    1;

-- name: GetRadusergroups :many
SELECT
    *
FROM
    radusergroup
LIMIT
    ?
OFFSET
    ?;