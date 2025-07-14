-- name: GetPasswordResetRequest :one
SELECT
    *
FROM
    PasswordResetRequests
WHERE
    id = ?
LIMIT
    1;

-- name: GetPasswordResetRequests :many
SELECT
    *
FROM
    PasswordResetRequests
LIMIT
    ?
OFFSET
    ?;