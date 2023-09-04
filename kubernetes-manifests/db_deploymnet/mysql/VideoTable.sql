CREATE TABLE `video` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'PrimaryKey', -- 视频唯一标识
  `author_id` bigint NOT NULL COMMENT 'UserId', -- 视频作者信息
  `play_url` varchar(128) NOT NULL COMMENT 'VideoPlayUrl', -- 视频播放地址
  `cover_url` varchar(128) NOT NULL COMMENT 'VideoCoverUrl', -- 视频封面地址,
  `favorite_count` bigint NOT NULL DEFAULT 0 COMMENT 'VideoFavoriteCount', -- 视频的点赞总数
  `comment_count` bigint NOT NULL DEFAULT 0 COMMENT 'VideoCommentCount', -- 视频的评论总数
  `is_favorite` boolean NOT NULL DEFAULT 0 COMMENT 'IsFavorite', -- true-已点赞，false-未点赞
  `title` varchar(20) NOT NULL DEFAULT '' COMMENT 'VideoTitle', -- 视频标题
  `publish_time` bigint NOT NULL COMMENT 'PublishTime', -- 发布时间
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'video create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'video update time',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'video delete time',
  PRIMARY KEY (`id`),
  KEY          `video_user_id` (`author_id`) COMMENT 'VideoOfUserId index',
  KEY          `video_publish_time` (`publish_time`) COMMENT 'PublishTimeOfVideo index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Video table';