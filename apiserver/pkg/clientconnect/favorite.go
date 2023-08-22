package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/favoriteservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var FavoriteActionChan chan favoriteservice.FavoriteActionServiceClient
var FavoriteListChan chan favoriteservice.FavoriteListServiceClient

var favoriteAddr = "favorite:8080"

func init() {
	FavoriteActionChan = make(chan favoriteservice.FavoriteActionServiceClient, 10)
	FavoriteListChan = make(chan favoriteservice.FavoriteListServiceClient, 10)
	for i := 0; i < 10; i++ {
		conn1, _ := grpc.Dial(favoriteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn2, _ := grpc.Dial(favoriteAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		FavoriteActionChan <- favoriteservice.NewFavoriteActionServiceClient(conn1)
		FavoriteListChan <- favoriteservice.NewFavoriteListServiceClient(conn2)
	}
}