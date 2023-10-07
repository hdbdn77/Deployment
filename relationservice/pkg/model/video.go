package model

import (
	_ "errors"
	"time"
	"fmt"

	_"github.com/micro/simplifiedTikTok/relationservice/pkg/dao"
	"gorm.io/gorm"
)

type Video struct {
    Id            int64  `gorm:"primaryKey;autoIncrement"` 
    AuthorId      int64  `gorm:"not null"`
    PlayUrl       string `gorm:"not null"`
    CoverUrl      string `gorm:"default:''"`
    FavoriteCount int64  `gorm:"default:0"`
    CommentCount  int64  `gorm:"default:0"`
    IsFavorite    bool   `gorm:"default:false"`
    Title         string `gorm:"default:'';size:20"`
	PublishTime   int64
}

func (Video) TableName() string {
    return "video" 
}

func CreateVideo(video *Video, tx *gorm.DB) (*Video, error) {
	// 创建
	video.PublishTime = time.Now().Unix()
	err := tx.Create(video).Error
	return video, err
}

func ListVideoByAuthorId(authorId int64, tx *gorm.DB) (*[]Video, error) {
	//查询
	var videos []Video
	err := tx.Model(&Video{}).Where("author_id = ?", authorId).Find(&videos).Error
	return &videos, err

}

func GetVideoById(id int64, tx *gorm.DB) (*Video, error) {
	var video Video
	err := tx.Model(&Video{}).Where("id = ?", id).Take(&video).Error
	return &video, err
}

func ListVideoByTime(time int64, tx *gorm.DB) (*[]Video, error) {
	//查询
	var videos []Video
	err := tx.Model(&Video{}).Where("publish_time < ?", time).Order("publish_time desc").Limit(30).Find(&videos).Error
	return &videos, err
}

func AddViseoFavoriteCount(video *Video, tx *gorm.DB) (*Video, error) {
	err := tx.Model(&Video{}).Where("id = ?", video.Id).Take(video).Error
	if err != nil {
		fmt.Println("增加视频获赞总数时查找视频失败：", err)
		return nil, err
	}
	err = tx.Model(video).Update("favorite_count", video.FavoriteCount + 1).Error
	if err != nil {
		fmt.Println("增加视频获赞总数失败：", err)
		return nil, err
	}
	return video, nil
}

func MinusViseoFavoriteCount(video *Video, tx *gorm.DB) (*Video, error) {
	err := tx.Model(&Video{}).Where("id = ?", video.Id).Take(video).Error
	if err != nil {
		fmt.Println("减少视频获赞总数时查找视频失败：", err)
		return nil, err
	}
	err = tx.Model(video).Update("favorite_count", video.FavoriteCount - 1).Error
	if err != nil {
		fmt.Println("减少视频获赞总数失败：", err)
		return nil, err
	}
	return video, nil
}

func AddVideoCommentCount(video *Video, tx *gorm.DB) (*Video, error) {
	err := tx.Model(&Video{}).Where("id = ?", video.Id).Take(video).Error
	if err != nil {
		fmt.Println("增加视频评论总数时查找视频失败：", err)
		return nil, err
	}
	err = tx.Model(video).Update("comment_count", video.CommentCount + 1).Error
	if err != nil {
		fmt.Println("增加视频评论总数失败：", err)
		return nil, err
	}
	return video, nil
}

func MinusVideoCommentCount(video *Video, tx *gorm.DB) (*Video, error) {
	err := tx.Model(&Video{}).Where("id = ?", video.Id).Take(video).Error
	if err != nil {
		fmt.Println("减少视频评论总数时查找视频失败：", err)
		return nil, err
	}
	err = tx.Model(video).Update("comment_count", video.CommentCount - 1).Error
	if err != nil {
		fmt.Println("减少视频评论总数失败：", err)
		return nil, err
	}
	return video, nil
}