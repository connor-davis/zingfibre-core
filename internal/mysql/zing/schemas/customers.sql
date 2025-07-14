CREATE TABLE
    `Customers` (
        `Id` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci NOT NULL,
            `FirstName` longtext DEFAULT NULL,
            `Surname` longtext DEFAULT NULL,
            `Password` longtext DEFAULT NULL,
            `PasswordSalt` longtext DEFAULT NULL,
            `Email` longtext DEFAULT NULL,
            `PhoneNumber` longtext DEFAULT NULL,
            `IdNumber` longtext DEFAULT NULL,
            `RadiusUsername` longtext DEFAULT NULL,
            `PreferEmailCommunication` tinyint (1) NOT NULL,
            `Language` longtext DEFAULT NULL,
            `RegistrationApproved` tinyint (1) NOT NULL,
            `RegistrationDeclined` tinyint (1) NOT NULL,
            `SetOwnPassword` tinyint (1) NOT NULL,
            `SubscriptionToken` longtext DEFAULT NULL,
            `ProofOfAddressDocumentId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `IDBookDocumentId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `ApprovedByUserId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `AddressId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `PotentialAddress` longtext DEFAULT NULL,
            `SalesAgentId` char(36) CHARACTER
        SET
            ascii COLLATE ascii_general_ci DEFAULT NULL,
            `DateCreated` datetime (6) NOT NULL,
            `Deleted` tinyint (1) NOT NULL,
            PRIMARY KEY (`Id`),
            KEY `IX_Customers_AddressId` (`AddressId`),
            KEY `IX_Customers_ApprovedByUserId` (`ApprovedByUserId`),
            KEY `IX_Customers_IDBookDocumentId` (`IDBookDocumentId`),
            KEY `IX_Customers_ProofOfAddressDocumentId` (`ProofOfAddressDocumentId`),
            KEY `IX_Customers_SalesAgentId` (`SalesAgentId`),
            CONSTRAINT `FK_Customers_Addresses_AddressId` FOREIGN KEY (`AddressId`) REFERENCES `Addresses` (`Id`),
            CONSTRAINT `FK_Customers_Documents_IDBookDocumentId` FOREIGN KEY (`IDBookDocumentId`) REFERENCES `Documents` (`Id`),
            CONSTRAINT `FK_Customers_Documents_ProofOfAddressDocumentId` FOREIGN KEY (`ProofOfAddressDocumentId`) REFERENCES `Documents` (`Id`),
            CONSTRAINT `FK_Customers_SalesAgents_SalesAgentId` FOREIGN KEY (`SalesAgentId`) REFERENCES `SalesAgents` (`Id`),
            CONSTRAINT `FK_Customers_Users_ApprovedByUserId` FOREIGN KEY (`ApprovedByUserId`) REFERENCES `Users` (`Id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;