

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";




-- --------------------------------------------------------

--
-- 表的结构 `device`
--

CREATE TABLE `device` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
  `uid` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户主键',
  `ip` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'ip地址',
  `created_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '注册时间',
  `client` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `trace`
--

CREATE TABLE `trace` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
  `uid` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户主键',
  `ip` int(10) UNSIGNED NOT NULL COMMENT 'ip',
  `created_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间'

) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- 表的结构 `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
  `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `mobile` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `passwd` varchar(40) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `created_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated_time` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '修改时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- 转储表的索引
--

--
-- 表的索引 `device`
--
ALTER TABLE `device`
  ADD PRIMARY KEY (`id`),
  ADD KEY `uid` (`uid`);

--
-- 表的索引 `trace`
--
ALTER TABLE `trace`
  ADD PRIMARY KEY (`id`),
  ADD KEY `UT` (`uid`) USING BTREE;

--
-- 表的索引 `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD KEY `create_time` (`created_time`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `device`
--
ALTER TABLE `device`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键';

--
-- 使用表AUTO_INCREMENT `trace`
--
ALTER TABLE `trace`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键', AUTO_INCREMENT=2;

--
-- 使用表AUTO_INCREMENT `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键';
COMMIT;
