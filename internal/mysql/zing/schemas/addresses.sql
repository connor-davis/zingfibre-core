CREATE TABLE
    `Addresses` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `ERF` longtext DEFAULT NULL,
            `StreetAddress` longtext DEFAULT NULL,
            `MDUName` longtext DEFAULT NULL,
            `MDUUnitNumber` longtext DEFAULT NULL,
            `MDUBlock` longtext DEFAULT NULL,
            `Township` longtext DEFAULT NULL,
            `PropertyType` longtext DEFAULT NULL,
            `POP` longtext DEFAULT NULL,
            `InstallDate` datetime (6) DEFAULT NULL,
            `RadiusUsername` longtext DEFAULT NULL,
            `RadiusPassword` longtext DEFAULT NULL,
            `InstallComplete` tinyint (1) NOT NULL,
            `W3W` longtext DEFAULT NULL,
            `InstallState` tinyint (3) unsigned DEFAULT NULL,
            `PoleNumber` longtext DEFAULT NULL,
            `ServiceID` bigint (20) NOT NULL AUTO_INCREMENT,
            `BuildId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `SalesAgentId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            UNIQUE KEY `ServiceID_UNIQUE` (`ServiceID`),
            UNIQUE KEY `Id_UNIQUE` (`Id`),
            KEY `IX_Addresses_BuildId` (`BuildId`),
            KEY `IX_Addresses_SalesAgentId` (`SalesAgentId`),
            CONSTRAINT `FK_Addresses_Builds_BuildId` FOREIGN KEY (`BuildId`) REFERENCES `Builds` (`Id`),
            CONSTRAINT `FK_Addresses_SalesAgents_SalesAgentId` FOREIGN KEY (`SalesAgentId`) REFERENCES `SalesAgents` (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 6102 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;