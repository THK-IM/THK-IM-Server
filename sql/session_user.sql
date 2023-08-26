CREATE TABLE IF NOT EXISTS `session_user_%s`
(
    `session_id`  BIGINT  NOT NULL,
    `user_id`     BIGINT  NOT NULL,
    `type`        INT     NOT NULL DEFAULT 1 COMMENT '1单聊/2群聊',
    `role`        INT     NOT NULL DEFAULT 1 COMMENT '4拥有者/3超级管理员/2管理员/1成员',
    `mute`        INT     NOT NULL DEFAULT 0 COMMENT '2^0(全员被禁言) 2^1(自己被禁言)',
    `status`      INT     NOT NULL DEFAULT 0 COMMENT '2^1(不接收消息) 2^2(静音)',
    `update_time` BIGINT  NOT NULL DEFAULT 0 COMMENT '更新时间',
    `create_time` BIGINT  NOT NULL DEFAULT 0 COMMENT '创建时间',
    `deleted`     TINYINT NOT NULL DEFAULT 0 COMMENT '会话删除状态',
    INDEX `SESSION_USER_S_IDX` (`session_id`),
    UNIQUE INDEX `SESSION_USER_IDX` (`session_id`, `user_id`, `type`)
);