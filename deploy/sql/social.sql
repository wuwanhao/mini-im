CREATE TABLE `friends`
(
    `id`         int(11) unsigned                       NOT NULL AUTO_INCREMENT,
    `user_id`    varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '本用户ID',
    `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '好友ID',
    `remark`     varchar(255)                                    DEFAULT NULL COMMENT '备注',
    `add_source` tinyint COLLATE utf8mb4_unicode_ci              DEFAULT NULL COMMENT '添加来源',
    `created_at` int(11)                                NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='好友列表';

CREATE TABLE `friend_requests`
(
    `id`            int(11) unsigned                       NOT NULL AUTO_INCREMENT,
    `user_id`       varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '要添加的好友',
    `req_uid`       varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '发起人',
    `req_msg`       varchar(255)                                    DEFAULT NULL COMMENT '验证信息',
    `req_time`      int(11)                                NOT NULL DEFAULT 0 COMMENT '请求时间',
    `handle_result` tinyint COLLATE utf8mb4_unicode_ci              DEFAULT NULL COMMENT '处理结果',
    `handle_msg`    varchar(255)                                    DEFAULT NULL COMMENT '处理信息',
    `handled_at`    int(11)                                NOT NULL DEFAULT 0 COMMENT '处理时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='好友请求列表';

CREATE TABLE `groups`
(
    `id`               varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL,
    `name`             varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群名',
    `icon`             varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群头像',
    `status`           tinyint COLLATE utf8mb4_unicode_ci               DEFAULT NULL COMMENT '状态',
    `creator_uid`      varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL COMMENT '创建者ID',
    `group_type`       int(11)                                 NOT NULL COMMENT '群类型',
    `is_verify`        boolean                                 NOT NULL COMMENT '是否需要验证',
    `notification`     varchar(255)                                     DEFAULT NULL COMMENT '系统通知',
    `notification_uid` varchar(64)                                      DEFAULT NULL COMMENT '系统通知接收者',
    `created_at`       int(11)                                 NOT NULL DEFAULT 0 COMMENT '创建时间',
    `updated_at`       int(11)                                 NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='群组';

CREATE TABLE `group_members`
(
    `id`           int(11) unsigned                       NOT NULL AUTO_INCREMENT,
    `group_id`     varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群ID',
    `user_id`      varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户ID',
    `role_level`   tinyint COLLATE utf8mb4_unicode_ci     NOT NULL COMMENT '角色等级',
    `join_time`    int(11)                                NOT NULL DEFAULT 0 COMMENT '入群时间',
    `join_source`  tinyint COLLATE utf8mb4_unicode_ci              DEFAULT NULL COMMENT '入群方式',
    `inviter_uid`  varchar(64)                                     DEFAULT NULL COMMENT '邀请者ID',
    `operator_uid` varchar(64)                                     DEFAULT NULL COMMENT '操作者ID',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='群-群成员关系映射表';

CREATE TABLE `group_requests`
(
    `id`              int(11) unsigned                       NOT NULL AUTO_INCREMENT,
    `req_id`          varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求ID',
    `group_id`        varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标群ID',
    `req_msg`         varchar(255)                                    DEFAULT NULL COMMENT '申请信息',
    `req_time`        int(11)                                NOT NULL DEFAULT 0 COMMENT '申请时间',
    `join_source`     tinyint COLLATE utf8mb4_unicode_ci              DEFAULT NULL COMMENT '入群方式',
    `inviter_user_id` varchar(64)                                     DEFAULT NULL COMMENT '邀请人 ID',
    `handle_user_id`  varchar(64)                                     DEFAULT NULL COMMENT '处理人 ID',
    `handle_time`     int(11)                                NOT NULL DEFAULT 0 COMMENT '处理时间',
    `handle_result`   tinyint COLLATE utf8mb4_unicode_ci              DEFAULT NULL COMMENT '处理结果',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='入群请求列表';

