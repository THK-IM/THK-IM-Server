CREATE TABLE IF NOT EXISTS `group_member_apply_%s`
(
    `id`          BIGINT NOT NULL COMMENT '申请id',
    `group_id`    BIGINT NOT NULL COMMENT '群id',
    `member_id`   BIGINT NOT NULL COMMENT '成员id',
    `channel`     INT    NOT NULL COMMENT '渠道:1群号,2群二维码,3分享',
    `content`     TEXT   NOT NULL COMMENT '申请内容',
    `status`      INT    NOT NULL DEFAULT 0 COMMENT '申请状态 0 申请中，1 拒绝，2 通过',
    `update_time` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `Group_IDX` (`id`)
);