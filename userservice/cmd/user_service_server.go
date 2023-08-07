package main

import (
	"fmt"
	"log"
	"net"

	"github.com/micro/simplifiedTikTok/userservice/pkg/userservice"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()

	userservice.RegisterRegisterServiceServer(server, userservice.RegisterService)
	userservice.RegisterLoginServiceServer(server, userservice.LoginService)
	userservice.RegisterUserServiceServer(server, userservice.UserService)

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("启动监听出错: %v", err)
	}

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("启动服务出错: %v", err)
	}

	fmt.Println("grpc服务器后台进行中...")

}
