CREATE TABLE
    `radreply` (
        `id` int (11) unsigned NOT NULL AUTO_INCREMENT,
        `username` varchar(64) NOT NULL DEFAULT '',
        `attribute` varchar(64) NOT NULL DEFAULT '',
        `op` char(2) NOT NULL DEFAULT '=',
        `value` varchar(253) NOT NULL DEFAULT '',
        PRIMARY KEY (`id`),
        KEY `username` (`username`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;