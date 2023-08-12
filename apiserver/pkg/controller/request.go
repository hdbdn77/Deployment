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

// favorite request
type FavoriteActionRequest struct {
	Token string `form:"token" binding:"required"`
	VideoId string `form:"video_id" binding:"required"`
	ActionType string `form:"action_type" binding:"required"`
}

type FavoriteListRequest struct {
	UserId string `form:"user_id" binding:"required"`
	Token string `form:"token" binding:"required"`
}

// comment request
type CommentActionRequest struct {
	Token string `form:"token" binding:"required"`
	VideoId string `form:"video_id" binding:"required"`
	ActionType string `form:"action_type" binding:"required"`
	CommentText string `form:"comment_text"`
	CommentId string `form:"comment_id"`
}

type CommentListRequest struct {
	Token string `form:"token" binding:"required"`
	VideoId string `form:"video_id" binding:"required"`
}