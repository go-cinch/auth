-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE role
(
    `id`     BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    `name`   VARCHAR(50)          NULL COMMENT 'name',
    `key`    VARCHAR(50)          NULL COMMENT 'keyword, must be unique, used as frontend display',
    `status` TINYINT(1) DEFAULT 1 NULL COMMENT 'status(0: disabled, 1: enable)'
);

CREATE UNIQUE INDEX idx_key ON role (`key`);