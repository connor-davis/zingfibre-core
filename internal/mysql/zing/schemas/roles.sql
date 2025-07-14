CREATE TABLE
    `Roles` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `Name` longtext NOT NULL,
            `CreatedByUserId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_Roles_CreatedByUserId` (`CreatedByUserId`),
            CONSTRAINT `FK_Roles_Users_CreatedByUserId` FOREIGN KEY (`CreatedByUserId`) REFERENCES `Users` (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;