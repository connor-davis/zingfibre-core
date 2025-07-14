CREATE TABLE
    `radusergroup` (
        `username` varchar(64) NOT NULL DEFAULT '',
        `groupname` varchar(64) NOT NULL DEFAULT '',
        `priority` int (11) NOT NULL DEFAULT '1',
        KEY `username` (`username`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;