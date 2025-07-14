-- name: GetRmDailyacct :one
SELECT
    *
FROM
    rm_dailyacct
WHERE
    radacctid = ?
LIMIT
    1;

-- name: GetRmDailyaccts :many
SELECT
    *
FROM
    rm_dailyacct
LIMIT
    ?
OFFSET
    ?;