CREATE TABLE
    `rm_phpsess` (
        `managername` varchar(64) NOT NULL,
        `ip` varchar(15) NOT NULL,
        `sessid` varchar(64) NOT NULL,
        `lastact` datetime NOT NULL,
        `closed` tinyint (1) DEFAULT NULL,
        KEY `managername` (`managername`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;