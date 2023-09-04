CREATE TABLE `favorite` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'PrimaryKey', -- 点赞唯一标识
  `user_id` bigint NOT NULL COMMENT 'UserId', -- 点赞者信息,
  `video_id` bigint NOT NULL COMMENT 'VideoId',  -- 视频信息
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'favorite create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'favorite update time',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'favorite delete time',
  PRIMARY KEY (`id`),
  KEY          `favorite_user_id` (`user_id`) COMMENT 'Favorite user id',
  KEY          `favorite_video_id` (`video_id`) COMMENT 'Favorite video id'
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Favorite table';