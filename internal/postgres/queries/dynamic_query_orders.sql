-- name: GetDynamicQueryOrder :one
SELECT
    *
FROM
    dynamic_query_orders
WHERE
    id = $1
LIMIT
    1;

-- name: CreateDynamicQueryOrder :one
INSERT INTO
    dynamic_query_orders (
        dynamic_query_id,
        order_column,
        order_type
    )
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: UpdateDynamicQueryOrder :one
UPDATE dynamic_query_orders
SET
    dynamic_query_id = $1,
    order_column = $2,
    order_type = $3
WHERE
    id = $4
RETURNING
    *;

-- name: DeleteDynamicQueryOrder :exec
DELETE FROM
    dynamic_query_orders
WHERE
    id = $1;