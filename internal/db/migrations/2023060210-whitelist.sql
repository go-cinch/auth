-- +migrate Up
CREATE TABLE `whitelist`
(
    `id`       BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `category` SMALLINT UNSIGNED NOT NULL COMMENT 'category(0:permission, 1:jwt, 2:idempotent)',
    `resource` LONGTEXT NOT NULL COMMENT 'resource array, split by break line str, example: GET|/user+\n+PUT,PATCH|/role/*+\n+GET|/action'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

-- +migrate Down
DROP TABLE `whitelist`;
