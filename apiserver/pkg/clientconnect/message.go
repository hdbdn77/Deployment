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
	conn, _ := grpc.Dial(messageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	MessageChatChan = make(chan messageservice.MessageChatServiceClient, 10)
	MessageActionChan = make(chan messageservice.MessageActionServiceClient, 10)
	for i := 0; i < 10; i++ {
		MessageChatChan <- messageservice.NewMessageChatServiceClient(conn)
		MessageActionChan <- messageservice.NewMessageActionServiceClient(conn)
	}
}