package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/videoservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var PublishActionChan chan videoservice.PublishActionServiceClient
var PublishListChan chan videoservice.PublishListServiceClient
var videoAddr = "121.41.123.39:30828"

func init() {
	coon, _ := grpc.Dial(videoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	PublishActionChan = make(chan videoservice.PublishActionServiceClient, 10)
	PublishListChan = make(chan videoservice.PublishListServiceClient, 10)
	for i := 0; i < 10; i++ {
		PublishActionChan <- videoservice.NewPublishActionServiceClient(coon)
		PublishListChan <- videoservice.NewPublishListServiceClient(coon)
	}
}