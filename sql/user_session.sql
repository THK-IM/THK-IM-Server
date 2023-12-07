CREATE TABLE IF NOT EXISTS `user_session_%s`
(
    `session_id`  BIGINT  NOT NULL,
    `user_id`     BIGINT  NOT NULL,
    `name`        TEXT    NOT NULL COMMENT '名称',
    `remark`      TEXT    NOT NULL COMMENT '描述',
    `type`        INT     NOT NULL DEFAULT 1 COMMENT '1单聊/2群聊',
    `role`        INT     NOT NULL DEFAULT 1 COMMENT '4拥有者/3超级管理员/2管理员/1成员',
    `entity_id`   BIGINT  NOT NULL COMMENT '用户id/群id',
    `top`         BIGINT  NOT NULL DEFAULT 0 COMMENT '置顶时间戳',
    `mute`        INT     NOT NULL DEFAULT 0 COMMENT '2^0(全员被禁言) 2^1(自己被禁言)',
    `status`      INT     NOT NULL DEFAULT 0 COMMENT '2^1(不接收消息) 2^2(静音)',
    `parent_id`   BIGINT COMMENT '父sessionId，用于合并session',
    `ext_data`    TEXT COMMENT '扩展字段',
    `update_time` BIGINT  NOT NULL DEFAULT 0 COMMENT '更新时间',
    `create_time` BIGINT  NOT NULL DEFAULT 0 COMMENT '创建时间',
    `deleted`     TINYINT NOT NULL DEFAULT 0 COMMENT '会话删除状态',
    INDEX `USER_SESSION_U_IDX` (`user_id`),
    INDEX `USER_SESSION_TIME_IDX` (`update_time`),
    UNIQUE INDEX `USER_SESSION_IDX` (`session_id`, `user_id`, `entity_id`, `type`)
);