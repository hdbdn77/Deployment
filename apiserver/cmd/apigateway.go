package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/route"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/utils"
)

func main() {
	r := gin.Default()

	route.InitRouter(r)
	utils.InitRequestRateLimit(r, 100)
	r.Use(gin.Recovery())
	
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}


