CREATE TABLE
    `rm_radacct` (
        `radacctid` bigint (21) NOT NULL,
        `acctuniqueid` varchar(32) NOT NULL,
        `username` varchar(64) NOT NULL,
        `acctstarttime` datetime NOT NULL,
        `acctstoptime` datetime NOT NULL,
        `acctsessiontime` int (12) NOT NULL,
        `acctsessiontimeratio` decimal(3, 2) NOT NULL,
        `dlbytesstart` bigint (20) NOT NULL,
        `dlbytesstop` bigint (20) NOT NULL,
        `dlbytes` bigint (20) NOT NULL,
        `dlratio` decimal(3, 2) NOT NULL,
        `ulbytesstart` bigint (20) NOT NULL,
        `ulbytesstop` bigint (20) NOT NULL,
        `ulbytes` bigint (20) NOT NULL,
        `ulratio` decimal(3, 2) NOT NULL,
        KEY `radacctid` (`radacctid`),
        KEY `acctuniqueid` (`acctuniqueid`),
        KEY `username` (`username`),
        KEY `acctstarttime` (`acctstarttime`),
        KEY `acctstoptime` (`acctstoptime`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;