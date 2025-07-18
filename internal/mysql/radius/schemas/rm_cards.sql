CREATE TABLE
    `rm_cards` (
        `id` bigint (20) NOT NULL,
        `cardnum` varchar(16) NOT NULL,
        `password` varchar(8) NOT NULL,
        `value` decimal(22, 2) NOT NULL,
        `expiration` date NOT NULL,
        `series` varchar(16) NOT NULL,
        `date` date NOT NULL,
        `owner` varchar(64) NOT NULL,
        `used` datetime NOT NULL,
        `cardtype` tinyint (1) NOT NULL,
        `revoked` tinyint (1) NOT NULL,
        `downlimit` bigint (20) NOT NULL,
        `uplimit` bigint (20) NOT NULL,
        `comblimit` bigint (20) NOT NULL,
        `uptimelimit` bigint (20) NOT NULL,
        `srvid` int (11) NOT NULL,
        `transid` varchar(32) NOT NULL,
        `active` tinyint (1) NOT NULL,
        `expiretime` bigint (20) NOT NULL,
        `timebaseexp` tinyint (1) NOT NULL,
        `timebaseonline` tinyint (1) NOT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY `cardnum` (`cardnum`),
        KEY `series` (`series`),
        KEY `used` (`used`),
        KEY `owner` (`owner`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;