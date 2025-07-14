CREATE TABLE
    `rm_cmts` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `ip` varchar(15) NOT NULL,
        `name` varchar(32) NOT NULL,
        `community` varchar(32) NOT NULL,
        `descr` varchar(200) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `ip` (`ip`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;