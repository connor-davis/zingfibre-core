CREATE TABLE
    `Products` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `Price` decimal(18, 2) NOT NULL,
            `Name` longtext DEFAULT NULL,
            `Category` longtext DEFAULT NULL,
            `Period` int (11) NOT NULL,
            `ServiceId` int (11) NOT NULL,
            `Months` int (11) DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;