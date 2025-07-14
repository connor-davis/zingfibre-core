-- name: GetRmSetting :one
SELECT
    *
FROM
    rm_settings
LIMIT
    1;

-- name: GetRmSettings :many
SELECT
    *
FROM
    rm_settings
LIMIT
    ?
OFFSET
    ?;