CREATE TABLE
    `rm_ap` (
        `id` int (11) NOT NULL AUTO_INCREMENT,
        `name` varchar(32) NOT NULL,
        `enable` tinyint (1) NOT NULL,
        `accessmode` tinyint (1) NOT NULL,
        `ip` varchar(15) NOT NULL,
        `community` varchar(32) NOT NULL,
        `apiusername` varchar(32) NOT NULL,
        `apipassword` varchar(32) NOT NULL,
        `apiver` tinyint (1) NOT NULL,
        `description` varchar(200) NOT NULL,
        PRIMARY KEY (`id`),
        KEY `ip` (`ip`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;