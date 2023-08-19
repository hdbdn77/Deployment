package controller

import (
	"context"
	_ "fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/clientconnect"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/messageservice"
)

func MessageChat(c *gin.Context){
	var messageChatRequest MessageChatRequest
	err := c.ShouldBindQuery(&messageChatRequest)
	if err != nil {
		c.JSON(http.StatusOK, MessageChatResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "聊天记录信息输入有误",
			},
		})
		return
	}

	toUserID, _ := strconv.ParseInt(messageChatRequest.ToUserID, 10, 64)
	preMsgTime, _ := strconv.ParseInt(messageChatRequest.PreMsgTime, 10, 64)
	messageChatClient := <- clientconnect.MessageChatChan
	messageChatResponse, err := messageChatClient.MessageChat(context.Background(), &messageservice.DouYinMessageChatRequest{Token: messageChatRequest.Token, ToUserId: toUserID, PreMsgTime: preMsgTime})
	clientconnect.MessageChatChan <- messageChatClient

	if (messageChatResponse == nil) || (err != nil) {
		c.JSON(http.StatusOK, MessageChatResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: " message chat failed",
			},
		})
		return
	}

	if messageChatResponse.StatusCode != 0  {
		c.JSON(http.StatusOK, MessageChatResponse{
			Response : Response{
				StatusCode: messageChatResponse.StatusCode, 
				StatusMsg: messageChatResponse.StatusMsg,
			},
		})
		return
	}

	var messageList []Message
	for _, message := range messageChatResponse.MessageList {
		createTime := time.Unix(message.CreateTime, 0).Format("2006-01-02 15:04:05")
		messageList = append(messageList, Message{
			Id: message.Id,
			ToUserId: message.ToUserId,
			FromUserId: message.FromUserId,
			Content: message.Content,
			CreateTime: createTime,
		})
	}
	c.JSON(http.StatusOK, MessageChatResponse{
		Response : Response{
			StatusCode: messageChatResponse.StatusCode, 
			StatusMsg: messageChatResponse.StatusMsg,
		},
		MessageList: messageList,
	})
}

func MessageAction(c *gin.Context){
	var messageActionRequest MessageActionRequest
	err := c.ShouldBind(&messageActionRequest)
	if err != nil {
		c.JSON(http.StatusOK, MessageActionResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "聊天记录信息输入有误",
			},
		})
		return
	}

	toUserID, _ := strconv.ParseInt(messageActionRequest.ToUserID, 10, 64)
	actionType, _ := strconv.ParseInt(messageActionRequest.ActionType, 10, 64)
	messageActionClient := <- clientconnect.MessageActionChan
	messageActionResponse, err := messageActionClient.MessageAction(context.Background(), &messageservice.DouYinMessageActionRequest{Token: messageActionRequest.Token, ToUserId: toUserID, ActionType: int32(actionType) , Content: messageActionRequest.Content})
	clientconnect.MessageActionChan <- messageActionClient

	if (messageActionResponse == nil) || (err != nil) {
		c.JSON(http.StatusOK, MessageActionResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "message action faild",
			},
		})
		return
	}

	c.JSON(http.StatusOK, MessageActionResponse{
		Response : Response{
			StatusCode: messageActionResponse.StatusCode, 
			StatusMsg: messageActionResponse.StatusMsg,
		},
	})
}