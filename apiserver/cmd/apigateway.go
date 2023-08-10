package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/route"

	"github.com/redis/go-redis/v9"
	"context"
)

func main() {
	r := gin.Default()

	route.InitRouter(r)
	
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func Refresh() {
	rdb:= redis.NewClient(&redis.Options{
		Addr:	  "redis:6379",
		Password: "",
	})
	rdb.LTrim(context.Background(),"video", 0, -1)
}