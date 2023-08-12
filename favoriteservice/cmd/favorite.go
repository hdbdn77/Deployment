package main

import (
	"fmt"
	"log"
	"net"

	"github.com/micro/simplifiedTikTok/favoriteservice/pkg/favoriteservice"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()

	favoriteservice.RegisterFavoriteActionServiceServer(server, favoriteservice.FavoriteActionService)
	favoriteservice.RegisterFavoriteListServiceServer(server, favoriteservice.FavoriteListService)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("启动监听出错: %v", err)
	}

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("启动服务出错: %v", err)
	}

	fmt.Println("grpc服务器后台进行中...")

}