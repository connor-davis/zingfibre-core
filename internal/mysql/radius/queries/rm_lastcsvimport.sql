-- name: GetRmLastcsvimport :one
SELECT
    *
FROM
    rm_lastcsvimport
WHERE
    username = ?
LIMIT
    1;

-- name: GetRmLastcsvimports :many
SELECT
    *
FROM
    rm_lastcsvimport
LIMIT
    ?
OFFSET
    ?;