create database if not exists agent_user;
use agent_user;



CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '软删除时间',
  `username` VARCHAR(255) NOT NULL COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码',
  `nickname` VARCHAR(255) NULL DEFAULT NULL COMMENT '昵称',
  `email` VARCHAR(255) NULL DEFAULT NULL COMMENT '邮箱',
  `status` INT NULL DEFAULT 1 COMMENT '状态：1正常，0禁用',
  `role` VARCHAR(255) NULL DEFAULT 'user' COMMENT '角色 user/admin',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;