package favoriteservice

import (
	"context"
	"fmt"

	"github.com/micro/simplifiedTikTok/favoriteservice/pkg/model"
	"github.com/micro/simplifiedTikTok/favoriteservice/pkg/utils"
	"github.com/micro/simplifiedTikTok/favoriteservice/pkg/dao"
)

var FavoriteActionService = &favoriteActionService{}
var FavoriteListService = &favoriteListService{}

type favoriteActionService struct {}

type favoriteListService struct {}

func (fA *favoriteActionService) FavoriteAction(context context.Context, request *DouYinFavoriteActionRequest) (*DouYinFavoriteActionResponse, error) {
	//实现具体的业务逻辑
	// 检查token是否有效
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) {
		return &DouYinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
		}, nil
	}

	db := dao.GetDB()
	video, err := model.GetVideoById(request.VideoId, db)
	if (err != nil) || (video.PlayUrl == "")  {
		return &DouYinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg: "videoId无效",
		}, nil
	}

	if (request.ActionType != 1) && (request.ActionType != 2) {
		return &DouYinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg: "点赞操作代码无效",
		}, nil
	}

	tx := db.Begin()
	if (request.ActionType == 1) {
		// 添加点赞表记录
		_, err := model.Favorited(&model.Favorite{UserID: claims.ID, VideoID: request.VideoId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "点赞失败",
			}, nil
		}
		// 添加视频表记录
		_, err = model.AddViseoFavoriteCount(video, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "点赞失败",
			}, nil
		}
		// 添加用户表记录
		_, err = model.AddTotalFavorited(&model.User{Id: video.AuthorId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "点赞失败",
			}, nil
		}
		_, err = model.AddUserFavoriteCount(&model.User{Id: claims.ID}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "点赞失败",
			}, nil
		}

		err = tx.Commit().Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "点赞失败",
			}, nil
		}
		return &DouYinFavoriteActionResponse{
			StatusCode: 0,
			StatusMsg: "点赞成功",
		}, nil
	}else {
		_, err := model.Unfavorited(&model.Favorite{UserID: claims.ID, VideoID: request.VideoId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "取消点赞失败",
			}, nil
		}
		// 添加视频表记录
		video, err := model.MinusViseoFavoriteCount(&model.Video{Id: request.VideoId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "取消点赞失败",
			}, nil
		}
		// 添加用户表记录
		_, err = model.MinusTotalFavorited(&model.User{Id: video.AuthorId}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "取消点赞失败",
			}, nil
		}
		_, err = model.MinusUserFavoriteCount(&model.User{Id: claims.ID}, tx)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "取消点赞失败",
			}, nil
		}

		err = tx.Commit().Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return &DouYinFavoriteActionResponse{
				StatusCode: -2,
				StatusMsg: "取消点赞失败",
			}, nil
		}
		return &DouYinFavoriteActionResponse{
			StatusCode: 0,
			StatusMsg: "取消点赞成功",
		}, nil
	}

}

func (fL *favoriteListService) FavoriteList(context context.Context, request *DouYinFavoriteListRequest) (*DouYinFavoriteListResponse, error) {
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) || (claims.ID != request.UserId)  {
		return &DouYinFavoriteListResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
			VideoList: nil,
		}, nil
	}

	db := dao.GetDB()

	favorites, err := model.GetFavoritesByUserID(request.UserId, db)
	if err != nil {
		return &DouYinFavoriteListResponse{
			StatusCode: -2,
			StatusMsg: "获取点赞列表失败",
			VideoList: nil,
		}, nil
	}
	
	var videoList []*Video
	for _, favorite := range favorites {
		video, err := model.GetVideoById(favorite.VideoID, db)
		if err != nil {
			return &DouYinFavoriteListResponse{
				StatusCode: -2,
				StatusMsg: "获取点赞列表视频失败",
				VideoList: nil,
			}, nil
		}
		user , err := model.FindUserById(&model.User{Id: video.AuthorId}, db)
		if err != nil {
			return &DouYinFavoriteListResponse{
				StatusCode: -2,
				StatusMsg: "获取点赞列表视频作者失败",
				VideoList: nil,
			}, nil
		}
		videoList = append(videoList, &Video{
			Id: video.Id,
			Author: &User{
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
			PlayUrl: video.PlayUrl,
			CoverUrl: video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount: video.CommentCount,
			IsFavorite: true,
			Title: video.Title,
		})
	}

	return &DouYinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg: "获取点赞列表成功",
		VideoList: videoList,
	}, nil

}

func (fA *favoriteActionService) mustEmbedUnimplementedFavoriteActionServiceServer() {}

func (fL *favoriteListService) mustEmbedUnimplementedFavoriteListServiceServer() {}