package clientconnect

import (
	"github.com/micro/simplifiedTikTok/apiserver/pkg/userservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UserRegisterChan chan userservice.RegisterServiceClient
var UserLoginChan chan userservice.LoginServiceClient
var UserChan chan userservice.UserServiceClient
var userAddr = "user:8080"

func init() {
	UserRegisterChan = make(chan userservice.RegisterServiceClient, 10)
	UserLoginChan = make(chan userservice.LoginServiceClient, 10)
	UserChan = make(chan userservice.UserServiceClient, 10)
	for i := 0; i < 10; i++ {
		conn1, _ := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn2, _ := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		conn3, _ := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		UserRegisterChan <- userservice.NewRegisterServiceClient(conn1)
		UserLoginChan <- userservice.NewLoginServiceClient(conn2)
		UserChan <- userservice.NewUserServiceClient(conn3)
	}
}