CREATE TABLE `user` (
  `user_id` VARCHAR(32) NOT NULL COMMENT '用户id',
  `user_name` VARCHAR(128) NOT NULL COMMENT '用户名',
  `alias` varchar(128) COLLATE utf8mb4_unicode_ci COMMENT '别名',
  `user_passwd` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `hashed_passwd` char(80) NOT NULL COMMENT '密码',
  `passwd_level` int NOT NULL DEFAULT 0 NULL COMMENT '密码强度',
  `last_update_passwd` datetime NOT NULL DEFAULT '2018-06-01 00:00:01'  COMMENT '最近修改密码的时间',
  `register_source` varchar(256) NOT NULL DEFAULT '' COMMENT '注册来源',
  `portrait` varchar(258) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '头像',
  `tel` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '电话',
  `email` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱',
  `qq` varchar(16) COLLATE utf8mb4_unicode_ci COMMENT 'qq',
  `is_member` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否会员',
  `member_start_time` datetime NOT NULL DEFAULT '2018-06-01 00:00:01' COMMENT '会员开始时间',
  `member_duration` smallint(6) NOT NULL DEFAULT 0 COMMENT '会员时长',
  `role` tinyint(1) NOT NULL DEFAULT 0 COMMENT '角色',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (user_id),
  UNIQUE KEY (`tel`),
  UNIQUE KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';

