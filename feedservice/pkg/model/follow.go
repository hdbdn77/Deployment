package model

import (
	"errors"
	_"fmt"
	_"time"

	"gorm.io/gorm"
)

type Follow struct {
	Id              int64  `gorm:"primaryKey;autoIncrement;comment:'PrimaryKey'"`
	FollowerUserID  int64  `gorm:"not null;comment:'关注者ID'"`
	FollowedUserID  int64  `gorm:"not null;comment:'被关注者ID'"`
}

func (Follow) TableName() string {
	return "follow"
}

func (f *Follow) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Model(f).Where("follower_user_id = ? AND followed_user_id = ?", f.FollowerUserID, f.FollowedUserID).Count(&count)
	if count > 0 {
		return errors.New("follow already existed")
	}

	return nil
}

func (f *Follow) BeforeDelete(tx *gorm.DB) error {
	var count int64
	tx.Model(f).Where("follower_user_id = ? AND followed_user_id = ?", f.FollowerUserID, f.FollowedUserID).Count(&count)
	if count == 0 {
		return errors.New("follow not exist")
	}

	return nil
}

func Followed(follow *Follow, tx *gorm.DB) (*Follow, error) {
	// 创建
	err := tx.Create(follow).Error
	return follow, err
}

func Unfollowed(follow *Follow, tx *gorm.DB) (*Follow, error) {
	// 删除
	tx.Model(&Follow{}).Where("follower_user_id = ? AND followed_user_id = ?", follow.FollowerUserID, follow.FollowedUserID).Find(&follow)
	err := tx.Delete(follow).Error
	return follow, err
}

func IsFollowed(follow *Follow, tx *gorm.DB) bool {
	var count int64
	tx.Model(follow).Where("follower_user_id = ? AND followed_user_id = ?", follow.FollowerUserID, follow.FollowedUserID).Count(&count)
	return count > 0
}

func FollowList(useID int64, tx *gorm.DB) ([]*Follow, error) {
	var followList []*Follow
	err := tx.Model(&Follow{}).Where("follower_user_id = ?", useID).Find(&followList).Error
	return followList, err
}

func FollowerList(useID int64, tx *gorm.DB) ([]*Follow, error) {
	var followList []*Follow
	err := tx.Model(&Follow{}).Where("followed_user_id = ?", useID).Find(&followList).Error
	return followList, err
}