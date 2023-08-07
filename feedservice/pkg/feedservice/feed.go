package feedservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/micro/simplifiedTikTok/feedservice/pkg/model"
	"github.com/micro/simplifiedTikTok/feedservice/pkg/utils"
)

var FeedService = &feedService{}

type feedService struct {}

func (f *feedService) Feed(context context.Context, request *DouYinFeedRequest) (*DouYinFeedResponse, error) {
	//实现具体的业务逻辑
	// 检查token是否有效
	if request.Token != "" {
		claims, _ := utils.ParseToken(request.Token)
		if (claims == nil) {
			return &DouYinFeedResponse{
				StatusCode: -1,
				StatusMsg: "token无效",
				VideoList: nil,
				NextTime: time.Now().Unix(),
			}, nil
		}
	}

	//从redis中获取video缓存列表
	if request.LatestTime == 0 {
		len, err := model.GetVideoListSize("video")
		if len == 0 || err != nil {
			return &DouYinFeedResponse{
				StatusCode: -2,
				StatusMsg: "获取redis缓存失败",
				VideoList: nil,
				NextTime: time.Now().Unix(),
			}, nil
		}

		videos, err := model.GetVideoList("video")
		if len == 0 || err != nil {
			return &DouYinFeedResponse{
				StatusCode: -2,
				StatusMsg: "获取redis视频列表失败",
				VideoList: nil,
				NextTime: time.Now().Unix(),
			}, nil
		}
		var nextTime int64
		var videoList []*Video
		for i := 0; i < int(len); i++ {
			var video Video
			json.Unmarshal([]byte(videos[i]), &video)
			videoList = append(videoList, &video)

			if i == 0 {
				newVideo , err:= model.GetVideoById(video.Id)
				if err != nil {
					return &DouYinFeedResponse{
						StatusCode: -2,
						StatusMsg: "获取视频投稿时间失败",
						VideoList: nil,
						NextTime: time.Now().Unix(),
					}, nil
				}
				nextTime = newVideo.PublishTime
			}
		}

		return &DouYinFeedResponse{
			StatusCode: 0,
			StatusMsg: "获取视频列表成功",
			VideoList: videoList,
			NextTime: nextTime,
		}, nil
	}

	//查询mysql
	videos, err := model.ListVideoByTime(request.LatestTime)
	if err != nil {
		return &DouYinFeedResponse{
			StatusCode: -2,
			StatusMsg: "获取feed视频列表失败",
			VideoList: nil,
			NextTime: time.Now().Unix(),
		}, nil
	}
	var videoList []*Video
	for _, video := range *videos {
		user , err := model.FindUserById(&model.User{Id: video.AuthorId})
		if err != nil {
			return &DouYinFeedResponse{
				StatusCode: -2,
				StatusMsg: "获取feed视频列表时用户信息查询失败",
				VideoList: nil,
				NextTime: time.Now().Unix(),
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
			IsFavorite: video.IsFavorite,
			Title: video.Title,
		})
	}
	var nextTime int64
	if len(*videos) != 0 {
		nextTime = (*videos)[0].PublishTime
	}
	

	return &DouYinFeedResponse{
		StatusCode: 0,
		StatusMsg: "获取视频列表成功",
		VideoList: videoList,
		NextTime: nextTime,
	}, nil

}

func (f *feedService) mustEmbedUnimplementedFeedServiceServer() {}