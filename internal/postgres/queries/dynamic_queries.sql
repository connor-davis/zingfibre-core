-- name: GetDynamicQuery :one
SELECT
    *
FROM
    dynamic_queries
WHERE
    id = $1
LIMIT
    1;

-- name: GetDynamicQueries :many
SELECT
    *
FROM
    dynamic_queries
ORDER BY
    query_name ASC;

-- name: GetDynamicQueryChildQueries :many
SELECT
    t1.*
FROM
    dynamic_queries t1
LEFT JOIN dynamic_query_child_queries t2 ON t2.child_query_id = t1.id
WHERE
    t2.dynamic_query_id = $1
ORDER BY
    t1.query_name ASC;

-- name: GetDynamicQueryColumns :many
SELECT
    t1.*
FROM
    dynamic_query_columns t1
WHERE
    t1.dynamic_query_id = $1
ORDER BY
    t1.column_name ASC,
    t1.column_table ASC,
    t1.column_label ASC;

-- name: GetDynamicQueryJoins :many
SELECT
    t1.*
FROM
    dynamic_query_joins t1
WHERE
    t1.dynamic_query_id = $1
ORDER BY
    t1.join_left_table ASC,
    t1.join_left_column ASC,
    t1.join_right_table ASC,
    t1.join_right_column ASC;

-- name: GetDynamicQueryClauses :many
SELECT
    t1.*
FROM
    dynamic_query_clauses t1
WHERE
    t1.dynamic_query_id = $1
ORDER BY
    t1.clause_column ASC,
    t1.clause_operator ASC,
    t1.clause_value ASC;

-- name: GetDynamicQueryOrders :many
SELECT
    t1.*
FROM
    dynamic_query_orders t1
WHERE
    t1.dynamic_query_id = $1
ORDER BY
    t1.order_column ASC,
    t1.order_type ASC;


-- name: CreateDynamicQuery :one
INSERT INTO
    dynamic_queries (
        query_name,
        is_primary
    )
VALUES
    ($1, $2)
RETURNING
    *;

-- name: UpdateDynamicQuery :one
UPDATE dynamic_queries
SET
    query_name = $1,
    is_primary = $2,
    updated_at = now()
WHERE
    id = $3
RETURNING
    *;

-- name: DeleteDynamicQuery :exec
DELETE FROM
    dynamic_queries
WHERE
    id = $1;