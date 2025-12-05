-- name: GetDynamicQuery :one
SELECT
    *
FROM
    dynamic_queries
WHERE
    id = $1
LIMIT
    1;

-- name: GetTotalDynamicQueries :one
SELECT
    COUNT(*) AS total
FROM
    dynamic_queries
LIMIT
    1;

-- name: GetDynamicQueries :many
SELECT
    *
FROM
    dynamic_queries
WHERE
    TRIM(LOWER(name)) ILIKE '%' || TRIM(LOWER(sqlc.arg(search_term)::text)) || '%'
ORDER BY
    name ASC
LIMIT $1
OFFSET $2;

-- name: CreateDynamicQuery :one
INSERT INTO
    dynamic_queries (
        name,
        query,
        response_id,
        status
    )
VALUES
    ($1, $2, $3, $4) RETURNING *;

-- name: UpdateDynamicQuery :one
UPDATE dynamic_queries
SET
    name = $1,
    query = $2,
    response_id = $3,
    status = $4,
    updated_at = NOW()
WHERE
    id = $5 RETURNING *;

-- name: DeleteDynamicQuery :one
DELETE FROM dynamic_queries
WHERE
    id = $1 RETURNING *;
