-- name: GetReportExportsCustomers :many
SELECT
    CONCAT(t1.FirstName, ' ', t1.Surname) AS full_name,
    t1.Email AS email,
    t2.RadiusUsername AS radius_username,
    t1.PhoneNumber AS phone_number
FROM Customers t1
LEFT JOIN Addresses t2 ON t1.AddressId = t2.Id
WHERE
    TRIM(LOWER(t2.POP)) LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('poi'))), '%')
ORDER BY
    CONCAT(t1.FirstName, ' ', t1.Surname) ASC,
    t1.Email ASC;

-- name: GetReportExportsExpiringCustomers :many
SELECT
    CONCAT(t1.FirstName, ' ', t1.Surname) AS full_name,
    t1.Email AS email,
    t1.PhoneNumber AS phone_number,
    t4.RadiusUsername AS radius_username,
    t3.Name AS last_purchase_duration,
    t3.Category AS last_purchase_speed,
    t4.StreetAddress AS Address,
    CONCAT('') AS expiration
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
    TRIM(LOWER(t4.POP)) LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND (
        t1.FirstName LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
        OR t1.Surname LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
        OR t1.Email LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
        OR t1.PhoneNumber LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
        OR t4.RadiusUsername LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
        OR t3.Name LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
        OR t3.Category LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('search'))), '%')
    )
    AND t4.RadiusUsername IN (sqlc.slice('radius_usernames'))
ORDER BY
    CONCAT(t1.FirstName, ' ', t1.Surname) ASC,
    t1.Email ASC;

-- name: GetReportExportsRecharges :many
SELECT
    t1.DateCreated AS date_created,
    t2.Email AS email,
    CONCAT(t2.FirstName, ' ', t2.Surname) AS full_name,
    CASE 
        WHEN t3.Category IS NULL OR t3.Name IS NULL THEN 'Intro Package'
        ELSE CONCAT(t3.Category, ' ', t3.Name, ' Access')
    END AS item_name,
    t1.PaymentAmount AS amount,
    t1.Method AS method,
    t1.RechargeSuccessful AS successful,
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
    TRIM(LOWER(t4.POP)) LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND CAST(t1.DateCreated AS DATE) >= sqlc.arg('start_date')
    AND CAST(t1.DateCreated AS DATE) <= sqlc.arg('end_date')
ORDER BY
    t1.DateCreated DESC;

-- name: GetReportExportsRechargesSummary :many
SELECT
    t1.DateCreated AS date_created,
    t2.Email AS email,
    CONCAT(t2.FirstName, ' ', t2.Surname) AS full_name,
    CASE 
        WHEN t3.Category IS NULL OR t3.Name IS NULL THEN 'Intro Package'
        ELSE CONCAT(t3.Category, ' ', t3.Name, ' Access')
    END AS item_name,
    t1.PaymentAmount AS amount,
    t1.Method AS method,
    t1.RechargeSuccessful AS successful,
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
    TRIM(LOWER(t4.POP)) LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND t1.DateCreated >= DATE_FORMAT(NOW(), '%Y-%m-01')
ORDER BY
    t1.DateCreated DESC;

-- name: GetReportExportsSummary :many
SELECT
    t2.DateCreated AS date_created,
    CASE 
        WHEN t3.Category IS NULL OR t3.Name IS NULL THEN 'Intro Package'
        ELSE CONCAT(t3.Category, ' ', t3.Name, ' Access')
    END AS item_name,
    t4.RadiusUsername AS radius_username,
    t2.Method AS method,
    t2.PaymentAmount AS amount,
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
    TRIM(LOWER(t4.POP)) LIKE CONCAT('%', TRIM(LOWER(sqlc.arg('poi'))), '%')
    AND t2.DateCreated >= 
        CASE 
            WHEN sqlc.arg('months') = 1 THEN DATE_FORMAT(NOW(), '%Y-%m-01 00:00:00')
            ELSE DATE_FORMAT(DATE_SUB(DATE_FORMAT(NOW(), '%Y-%m-01'), INTERVAL (sqlc.arg('months') - 1) MONTH), '%Y-%m-01 00:00:00')
        END
    AND t2.RechargeSuccessful = 1
ORDER BY
    t2.DateCreated DESC;