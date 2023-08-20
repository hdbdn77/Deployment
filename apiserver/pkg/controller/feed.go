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
	"github.com/micro/simplifiedTikTok/apiserver/pkg/feedservice"
)

func Feed(c *gin.Context) {
	var feedRequest FeedRequest
	err := c.ShouldBindQuery(&feedRequest)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "视频流信息输入有误",
			},
		})
	}

	latestTime, _ := strconv.ParseInt(feedRequest.LatestTime, 10, 64)
	feedServiceClient := <- clientconnect.FeedChan
	// feedResponse, err := feedServiceClient.Feed(context.Background(), &feedservice.DouYinFeedRequest{LatestTime: latestTime, Token: feedRequest.Token})
	// 超时重试
	var feedResponse *feedservice.DouYinFeedResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		feedResponse, err = feedServiceClient.Feed(ctx, &feedservice.DouYinFeedRequest{LatestTime: latestTime, Token: feedRequest.Token})
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

	clientconnect.FeedChan <- feedServiceClient

	if (feedResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "feed failed",
			},
		})
		return
	}
	if feedResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: feedResponse.StatusCode,
				StatusMsg: feedResponse.StatusMsg,
			},
			NextTime: feedResponse.NextTime,
		})
		return
	}
	var videoList []Video
	for _, video := range feedResponse.VideoList {
		if strings.HasPrefix(video.PlayUrl, "static/") {
			video.PlayUrl = "http://121.41.85.100:30808/" + video.PlayUrl
		}
		videoList = append(videoList, Video{
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
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: feedResponse.StatusCode,
			StatusMsg: feedResponse.StatusMsg,
		},
		NextTime: feedResponse.NextTime,
		VideoList: videoList,
	})
}