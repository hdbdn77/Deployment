package model

import (
	"errors"
	_"fmt"
	"time"

	_"github.com/micro/simplifiedTikTok/favoriteservice/pkg/dao"
	"gorm.io/gorm"
)

type Comment struct {
	Id              int64  `gorm:"primaryKey;autoIncrement;comment:'PrimaryKey'"`
	UserID          int64  `gorm:"not null;comment:'用户ID'"`
	VideoID         int64  `gorm:"not null;comment:'视频ID'"`
	Content         string `gorm:"not null;comment:'评论内容'"`
	CreateDate      int64  `gorm:"not null;comment:'评论发布日期'"`
}

func (Comment) TableName() string {
	return "comment"
}

func (c *Comment) BeforeDelete(tx *gorm.DB) error {
	// 自定义Comment存在性校验
	var count int64
	tx.Model(c).Where("id = ?", c.Id).Count(&count)
	if count == 0 {
		return errors.New("Comment not exist")
	}

	return nil
}

func Commented(comment *Comment, tx *gorm.DB) (*Comment, error) {
	// 创建
	comment.CreateDate = time.Now().Unix()
	err := tx.Create(comment).Error
	return comment, err
}

func UnCommented(comment *Comment, tx *gorm.DB) (*Comment, error) {
	// 删除
	tx.Model(&Comment{}).Where("id = ?", comment.Id).Find(&comment)
	err := tx.Delete(comment).Error
	return comment, err
}

func ListComment(comment *Comment, tx *gorm.DB) (*[]Comment, error) {
	var comments []Comment
	err := tx.Model(&Comment{}).Where("video_id = ?", comment.VideoID).Order("create_date desc").Find(&comments).Error
	return &comments, err
}