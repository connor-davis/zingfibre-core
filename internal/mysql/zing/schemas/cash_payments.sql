CREATE TABLE
    `CashPayments` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `PaymentCode` bigint (20) NOT NULL AUTO_INCREMENT,
            `CustomerId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `ProductId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `RechargeId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `DateCompleted` datetime (6) DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            UNIQUE KEY `PaymentCode_UNIQUE` (`PaymentCode`),
            UNIQUE KEY `Id_UNIQUE` (`Id`),
            KEY `IX_CashPayments_CustomerId` (`CustomerId`),
            KEY `IX_CashPayments_ProductId` (`ProductId`),
            KEY `IX_CashPayments_RechargeId` (`RechargeId`),
            CONSTRAINT `FK_CashPayments_Customers_CustomerId` FOREIGN KEY (`CustomerId`) REFERENCES `Customers` (`Id`),
            CONSTRAINT `FK_CashPayments_Products_ProductId` FOREIGN KEY (`ProductId`) REFERENCES `Products` (`Id`),
            CONSTRAINT `FK_CashPayments_Recharges_RechargeId` FOREIGN KEY (`RechargeId`) REFERENCES `Recharges` (`Id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 8383 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;