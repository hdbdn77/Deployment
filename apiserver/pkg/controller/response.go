package controller

//通用回复
type Response struct { 
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
}

type FriendUser struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	Message         string `json:"message,omitempty"`
	MsgType         int64  `json:"msgType,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

// type MessageSendEvent struct {
// 	UserId     int64  `json:"user_id,omitempty"`
// 	ToUserId   int64  `json:"to_user_id,omitempty"`
// 	MsgContent string `json:"msg_content,omitempty"`
// }

// type MessagePushEvent struct {
// 	FromUserId int64  `json:"user_id,omitempty"`
// 	MsgContent string `json:"msg_content,omitempty"`
// }

// 专属回复
// user
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// video
type PublishResponse struct {
	Response
}

type PublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// feed
type FeedResponse struct {
	Response
	NextTime int64 `json:"next_time,omitempty"`
	VideoList []Video `json:"video_list"`
}

// favorite
type FavoriteActionResponse struct {
	Response
}

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// comment
type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}

// relation
type RelationActionResponse struct {
	Response
}

type FollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type FollowerListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type FriendListResponse struct {
	Response
	FriendUserList []FriendUser `json:"user_list"`
}

// message
type MessageChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

type MessageActionResponse struct {
	Response
}