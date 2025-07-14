CREATE TABLE
    `PaymentRequests` (
        `Id` bigint (20) NOT NULL AUTO_INCREMENT,
        `CustomerId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `ProductId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `DateCreated` datetime (6) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_PaymentRequests_CustomerId` (`CustomerId`),
            KEY `IX_PaymentRequests_ProductId` (`ProductId`),
            CONSTRAINT `FK_PaymentRequests_Customers_CustomerId` FOREIGN KEY (`CustomerId`) REFERENCES `Customers` (`Id`),
            CONSTRAINT `FK_PaymentRequests_Products_ProductId` FOREIGN KEY (`ProductId`) REFERENCES `Products` (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 24330 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;