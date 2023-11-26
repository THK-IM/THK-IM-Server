CREATE TABLE IF NOT EXISTS `object_%s`
(
    `id`          BIGINT      NOT NULL COMMENT '对象id',
    `s_id`        BIGINT      NOT NULL COMMENT '所属会话id',
    `engine`      varchar(10) NOT NULL COMMENT '存储引擎 minio/oss/obs/s3',
    `key`         text        NOT NULL COMMENT '对象key',
    `create_time` BIGINT      NOT NULL DEFAULT 0 COMMENT '创建时间 毫秒',
    UNIQUE INDEX `Object_IDX` (`id`, `s_id`)
);