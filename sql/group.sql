CREATE TABLE IF NOT EXISTS `group_%s`
(
    `id`           BIGINT       NOT NULL COMMENT '群id',
    `s_id`         BIGINT       NOT NULL COMMENT '会话id',
    `owner_id`     BIGINT       NOT NULL COMMENT '群主',
    `name`         varchar(100) NOT NULL COMMENT '群名称',
    `avatar`       TEXT COMMENT '群头像',
    `announce`     TEXT COMMENT '公告',
    `qr_code`      TEXT COMMENT '二维码',
    `members`      INT COMMENT  '群成员数量',
    `enter_flag`   INT          NOT NULL default 0 COMMENT '进群条件，1 扫码或通过群id随意进群，2 申请通过后可以进入 4 管理员邀请 ',
    `update_time`  BIGINT       NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    `create_time`  BIGINT       NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `Group_IDX` (`id`)
);