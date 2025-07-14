CREATE TABLE
    `CustomerNotes` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `CustomerId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `Note` longtext DEFAULT NULL,
            `CreatedByUserId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_CustomerNotes_CreatedByUserId` (`CreatedByUserId`),
            KEY `IX_CustomerNotes_CustomerId` (`CustomerId`),
            CONSTRAINT `FK_CustomerNotes_Customers_CustomerId` FOREIGN KEY (`CustomerId`) REFERENCES `Customers` (`Id`),
            CONSTRAINT `FK_CustomerNotes_Users_CreatedByUserId` FOREIGN KEY (`CreatedByUserId`) REFERENCES `Users` (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;