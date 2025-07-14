CREATE TABLE
    `rm_allowednases` (
        `srvid` int (11) NOT NULL,
        `nasid` int (11) NOT NULL,
        KEY `srvid` (`srvid`),
        KEY `nasid` (`nasid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;