CREATE TABLE
    `rm_changesrv` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `username` varchar(64) NOT NULL,
        `newsrvid` int (11) NOT NULL,
        `newsrvname` varchar(50) NOT NULL,
        `scheduledate` date NOT NULL,
        `requestdate` date NOT NULL,
        `status` tinyint (1) NOT NULL,
        `transid` varchar(32) NOT NULL,
        `requested` varchar(64) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `requestdate` (`requestdate`),
        KEY `scheduledate` (`scheduledate`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 2113 DEFAULT CHARSET = utf8;