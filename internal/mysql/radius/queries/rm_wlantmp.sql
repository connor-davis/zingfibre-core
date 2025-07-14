-- name: GetRmWlantmp :one
SELECT
    *
FROM
    rm_wlantmp
WHERE
    maccpe = ?
LIMIT
    1;

-- name: GetRmWlantmps :many
SELECT
    *
FROM
    rm_wlantmp
LIMIT
    ?
OFFSET
    ?;