

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


--
-- 表的结构 `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
  `uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'用户名',
  `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '唯一标识符',
  `mobile` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `passwd` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `user_state` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户状态',
  `created_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `rooms` (
    `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
    `uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'用户唯一标识符',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间名称',
    `host_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '主持人标识符',
    `type` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间类型',
    `created_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '修改时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `participants` (
    `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
    `uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'参会者唯一标识符',
    `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '参会者ID',
    `room_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '房间ID',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '参会者姓名',
    `permission` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户权限',
    `created_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- 转储表的索引
--

-- 表的索引 `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD KEY `uuid` (`uuid`),
  ADD KEY `mobile` (`mobile`);


ALTER TABLE `rooms`
    ADD PRIMARY KEY (`id`),
    ADD KEY `uuid` (`uuid`);

ALTER TABLE `participants`
    ADD PRIMARY KEY (`id`),
    ADD KEY `uuid` (`uuid`);

--
-- 使用表AUTO_INCREMENT `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键';
COMMIT;

-- 使用表AUTO_INCREMENT `rooms`
--
ALTER TABLE `rooms`
    MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键';
COMMIT;

-- 使用表AUTO_INCREMENT `participants`
--
ALTER TABLE `participants`
    MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键';
COMMIT;

