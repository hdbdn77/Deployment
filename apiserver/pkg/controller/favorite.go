package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/clientconnect"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/favoriteservice"
)

func FavoriteAction(c *gin.Context) {
	var favoriteActionRequest FavoriteActionRequest
	err := c.ShouldBind(&favoriteActionRequest)
	if err != nil {
		c.JSON(http.StatusOK, FavoriteActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "点赞信息输入有误",
			},
		})
		return
	}

	videoId, _ := strconv.ParseInt(favoriteActionRequest.VideoId, 10, 64)
	actionType, _ := strconv.ParseInt(favoriteActionRequest.ActionType, 10, 64)
	favoriteActionClient := <- clientconnect.FavoriteActionChan
	favoriteActionResponse, err := favoriteActionClient.FavoriteAction(context.Background(), &favoriteservice.DouYinFavoriteActionRequest{Token: favoriteActionRequest.Token, VideoId: videoId, ActionType: int32(actionType)})
	clientconnect.FavoriteActionChan <- favoriteActionClient

	if (favoriteActionResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, FavoriteActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "favorite action faild",
			},
		})
		return
	}

	c.JSON(http.StatusOK, FavoriteActionResponse{
		Response: Response{
			StatusCode: favoriteActionResponse.StatusCode,
			StatusMsg: favoriteActionResponse.StatusMsg,
		},
	})
}
func FavoriteList(c *gin.Context) {
	var favoriteListRequest FavoriteListRequest
	err := c.ShouldBindQuery(&favoriteListRequest)
	if err != nil {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "用户身份信息输入有误",
			},
		})
		return
	}

	userId, _ := strconv.ParseInt(favoriteListRequest.UserId, 10, 64)
	favoriteListClient := <- clientconnect.FavoriteListChan
	favoriteListResponse, err := favoriteListClient.FavoriteList(context.Background(), &favoriteservice.DouYinFavoriteListRequest{UserId: userId, Token: favoriteListRequest.Token})
	clientconnect.FavoriteListChan <- favoriteListClient

	if (favoriteListResponse == nil) || (err != nil) {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "favorite list failed",
			},
		})
		return
	}

	if favoriteListResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response : Response{
				StatusCode: favoriteListResponse.StatusCode, 
				StatusMsg: favoriteListResponse.StatusMsg,
			},
		})
		return
	}
	
	var favoriteList []Video
	for _, video:= range favoriteListResponse.VideoList {
		favoriteList = append(favoriteList, Video{
			Id: video.Id,
			Author: User{
				Id: video.Author.Id,
				Name: video.Author.Name,
				FollowCount: video.Author.FollowCount,
				FollowerCount: video.Author.FollowerCount,
				IsFollow: video.Author.IsFollow,
				Avatar: video.Author.Avatar,
				BackgroundImage: video.Author.BackgroundImage,
				Signature: video.Author.Signature,
				TotalFavorited: video.Author.TotalFavorited,
				WorkCount: video.Author.WorkCount,
				FavoriteCount: video.Author.FavoriteCount,
			},
			PlayUrl: "http://localhost:8080/" + video.PlayUrl,
			CoverUrl: "http://5b0988e595225.cdn.sohucs.com/images/20180430/fcf555aed1804ad586b24b3aeda6c031.jpeg",
			FavoriteCount: video.FavoriteCount,
			CommentCount: video.CommentCount,
			IsFavorite: video.IsFavorite,
			Title: video.Title,
		})
	}
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response : Response{
			StatusCode: favoriteListResponse.StatusCode, 
			StatusMsg: favoriteListResponse.StatusMsg,
		},
		VideoList: favoriteList,
	})
}