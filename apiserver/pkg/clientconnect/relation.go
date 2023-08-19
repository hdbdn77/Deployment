package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/relationservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var RelationActionChan chan relationservice.RelationActionServiceClient
var FollowListChan chan relationservice.RelationFollowListServiceClient
var FollowerListChan chan relationservice.RelationFollowerListServiceClient
var FriendListChan chan relationservice.RelationFriendListServiceClient

var relationAddr = "relation:8080"

func init() {
	conn, _ := grpc.Dial(relationAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	RelationActionChan = make(chan relationservice.RelationActionServiceClient, 10)
	FollowListChan = make(chan relationservice.RelationFollowListServiceClient, 10)
	FollowerListChan = make(chan relationservice.RelationFollowerListServiceClient, 10)
	FriendListChan = make(chan relationservice.RelationFriendListServiceClient, 10)
	for i := 0; i < 10; i++ {
		RelationActionChan <- relationservice.NewRelationActionServiceClient(conn)
		FollowListChan <- relationservice.NewRelationFollowListServiceClient(conn)
		FollowerListChan <- relationservice.NewRelationFollowerListServiceClient(conn)
		FriendListChan <- relationservice.NewRelationFriendListServiceClient(conn)
	}
}