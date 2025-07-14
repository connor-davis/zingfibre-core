CREATE TABLE
    `rm_colsetlistdocsis` (
        `managername` varchar(64) NOT NULL,
        `colname` varchar(32) NOT NULL,
        KEY `managername` (`managername`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;