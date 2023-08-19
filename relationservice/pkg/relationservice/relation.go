package relationservice

import (
	"context"
	"fmt"

	"github.com/micro/simplifiedTikTok/relationservice/pkg/dao"
	"github.com/micro/simplifiedTikTok/relationservice/pkg/model"
	"github.com/micro/simplifiedTikTok/relationservice/pkg/utils"
)

var RelationActionService = &relationActionService{}
var RelationFollowListService = &relationFollowListService{}
var RelationFollowerListService = &relationFollowerListService{}
var RelationFriendListService = &relationFriendListService{}

type relationActionService struct {}
type relationFollowListService struct {}
type relationFollowerListService struct {}
type relationFriendListService struct {}

func FriendID(userAID int64, userBID int64) string {
	if userAID < userBID {
		return fmt.Sprintf("_%d_%d_", userAID, userBID)
	}else {
		return fmt.Sprintf("_%d_%d_", userBID, userAID)
	}
}

func (rA *relationActionService) RelationAction(context context.Context, request *DouYinRelationActionRequest) (*DouYinRelationActionResponse, error) {
	claims, _ := utils.ParseToken(request.Token)
	if claims == nil {
		return &DouYinRelationActionResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
		}, nil
	}
	if claims.ID == request.ToUserId {
		return &DouYinRelationActionResponse{
			StatusCode: -1,
			StatusMsg: "关注者ID无效",
		}, nil
	}

	db := dao.GetDB()
	user , err := model.FindUserById(&model.User{Id: request.ToUserId}, db)
	if (err != nil) || (user.Username == "") {
		return &DouYinRelationActionResponse{
			StatusCode: -1,
			StatusMsg: "关注者ID无效",
		}, nil
	}

	if (request.ActionType != 1) && (request.ActionType != 2) {
		return &DouYinRelationActionResponse{
			StatusCode: -1,
			StatusMsg: "关注操作代码无效",
		}, nil
	}

	tx := db.Begin()
	if (request.ActionType == 1) {
		// 添加关注记录
		_, err := model.Followed(&model.Follow{FollowerUserID: claims.ID, FollowedUserID: request.ToUserId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "关注失败",
			}, nil
		}
		// 若相互关注，则添加朋友表记录
		mutualFollow := model.IsFollowed(&model.Follow{FollowerUserID: request.ToUserId, FollowedUserID: claims.ID}, tx)
		if mutualFollow {
			friendID := FriendID(claims.ID, request.ToUserId)
			_, err = model.Friended(&model.Friend{Id: friendID, FromUserID: claims.ID, ToUserID: request.ToUserId}, tx)
			if err != nil {
				fmt.Println(err)
				tx.Rollback()
				return &DouYinRelationActionResponse{
					StatusCode: -2,
					StatusMsg: "关注失败",
				}, nil
			}
		}
		// 添加用户表记录
		_, err = model.AddFollowCount(&model.User{Id: claims.ID}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "关注失败",
			}, nil
		}
		_, err = model.AddFollowerCount(&model.User{Id: request.ToUserId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "关注失败",
			}, nil
		}

		err = tx.Commit().Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "关注失败",
			}, nil
		}
		return &DouYinRelationActionResponse{
			StatusCode: 0,
			StatusMsg: "关注成功",
		}, nil
	}else {
		// 添加关注记录
		_, err := model.Unfollowed(&model.Follow{FollowerUserID: claims.ID, FollowedUserID: request.ToUserId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "取消关注失败",
			}, nil
		}
		// 若相互关注，则添加朋友表记录
		mutualFollow := model.IsFollowed(&model.Follow{FollowerUserID: request.ToUserId, FollowedUserID: claims.ID}, tx)
		if mutualFollow {
			friendID := FriendID(claims.ID, request.ToUserId)
			_, err = model.Unfriended(&model.Friend{Id: friendID, FromUserID: claims.ID, ToUserID: request.ToUserId}, tx)
			if err != nil {
				fmt.Println(err)
				tx.Rollback()
				return &DouYinRelationActionResponse{
					StatusCode: -2,
					StatusMsg: "取消关注失败",
				}, nil
			}
		}
		// 添加用户表记录
		_, err = model.MinusFollowCount(&model.User{Id: claims.ID}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "取消关注失败",
			}, nil
		}
		_, err = model.MinusFollowerCount(&model.User{Id: request.ToUserId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "取消关注失败",
			}, nil
		}

		err = tx.Commit().Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinRelationActionResponse{
				StatusCode: -2,
				StatusMsg: "取消关注失败",
			}, nil
		}
		return &DouYinRelationActionResponse{
			StatusCode: 0,
			StatusMsg: "取消关注成功",
		}, nil
	}

}

func (rFW *relationFollowListService) RelationFollowList(context context.Context, request *DouYinRelationFollowListRequest) (*DouYinRelationFollowListResponse, error) {
	//实现具体的业务逻辑
	// 检查token是否有效
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) || (claims.ID != request.UserId) {
		return &DouYinRelationFollowListResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
			UserList: nil,
		}, nil
	}

	db := dao.GetDB()
	follows, err := model.FollowList(request.UserId, db)
	if err != nil {
		return &DouYinRelationFollowListResponse{
			StatusCode: -2,
			StatusMsg: "获取关注列表失败",
			UserList: nil,
		}, nil
	}

	var followList []*User
	for _, follow := range follows {
		user , err := model.FindUserById(&model.User{Id: follow.FollowedUserID}, db)
		if err != nil {
			return &DouYinRelationFollowListResponse{
				StatusCode: -2,
				StatusMsg: "获取关注列表用户失败",
				UserList: nil,
			}, nil
		}
		followList = append(followList, &User{
			Id: user.Id,
			Name: user.Username,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: true,
			Avatar: user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
		})
	}
	return &DouYinRelationFollowListResponse{
		StatusCode: 0,
		StatusMsg: "获取关注列表成功",
		UserList: followList,
	}, nil
}

func (rFR *relationFollowerListService) RelationFollowerList(context context.Context, request *DouYinRelationFollowerListRequest) (*DouYinRelationFollowerListResponse, error) {
	//实现具体的业务逻辑
	// 检查token是否有效
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) || (claims.ID != request.UserId) {
		return &DouYinRelationFollowerListResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
			UserList: nil,
		}, nil
	}

	db := dao.GetDB()
	followers, err := model.FollowerList(request.UserId, db)
	if err != nil {
		return &DouYinRelationFollowerListResponse{
			StatusCode: -2,
			StatusMsg: "获取粉丝列表失败",
			UserList: nil,
		}, nil
	}

	var followerList []*User
	for _, follow := range followers {
		user , err := model.FindUserById(&model.User{Id: follow.FollowerUserID}, db)
		if err != nil {
			return &DouYinRelationFollowerListResponse{
				StatusCode: -2,
				StatusMsg: "获取粉丝列表用户失败",
				UserList: nil,
			}, nil
		}
		isFollow := model.IsFollowed(&model.Follow{FollowerUserID: claims.ID, FollowedUserID: follow.FollowerUserID}, db)
		followerList = append(followerList, &User{
			Id: user.Id,
			Name: user.Username,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: isFollow,
			Avatar: user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
		})
	}
	return &DouYinRelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg: "获取粉丝列表成功",
		UserList: followerList,
	}, nil
}

func (rFD *relationFriendListService) RelationFriendList(context context.Context, request *DouYinRelationFriendListRequest) (*DouYinRelationFriendListResponse, error) {
	//实现具体的业务逻辑
	// 检查token是否有效
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) || (claims.ID != request.UserId) {
		return &DouYinRelationFriendListResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
			UserList: nil,
		}, nil
	}

	db := dao.GetDB()
	friendKey := fmt.Sprintf("_%d_", request.UserId)
	friends, err := model.FriendList(friendKey, db)
	if err != nil {
		return &DouYinRelationFriendListResponse{
			StatusCode: -2,
			StatusMsg: "获取好友列表失败",
			UserList: nil,
		}, nil
	}
	var friendList []*FriendUser
	for _, friend := range friends {
		var msgType int64
		var friendID int64
		if friend.FromUserID == request.UserId {
			friendID = friend.ToUserID
			msgType = 1
		}else {
			friendID = friend.FromUserID
			msgType = 0
		}
		user , err := model.FindUserById(&model.User{Id: friendID}, db)
		if err != nil {
			return &DouYinRelationFriendListResponse{
				StatusCode: -2,
				StatusMsg: "获取好友列表用户失败",
				UserList: nil,
			}, nil
		}
		friendList = append(friendList, &FriendUser{
			Id: user.Id,
			Name: user.Username,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: true,
			Avatar: user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
			Message: friend.LatestMessage,
			MsgType: msgType,
		})
	}
	return &DouYinRelationFriendListResponse{
		StatusCode: 0,
		StatusMsg: "获取朋友列表成功",
		UserList: friendList,
	}, nil
}

func (rA *relationActionService) mustEmbedUnimplementedRelationActionServiceServer() {}
func (rFW *relationFollowListService) mustEmbedUnimplementedRelationFollowListServiceServer() {}
func (rFR *relationFollowerListService) mustEmbedUnimplementedRelationFollowerListServiceServer() {}
func (rFD *relationFriendListService) mustEmbedUnimplementedRelationFriendListServiceServer() {}





