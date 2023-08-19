package dao

import (
	"github.com/micro/simplifiedTikTok/relationservice/configs"
	"github.com/redis/go-redis/v9"
)

var (
	redisConfig = configs.Redis
	client 		*redis.Client
)

func init() {
	// 连接数据库
	rdb:= redis.NewClient(&redis.Options{
		Addr:	  redisConfig.Addr,
		Password: redisConfig.Password,
	})

	client = rdb
}

func GetClient() *redis.Client { 
	return client
}