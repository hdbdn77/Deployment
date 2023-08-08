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
	conn, _ := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	UserRegisterChan = make(chan userservice.RegisterServiceClient, 10)
	UserLoginChan = make(chan userservice.LoginServiceClient, 10)
	UserChan = make(chan userservice.UserServiceClient, 10)
	for i := 0; i < 10; i++ {
		UserRegisterChan <- userservice.NewRegisterServiceClient(conn)
		UserLoginChan <- userservice.NewLoginServiceClient(conn)
		UserChan <- userservice.NewUserServiceClient(conn)
	}
}