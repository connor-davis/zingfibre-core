-- name: GetAnalyticsMonthlyRevenueStatistics :one
SELECT
	CAST(
        SUM(
            CASE
                WHEN 
                    t1.DateCreated >= DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
                THEN t1.PaymentAmount
                ELSE 0
            END
        ) AS SIGNED
    ) AS revenue,
	CAST(
        (
            SUM(
                CASE
                    WHEN 
                        t1.DateCreated >= DATE_FORMAT(CURDATE(), '%Y-%m-01 00:00:00')
                    THEN t1.PaymentAmount
                    ELSE 0
                END
            )
            -
            SUM(
                CASE
                    WHEN 
                        t1.DateCreated >= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-01 00:00:00')
                        AND t1.DateCreated <= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-%d 23:59:59')
                    THEN t1.PaymentAmount
                    ELSE 0
                END
            )
        ) AS SIGNED
    ) AS revenue_growth_amount,
	ROUND(
        ABS(
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
                            AND t1.DateCreated <= DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 1 MONTH), '%Y-%m-%d 23:59:59')
                        THEN PaymentAmount
                        ELSE 0
                    END
                ), 0) - 1
            )
        ), 4
    ) AS revenue_growth_percentage
FROM Recharges t1
LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
LEFT JOIN Addresses t3 ON t2.AddressId = t3.Id
WHERE
    t1.RechargeSuccessful = 1
    AND TRIM(LOWER(t3.POP)) LIKE TRIM(LOWER(CONCAT(sqlc.arg('poi'), '%')));

-- name: GetAnalyticsMonthlyUniquePurchasers :one
SELECT
    CAST(
        COUNT(DISTINCT t2.RadiusUsername) AS SIGNED
    ) AS unique_purchasers
FROM
    Recharges t1
LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
LEFT JOIN Addresses t3 ON t2.AddressId = t3.Id
WHERE
    t1.RechargeSuccessful = 1
    AND TRIM(LOWER(t3.POP)) LIKE TRIM(LOWER(CONCAT(sqlc.arg('poi'), '%')))
    AND t1.DateCreated >= DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')