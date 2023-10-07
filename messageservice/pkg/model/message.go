package model

import (
	_"errors"
	_"fmt"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	Id              int64  `gorm:"primaryKey;autoIncrement;comment:'PrimaryKey'"`
	ChatID          string  `gorm:"not null;size:256;comment:'聊天ID'"`
	FromUserID      int64  `gorm:"not null;comment:'发送者ID'"`
	ToUserID  	    int64  `gorm:"not null;comment:'接收者ID'"`
	Content         string `gorm:"not null;comment:'消息内容'"`
	CreateTime      int64  `gorm:"not null;comment:'消息发送日期'"`
}

func (Message) TableName() string {
	return "message"
}

func SendMessage(message *Message, tx *gorm.DB) (*Message, error) {
	// 创建
	message.CreateTime = time.Now().Unix()
	err := tx.Create(message).Error
	return message, err
}

func ListMessageByTime(message *Message, tx *gorm.DB) ([]*Message, error) {
	var messages []*Message
	err := tx.Model(&Message{}).Where("chat_id = ? AND create_time > ?", message.ChatID, message.CreateTime).Find(&messages).Error
	return messages, err
}