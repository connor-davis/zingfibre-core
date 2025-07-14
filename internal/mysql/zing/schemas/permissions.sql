CREATE TABLE
    `Permissions` (
        `Id` bigint (20) NOT NULL AUTO_INCREMENT,
        `Name` longtext DEFAULT NULL,
        `Section` longtext DEFAULT NULL,
        `Type` longtext DEFAULT NULL,
        `Order` bigint (20) NOT NULL,
        `Deleted` tinyint (1) NOT NULL,
        PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 76 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;