package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/clientconnect"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/commentservice"
)

var MaxRetry = 3

func CommentAction(c *gin.Context){
	var commentActionRequest CommentActionRequest
	err := c.ShouldBind(&commentActionRequest)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "评论信息输入有误",
			},
		})
		return
	}

	videoId, _ := strconv.ParseInt(commentActionRequest.VideoId, 10, 64)
	actionType, _ := strconv.ParseInt(commentActionRequest.ActionType, 10, 64)
	commentId, _ := strconv.ParseInt(commentActionRequest.CommentId, 10, 64)
	commentActionClient := <- clientconnect.CommentActionChan
	// commentActionResponse, err := commentActionClient.CommentAction(context.Background(), &commentservice.DouYinCommentActionRequest{Token: commentActionRequest.Token, VideoId: videoId, ActionType: int32(actionType), CommentId: commentId, CommentText: commentActionRequest.CommentText})
	// 超时重试
	var commentActionResponse *commentservice.DouYinCommentActionResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		commentActionResponse, err = commentActionClient.CommentAction(ctx, &commentservice.DouYinCommentActionRequest{Token: commentActionRequest.Token, VideoId: videoId, ActionType: int32(actionType), CommentId: commentId, CommentText: commentActionRequest.CommentText})
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

	clientconnect.CommentActionChan <- commentActionClient

	if (commentActionResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "comment action faild",
			},
		})
		return
	}

	if commentActionResponse.StatusMsg == "评论成功" {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: commentActionResponse.StatusCode,
				StatusMsg: commentActionResponse.StatusMsg,
			},
			Comment: Comment{
				Id: commentActionResponse.Comment.Id,
				User: User{
					Id: commentActionResponse.Comment.User.Id,
					Name: commentActionResponse.Comment.User.Name,
					FollowCount: commentActionResponse.Comment.User.FollowCount,
					FollowerCount: commentActionResponse.Comment.User.FollowerCount,
					IsFollow: commentActionResponse.Comment.User.IsFollow,
					Avatar: commentActionResponse.Comment.User.Avatar,
					BackgroundImage: commentActionResponse.Comment.User.BackgroundImage,
					Signature: commentActionResponse.Comment.User.Signature,
					TotalFavorited: commentActionResponse.Comment.User.TotalFavorited,
					WorkCount: commentActionResponse.Comment.User.WorkCount,
					FavoriteCount: commentActionResponse.Comment.User.FavoriteCount,
				},
				Content: commentActionResponse.Comment.Content,
				CreateDate: commentActionResponse.Comment.CreateDate,
			},
		})
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: Response{
			StatusCode: commentActionResponse.StatusCode,
			StatusMsg: commentActionResponse.StatusMsg,
		},
	})

}

func CommentList(c *gin.Context){
	var commentListRequest CommentListRequest
	err := c.ShouldBindQuery(&commentListRequest)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "用户身份信息输入有误",
			},
		})
		return
	}

	videoId, _ := strconv.ParseInt(commentListRequest.VideoId, 10, 64)
	commentListClient := <- clientconnect.CommentListChan
	// commentListResponse, err := commentListClient.CommentList(context.Background(), &commentservice.DouYinCommentListRequest{Token: commentListRequest.Token, VideoId: videoId})
	// 超时重试
	var commentListResponse *commentservice.DouYinCommentListResponse
	for try := 0; try < MaxRetry; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		commentListResponse, err = commentListClient.CommentList(ctx, &commentservice.DouYinCommentListRequest{Token: commentListRequest.Token, VideoId: videoId})
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
	
	clientconnect.CommentListChan <- commentListClient

	if (commentListResponse == nil) || (err != nil) {
		c.JSON(http.StatusOK, CommentListResponse{
			Response : Response{
				StatusCode: -1, 
				StatusMsg: "comment list failed",
			},
		})
		return
	}

	if commentListResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, CommentListResponse{
			Response : Response{
				StatusCode: commentListResponse.StatusCode, 
				StatusMsg: commentListResponse.StatusMsg,
			},
		})
		return
	}

	var commentList []Comment
	for _, comment := range commentListResponse.CommentList {
		commentList = append(commentList, Comment{
			Id: comment.Id,
			User: User{
				Id: comment.User.Id,
				Name: comment.User.Name,
				FollowCount: comment.User.FollowCount,
				FollowerCount: comment.User.FollowerCount,
				IsFollow: comment.User.IsFollow,
				Avatar: comment.User.Avatar,
				BackgroundImage: comment.User.BackgroundImage,
				Signature: comment.User.Signature,
				TotalFavorited: comment.User.TotalFavorited,
				WorkCount: comment.User.WorkCount,
				FavoriteCount: comment.User.FavoriteCount,
			},
			Content: comment.Content,
			CreateDate: comment.CreateDate,
		})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response : Response{
			StatusCode: commentListResponse.StatusCode, 
			StatusMsg: commentListResponse.StatusMsg,
		},
		CommentList: commentList,
	})
}