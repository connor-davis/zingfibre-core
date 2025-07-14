-- name: GetRmWlan :one
SELECT
    *
FROM
    rm_wlan
WHERE
    maccpe = ?
LIMIT
    1;

-- name: GetRmWlans :many
SELECT
    *
FROM
    rm_wlan
LIMIT
    ?
OFFSET
    ?;