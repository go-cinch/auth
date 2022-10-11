-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE action
(
    `id`   BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `name` VARCHAR(50) NULL COMMENT 'name',
    `code` CHAR(8)     NOT NULL COMMENT 'code',
    `key`  VARCHAR(50) NULL COMMENT 'keyword, must be unique, used as frontend display',
    `path` LONGTEXT    NULL COMMENT 'url path array, split by break line str, example: GET,/user+\n+POST,/role+\n+GET,/action'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE UNIQUE INDEX idx_key ON action (`key`);
CREATE UNIQUE INDEX idx_code ON action (`code`);