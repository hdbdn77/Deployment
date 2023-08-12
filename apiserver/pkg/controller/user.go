package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/clientconnect"
	"github.com/micro/simplifiedTikTok/apiserver/pkg/userservice"
)

func Register(c *gin.Context) {
	var userLoginRequest UserLoginRequest
	err := c.ShouldBind(&userLoginRequest)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "用户名或密码输入有误",
			},
		})
		return
	}
	if len(userLoginRequest.Username) >32 || len(userLoginRequest.Password) >32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "用户名或密码长度超过32",
			},
		})
		return
	}

	registerServiceClient := <- clientconnect.UserRegisterChan
	registerResponse, err := registerServiceClient.Register(context.Background(), &userservice.DouYinUserRegisterRequest{Username: userLoginRequest.Username, Password: userLoginRequest.Password})
	clientconnect.UserRegisterChan <- registerServiceClient
	
	if (registerResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "register failed",
			},
		})
		return
	}
	if registerResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: registerResponse.StatusCode,
				StatusMsg: registerResponse.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: registerResponse.StatusCode, 
			StatusMsg: registerResponse.StatusMsg,
		},
		UserId:   registerResponse.UserId,
		Token:    registerResponse.Token,
	})

}

func Login(c *gin.Context) {
	var userLoginRequest UserLoginRequest
	err := c.ShouldBind(&userLoginRequest)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "用户名或密码输入有误",
			},
		})
		return
	}

	loginServiceClient := <- clientconnect.UserLoginChan
	loginResponse, err := loginServiceClient.Login(context.Background(), &userservice.DouYinUserLoginRequest{Username: userLoginRequest.Username, Password: userLoginRequest.Password})
	clientconnect.UserLoginChan <- loginServiceClient
	
	if (loginResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "login failed",
			},
		})
		return
	}
	if loginResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: loginResponse.StatusCode,
				StatusMsg: loginResponse.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: loginResponse.StatusCode,
			StatusMsg: loginResponse.StatusMsg,
		},
		UserId:   loginResponse.UserId,
		Token:    loginResponse.Token,
	})
}

func UserInfo(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBindQuery(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "用户身份信息输入有误",
			},
		})
		return
	}

	userId, _ := strconv.ParseInt(userInfoRequest.UserId, 10, 64)
	userServiceClient := <- clientconnect.UserChan
	userResponse, err := userServiceClient.Find(context.Background(), &userservice.DouYinUserRequest{UserId: userId, Token: userInfoRequest.Token})
	clientconnect.UserChan <- userServiceClient
	
	if (userResponse == nil) || (err != nil) {
		fmt.Println(err)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg: "user info failed",
			},
		})
		return
	}
	if userResponse.StatusCode != 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: userResponse.StatusCode,
				StatusMsg: userResponse.StatusMsg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{
			StatusCode: userResponse.StatusCode,
			StatusMsg: userResponse.StatusMsg,
		},
		User: User{
			Id: userResponse.User.Id,
			Name: userResponse.User.Name,
			FollowCount: userResponse.User.FollowCount,
			FollowerCount: userResponse.User.FollowerCount,
			IsFollow: userResponse.User.IsFollow,
			Avatar: userResponse.User.Avatar,
			BackgroundImage: userResponse.User.BackgroundImage,
			Signature: userResponse.User.Signature,
			TotalFavorited: userResponse.User.TotalFavorited,
			WorkCount: userResponse.User.WorkCount,
			FavoriteCount: userResponse.User.FavoriteCount,
		},
	})
}