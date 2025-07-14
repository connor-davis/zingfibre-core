-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT
    1;

-- name: GetUsers :many
SELECT
    *
FROM
    users
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
INSERT INTO
    users (
        email,
        password,
        mfa_secret,
        mfa_enabled,
        mfa_verified
    )
VALUES
    ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    email = $1,
    password = $2,
    mfa_secret = $3,
    mfa_enabled = $4,
    mfa_verified = $5
WHERE
    id = $6 RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE
    id = $1 RETURNING *;