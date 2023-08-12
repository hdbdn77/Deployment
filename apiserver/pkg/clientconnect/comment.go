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
	conn, _ := grpc.Dial(commentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	CommentActionChan = make(chan commentservice.CommentActionServiceClient, 10)
	CommentListChan = make(chan commentservice.CommentListServiceClient, 10)
	for i := 0; i < 10; i++ {
		CommentActionChan <- commentservice.NewCommentActionServiceClient(conn)
		CommentListChan <- commentservice.NewCommentListServiceClient(conn)
	}
}