-- name: CreateDynamicQueryChildQuery :one
INSERT INTO
    dynamic_query_child_queries (
        dynamic_query_id,
        child_query_id
    )
VALUES
    ($1, $2)
RETURNING *;