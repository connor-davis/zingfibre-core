-- name: GetDynamicQueryColumn :one
SELECT
    *
FROM
    dynamic_query_columns
WHERE
    id = $1
LIMIT
    1;

-- name: GetDynamicQueryColumnAggregates :many
SELECT
    t1.*
FROM
    dynamic_query_aggregates t1
WHERE
    t1.dynamic_query_column_id = $1
ORDER BY
    t1.aggregate_column ASC,
    t1.aggregate_type ASC;

-- name: CreateDynamicQueryColumn :one
INSERT INTO
    dynamic_query_columns (
        dynamic_query_id,
        column_table,
        column_name,
        column_label
    )
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateDynamicQueryColumn :one
UPDATE dynamic_query_columns
SET
    dynamic_query_id = $1,
    column_table = $2,
    column_name = $3,
    column_label = $4
WHERE
    id = $5
RETURNING
    *;

-- name: DeleteDynamicQueryColumn :exec
DELETE FROM
    dynamic_query_columns
WHERE
    id = $1;