CREATE TABLE
    `rm_mpesa` (
        `merchantrequestid` varchar(50) NOT NULL,
        `checkoutrequestid` varchar(50) NOT NULL,
        `resultcode` tinyint (4) NOT NULL,
        `resultdesc` varchar(100) NOT NULL,
        `amount` int (11) NOT NULL,
        `receiptnumber` varchar(20) NOT NULL,
        `transactiondate` datetime NOT NULL,
        `phonenumber` varchar(15) NOT NULL,
        UNIQUE KEY `merchantrequestid` (`merchantrequestid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;