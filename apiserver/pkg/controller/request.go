package controller

import(
	"mime/multipart"
)

type Request struct {
	UserId string `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

// user request
type UserLoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserInfoRequest struct {
	Request
}

//video request
type PublishRequest struct {
	Data *multipart.FileHeader `form:"data" binding:"required"`
	Token string `form:"token" binding:"required"`
	Title string `form:"title" binding:"required"`
}

type PublishListRequest struct {
	Request
}

// feed request
type FeedRequest struct {
	LatestTime string `form:"latest_time"`
	Token string `form:"token"`
}