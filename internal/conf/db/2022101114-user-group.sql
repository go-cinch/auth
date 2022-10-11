-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user_group
(
    `id`     BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `name`   VARCHAR(50) NULL COMMENT 'name',
    `key`    VARCHAR(50) NULL COMMENT 'keyword, must be unique, used as frontend display',
    `action` LONGTEXT    NULL COMMENT 'user group action code array'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE UNIQUE INDEX idx_key ON user_group (`key`);

CREATE TABLE `user_user_group_relation`
(
    `user_id`       BIGINT UNSIGNED NOT NULL COMMENT 'auto increment id',
    `user_group_id` BIGINT UNSIGNED NOT NULL COMMENT 'auto increment id',
    PRIMARY KEY (`user_id`, `user_group_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;