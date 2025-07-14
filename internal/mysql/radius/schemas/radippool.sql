CREATE TABLE
    `radippool` (
        `id` int (11) unsigned NOT NULL AUTO_INCREMENT,
        `pool_name` varchar(30) NOT NULL,
        `framedipaddress` varchar(15) NOT NULL,
        `nasipaddress` varchar(15) NOT NULL,
        `calledstationid` varchar(30) NOT NULL,
        `callingstationid` varchar(30) NOT NULL,
        `expiry_time` datetime DEFAULT NULL,
        `username` varchar(64) NOT NULL,
        `pool_key` varchar(30) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;