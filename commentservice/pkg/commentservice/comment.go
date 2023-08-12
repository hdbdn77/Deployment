package commentservice

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/simplifiedTikTok/commentservice/pkg/dao"
	"github.com/micro/simplifiedTikTok/commentservice/pkg/model"
	"github.com/micro/simplifiedTikTok/commentservice/pkg/utils"
)

var CommentActionService = &commentActionService{}
var CommentListService = &commentListService{}

type commentActionService struct {}
type commentListService struct {}

func (cA *commentActionService) CommentAction(context context.Context, request *DouYinCommentActionRequest) (*DouYinCommentActionResponse, error) {
	//实现具体的业务逻辑
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) {
		return &DouYinCommentActionResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
		}, nil
	}

	db := dao.GetDB()
	video, err := model.GetVideoById(request.VideoId, db)
	if (err != nil) || (video.PlayUrl == "")  {
		return &DouYinCommentActionResponse{
			StatusCode: -1,
			StatusMsg: "videoId无效",
		}, nil
	}

	if (request.ActionType != 1) && (request.ActionType != 2) {
		return &DouYinCommentActionResponse{
			StatusCode: -1,
			StatusMsg: "评论操作代码无效",
		}, nil
	}

	tx := db.Begin()
	if (request.ActionType == 1) {
		// 添加评论表记录
		comment, err := model.Commented(&model.Comment{UserID: claims.ID, VideoID: request.VideoId, Content: request.CommentText}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinCommentActionResponse{
				StatusCode: -2,
				StatusMsg: "评论失败",
			}, nil
		}
		// 添加视频表记录
		_, err = model.AddVideoCommentCount(video, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinCommentActionResponse{
				StatusCode: -2,
				StatusMsg: "评论失败",
			}, nil
		}

		err = tx.Commit().Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinCommentActionResponse{
				StatusCode: -2,
				StatusMsg: "评论失败",
			}, nil
		}
		user, _ := model.FindUserById(&model.User{Id: claims.ID}, db)
		createDate := time.Unix(comment.CreateDate, 0).Format("01-02")
		return &DouYinCommentActionResponse{
			StatusCode: 0,
			StatusMsg: "评论成功",
			Comment: &Comment{
				Id: comment.Id,
				User: &User{
					Id: user.Id,
					Name: user.Username,
					FollowCount: user.FollowCount,
					FollowerCount: user.FollowerCount,
					IsFollow: user.IsFollow,
					Avatar: user.Avatar,
					BackgroundImage: user.BackgroundImage,
					Signature: user.Signature,
					TotalFavorited: user.TotalFavorited,
					WorkCount: user.WorkCount,
					FavoriteCount: user.FavoriteCount,
					},
				Content: comment.Content,
				CreateDate: createDate,
			},
		}, nil
	}else {
		_, err := model.UnCommented(&model.Comment{Id: request.CommentId, UserID: claims.ID, VideoID: request.VideoId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinCommentActionResponse{
				StatusCode: -2,
				StatusMsg: "删除评论失败",
			}, nil
		}
		// 添加视频表记录
		_, err = model.MinusVideoCommentCount(video, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinCommentActionResponse{
				StatusCode: -2,
				StatusMsg: "删除评论失败",
			}, nil
		}

		err = tx.Commit().Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinCommentActionResponse{
				StatusCode: -2,
				StatusMsg: "删除评论失败",
			}, nil
		}
		return &DouYinCommentActionResponse{
			StatusCode: 0,
			StatusMsg: "删除评论成功",
		}, nil
	}
}

func (cL *commentListService) CommentList(context context.Context, request *DouYinCommentListRequest) (*DouYinCommentListResponse, error) {
	//实现具体的业务逻辑
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) {
		return &DouYinCommentListResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
		}, nil
	}

	db := dao.GetDB()
	video, err := model.GetVideoById(request.VideoId, db)
	if (err != nil) || (video.PlayUrl == "")  {
		return &DouYinCommentListResponse{
			StatusCode: -1,
			StatusMsg: "videoId无效",
		}, nil
	}

	comments, err := model.ListComment(&model.Comment{VideoID: request.VideoId}, db)
	if err != nil {
		return &DouYinCommentListResponse{
			StatusCode: -2,
			StatusMsg: "查询评论列表失败",
		}, nil
	}

	var commentList []*Comment
	for _, comment := range *comments {
		user, err := model.FindUserById(&model.User{Id: comment.UserID}, db)
		if err != nil {
			return &DouYinCommentListResponse{
				StatusCode: -2,
				StatusMsg: "查询评论列表时获取用户失败",
			}, nil
		} 
		commentList = append(commentList, &Comment{
			Id: comment.Id,
			User: &User{
				Id: user.Id,
				Name: user.Username,
				FollowCount: user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow: user.IsFollow,
				Avatar: user.Avatar,
				BackgroundImage: user.BackgroundImage,
				Signature: user.Signature,
				TotalFavorited: user.TotalFavorited,
				WorkCount: user.WorkCount,
				FavoriteCount: user.FavoriteCount,
			},
			Content: comment.Content,
			CreateDate: time.Unix(comment.CreateDate, 0).Format("01-02"),
		})
	}
	
	return &DouYinCommentListResponse{
		StatusCode: 0,
		StatusMsg: "查询评论列表成功",
		CommentList: commentList,
	}, nil
	
}

func (cA *commentActionService) mustEmbedUnimplementedCommentActionServiceServer() {}

func (cL *commentListService) mustEmbedUnimplementedCommentListServiceServer() {}