CREATE TABLE
    `rm_invoices` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `invgroup` tinyint (1) NOT NULL,
        `invnum` varchar(16) NOT NULL,
        `managername` varchar(64) NOT NULL,
        `username` varchar(64) NOT NULL,
        `date` date NOT NULL,
        `bytesdl` bigint (20) NOT NULL,
        `bytesul` bigint (20) NOT NULL,
        `bytescomb` bigint (20) NOT NULL,
        `downlimit` bigint (20) NOT NULL,
        `uplimit` bigint (20) NOT NULL,
        `comblimit` bigint (20) NOT NULL,
        `time` int (11) NOT NULL,
        `uptimelimit` bigint (20) NOT NULL,
        `days` int (6) NOT NULL,
        `expiration` date NOT NULL,
        `capdl` tinyint (1) NOT NULL,
        `capul` tinyint (1) NOT NULL,
        `captotal` tinyint (1) NOT NULL,
        `captime` tinyint (1) NOT NULL,
        `capdate` tinyint (1) NOT NULL,
        `service` varchar(60) NOT NULL,
        `comment` varchar(200) NOT NULL,
        `transid` varchar(32) NOT NULL,
        `amount` decimal(13, 2) NOT NULL,
        `address` varchar(50) NOT NULL,
        `city` varchar(50) NOT NULL,
        `zip` varchar(8) NOT NULL,
        `country` varchar(50) NOT NULL,
        `state` varchar(50) NOT NULL,
        `fullname` varchar(100) NOT NULL,
        `taxid` varchar(40) NOT NULL,
        `contractid` varchar(50) NOT NULL,
        `paymentopt` date NOT NULL,
        `invtype` tinyint (1) NOT NULL,
        `paymode` tinyint (4) NOT NULL,
        `paid` date NOT NULL,
        `price` decimal(25, 6) NOT NULL,
        `tax` decimal(25, 6) NOT NULL,
        `advtax` decimal(25, 6) NOT NULL,
        `vatpercent` decimal(4, 2) NOT NULL,
        `advtaxpercent` decimal(4, 2) NOT NULL,
        `remark` varchar(400) NOT NULL,
        `balance` decimal(20, 2) NOT NULL,
        `gwtransid` varchar(255) NOT NULL,
        `phone` varchar(15) NOT NULL,
        `mobile` varchar(15) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `invnum` (`invnum`),
        KEY `username` (`username`),
        KEY `managername` (`managername`),
        KEY `date` (`date`),
        KEY `gwtransid` (`gwtransid`),
        KEY `comment` (`comment`),
        KEY `paymode` (`paymode`),
        KEY `invgroup` (`invgroup`),
        KEY `paid` (`paid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;