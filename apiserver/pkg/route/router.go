package route

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/controller"
)

func InitRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
  
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.GET("/user/", controller.UserInfo)
	
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)
  
	// extra apis - I
	// apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	// apiRouter.GET("/favorite/list/", controller.FavoriteList)
	// apiRouter.POST("/comment/action/", controller.CommentAction)
	// apiRouter.GET("/comment/list/", controller.CommentList)
  
	// extra apis - II
	// apiRouter.POST("/relation/action/", controller.RelationAction)
	// apiRouter.GET("/relation/follow/list/", controller.FollowList)
	// apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	// apiRouter.GET("/relation/friend/list/", controller.FriendList)
	// apiRouter.GET("/message/chat/", controller.MessageChat)
	// apiRouter.POST("/message/action/", controller.MessageAction)
  }