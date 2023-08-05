CREATE TABLE IF NOT EXISTS `session_user_%s`
(
    `session_id`  BIGINT  NOT NULL,
    `user_id`     BIGINT  NOT NULL,
    `type`        INT     NOT NULL           COMMENT '1单聊/2群聊',
    `status`      INT     NOT NULL DEFAULT 0 COMMENT 'session状态 2^0(0未禁言 1禁言) 2^1(0可以接收 1不可接收) 2^2(0提醒 1静音) ',
    `update_time` BIGINT           DEFAULT 0 COMMENT '更新时间',
    `create_time` BIGINT           DEFAULT 0 COMMENT '创建时间',
    `deleted`     TINYINT NOT NULL DEFAULT 0 COMMENT '会话删除状态',
    INDEX `SESSION_USER_S_IDX` (`session_id`),
    UNIQUE INDEX `SESSION_USER_IDX` (`session_id`, `user_id`)
);