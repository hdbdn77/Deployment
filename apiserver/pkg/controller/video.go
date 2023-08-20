package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"io"
	"bytes"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/clientconnect"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/videoservice"
)

func Publish(c *gin.Context) {
	var publishRequest PublishRequest
	err := c.ShouldBind(&publishRequest)
	if err != nil {
		c.JSON(http.StatusOK, PublishResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "投稿信息输入有误",
			},
		})
		return
	}

	dataBytes, _ := publishRequest.Data.Open()
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, dataBytes)
	bytesData := buf.Bytes()
	publishActionServiceClient := <- clientconnect.PublishActionChan
	// publishActionResponse, err := publishActionServiceClient.PublishAction(context.Background(), &videoservice.DouYinPublishActionRequest{Token: publishRequest.Token, Data: bytesData, Title: publishRequest.Title})
	// 超时重试
	var publishActionResponse *videoservice.DouYinPublishActionResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		publishActionResponse, err = publishActionServiceClient.PublishAction(ctx, &videoservice.DouYinPublishActionRequest{Token: publishRequest.Token, Data: bytesData, Title: publishRequest.Title})
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
	
	clientconnect.PublishActionChan <- publishActionServiceClient

	if (publishActionResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, PublishResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "publish action failed",
			},
		})
		return
	}
	if publishActionResponse.StatusCode == 0 {
		fileName := strings.Replace(publishActionResponse.StatusMsg, "保存视频成功:", "", 1)
		go c.SaveUploadedFile(publishRequest.Data, "static/" + fileName)
	}
	
	c.JSON(http.StatusOK, Response{
		StatusCode: publishActionResponse.StatusCode,
		StatusMsg: publishActionResponse.StatusMsg,
	})

}

func PublishList(c *gin.Context) {
	var publishListRequest PublishListRequest
	err := c.ShouldBindQuery(&publishListRequest)
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "用户身份信息输入有误",
			},
		})
		return
	}

	userId, _ := strconv.ParseInt(publishListRequest.UserId, 10, 64)
	publishListServiceClient := <- clientconnect.PublishListChan
	// publishListResponse, err := publishListServiceClient.PublishList(context.Background(), &videoservice.DouYinPublishListRequest{Token: publishListRequest.Token, UserId: userId})
	// 超时重试
	var publishListResponse *videoservice.DouYinPublishListResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		publishListResponse, err = publishListServiceClient.PublishList(ctx, &videoservice.DouYinPublishListRequest{Token: publishListRequest.Token, UserId: userId})
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
	
	clientconnect.PublishListChan <- publishListServiceClient

	if (publishListResponse == nil) || (err != nil) {
		c.JSON(http.StatusOK, PublishListResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "publish list failed",
			},
		})
		return
	}

	if publishListResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, PublishListResponse{
			Response : Response{
				StatusCode: publishListResponse.StatusCode, 
				StatusMsg: publishListResponse.StatusMsg,
			},
		})
		return
	}

	var publishList []Video
	for _, video := range publishListResponse.VideoList {
		if strings.HasPrefix(video.PlayUrl, "static/") {
			video.PlayUrl = "http://121.41.85.100:30808/" + video.PlayUrl
		}
		publishList = append(publishList, Video{
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
	c.JSON(http.StatusOK, PublishListResponse{
		Response : Response{
			StatusCode: publishListResponse.StatusCode, 
			StatusMsg: publishListResponse.StatusMsg,
		},
		VideoList: publishList,
	})
}