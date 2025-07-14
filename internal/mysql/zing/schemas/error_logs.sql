CREATE TABLE
    `ErrorLogs` (
        `Id` bigint (20) NOT NULL AUTO_INCREMENT,
        `ExceptionMessage` longtext DEFAULT NULL,
        `InnerExceptionMessage` longtext DEFAULT NULL,
        `StackTrace` longtext DEFAULT NULL,
        `DateCreated` datetime (6) NOT NULL,
        PRIMARY KEY (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 1505 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;