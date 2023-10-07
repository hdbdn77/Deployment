package model

import (
	"errors"
	"fmt"

	_"github.com/micro/simplifiedTikTok/messageservice/pkg/dao"
	"gorm.io/gorm"
)

type User struct {
	Id              int64  `gorm:"primaryKey;autoIncrement;comment:'PrimaryKey'"`
	Username        string `gorm:"size:32;not null;default:'';comment:'Username'"`
	Password        string `gorm:"size:32;not null;default:'';comment:'Password'"`
	FollowCount     int64  `gorm:"not null;default:0;comment:'FollowCount'"`
	FollowerCount   int64  `gorm:"not null;default:0;comment:'FollowerCount'"`
	IsFollow        bool   `gorm:"not null;default:false;comment:'IsFollow'"`
	Avatar          string `gorm:"size:128;not null;default:'';comment:'Avatar'"`
	BackgroundImage string `gorm:"size:128;not null;default:'';comment:'BackgroundImage'"`
	Signature       string `gorm:"size:256;not null;default:'';comment:'Signature'"`
	TotalFavorited  int64  `gorm:"not null;default:0;comment:'TotalFavorited'"`
	WorkCount       int64  `gorm:"not null;default:0;comment:'WorkCount'"`
	FavoriteCount   int64  `gorm:"not null;default:0;comment:'FavoriteCount'"`
}

// 用户名索引
func (User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 自定义Username唯一性校验
	var count int64
	tx.Model(u).Where("username = ?", u.Username).Count(&count)
	if count > 0 {
		return errors.New("username already existed")
	}

	return nil
}

func Register(user *User, tx *gorm.DB) (*User, error) {
	// 创建
	err := tx.Create(user).Error
	return user, err
}

func FindUserByUsername(user *User, tx *gorm.DB) (*User, error) {
	// 查询
	err := tx.Model(&User{}).Where("username = ?", user.Username).Take(&user).Error
	return user, err
}

func FindUserById(user *User, tx *gorm.DB) (*User, error) {
	// 查询
	err := tx.Model(&User{}).Where("id = ?", user.Id).Take(&user).Error
	return user, err
}

func AddTotalFavorited(user *User, tx *gorm.DB) (*User, error) {
	err := tx.Model(&User{}).Where("id = ?", user.Id).Take(&user).Error
	if err != nil {
		fmt.Println("增加获赞总数时查找用户失败：", err)
		return nil, err
	}

	err = tx.Model(user).Update("total_favorited", user.TotalFavorited + 1).Error
	if err != nil {
		fmt.Println("增加获赞总数失败：", err)
		return nil, err
	}
	return user, nil
}

func MinusTotalFavorited(user *User, tx *gorm.DB) (*User, error) {
	err := tx.Model(&User{}).Where("id = ?", user.Id).Take(&user).Error
	if err != nil {
		fmt.Println("减少获赞总数时查找用户失败：", err)
		return nil, err
	}

	err = tx.Model(user).Update("total_favorited", user.TotalFavorited - 1).Error
	if err != nil {
		fmt.Println("减少获赞总数失败：", err)
		return nil, err
	}
	return user, nil
}

func AddWorkCount(user *User, tx *gorm.DB) (*User, error) {
	err := tx.Model(&User{}).Where("id = ?", user.Id).Take(&user).Error
	if err != nil {
		fmt.Println("更新作品总数时查找用户失败：", err)
		return nil, err
	}

	err = tx.Model(user).Update("work_count", user.WorkCount + 1).Error
	if err != nil {
		fmt.Println("更新作品总数失败：", err)
		return nil, err
	}
	return user, nil
}

func AddUserFavoriteCount(user *User, tx *gorm.DB) (*User, error) {
	err := tx.Model(&User{}).Where("id = ?", user.Id).Take(&user).Error
	if err != nil {
		fmt.Println("增加点赞总数时查找用户失败：", err)
		return nil, err
	}

	err = tx.Model(user).Update("favorite_count", user.FavoriteCount + 1).Error
	if err != nil {
		fmt.Println("增加点赞总数失败：", err)
		return nil, err
	}
	return user, nil
}

func MinusUserFavoriteCount(user *User, tx *gorm.DB) (*User, error) {
	err := tx.Model(&User{}).Where("id = ?", user.Id).Take(&user).Error
	if err != nil {
		fmt.Println("减少点赞总数时查找用户失败：", err)
		return nil, err
	}

	err = tx.Model(user).Update("favorite_count", user.FavoriteCount - 1).Error
	if err != nil {
		fmt.Println("减少点赞总数失败：", err)
		return nil, err
	}
	return user, nil
}