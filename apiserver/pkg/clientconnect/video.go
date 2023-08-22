package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/videoservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var PublishActionChan chan videoservice.PublishActionServiceClient
var PublishListChan chan videoservice.PublishListServiceClient
var videoAddr = "video:8080"

func init() {
	PublishActionChan = make(chan videoservice.PublishActionServiceClient, 10)
	PublishListChan = make(chan videoservice.PublishListServiceClient, 10)
	for i := 0; i < 10; i++ {
		coon1, _ := grpc.Dial(videoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		coon2, _ := grpc.Dial(videoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		PublishActionChan <- videoservice.NewPublishActionServiceClient(coon1)
		PublishListChan <- videoservice.NewPublishListServiceClient(coon2)
	}
}