CREATE TABLE
    `rm_wlantmp` (
        `maccpe` varchar(17) DEFAULT NULL,
        `uspwr` smallint (6) DEFAULT NULL,
        `ccq` smallint (6) DEFAULT NULL,
        `snr` smallint (6) DEFAULT NULL,
        `apip` varchar(15) DEFAULT NULL,
        `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        KEY `maccpe` (`maccpe`)
    ) ENGINE = MEMORY DEFAULT CHARSET = utf8;