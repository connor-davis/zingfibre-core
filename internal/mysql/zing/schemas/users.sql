CREATE TABLE
    `Users` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `EmailAddress` longtext NOT NULL,
            `Password` longtext NOT NULL,
            `PasswordSalt` longtext DEFAULT NULL,
            `FirstName` longtext NOT NULL,
            `LastName` longtext DEFAULT NULL,
            `Active` tinyint (1) NOT NULL,
            `RoleId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `TwoFactorKey` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;