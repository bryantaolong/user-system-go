-- MySQL 版本表结构初始化脚本
-- 注意：MySQL 8.0+ 支持以下语法

-- sys_user 用户表
CREATE TABLE IF NOT EXISTS `sys_user`
(
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，主键，自增长',
    username            VARCHAR(255) NOT NULL COMMENT '用户名，用于登录的唯一标识',
    `password`          VARCHAR(255) NOT NULL COMMENT '加密后的用户密码',
    phone               VARCHAR(50) COMMENT '用户手机号码',
    email               VARCHAR(255) COMMENT '用户电子邮箱',
    `status`            INT DEFAULT 0 COMMENT '用户状态(0-正常 1-禁用 2-锁定)',
    roles               VARCHAR(255) COMMENT '用户角色，多个角色用逗号分隔',
    last_login_at       DATETIME COMMENT '最后一次登录时间',
    last_login_ip       VARCHAR(255) COMMENT '最后一次登录IP地址',
    last_login_device   VARCHAR(255) COMMENT '最后一次登录设备信息',
    password_reset_at   DATETIME COMMENT '密码重置时间',
    login_fail_count    INT DEFAULT 0 COMMENT '登录失败次数',
    locked_at           DATETIME COMMENT '账户锁定时间',
    deleted             INT DEFAULT 0 NOT NULL COMMENT '软删除标记(0-未删除 1-已删除)',
    `version`           INT DEFAULT 0 NOT NULL COMMENT '乐观锁版本号',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '记录创建时间',
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    created_by          VARCHAR(255) COMMENT '记录创建人',
    updated_by          VARCHAR(255) COMMENT '记录更新人'
) COMMENT='用户表，存储系统用户的基本信息、认证信息和状态';

-- 用户名索引
CREATE INDEX idx_user_username ON `sys_user` (username);

-- user_role 用户角色表
CREATE TABLE IF NOT EXISTS `user_role`
(
    id          INT NOT NULL PRIMARY KEY COMMENT '角色ID，主键',
    role_name   VARCHAR(50) NOT NULL COMMENT '角色名称',
    is_default  TINYINT(1) DEFAULT 0 NOT NULL COMMENT '是否为默认角色(0-否 1-是)',
    deleted     INT DEFAULT 0 NOT NULL COMMENT '软删除标记(0-未删除 1-已删除)',
    `version`   INT DEFAULT 0 NOT NULL COMMENT '乐观锁版本号',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '记录创建时间',
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    created_by  VARCHAR(255) COMMENT '记录创建人',
    updated_by  VARCHAR(255) COMMENT '记录更新人'
) COMMENT='用户角色表，存储角色权限';

-- 默认角色唯一索引（只允许一个默认角色）
-- MySQL 8.0.13+ 支持函数索引，使用如下语法
CREATE UNIQUE INDEX uk_user_role_default_true ON `user_role` ((CASE WHEN is_default = 1 THEN is_default END));

-- 插入管理员角色
INSERT INTO `user_role` (id, role_name, is_default, deleted, `version`, created_at, updated_at)
VALUES (1, 'ROLE_ADMIN', 0, 0, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE role_name = role_name;

-- 插入普通用户角色
INSERT INTO `user_role` (id, role_name, is_default, deleted, `version`, created_at, updated_at)
VALUES (2, 'ROLE_USER', 1, 0, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE role_name = role_name;

-- user_profile 用户资料表
CREATE TABLE IF NOT EXISTS `user_profile`
(
    user_id     BIGINT PRIMARY KEY COMMENT '用户ID，关联sys_user表的主键',
    real_name   VARCHAR(255) COMMENT '用户真实姓名',
    gender      INT COMMENT '性别(0-女 1-男)',
    birthday    DATE COMMENT '用户生日',
    avatar      VARCHAR(255) COMMENT '用户头像URL',
    deleted     INT DEFAULT 0 NOT NULL COMMENT '软删除标记(0-未删除 1-已删除)',
    `version`   INT DEFAULT 0 NOT NULL COMMENT '乐观锁版本号',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '记录创建时间',
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    created_by  VARCHAR(255) COMMENT '记录创建人',
    updated_by  VARCHAR(255) COMMENT '记录更新人'
) COMMENT='用户资料表，存储用户的详细信息';
