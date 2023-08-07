package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/micro/simplifiedTikTok/userservice/configs"
	"time"
)

var privKey = []byte(configs.Jwt.PrivKey)

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenToken(id int64, username string) (string, error) {
	// 创建一个我们自己的声明的数据
	claims := Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "test",
			Audience: []string{"zyx"},
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(privKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return privKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 令牌有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
