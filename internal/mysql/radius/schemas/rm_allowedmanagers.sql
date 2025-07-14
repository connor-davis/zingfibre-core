CREATE TABLE
    `rm_allowedmanagers` (
        `srvid` int (11) NOT NULL,
        `managername` varchar(64) NOT NULL,
        KEY `srvid` (`srvid`),
        KEY `managername` (`managername`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;