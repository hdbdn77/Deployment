package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/clientconnect"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/relationservice"
)

func RelationAction(c *gin.Context){
	var relationActionRequest RelationActionRequest
	err := c.ShouldBind(&relationActionRequest)
	if err != nil {
		c.JSON(http.StatusOK, RelationActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "关注信息输入有误",
			},
		})
		return
	}

	toUserID, _ := strconv.ParseInt(relationActionRequest.ToUserID, 10, 64)
	actionType, _ := strconv.ParseInt(relationActionRequest.ActionType, 10, 64)
	relationActionClient := <- clientconnect.RelationActionChan
	// relationActionResponse, err := relationActionClient.RelationAction(context.Background(), &relationservice.DouYinRelationActionRequest{Token: relationActionRequest.Token, ToUserId: toUserID, ActionType: int32(actionType)})
	// 超时重试
	var relationActionResponse *relationservice.DouYinRelationActionResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		relationActionResponse, err = relationActionClient.RelationAction(ctx, &relationservice.DouYinRelationActionRequest{Token: relationActionRequest.Token, ToUserId: toUserID, ActionType: int32(actionType)})
		if err != nil {
			if err == context.DeadlineExceeded {
			  // 超时,可以重试继续
			  continue
			} else {
			  // 其他错误,不重试
			  break 
			}   
		}else {
			break
		}
	}
	
	clientconnect.RelationActionChan <- relationActionClient

	if (relationActionResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, RelationActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "relation action faild",
			},
		})
		return
	}

	c.JSON(http.StatusOK, RelationActionResponse{
		Response: Response{
			StatusCode: relationActionResponse.StatusCode,
			StatusMsg: relationActionResponse.StatusMsg,
		},
	})
}

func FollowList(c *gin.Context){
	var followListRequest FollowListRequest
	err := c.ShouldBindQuery(&followListRequest)
	if err != nil {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "关注列表信息输入有误",
			},
		})
		return
	}

	UserID, _ := strconv.ParseInt(followListRequest.UserId, 10, 64)
	followListClient := <- clientconnect.FollowListChan
	// followListResponse, err := followListClient.RelationFollowList(context.Background(), &relationservice.DouYinRelationFollowListRequest{UserId: UserID, Token: followListRequest.Token})
	// 超时重试
	var followListResponse *relationservice.DouYinRelationFollowListResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		followListResponse, err = followListClient.RelationFollowList(ctx, &relationservice.DouYinRelationFollowListRequest{UserId: UserID, Token: followListRequest.Token})
		if err != nil {
			if err == context.DeadlineExceeded {
			  // 超时,可以重试继续
			  continue
			} else {
			  // 其他错误,不重试
			  break 
			}   
		}else {
			break
		}
	}
	
	clientconnect.FollowListChan <- followListClient

	if (followListResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "follow list faild",
			},
		})
		return
	}

	if followListResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, FollowListResponse{
			Response: Response{
				StatusCode: followListResponse.StatusCode,
				StatusMsg: followListResponse.StatusMsg,
			},
		})
		return
	}

	var userList []User
	for _, user := range followListResponse.UserList {
		userList = append(userList, User{
			Id: user.Id,
			Name: user.Name,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: user.IsFollow,
			Avatar: user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
		})
	}
	c.JSON(http.StatusOK, FollowListResponse{
		Response: Response{
			StatusCode: followListResponse.StatusCode,
			StatusMsg: followListResponse.StatusMsg,
		},
		UserList: userList,
	})
}

func FollowerList(c *gin.Context){
	var followerListRequest FollowerListRequest
	err := c.ShouldBindQuery(&followerListRequest)
	if err != nil {
		c.JSON(http.StatusOK, FollowerListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "粉丝列表信息输入有误",
			},
		})
		return
	}

	UserID, _ := strconv.ParseInt(followerListRequest.UserId, 10, 64)
	followerListClient := <- clientconnect.FollowerListChan
	// followerListResponse, err := followerListClient.RelationFollowerList(context.Background(), &relationservice.DouYinRelationFollowerListRequest{UserId: UserID, Token: followerListRequest.Token})
	// 超时重试
	var followerListResponse *relationservice.DouYinRelationFollowerListResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		followerListResponse, err = followerListClient.RelationFollowerList(ctx, &relationservice.DouYinRelationFollowerListRequest{UserId: UserID, Token: followerListRequest.Token})
		if err != nil {
			if err == context.DeadlineExceeded {
			  // 超时,可以重试继续
			  continue
			} else {
			  // 其他错误,不重试
			  break 
			}   
		}else {
			break
		}
	}
	
	clientconnect.FollowerListChan <- followerListClient

	if (followerListResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, FollowerListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "follower list faild",
			},
		})
		return
	}

	if followerListResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, FollowerListResponse{
			Response: Response{
				StatusCode: followerListResponse.StatusCode,
				StatusMsg: followerListResponse.StatusMsg,
			},
		})
		return
	}

	var userList []User
	for _, user := range followerListResponse.UserList {
		userList = append(userList, User{
			Id: user.Id,
			Name: user.Name,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: user.IsFollow,
			Avatar: user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
		})
	}
	c.JSON(http.StatusOK, FollowerListResponse{
		Response: Response{
			StatusCode: followerListResponse.StatusCode,
			StatusMsg: followerListResponse.StatusMsg,
		},
		UserList: userList,
	})
}

func FriendList(c *gin.Context){
	var friendListRequest FriendListRequest
	err := c.ShouldBindQuery(&friendListRequest)
	if err != nil {
		c.JSON(http.StatusOK, FriendListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "关注列表信息输入有误",
			},
		})
		return
	}

	UserID, _ := strconv.ParseInt(friendListRequest.UserId, 10, 64)
	friendListClient := <- clientconnect.FriendListChan
	// friendListResponse, err := friendListClient.RelationFriendList(context.Background(), &relationservice.DouYinRelationFriendListRequest{UserId: UserID, Token: friendListRequest.Token})
	// 超时重试
	var friendListResponse *relationservice.DouYinRelationFriendListResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		friendListResponse, err = friendListClient.RelationFriendList(ctx, &relationservice.DouYinRelationFriendListRequest{UserId: UserID, Token: friendListRequest.Token})
		if err != nil {
			if err == context.DeadlineExceeded {
			  // 超时,可以重试继续
			  continue
			} else {
			  // 其他错误,不重试
			  break 
			}   
		}else {
			break
		}
	}
	
	clientconnect.FriendListChan <- friendListClient

	if (friendListResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, FriendListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "friend list faild",
			},
		})
		return
	}

	if friendListResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, FriendListResponse{
			Response: Response{
				StatusCode: friendListResponse.StatusCode,
				StatusMsg: friendListResponse.StatusMsg,
			},
		})
		return
	}

	var friendUserList []FriendUser
	for _, user := range friendListResponse.UserList {
		friendUserList = append(friendUserList, FriendUser{
			Id: user.Id,
			Name: user.Name,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: user.IsFollow,
			Avatar: user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
			Message: user.Message,
			MsgType: user.MsgType,
		})
	}
	c.JSON(http.StatusOK, FriendListResponse{
		Response: Response{
			StatusCode: friendListResponse.StatusCode,
			StatusMsg: friendListResponse.StatusMsg,
		},
		FriendUserList: friendUserList,
	})
}