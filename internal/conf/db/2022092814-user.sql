-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE user
(
    id           BIGINT UNSIGNED AUTO_INCREMENT COMMENT 'auto increment id' PRIMARY KEY,
    created_at   DATETIME(3)          NULL COMMENT 'create time',
    updated_at   DATETIME(3)          NULL COMMENT 'update time',
    username     VARCHAR(191)         NULL COMMENT 'user login name',
    user_code    CHAR(8)              NOT NULL COMMENT 'user code',
    password     LONGTEXT             NULL COMMENT 'password',
    mobile       LONGTEXT             NULL COMMENT 'mobile number',
    avatar       LONGTEXT             NULL COMMENT 'avatar url',
    nickname     LONGTEXT             NULL COMMENT 'nickname',
    introduction LONGTEXT             NULL COMMENT 'introduction',
    status       TINYINT(1) DEFAULT 1 NULL COMMENT 'status(0: disabled, 1: enable)',
    last_login   DATETIME(3)          NULL COMMENT 'last login time',
    locked       TINYINT(1) DEFAULT 0 NULL COMMENT 'locked(0: unlock, 1: locked)',
    lock_expire  BIGINT UNSIGNED      NULL COMMENT 'lock expiration time',
    wrong        BIGINT UNSIGNED      NULL COMMENT 'type wrong password count'
);

CREATE UNIQUE INDEX idx_username ON user (username);
CREATE UNIQUE INDEX idx_user_code ON user (user_code);