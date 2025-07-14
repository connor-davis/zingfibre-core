CREATE TABLE
    `Documents` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `Base64Data` longtext DEFAULT NULL,
            `FileName` longtext DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;