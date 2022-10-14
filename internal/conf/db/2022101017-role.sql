-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE role
(
    `id`     BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `name`   VARCHAR(50)          NULL COMMENT 'name',
    `word`   VARCHAR(50)          NULL COMMENT 'keyword, must be unique, used as frontend display',
    `status` TINYINT(1) DEFAULT 1 NULL COMMENT 'status(0: disabled, 1: enable)',
    `action` LONGTEXT             NULL COMMENT 'role action code array'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE UNIQUE INDEX idx_word ON role (`word`);