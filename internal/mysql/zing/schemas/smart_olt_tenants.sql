CREATE TABLE
    `SmartOLTTenants` (
        `Id` smallint (6) NOT NULL AUTO_INCREMENT,
        `APIUrl` longtext DEFAULT NULL,
        `APIKey` longtext DEFAULT NULL,
        PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 3 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;