CREATE TABLE
    `rm_ippools` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `type` tinyint (1) NOT NULL,
        `name` varchar(32) NOT NULL,
        `fromip` varchar(15) NOT NULL,
        `toip` varchar(15) NOT NULL,
        `descr` varchar(200) NOT NULL,
        `nextpoolid` int (11) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `name` (`name`),
        KEY `nextid` (`nextpoolid`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 23 DEFAULT CHARSET = utf8;