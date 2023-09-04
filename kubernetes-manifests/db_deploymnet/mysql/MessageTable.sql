CREATE TABLE `message` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'PrimaryKey',
  `chat_id` varchar(256) NOT NULL COMMENT 'ChatID',
  `from_user_id` bigint NOT NULL COMMENT 'FromUserID',
  `to_user_id` bigint NOT NULL COMMENT 'ToUserID',
  `content` TEXT NOT NULL COMMENT 'Content',
  `create_time` bigint NOT NULL COMMENT 'CreateTime',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Message create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Message update time',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Message delete time',
  PRIMARY KEY (`id`),
  KEY          `chat_id` (`chat_id`) COMMENT 'ChatID',
  KEY          `create_time` (`create_time`) COMMENT 'CreateTime'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Message table';