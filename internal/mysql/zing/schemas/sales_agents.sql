CREATE TABLE
    `SalesAgents` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `Name` longtext DEFAULT NULL,
            `Code` longtext DEFAULT NULL,
            PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;