-- name: GetRmOnlinecm :one
SELECT
    *
FROM
    rm_onlinecm
WHERE
    username = ?
LIMIT
    1;

-- name: GetRmOnlinecms :many
SELECT
    *
FROM
    rm_onlinecm
LIMIT
    ?
OFFSET
    ?;