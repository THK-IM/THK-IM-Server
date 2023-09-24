CREATE TABLE IF NOT EXISTS `user_online_status_%s`
(
    `user_id`     BIGINT      NOT NULL,
    `platform`    VARCHAR(10) NOT NULL,
    `conn_id`     BIGINT      NOT NULL,
    `online_time` BIGINT      NOT NULL,
    UNIQUE INDEX `USER_ONLINE_STATUS_U_IDX` (`user_id`, `platform`)
);