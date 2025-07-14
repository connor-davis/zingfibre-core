CREATE TABLE
    `nas` (
        `id` int (10) NOT NULL AUTO_INCREMENT,
        `nasname` varchar(128) NOT NULL,
        `shortname` varchar(32) DEFAULT NULL,
        `type` varchar(30) DEFAULT 'other',
        `ports` int (5) DEFAULT NULL,
        `secret` varchar(60) NOT NULL DEFAULT 'secret',
        `community` varchar(50) DEFAULT NULL,
        `description` varchar(200) DEFAULT 'RADIUS Client',
        `starospassword` varchar(32) NOT NULL,
        `ciscobwmode` tinyint (1) NOT NULL,
        `apiusername` varchar(32) NOT NULL,
        `apipassword` varchar(32) NOT NULL,
        `apiver` tinyint (1) NOT NULL,
        `coamode` tinyint (1) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `nasname` (`nasname`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 37 DEFAULT CHARSET = utf8;