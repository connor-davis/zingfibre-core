CREATE TABLE
    `rm_actsrv` (
        `id` bigint (20) NOT NULL AUTO_INCREMENT,
        `datetime` datetime NOT NULL,
        `username` varchar(64) NOT NULL,
        `srvid` int (11) NOT NULL,
        `dailynextsrvactive` tinyint (1) NOT NULL,
        UNIQUE KEY `id` (`id`),
        KEY `datetime` (`datetime`),
        KEY `username` (`username`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 1311748 DEFAULT CHARSET = utf8;