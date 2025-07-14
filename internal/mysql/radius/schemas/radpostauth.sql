CREATE TABLE
    `radpostauth` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `username` varchar(64) NOT NULL DEFAULT '',
        `pass` varchar(64) NOT NULL DEFAULT '',
        `reply` varchar(32) NOT NULL DEFAULT '',
        `authdate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `nasipaddress` varchar(15) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `username` (`username`),
        KEY `authdate` (`authdate`),
        KEY `nasipaddress` (`nasipaddress`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 1065805 DEFAULT CHARSET = utf8;