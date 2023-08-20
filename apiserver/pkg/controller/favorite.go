package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	// favoriteActionResponse, err := favoriteActionClient.FavoriteAction(context.Background(), &favoriteservice.DouYinFavoriteActionRequest{Token: favoriteActionRequest.Token, VideoId: videoId, ActionType: int32(actionType)})
	// 超时重试
	var favoriteActionResponse *favoriteservice.DouYinFavoriteActionResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		favoriteActionResponse, err = favoriteActionClient.FavoriteAction(ctx, &favoriteservice.DouYinFavoriteActionRequest{Token: favoriteActionRequest.Token, VideoId: videoId, ActionType: int32(actionType)})
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
	// favoriteListResponse, err := favoriteListClient.FavoriteList(context.Background(), &favoriteservice.DouYinFavoriteListRequest{UserId: userId, Token: favoriteListRequest.Token})
	// 超时重试
	var favoriteListResponse *favoriteservice.DouYinFavoriteListResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		favoriteListResponse, err = favoriteListClient.FavoriteList(ctx, &favoriteservice.DouYinFavoriteListRequest{UserId: userId, Token: favoriteListRequest.Token})
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
		if strings.HasPrefix(video.PlayUrl, "static/") {
			video.PlayUrl = "http://121.41.85.100:30808/" + video.PlayUrl
		}
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
			PlayUrl: video.PlayUrl,
			CoverUrl: video.CoverUrl,
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