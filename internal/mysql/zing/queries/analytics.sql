-- name: GetAnalyticsMonthlyRevenueStatistics :one
SELECT
	SUM(
        CASE
            WHEN 
                t1.DateCreated >= DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
            THEN t1.PaymentAmount
            ELSE 0
        END
    ) - 0 AS revenue,
	COALESCE(SUM(
		CASE
			WHEN 
				t1.DateCreated >= DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
			THEN t1.PaymentAmount
			ELSE 0
		END
	), 0)
	-
	COALESCE(SUM(
		CASE
			WHEN 
				t1.DateCreated >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01 00:00:00')
				AND t1.DateCreated < DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
			THEN t1.PaymentAmount
			ELSE 0
		END
	), 0) AS revenue_growth_amount,
	ROUND(
	(
		(
			SUM(
				CASE
					WHEN 
						t1.DateCreated >= DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
					THEN PaymentAmount
					ELSE 0
				END
			)
			/
			NULLIF(SUM(
				CASE
					WHEN 
						t1.DateCreated >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01 00:00:00')
						AND t1.DateCreated < DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
					THEN PaymentAmount
					ELSE 0
				END
			), 0)) - 1
		) * 100, 2
	) AS revenue_growth_percentage
FROM Recharges t1
LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
LEFT JOIN Addresses t3 ON t2.AddressId = t3.Id
WHERE
    t1.RechargeSuccessful = 1
    AND TRIM(LOWER(t3.POP)) LIKE TRIM(LOWER(CONCAT(sql.arg('pop'), '%')));

-- name: GetAnalyticsMonthlyUniquePurchasers :one
SELECT
    COUNT(DISTINCT t2.RadiusUsername) - 0 AS unique_purchasers
FROM
    Recharges t1
LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
LEFT JOIN Addresses t3 ON t2.AddressId = t3.Id
WHERE
    t1.RechargeSuccessful = 1
    AND TRIM(LOWER(t3.POP)) LIKE TRIM(LOWER(CONCAT(sql.arg('pop'), '%')))
    AND t1.DateCreated >= DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')