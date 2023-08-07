package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/feedservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var FeedChan chan feedservice.FeedServiceClient

var feedAddr = ":8004"

func init() {
	coon, _ := grpc.Dial(feedAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	FeedChan = make(chan feedservice.FeedServiceClient, 10)
	for i := 0; i < 10; i++ {
		FeedChan <- feedservice.NewFeedServiceClient(coon)
	}
}
