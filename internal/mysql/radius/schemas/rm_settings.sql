CREATE TABLE
    `rm_settings` (
        `currency` varchar(15) NOT NULL,
        `unixacc` tinyint (1) NOT NULL,
        `diskquota` tinyint (1) NOT NULL,
        `quotatpl` varchar(30) NOT NULL,
        `paymentopt` int (11) NOT NULL,
        `changesrv` tinyint (1) NOT NULL,
        `vatpercent` decimal(4, 2) NOT NULL,
        `advtaxpercent` decimal(4, 2) NOT NULL,
        `disablenotpaid` tinyint (1) NOT NULL,
        `disableexpcont` tinyint (1) NOT NULL,
        `resetctr` tinyint (1) NOT NULL,
        `newnasallsrv` tinyint (1) NOT NULL,
        `newmanallsrv` tinyint (1) NOT NULL,
        `disconnmethod` tinyint (1) NOT NULL,
        `warndl` bigint (20) NOT NULL,
        `warndlpercent` int (3) NOT NULL,
        `warnul` bigint (20) NOT NULL,
        `warnulpercent` int (3) NOT NULL,
        `warncomb` bigint (20) NOT NULL,
        `warncombpercent` int (3) NOT NULL,
        `warnuptime` bigint (20) NOT NULL,
        `warnuptimepercent` int (3) NOT NULL,
        `warnexpiry` int (11) NOT NULL,
        `expalertmode` tinyint (1) NOT NULL,
        `emailselfregman` tinyint (1) NOT NULL,
        `emailwelcome` tinyint (1) NOT NULL,
        `emailnewsrv` tinyint (1) NOT NULL,
        `smsnewsrv` tinyint (1) NOT NULL,
        `emailrenew` tinyint (1) NOT NULL,
        `smsrenew` tinyint (1) NOT NULL,
        `emailexpiry` tinyint (1) NOT NULL,
        `smswelcome` tinyint (1) NOT NULL,
        `smsexpiry` tinyint (1) NOT NULL,
        `warnmode` tinyint (1) NOT NULL,
        `selfreg` tinyint (1) NOT NULL,
        `edituserdata` tinyint (1) NOT NULL,
        `hidelimits` tinyint (1) NOT NULL,
        `pm_internal` tinyint (1) NOT NULL,
        `pm_paypalstd` tinyint (1) NOT NULL,
        `pm_paypalpro` tinyint (1) NOT NULL,
        `pm_paypalexp` tinyint (1) NOT NULL,
        `pm_netcash` tinyint (1) NOT NULL,
        `pm_authorizenet` tinyint (1) NOT NULL,
        `pm_dps` tinyint (1) NOT NULL,
        `pm_2co` tinyint (1) NOT NULL,
        `pm_payfast` tinyint (1) NOT NULL,
        `pm_payu` tinyint (1) NOT NULL,
        `pm_paytm` tinyint (1) NOT NULL,
        `pm_bkash` tinyint (1) NOT NULL,
        `pm_flutterwave` tinyint (1) NOT NULL,
        `pm_easypay` tinyint (1) NOT NULL,
        `pm_mpesa` tinyint (1) NOT NULL,
        `pm_custom` tinyint (1) NOT NULL,
        `unixhost` tinyint (1) NOT NULL,
        `remotehostname` varchar(100) NOT NULL,
        `maclock` tinyint (1) NOT NULL,
        `billingstart` tinyint (2) NOT NULL,
        `disconnpostpaid` tinyint (1) NOT NULL,
        `renewday` tinyint (2) NOT NULL,
        `changepswucp` tinyint (1) NOT NULL,
        `redeemucp` tinyint (1) NOT NULL,
        `buycreditsucp` tinyint (1) NOT NULL,
        `buydepositucp` tinyint (1) NOT NULL,
        `reg_firstname` tinyint (1) NOT NULL,
        `reg_lastname` tinyint (1) NOT NULL,
        `reg_address` tinyint (1) NOT NULL,
        `reg_city` tinyint (1) NOT NULL,
        `reg_zip` tinyint (1) NOT NULL,
        `reg_country` tinyint (1) NOT NULL,
        `reg_state` tinyint (1) NOT NULL,
        `reg_phone` tinyint (1) NOT NULL,
        `reg_mobile` tinyint (1) NOT NULL,
        `reg_email` tinyint (1) NOT NULL,
        `reg_vatid` tinyint (1) NOT NULL,
        `reg_cnic` tinyint (1) NOT NULL,
        `selfreg_firstname` tinyint (1) NOT NULL,
        `selfreg_lastname` tinyint (1) NOT NULL,
        `selfreg_address` tinyint (1) NOT NULL,
        `selfreg_city` tinyint (1) NOT NULL,
        `selfreg_zip` tinyint (1) NOT NULL,
        `selfreg_country` tinyint (1) NOT NULL,
        `selfreg_state` tinyint (1) NOT NULL,
        `selfreg_phone` tinyint (1) NOT NULL,
        `selfreg_mobile` tinyint (1) NOT NULL,
        `selfreg_email` tinyint (1) NOT NULL,
        `selfreg_mobactsms` tinyint (1) NOT NULL,
        `selfreg_nameactemail` tinyint (1) NOT NULL,
        `selfreg_nameactsms` tinyint (1) NOT NULL,
        `selfreg_endupemail` tinyint (1) NOT NULL,
        `selfreg_endupmobile` tinyint (1) NOT NULL,
        `selfreg_vatid` tinyint (1) NOT NULL,
        `ias_email` tinyint (1) NOT NULL,
        `ias_mobile` tinyint (1) NOT NULL,
        `ias_verify` tinyint (1) NOT NULL,
        `ias_endupemail` tinyint (1) NOT NULL,
        `ias_endupmobile` tinyint (1) NOT NULL,
        `simuseselfreg` int (11) NOT NULL,
        `defgrpid` int (11) NOT NULL,
        `captcha` tinyint (1) NOT NULL,
        `discontime` time NOT NULL
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8;