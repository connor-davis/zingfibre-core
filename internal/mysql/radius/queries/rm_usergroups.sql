-- name: GetRmUsergroup :one
SELECT
    *
FROM
    rm_usergroups
WHERE
    groupid = ?
LIMIT
    1;

-- name: GetRmUsergroups :many
SELECT
    *
FROM
    rm_usergroups
LIMIT
    ?
OFFSET
    ?;