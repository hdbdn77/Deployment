CREATE TABLE `comment` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'PrimaryKey', -- 评论唯一标识
  `user_id` bigint NOT NULL COMMENT 'UserId', -- 评论者信息
  `video_id` bigint NOT NULL COMMENT 'VideoId', -- 视频信息
  `content` TEXT NOT NULL COMMENT 'Content', -- 评论内容
  `create_date` bigint NOT NULL COMMENT 'CreateDate', -- 评论发布日期
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'comment create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'comment update time',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'comment delete time',
  PRIMARY KEY (`id`),
  KEY          `comment_user_id` (`user_id`) COMMENT 'Comment user id',
  KEY          `comment_video_id` (`video_id`) COMMENT 'Comment video id',
  KEY          `comment_create_date` (`create_date`) COMMENT 'Comment Create Date'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Comment table';