CREATE TABLE IF NOT EXISTS `session_%s`
(
    `id`          BIGINT PRIMARY KEY NOT NULL,
    `name`        TEXT               NOT NULL COMMENT '名称',
    `remark`      TEXT               NOT NULL COMMENT '描述',
    `type`        INT                NOT NULL COMMENT '1单聊/2群聊',
    `status`      INT                NOT NULL DEFAULT 0 COMMENT '会话状态',
    `update_time` BIGINT                      DEFAULT 0 COMMENT '更新时间',
    `create_time` BIGINT                      DEFAULT 0 COMMENT '创建时间',
    `deleted`     TINYINT            NOT NULL DEFAULT 0 COMMENT '会话删除状态',
    UNIQUE INDEX `SESSION_IDX` (`id`)
);