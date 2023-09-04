CREATE TABLE `friend` (
  `id` varchar(256) NOT NULL COMMENT 'PrimaryKey',
  `from_user_id` bigint NOT NULL COMMENT 'FromUserID',
  `to_user_id` bigint NOT NULL COMMENT 'ToUserID',
  `latest_message` TEXT NOT NULL COMMENT 'Latest Message',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Friend create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Friend update time',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Friend delete time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Friend table';