-- name: GetDynamicQueryClause :one
SELECT
    *
FROM
    dynamic_query_clauses
WHERE
    id = $1
LIMIT
    1;

-- name: CreateDynamicQueryClause :one
INSERT INTO
    dynamic_query_clauses (
        dynamic_query_id,
        clause_column,
        clause_operator,
        clause_value
    )
VALUES
    ($1, $2, $3, $4)
RETURNING
    *;

-- name: UpdateDynamicQueryClause :one
UPDATE dynamic_query_clauses
SET
    dynamic_query_id = $1,
    clause_column = $2,
    clause_operator = $3,
    clause_value = $4
WHERE
    id = $5
RETURNING
    *;

-- name: DeleteDynamicQueryClause :exec
DELETE FROM
    dynamic_query_clauses
WHERE
    id = $1;