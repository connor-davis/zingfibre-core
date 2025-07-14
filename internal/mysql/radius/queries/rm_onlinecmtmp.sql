-- name: GetRmOnlinecmtmp :one
SELECT
    *
FROM
    rm_onlinecmtmp
WHERE
    username = ?
LIMIT
    1;

-- name: GetRmOnlinecmtmps :many
SELECT
    *
FROM
    rm_onlinecmtmp
LIMIT
    ?
OFFSET
    ?;