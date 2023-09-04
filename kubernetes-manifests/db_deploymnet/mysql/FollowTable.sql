CREATE TABLE `follow` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'PrimaryKey',
  `follower_user_id` bigint NOT NULL COMMENT 'Follower user id',
  `followed_user_id` bigint NOT NULL COMMENT 'Followed user id',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Follow create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Follow update time',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Follow delete time',
  PRIMARY KEY (`id`),
  KEY          `follower_user_id` (`follower_user_id`) COMMENT 'Follower index',
  KEY          `followed_user_id` (`followed_user_id`) COMMENT 'Followed index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Follow table';