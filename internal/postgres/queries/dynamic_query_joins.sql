-- name: GetDynamicQueryJoin :one
SELECT
    *
FROM
    dynamic_query_joins
WHERE
    id = $1
LIMIT
    1;

-- name: CreateDynamicQueryJoin :one
INSERT INTO
    dynamic_query_joins (
        dynamic_query_id,
        join_left_table,
        join_left_column,
        join_right_table,
        join_right_column,
        join_type
    )
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: UpdateDynamicQueryJoin :one
UPDATE dynamic_query_joins
SET
    dynamic_query_id = $1,
    join_left_table = $2,
    join_left_column = $3,
    join_right_table = $4,
    join_right_column = $5,
    join_type = $6
WHERE
    id = $7
RETURNING
    *;

-- name: DeleteDynamicQueryJoin :exec
DELETE FROM
    dynamic_query_joins
WHERE
    id = $1;