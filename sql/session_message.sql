CREATE TABLE IF NOT EXISTS `session_message_%s`
(
    `msg_id`       BIGINT  NOT NULL,
    `client_id`    BIGINT  NOT NULL,
    `session_id`   BIGINT  NOT NULL,
    `from_user_id` BIGINT  NOT NULL           COMMENT '发送者id',
    `msg_type`     INT     NOT NULL           COMMENT '消息类型',
    `msg_content`  TEXT    NOT NULL           COMMENT '消息内容',
    `at_users`     TEXT                       COMMENT '@谁, uid数据',
    `reply_msg_id` BIGINT                     COMMENT '回复消息id',
    `create_time`  BIGINT           DEFAULT 0 COMMENT '创建时间',
    `update_time`  BIGINT           DEFAULT 0 COMMENT '更新时间',
    `deleted`      TINYINT NOT NULL DEFAULT 0 COMMENT '消息删除状态',
    INDEX `SESSION_MESSAGE_S_IDX` (`session_id`),
    UNIQUE INDEX `SESSION_MESSAGE_IDX` (`session_id`, `msg_id`),
    UNIQUE INDEX `SESSION_CLIENT_MESSAGE_IDX` (`session_id`, `from_user_id`, `client_id`)
);