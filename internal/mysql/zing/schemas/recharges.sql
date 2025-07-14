CREATE TABLE
    `Recharges` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `CustomerId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `ProductId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `Method` longtext DEFAULT NULL,
            `PaymentServicePaymentId` longtext DEFAULT NULL,
            `PaymentServicePayload` longtext DEFAULT NULL,
            `PaymentServiceQueryParams` longtext DEFAULT NULL,
            `RechargeSuccessful` tinyint (1) NOT NULL,
            `FailureReason` longtext DEFAULT NULL,
            `PaymentAmount` decimal(18, 2) DEFAULT NULL,
            `ExpiryDate` datetime (6) DEFAULT NULL,
            `PreviousRMExpiryDate` datetime (6) DEFAULT NULL,
            `UserId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `FromRMSvcID` int (11) DEFAULT NULL,
            `ToRMSvcID` int (11) DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_Recharges_CustomerId` (`CustomerId`),
            KEY `IX_Recharges_ProductId` (`ProductId`),
            KEY `IX_Recharges_UserId` (`UserId`),
            CONSTRAINT `FK_Recharges_Customers_CustomerId` FOREIGN KEY (`CustomerId`) REFERENCES `Customers` (`Id`),
            CONSTRAINT `FK_Recharges_Products_ProductId` FOREIGN KEY (`ProductId`) REFERENCES `Products` (`Id`),
            CONSTRAINT `FK_Recharges_Users_UserId` FOREIGN KEY (`UserId`) REFERENCES `Users` (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;