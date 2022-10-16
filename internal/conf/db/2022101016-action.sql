-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE action
(
    `id`       BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `name`     VARCHAR(50) NULL COMMENT 'name',
    `code`     CHAR(8)     NOT NULL COMMENT 'code',
    `word`     VARCHAR(50) NULL COMMENT 'keyword, must be unique, used as frontend display',
    `resource` LONGTEXT    NULL COMMENT 'resource array, split by break line str, example: GET,/user+\n+POST,/role+\n+GET,/action',
    `menu`     LONGTEXT    NULL COMMENT 'menu array, split by break line str',
    `btn`      LONGTEXT    NULL COMMENT 'btn array, split by break line str'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

CREATE UNIQUE INDEX idx_word ON action (`word`);
CREATE UNIQUE INDEX idx_code ON action (`code`);