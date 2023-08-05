CREATE TABLE IF NOT EXISTS `user_session_%s`
(
    `session_id`  BIGINT  NOT NULL,
    `user_id`     BIGINT  NOT NULL,
    `type`        INT     NOT NULL           COMMENT '1单聊/2群聊',
    `entity_id`   BIGINT  NOT NULL           COMMENT '用户id/群id',
    `top`         BIGINT  NOT NULL DEFAULT 0 COMMENT '置顶时间戳',
    `status`      INT     NOT NULL DEFAULT 0 COMMENT 'session状态 2^0(0未禁言 1禁言) 2^1(0可以接收 1不可接收) 2^2(0提醒 1静音)',
    `update_time` BIGINT           DEFAULT 0 COMMENT '更新时间',
    `create_time` BIGINT           DEFAULT 0 COMMENT '创建时间',
    `deleted`     TINYINT NOT NULL DEFAULT 0 COMMENT '会话删除状态',
    INDEX `USER_SESSION_U_IDX` (`user_id`),
    INDEX `USER_SESSION_TIME_IDX` (`update_time`),
    UNIQUE INDEX `USER_SESSION_IDX` (`session_id`, `user_id`),
    UNIQUE INDEX `USER_SESSION_ENTITY_TYPE` (`user_id`, `entity_id`, `type`)
);