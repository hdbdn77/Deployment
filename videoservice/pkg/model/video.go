package model

import (
	_ "errors"
	"time"

	"github.com/micro/simplifiedTikTok/videoservice/pkg/dao"
	_ "gorm.io/gorm"
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

func CreateVideo(video *Video) (*Video, error) {
	// 获取数据库连接
	db := dao.GetDB()
	// 迁移模型
	db.AutoMigrate(&Video{})

	// 创建
	video.PublishTime = time.Now().Unix()
	err := db.Create(video).Error
	return video, err
}

func ListVideoByAuthorId(authorId int64) (*[]Video, error) {
	// 获取数据库连接
	db := dao.GetDB()
	// 迁移模型
	db.AutoMigrate(&Video{})

	//查询
	var videos []Video
	err := db.Where("author_id = ?", authorId).Find(&videos).Error
	return &videos, err

}

func GetVideoById(id int64) (*Video, error) {
	// 获取数据库连接
	db := dao.GetDB()
	// 迁移模型
	db.AutoMigrate(&Video{})

	var video Video
	err := db.Where("id = ?", id).Take(&video).Error
	return &video, err
}

func ListVideoByTime(time int64) (*[]Video, error) {
	// 获取数据库连接
	db := dao.GetDB()
	// 迁移模型
	db.AutoMigrate(&Video{})

	//查询
	var videos []Video
	err := db.Where("publish_time < ?", time).Order("publish_time asc").Limit(30).Find(&videos).Error
	return &videos, err
}