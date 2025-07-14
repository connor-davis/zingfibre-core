CREATE TABLE
    `rm_specperacnt` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `srvid` int (11) NOT NULL,
        `starttime` time NOT NULL,
        `endtime` time NOT NULL,
        `timeratio` decimal(3, 2) NOT NULL,
        `dlratio` decimal(3, 2) NOT NULL,
        `ulratio` decimal(3, 2) NOT NULL,
        `connallowed` tinyint (1) NOT NULL,
        `mon` tinyint (1) NOT NULL,
        `tue` tinyint (1) NOT NULL,
        `wed` tinyint (1) NOT NULL,
        `thu` tinyint (1) NOT NULL,
        `fri` tinyint (1) NOT NULL,
        `sat` tinyint (1) NOT NULL,
        `sun` tinyint (1) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `srvid` (`srvid`),
        KEY `fromtime` (`starttime`),
        KEY `totime` (`endtime`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;