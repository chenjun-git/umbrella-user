CREATE TABLE `user_token` (
  `token`      VARCHAR(255) NOT NULL,
  `name`       VARCHAR(50)  NOT NULL,
  `user_id`    VARCHAR(32)  NOT NULL,
  `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`token`),
  UNIQUE KEY `name_user` (`name`, `user_id`),
  INDEX `user_id` (`user_id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 2
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;