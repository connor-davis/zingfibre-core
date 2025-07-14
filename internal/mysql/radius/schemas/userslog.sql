CREATE TABLE
    `userslog` (
        `id` int (10) unsigned NOT NULL AUTO_INCREMENT,
        `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `username` varchar(255) NOT NULL,
        `oldexpiration` datetime DEFAULT NULL,
        `expiration` datetime DEFAULT NULL,
        `srvid` varchar(255) DEFAULT NULL,
        `oldsrvid` varchar(255) DEFAULT NULL,
        `action` varchar(255) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 8146 DEFAULT CHARSET = latin1;