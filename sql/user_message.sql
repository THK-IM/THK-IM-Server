CREATE TABLE IF NOT EXISTS `user_message_%s`
(
    `msg_id`       BIGINT  NOT NULL,
    `client_id`    BIGINT  NOT NULL,
    `user_id`      BIGINT  NOT NULL,
    `session_id`   BIGINT  NOT NULL,
    `from_user_id` BIGINT  NOT NULL           COMMENT '发送者id',
    `msg_type`     INT     NOT NULL           COMMENT '消息类型',
    `msg_content`  TEXT    NOT NULL           COMMENT '消息内容',
    `at_users`     TEXT                       COMMENT '@谁, uid数据',
    `reply_msg_id` BIGINT                     COMMENT '回复消息id',
    `status`       TINYINT          DEFAULT 0 COMMENT '用户消息状态:0:默认,2^0:已经发送给用户,2^1:用户已读, 2^2:用户撤回, 2^3:重新编辑',
    `create_time`  BIGINT           DEFAULT 0 COMMENT '创建时间',
    `update_time`  BIGINT           DEFAULT 0 COMMENT '更新时间',
    `deleted`      TINYINT NOT NULL DEFAULT 0 COMMENT '消息删除状态',
    INDEX `USER_MESSAGE_U_IDX` (`user_id`),
    INDEX `USER_MESSAGE_CTIME_IDX` (`create_time`),
    UNIQUE INDEX `USER_MESSAGE_IDX` (`user_id`, `session_id`, `msg_id`),
    UNIQUE INDEX `USER_MESSAGE_CLIENT_IDX` (`session_id`, `from_user_id`, `client_id`)
);