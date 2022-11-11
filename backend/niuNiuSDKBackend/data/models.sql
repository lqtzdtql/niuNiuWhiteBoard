

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


--
-- 表的结构 `whiteboards`
--

CREATE TABLE `whiteboards` (
  `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
  `uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'白板uuid',
  `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '白板名',
  `room_uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'白板所在房间',
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '白板创建时间',
  `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '白板修改时间',
  `deleted_time` datetime DEFAULT NULL COMMENT '白板删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `rooms` (
    `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
    `uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'用户唯一标识符',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间名称',
    `host_uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '主持人标识符',
    `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '建房时间',
    `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_time` datetime DEFAULT NULL COMMENT '删除时间'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `participants` (
    `id` bigint(20) UNSIGNED NOT NULL COMMENT '主键',
    `name` varchar(50) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '参会者用户名',
    `uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'参会者唯一标识符',
    `room_uuid` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'参会者所在房间标识符',
    `permission` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户权限',
    `current_whiteboard` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT'参会者最近使用的白板标识符',
    `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '进房时间',
    `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    `deleted_time` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- 转储表的索引
--

-- 表的索引 `whiteboards`
--
ALTER TABLE `whiteboards`
  ADD PRIMARY KEY (`id`),
  ADD KEY `uuid` (`uuid`);


ALTER TABLE `rooms`
    ADD PRIMARY KEY (`id`),
    ADD KEY `uuid` (`uuid`);

ALTER TABLE `participants`
    ADD PRIMARY KEY (`id`);

--
-- 使用表AUTO_INCREMENT `whiteboards`
--
ALTER TABLE `whiteboards`
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

