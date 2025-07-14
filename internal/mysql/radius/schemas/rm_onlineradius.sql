CREATE TABLE
    `rm_onlineradius` (
        `rtt` decimal(11, 1) DEFAULT NULL,
        `loss` int (3) DEFAULT NULL,
        `username` varchar(64) NOT NULL DEFAULT '',
        `cid` varchar(50) NOT NULL DEFAULT '',
        `cpeip` varchar(15) NOT NULL DEFAULT '',
        KEY `cpeip` (`cpeip`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;