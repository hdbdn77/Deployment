package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/messageservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var MessageChatChan chan messageservice.MessageChatServiceClient
var MessageActionChan chan messageservice.MessageActionServiceClient

var messageAddr = "message:8080"

func init() {
	MessageChatChan = make(chan messageservice.MessageChatServiceClient, 10)
	MessageActionChan = make(chan messageservice.MessageActionServiceClient, 10)
	for i := 0; i < 10; i++ {
		conn1, _ := grpc.Dial(messageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn2, _ := grpc.Dial(messageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		MessageChatChan <- messageservice.NewMessageChatServiceClient(conn1)
		MessageActionChan <- messageservice.NewMessageActionServiceClient(conn2)
	}
}