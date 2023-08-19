package model

import (
	"errors"
	_"fmt"

	_"github.com/micro/simplifiedTikTok/relationservice/pkg/dao"
	"gorm.io/gorm"
)

type Favorite struct {
	Id              int64  `gorm:"primaryKey;autoIncrement;comment:'PrimaryKey'"`
	UserID          int64  `gorm:"not null;comment:'用户ID'"`
	VideoID         int64  `gorm:"not null;comment:'视频ID'"`
}

func (Favorite) TableName() string {
	return "favorite"
}

func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	// 自定义Favorite唯一性校验
	var count int64
	tx.Model(f).Where("user_id = ? AND video_id = ?", f.UserID, f.VideoID).Count(&count)
	if count > 0 {
		return errors.New("favorite already existed")
	}

	return nil
}

func (f *Favorite) BeforeDelete(tx *gorm.DB) error {
	// 自定义Favorite存在性校验
	var count int64
	tx.Model(f).Where("user_id = ? AND video_id = ?", f.UserID, f.VideoID).Count(&count)
	if count == 0 {
		return errors.New("favorite not exist")
	}

	return nil
}

func Favorited(favorite *Favorite, tx *gorm.DB) (*Favorite, error) {
	// 迁移模型
	tx.AutoMigrate(&Favorite{})

	// 创建
	err := tx.Create(favorite).Error
	return favorite, err
}

func Unfavorited(favorite *Favorite, tx *gorm.DB) (*Favorite, error) {
	// 迁移模型
	tx.AutoMigrate(&Favorite{})

	// 删除
	tx.Where("user_id = ? AND video_id = ?", favorite.UserID, favorite.VideoID).Find(&favorite)
	err := tx.Delete(favorite).Error
	return favorite, err
}

func IsFavorited(favorite *Favorite, tx *gorm.DB) bool {
	// 迁移模型
	tx.AutoMigrate(&Favorite{})

	var count int64
	tx.Model(favorite).Where("user_id = ? AND video_id = ?", favorite.UserID, favorite.VideoID).Count(&count)
	return count > 0
}

func GetFavoritesByUserID(userId int64, tx *gorm.DB) ([]*Favorite, error) {
	// 迁移模型
	tx.AutoMigrate(&Favorite{})

	var favorites []*Favorite
	err := tx.Where("user_id = ?", userId).Find(&favorites).Error
	return favorites, err
}