package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/commentservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var CommentActionChan chan commentservice.CommentActionServiceClient
var CommentListChan chan commentservice.CommentListServiceClient

var commentAddr = "comment:8080"

func init() {
	CommentActionChan = make(chan commentservice.CommentActionServiceClient, 10)
	CommentListChan = make(chan commentservice.CommentListServiceClient, 10)
	for i := 0; i < 10; i++ {
		conn1, _ := grpc.Dial(commentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn2, _ := grpc.Dial(commentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		CommentActionChan <- commentservice.NewCommentActionServiceClient(conn1)
		CommentListChan <- commentservice.NewCommentListServiceClient(conn2)
	}
}