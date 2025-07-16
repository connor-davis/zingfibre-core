-- name: GetReportsRechargeTypeCounts :many
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
		WHERE
			TRIM(LOWER(t2.RadiusUsername)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%')
			AND(
                (
                    sqlc.arg('period') = 'weeks'
                    AND t1.DateCreated >= DATE_SUB(DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY), INTERVAL (sqlc.arg('count') - 1) WEEK)
                )
                OR(
                    sqlc.arg('period') = 'months'
                    AND t1.DateCreated >= DATE_SUB(DATE_SUB(CURDATE(), INTERVAL DAY(CURDATE()) - 1 DAY), INTERVAL (sqlc.arg('count') - 1) MONTH)
                )
			)
		GROUP BY
			recharge_name,
			recharge_period
	) AS sub
ORDER BY
	recharge_max_date ASC,
    recharge_count DESC;

-- name: GetReportsCustomers :many
SELECT
    t1.*,
    t2.RadiusUsername AS t2_radius_username
FROM Customers t1
LEFT JOIN Addresses t2 ON t1.AddressId = t2.Id
WHERE
    TRIM(LOWER(t2.POP)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%');

-- name: GetReportsExpiringCustomers :many
SELECT
    t1.*,
    t3.Name AS last_purchase_duration,
    t3.Category AS last_purchase_speed,
    t4.RadiusUsername AS t4_radius_username,
    CONCAT(sqlc.arg('expiration')) AS expiration,
    CONCAT(sqlc.arg('address')) AS address
FROM
    Customers t1
LEFT JOIN (
    SELECT
        CustomerID,
        MAX(DateCreated) AS LastRechargeDate
    FROM
        Recharges
    GROUP BY
        CustomerID
) AS latest_recharge ON t1.Id = latest_recharge.CustomerID
LEFT JOIN Recharges t2 ON latest_recharge.CustomerID = t2.CustomerID AND latest_recharge.LastRechargeDate = t2.DateCreated
LEFT JOIN Products t3 ON t2.ProductId = t3.Id
LEFT JOIN Addresses t4 ON t1.AddressId = t4.Id
WHERE
    TRIM(LOWER(t4.POP)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))),'%');

-- name: GetReportsRecharges :many
SELECT
    t2.RadiusUsername AS radius_username,
    t1.DateCreated AS date_created,
    CONCAT(t3.Category, ' ', t3.Name, ' Access') AS item_name,
    t1.*,
    t2.*,
    t4.ServiceId AS service_id,
    t5.Name AS build_name,
    t6.Name AS build_type
FROM
    Recharges t1
LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
LEFT JOIN Products t3 ON t1.ProductId = t3.Id
LEFT JOIN Addresses t4 ON t2.AddressId = t4.Id
LEFT JOIN Builds t5 ON t4.BuildId = t5.Id
LEFT JOIN BuildTypes t6 ON t5.BuildTypeId = t6.Id
WHERE
    TRIM(LOWER(t4.POP)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND CAST(t1.DateCreated AS DATE) >= sqlc.arg('start_date')
    AND CAST(t1.DateCreated AS DATE) <= sqlc.arg('end_date')
ORDER BY
    t1.DateCreated DESC;

-- name: GetReportsRechargesSummary :many
SELECT
    t2.RadiusUsername AS radius_username,
    t1.DateCreated AS date_created,
    CONCAT(t3.Category, ' ', t3.Name, ' Access') as item_name,
    t1.*,
    t2.*,
    t4.ServiceId AS service_id,
    t5.Name AS build_name,
    t6.Name AS build_type
FROM
    Recharges t1
LEFT JOIN Customers t2 ON t1.CustomerId = t2.Id
LEFT JOIN Products t3 ON t1.ProductId = t3.Id
LEFT JOIN Addresses t4 ON t2.AddressId = t4.Id
LEFT JOIN Builds t5 ON t4.BuildId = t5.Id
LEFT JOIN BuildTypes t6 ON t5.BuildTypeId = t6.Id
WHERE
    TRIM(LOWER(t4.POP)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND t1.DateCreated >= DATE_FORMAT(NOW(), '%Y-%m-01')
ORDER BY
    t1.DateCreated DESC;

-- name: GetReportsSummary :many
SELECT
    t1.RadiusUsername AS radius_username,
    t2.DateCreated AS date_created,

    CASE WHEN JSON_VALID(PaymentServicePayload) = 1
        THEN CONCAT(t3.Category, ' ', t3.Name, ' Access')
        ELSE 'Intro Package'
    END AS item_name,

    CASE WHEN JSON_VALID(PaymentServicePayload) = 1 
        THEN JSON_VALUE(PaymentServicePayload, '$.amount_gross') 
        ELSE '0' 
    END AS amount_gross,

    CASE WHEN JSON_VALID(PaymentServicePayload) = 1 
        THEN JSON_VALUE(PaymentServicePayload, '$.amount_fee') 
        ELSE '0.0' 
    END AS amount_fee,

    CASE WHEN JSON_VALID(PaymentServicePayload) = 1 
        THEN JSON_VALUE(PaymentServicePayload, '$.amount_net') 
        ELSE '0' 
    END AS amount_net,

    CASE WHEN JSON_VALID(PaymentServicePayload) = 1 
        THEN JSON_VALUE(PaymentServicePayload, '$.accountNumber') 
        ELSE '0'
    END AS cash_code,

    CASE 
        WHEN
            JSON_VALID(PaymentServicePayload) = 1 
            AND JSON_VALUE(PaymentServicePayload, '$.tenders[0].amount') REGEXP '^[0-9]+(\\.[0-9]+)?$'
        THEN 
            CAST(JSON_VALUE(PaymentServicePayload, '$.tenders[0].amount') AS DECIMAL(10,2)) / 100
        ELSE 0
    END AS cash_amount,

    t4.ServiceId AS service_id,
    t5.Name AS build_name,
    t6.Name AS build_type
FROM Customers t1
LEFT JOIN Recharges t2 ON t1.Id = t2.CustomerID
LEFT JOIN Products t3 ON t2.ProductId = t3.Id
LEFT JOIN Addresses t4 ON t1.AddressId = t4.Id
LEFT JOIN Builds t5 ON t4.BuildId = t5.Id
LEFT JOIN BuildTypes t6 ON t5.BuildTypeId = t6.Id
WHERE 
    TRIM(LOWER(t4.POP)) LIKE CONCAT(TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND t2.DateCreated >= DATE_FORMAT(NOW(), '%Y-%m-01')
ORDER BY
    t2.DateCreated DESC