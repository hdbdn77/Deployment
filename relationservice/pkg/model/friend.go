package model

import (
	"errors"
	_"fmt"
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
	// 创建
	err := tx.Create(friend).Error
	return friend, err
}

func Unfriended(friend *Friend, tx *gorm.DB) (*Friend, error) {
	// 删除
	tx.Model(&Friend{}).Where("id = ?", friend.Id).Find(&friend)
	err := tx.Delete(friend).Error
	return friend, err
}

func FriendList(friendKey string, tx *gorm.DB) ([]*Friend, error) {
	var friends []*Friend
	err := tx.Model(&Friend{}).Where("id LIKE ?", "%"+ friendKey + "%").Find(&friends).Error
	return friends, err
}