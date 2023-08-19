package model

import (
	"errors"
	"fmt"
	_"time"

	"gorm.io/gorm"
)

type Friend struct {
	Id              string  `gorm:"primaryKey;size:256;comment:'PrimaryKey'"`
	FromUserID      int64  `gorm:"not null;comment:'发送者ID'"`
	ToUserID  	    int64  `gorm:"not null;comment:'接收者ID'"`
	LatestMessage   string  `gorm:"not null;comment:'最新消息'"`
}

func (Friend) TableName() string {
	return "friend"
}

func (f *Friend) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Model(f).Where("id = ?", f.Id).Count(&count)
	if count > 0 {
		return errors.New("friend already existed")
	}

	return nil
}

func (f *Friend) BeforeDelete(tx *gorm.DB) error {
	var count int64
	tx.Model(f).Where("id = ?", f.Id).Count(&count)
	if count == 0 {
		return errors.New("friend not exist")
	}

	return nil
}

func Friended(friend *Friend, tx *gorm.DB) (*Friend, error) {
	// 迁移模型
	tx.AutoMigrate(&Friend{})

	// 创建
	err := tx.Create(friend).Error
	return friend, err
}

func Unfriended(friend *Friend, tx *gorm.DB) (*Friend, error) {
	// 迁移模型
	tx.AutoMigrate(&Friend{})

	// 删除
	tx.Where("id = ?", friend.Id).Find(&friend)
	err := tx.Delete(friend).Error
	return friend, err
}

func FriendList(friendKey string, tx *gorm.DB) ([]*Friend, error) {
	// 迁移模型
	tx.AutoMigrate(&Friend{})

	var friends []*Friend
	err := tx.Where("id LIKE ?", "%"+ friendKey + "%").Find(&friends).Error
	return friends, err
}

func IsFriend(friend *Friend, tx *gorm.DB) bool {
	// 迁移模型
	tx.AutoMigrate(&Friend{})

	var count int64
	tx.Model(friend).Where("id = ?", friend.Id).Count(&count)
	return count > 0
}

func UpdateLatestMessage(friendKey string, fromUserID int64, toUserID int64, latestMessage string, tx *gorm.DB) (*Friend, error) {
	// 迁移模型
	tx.AutoMigrate(&Friend{})

	var friend Friend
	err := tx.Where("id = ?", friendKey).Take(&friend).Error
	if err != nil {
		fmt.Println("更新朋友最新消息时查找用户失败：", err)
		return nil, err
	}

	err = tx.Model(&friend).Update("from_user_id", fromUserID).Error
	if err != nil {
		fmt.Println("更新朋友最新消息失败：", err)
		return nil, err
	}
	err = tx.Model(&friend).Update("to_user_id", toUserID).Error
	if err != nil {
		fmt.Println("更新朋友最新消息失败：", err)
		return nil, err
	}
	err = tx.Model(&friend).Update("latest_message", latestMessage).Error
	if err != nil {
		fmt.Println("更新朋友最新消息失败：", err)
		return nil, err
	}
	return &friend, nil
}