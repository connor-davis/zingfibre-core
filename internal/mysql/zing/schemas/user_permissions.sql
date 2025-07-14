CREATE TABLE
    `UserPermissions` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `UserId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `PermissionId` bigint (20) NOT NULL,
            `Allowed` tinyint (1) NOT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_UserPermissions_PermissionId` (`PermissionId`),
            KEY `IX_UserPermissions_UserId` (`UserId`),
            CONSTRAINT `FK_UserPermissions_Permissions_PermissionId` FOREIGN KEY (`PermissionId`) REFERENCES `Permissions` (`Id`),
            CONSTRAINT `FK_UserPermissions_Users_UserId` FOREIGN KEY (`UserId`) REFERENCES `Users` (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;