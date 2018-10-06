CREATE TABLE `access_token` (
  id                       bigint unsigned NOT NULL AUTO_INCREMENT                                COMMENT 'id',
  access_token             VARCHAR(255)    NOT NULL                                               COMMENT 'access token',
  refresh_token            VARCHAR(255)    NOT NULL                                               COMMENT 'refresh token',
  user_id                  VARCHAR(32)     NOT NULL                                               COMMENT '标准对象user id',
  device                   VARCHAR(128)    NOT NULL                                               COMMENT 'access token 给哪个source生成的',
  app                      VARCHAR(64)     NOT NULL                                               COMMENT 'access token 给哪个App生成的',
  access_token_expired_at  TIMESTAMP       NULL                                                   COMMENT 'access token过期时间',
  refresh_token_expired_at TIMESTAMP       NULL                                                   COMMENT 'refresh token过期时间',
  created_at               TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP                     COMMENT '创建时间',
  PRIMARY KEY (id),
  UNIQUE KEY idx_access_token (access_token),
  UNIQUE KEY idx_refresh_token (refresh_token),
  UNIQUE KEY idx_user_app_source (user_id, app, device),
  INDEX idx_user_access_token_expired_at (user_id, access_token_expired_at),
  INDEX idx_user_refresh_token_expired_at (user_id, refresh_token_expired_at)
)
ENGINE = InnoDB
AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin;