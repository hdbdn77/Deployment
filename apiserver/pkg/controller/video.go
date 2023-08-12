package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"io"
	"bytes"
	"strings"

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
	publishActionResponse, err := publishActionServiceClient.PublishAction(context.Background(), &videoservice.DouYinPublishActionRequest{Token: publishRequest.Token, Data: bytesData, Title: publishRequest.Title})
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
	publishListResponse, err := publishListServiceClient.PublishList(context.Background(), &videoservice.DouYinPublishListRequest{Token: publishListRequest.Token, UserId: userId})
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
			PlayUrl: "http://121.41.85.100:8080/" + video.PlayUrl,
			CoverUrl: "http://5b0988e595225.cdn.sohucs.com/images/20180430/fcf555aed1804ad586b24b3aeda6c031.jpeg",
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