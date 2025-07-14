-- name: GetRmOnlineradius :one
SELECT
    *
FROM
    rm_onlineradius
WHERE
    username = ?
LIMIT
    1;

-- name: GetRmOnlineradiuses :many
SELECT
    *
FROM
    rm_onlineradius
LIMIT
    ?
OFFSET
    ?;