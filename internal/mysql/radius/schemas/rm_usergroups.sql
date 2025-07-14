CREATE TABLE
    `rm_usergroups` (
        `groupid` int (11) NOT NULL AUTO_INCREMENT,
        `groupname` varchar(50) NOT NULL,
        `descr` varchar(200) NOT NULL,
        PRIMARY KEY (`groupid`)
    ) ENGINE = InnoDB AUTO_INCREMENT = 6 DEFAULT CHARSET = utf8;