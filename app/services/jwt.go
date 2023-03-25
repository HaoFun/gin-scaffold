package services

import (
	"context"
	"gin-scaffold/global"
	"gin-scaffold/locales"
	"gin-scaffold/utils"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type jwtService struct{}

type JwtUser interface {
	GetId() int64
	GetUid() string
	GetName() string
}

var JwtService = new(jwtService)

type CustomClaims struct {
	UserID   int64  `json:"id"`
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expire_in"`
	TokenType   string `json:"token_type"`
}

func (jwtService *jwtService) GetToken(GuardName string, user JwtUser) (tokenData TokenOutPut, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, CustomClaims{
			UserID:   user.GetId(),
			UserName: user.GetName(),
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.App.Config.Jwt.JwtTtl) * time.Second)),
				ID:        user.GetUid(),
				Issuer:    GuardName,
				NotBefore: jwt.NewNumericDate(time.Now().Add(-1000 * time.Second)),
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(global.App.Config.Jwt.Secret))

	tokenData = TokenOutPut{
		tokenStr,
		int(global.App.Config.Jwt.JwtTtl),
		TokenType,
	}
	return
}

func (jwtService *jwtService) GetUserInfo(GuardName string, id int64) (err error, user JwtUser) {
	switch GuardName {
	case AppGuardName:
		return UserService.GetUserInfo(id)
	default:
		err = &locales.CErrors{Code: locales.CodeGuardNameNotExist}
	}
	return
}

func (jwtService *jwtService) getBlackListKey(tokenStr string) string {
	return global.App.Config.Redis.Prefix + "_jwt_blacklist:" + utils.MD5([]byte(tokenStr))
}

func (jwtService *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*CustomClaims).ExpiresAt.Unix()-nowUnix) * time.Second

	err = global.App.Redis.SetNX(context.Background(), jwtService.getBlackListKey(token.Raw), nowUnix, timer).Err()
	return
}

func (jwtService *jwtService) IsInBlackList(tokenStr string) bool {
	blockedStr, err := global.App.Redis.Get(context.Background(), jwtService.getBlackListKey(tokenStr)).Result()
	blockedUnix, err := strconv.ParseInt(blockedStr, 10, 64)
	if blockedStr == "" || err != nil {
		return false
	}

	if time.Now().Unix()-blockedUnix < global.App.Config.Jwt.JwtBlacklistGracePeriod {
		return false
	}

	return true
}
