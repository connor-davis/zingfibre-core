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
			t3.Name AS recharge_name,
			COUNT(*) AS recharge_count,
			CASE
				WHEN sqlc.arg('period') = 'weeks' THEN CONCAT(
					FLOOR((DAY(t1.DateCreated) - 1) / 7) + 1,
					'-',
					MONTH(t1.DateCreated),
					'-',
					YEAR(t1.DateCreated)
				)
				WHEN sqlc.arg('period') = 'months' THEN CONCAT(MONTH(t1.DateCreated), '-', YEAR(t1.DateCreated))
			END AS recharge_period,
			MAX(t1.DateCreated) AS recharge_max_date
		FROM
			Recharges t1
			LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
			LEFT JOIN Products t3 ON t1.ProductId = t3.Id
            LEFT JOIN Addresses t4 ON t2.AddressId = t4.Id
		WHERE
			TRIM(LOWER(t4.POP)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%')
			AND(
                (
                    sqlc.arg('period') = 'weeks'
                    AND t1.DateCreated >= 
                        CASE 
                            WHEN sqlc.arg('count') = 1 THEN DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')
                            ELSE DATE_FORMAT(DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL (sqlc.arg('count') - 1) WEEK), '%Y-%m-01 00:00:00')
                        END
                )
                OR(
                    sqlc.arg('period') = 'months'
                    AND t1.DateCreated >= 
                        CASE 
                            WHEN sqlc.arg('count') = 1 THEN DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')
                            ELSE DATE_FORMAT(DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL (sqlc.arg('count') - 1) MONTH), '%Y-%m-01 00:00:00')
                        END
                )
			)
		GROUP BY
			recharge_name,
			recharge_period
	) AS sub
ORDER BY
	recharge_max_date ASC,
    recharge_count DESC;