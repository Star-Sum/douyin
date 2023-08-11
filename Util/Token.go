package Util

import (
	"douyin/Log"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// TokenExpireDuration token有效时长为1天
const TokenExpireDuration = time.Hour * 24

// Secret 密钥
var Secret = []byte("未来之星")

type BaseInfo struct {
	UserId int64
	jwt.StandardClaims
}

// MakeToken 根据用户Uid生成Token
func MakeToken(userId int64) string {
	baseInfo := BaseInfo{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "lite-douyin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, baseInfo)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		Log.ErrorLogWithoutPanic("Token creation failed!", err)
	}
	Log.NormalLog("Token created successfully!", err)
	return tokenString
}

// ParserToken 解析Token
func ParserToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &BaseInfo{},
		func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		})
	if err == nil && token.Valid {
		Log.NormalLog("Token is valid!", err)
		return token.Claims.(*BaseInfo).UserId, nil
	} else {
		Log.ErrorLogWithoutPanic("Invalid token!", err)
		return 0, err
	}
}

// TokenJudge 判断token是否与目标uid相等
func TokenJudge(tokenString string, UID int64) (bool, error) {
	uidJudge, err := ParserToken(tokenString)
	if err != nil {
		Log.ErrorLogWithoutPanic("Token parsing failed!", err)
		return false, err
	}
	Log.NormalLog("Token parsing successful!", err)
	if uidJudge == UID {
		return true, nil
	} else {
		return false, nil
	}
}
