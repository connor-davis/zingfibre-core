CREATE TABLE
    `rm_ias` (
        `iasid` int (11) NOT NULL AUTO_INCREMENT,
        `iasname` varchar(50) NOT NULL,
        `price` decimal(20, 2) NOT NULL,
        `downlimit` bigint (20) NOT NULL,
        `uplimit` bigint (20) NOT NULL,
        `comblimit` bigint (20) NOT NULL,
        `uptimelimit` bigint (20) NOT NULL,
        `expiretime` bigint (20) NOT NULL,
        `timebaseonline` tinyint (1) NOT NULL,
        `timebaseexp` tinyint (1) NOT NULL,
        `srvid` int (11) NOT NULL,
        `enableias` tinyint (1) NOT NULL,
        `expiremode` tinyint (1) NOT NULL,
        `expiration` date NOT NULL,
        `simuse` int (11) NOT NULL,
        PRIMARY KEY (`iasid`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 11 DEFAULT CHARSET = utf8;