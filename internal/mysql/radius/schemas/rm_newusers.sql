CREATE TABLE
    `rm_newusers` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `username` varchar(64) NOT NULL,
        `firstname` varchar(50) NOT NULL,
        `lastname` varchar(50) NOT NULL,
        `address` varchar(100) NOT NULL,
        `city` varchar(50) NOT NULL,
        `zip` varchar(8) NOT NULL,
        `country` varchar(50) NOT NULL,
        `state` varchar(50) NOT NULL,
        `phone` varchar(15) NOT NULL,
        `mobile` varchar(15) NOT NULL,
        `email` varchar(100) NOT NULL,
        `vatid` varchar(40) NOT NULL,
        `srvid` int (11) NOT NULL,
        `actcode` varchar(10) NOT NULL,
        `actcount` int (11) NOT NULL,
        `lang` varchar(30) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;