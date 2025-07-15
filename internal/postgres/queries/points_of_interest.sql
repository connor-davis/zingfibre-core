-- name: GetPointOfInterest :one
SELECT
    *
FROM
    points_of_interest
WHERE
    id = $1
LIMIT
    1;

-- name: GetTotalPointsOfInterest :one
SELECT
    COUNT(*) AS total
FROM
    points_of_interest
LIMIT
    1;

-- name: GetPointsOfInterest :many
SELECT
    *
FROM
    points_of_interest
LIMIT $1
OFFSET $2;

-- name: CreatePointOfInterest :one
INSERT INTO
    points_of_interest (
        name,
        key
    )
VALUES
    ($1, $2) RETURNING *;

-- name: UpdatePointOfInterest :one
UPDATE points_of_interest
SET
    name = $1,
    key = $2,
    updated_at = NOW()
WHERE
    id = $3 RETURNING *;

-- name: DeletePointOfInterest :one
DELETE FROM points_of_interest
WHERE
    id = $1 RETURNING *;