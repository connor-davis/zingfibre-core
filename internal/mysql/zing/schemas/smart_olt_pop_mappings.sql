CREATE TABLE
    `SmartOLTPOPMappings` (
        `Id` smallint (6) NOT NULL AUTO_INCREMENT,
        `SmartOLTTenantId` smallint (6) NOT NULL,
        `POP` longtext DEFAULT NULL,
        PRIMARY KEY (`Id`),
        KEY `IX_SmartOLTPOPMappings_SmartOLTTenantId` (`SmartOLTTenantId`),
        CONSTRAINT `FK_SmartOLTPOPMappings_SmartOLTTenants_SmartOLTTenantId` FOREIGN KEY (`SmartOLTTenantId`) REFERENCES `SmartOLTTenants` (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 8 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;