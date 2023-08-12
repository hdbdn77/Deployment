package userservice

import (
	"context"
	"fmt"

	"github.com/micro/simplifiedTikTok/userservice/pkg/model"
	"github.com/micro/simplifiedTikTok/userservice/pkg/utils"
	"github.com/micro/simplifiedTikTok/userservice/pkg/dao"
)

var RegisterService = &registerService{}
var LoginService = &loginService{}
var UserService = &userService{}

type registerService struct {}
type loginService struct {}
type userService struct {}

func (r *registerService) Register(context context.Context, request *DouYinUserRegisterRequest) (*DouYinUserRegisterResponse, error){
	//实现具体的业务逻辑
	db := dao.GetDB()
	user , err := model.Register(&model.User{Username: request.Username, Password: request.Password}, db)
	token, _ := utils.GenToken(user.Id, user.Username)
	if err != nil {
		fmt.Println(err)
		return &DouYinUserRegisterResponse{
			StatusCode: -1,
			StatusMsg: "注册失败",
			UserId: -1,
			Token: "",
		}, nil
	}

	return &DouYinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg: "注册成功",
		UserId: user.Id,
		Token: token,
	}, nil

}

func (l *loginService) Login(context context.Context, request *DouYinUserLoginRequest) (*DouYinUserLoginResponse, error){
	db := dao.GetDB()
	user , err := model.FindUserByUsername(&model.User{Username: request.Username, Password: request.Password}, db)
	if err != nil {
		fmt.Println(err)
		return &DouYinUserLoginResponse{
			StatusCode: -1,
			StatusMsg: "登陆失败",
			UserId: -1,
			Token: "",
		}, nil
	}
	token, _ := utils.GenToken(user.Id, user.Username)
	if request.Password != user.Password {
		fmt.Println(err)
		return &DouYinUserLoginResponse{
			StatusCode: -2,
			StatusMsg: "密码错误",
			UserId: -1,
			Token: "",
		}, nil
	}

	return &DouYinUserLoginResponse{
		StatusCode: 0,
		StatusMsg: "登陆成功",
		UserId: user.Id,
		Token: token,
	}, nil
}

func (u *userService) Find(context context.Context, request *DouYinUserRequest) (*DouYinUserResponse, error){
	claims, _ := utils.ParseToken(request.Token)
	if (claims == nil) || (claims.ID != request.UserId)  {
		return &DouYinUserResponse{
			StatusCode: -1,
			StatusMsg: "token无效",
			User: nil,
		}, nil
	}
	db := dao.GetDB()
	user , err := model.FindUserByUsername(&model.User{Username: claims.Username}, db)
	if err != nil {
		fmt.Println(err)
		return &DouYinUserResponse{
			StatusCode: -1,
			StatusMsg: "用户信息查询失败",
			User: nil,
			
		}, nil
	}
	return &DouYinUserResponse{
		StatusCode: 0,
		StatusMsg: "用户信息查询成功",
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
		
	}, nil

}

func (r *registerService) mustEmbedUnimplementedRegisterServiceServer() {}

func (l *loginService) mustEmbedUnimplementedLoginServiceServer() {}

func (u *userService) mustEmbedUnimplementedUserServiceServer() {}