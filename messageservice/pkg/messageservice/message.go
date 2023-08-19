package messageservice

import (
	"context"
	"fmt"

	"github.com/micro/simplifiedTikTok/messageservice/pkg/model"
	"github.com/micro/simplifiedTikTok/messageservice/pkg/utils"
	"github.com/micro/simplifiedTikTok/messageservice/pkg/dao"
)

var MessageChatService = &messageChatService{}
var MessageActionService = &messageActionService{}

type messageChatService struct {}
type messageActionService struct {}

func ChatID(userAID int64, userBID int64) string {
	if userAID < userBID {
		return fmt.Sprintf("_%d_%d_", userAID, userBID)
	}else {
		return fmt.Sprintf("_%d_%d_", userBID, userAID)
	}
}

func (mC *messageChatService) MessageChat(context context.Context, request *DouYinMessageChatRequest) (*DouYinMessageChatResponse, error) {
	//实现具体的业务逻辑
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) {
		return &DouYinMessageChatResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
			MessageList: nil,
		}, nil
	}
	if claims.ID == request.ToUserId {
		return &DouYinMessageChatResponse{
			StatusCode: -1,
			StatusMsg: "聊天对象ID无效",
		}, nil
	}

	db := dao.GetDB()
	user , err := model.FindUserById(&model.User{Id: request.ToUserId}, db)
	if (err != nil) || (user.Username == "") {
		return &DouYinMessageChatResponse{
			StatusCode: -1,
			StatusMsg: "聊天对象ID无效",
		}, nil
	}

	// 查询消息记录
	chatID := ChatID(claims.ID, request.ToUserId)
	messages, err := model.ListMessageByTime(&model.Message{ChatID: chatID, CreateTime: request.PreMsgTime}, db)
	if err != nil {
		return &DouYinMessageChatResponse{
			StatusCode: -2,
			StatusMsg: "查询聊天记录失败",
		}, nil
	}
	var messageList []*Message
	for _, message := range messages {
		messageList = append(messageList, &Message{
			Id: message.Id,
			ToUserId: message.ToUserID,
			FromUserId: message.FromUserID,
			Content: message.Content,
			CreateTime: message.CreateTime,
		})
	}

	return &DouYinMessageChatResponse{
		StatusCode: 0,
		StatusMsg: "查询聊天记录成功",
		MessageList: messageList,
	}, nil
}

func (mA *messageActionService) MessageAction(context context.Context, request *DouYinMessageActionRequest) (*DouYinMessageActionResponse, error) {
	//实现具体的业务逻辑
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) {
		return &DouYinMessageActionResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
		}, nil
	}
	if claims.ID == request.ToUserId {
		return &DouYinMessageActionResponse{
			StatusCode: -1,
			StatusMsg: "聊天对象ID无效",
		}, nil
	}
	if request.ActionType != 1 {
		return &DouYinMessageActionResponse{
			StatusCode: -1,
			StatusMsg: "聊天操作码无效",
		}, nil
	}

	db := dao.GetDB()
	user , err := model.FindUserById(&model.User{Id: request.ToUserId}, db)
	if (err != nil) || (user.Username == "") {
		return &DouYinMessageActionResponse{
			StatusCode: -1,
			StatusMsg: "聊天对象ID无效",
		}, nil
	}

	tx := db.Begin()
	// 添加消息记录
	chatID := ChatID(claims.ID, request.ToUserId)
	_, err = model.SendMessage(&model.Message{ChatID: chatID, FromUserID: claims.ID, ToUserID: request.ToUserId, Content: request.Content}, tx)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return &DouYinMessageActionResponse{
			StatusCode: -2,
			StatusMsg: "发送消息失败",
		}, nil
	}
	// 若为好友，则添加朋友表记录
	isFriend := model.IsFriend(&model.Friend{Id: chatID}, tx)
	if isFriend {
		_, err = model.UpdateLatestMessage(chatID, claims.ID, request.ToUserId, request.Content, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinMessageActionResponse{
				StatusCode: -2,
				StatusMsg: "发送消息失败",
			}, nil
		}
	}

	err = tx.Commit().Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return &DouYinMessageActionResponse{
			StatusCode: -2,
			StatusMsg: "发送消息失败",
		}, nil
	}
	return &DouYinMessageActionResponse{
		StatusCode: 0,
		StatusMsg: "发送消息成功",
	}, nil
}

func (mC *messageChatService) mustEmbedUnimplementedMessageChatServiceServer() {}
func (mA *messageActionService) mustEmbedUnimplementedMessageActionServiceServer() {}
