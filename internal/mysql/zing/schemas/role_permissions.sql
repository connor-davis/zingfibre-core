CREATE TABLE
    `RolePermissions` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `RoleId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `PermissionId` bigint (20) NOT NULL,
            `Allowed` tinyint (1) NOT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_RolePermissions_PermissionId` (`PermissionId`),
            KEY `IX_RolePermissions_RoleId` (`RoleId`),
            CONSTRAINT `FK_RolePermissions_Permissions_PermissionId` FOREIGN KEY (`PermissionId`) REFERENCES `Permissions` (`Id`),
            CONSTRAINT `FK_RolePermissions_Roles_RoleId` FOREIGN KEY (`RoleId`) REFERENCES `Roles` (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;