package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/route"
)

func main() {
	r := gin.Default()

	route.InitRouter(r)
	
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
