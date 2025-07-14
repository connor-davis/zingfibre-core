CREATE TABLE
    `rm_syslog` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `datetime` datetime NOT NULL,
        `ip` varchar(15) NOT NULL,
        `name` varchar(64) NOT NULL,
        `eventid` int (11) NOT NULL,
        `data1` varchar(64) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 20119 DEFAULT CHARSET = utf8;