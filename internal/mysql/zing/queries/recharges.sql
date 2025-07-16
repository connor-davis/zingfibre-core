-- name: GetRecharge :one
SELECT
  *
FROM
  Recharges
WHERE
  id = ?
LIMIT
  1;

-- name: GetRecharges :many
SELECT
  *
FROM
  Recharges
LIMIT
  ?
OFFSET
  ?;

-- name: GetRechargeTypeCounts :many
SELECT
  *
FROM
  (
    SELECT
      t3.Name AS product_name,
      COUNT(*) AS recharge_count,
      CASE
        WHEN sqlc.arg (period) = 'weeks' THEN CONCAT(
          FLOOR((DAY(t1.DateCreated) - 1) / 7) + 1,
          '-',
          MONTH(t1.DateCreated),
          '-',
          YEAR(t1.DateCreated)
        )
        WHEN sqlc.arg (period) = 'months' THEN CONCAT(MONTH(t1.DateCreated), '-', YEAR(t1.DateCreated))
      END AS period,
      MAX(t1.DateCreated) AS max_date
    FROM
      Recharges t1
      LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
      LEFT JOIN Products t3 ON t1.ProductId = t3.Id
    WHERE
      TRIM(LOWER(t2.RadiusUsername)) LIKE CONCAT(TRIM(LOWER(sqlc.arg (poi))), '%')
      AND (
        (
          sqlc.arg (period) = 'weeks'
          AND t1.DateCreated >= DATE_SUB(NOW(), INTERVAL sqlc.arg (COUNT) WEEK)
        )
        OR (
          sqlc.arg (period) = 'months'
          AND t1.DateCreated >= DATE_SUB(NOW(), INTERVAL sqlc.arg (COUNT) MONTH)
        )
      )
    GROUP BY
      t3.Name,
      CASE
        WHEN sqlc.arg (period) = 'weeks' THEN CONCAT(
          FLOOR((DAY(t1.DateCreated) - 1) / 7) + 1,
          '-',
          MONTH(t1.DateCreated),
          '-',
          YEAR(t1.DateCreated)
        )
        WHEN sqlc.arg (period) = 'months' THEN CONCAT(MONTH(t1.DateCreated), '-', YEAR(t1.DateCreated))
      END
  ) AS Sub
ORDER BY
  MaxDate ASC;