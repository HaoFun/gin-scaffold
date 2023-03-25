package middlewares

import (
	"gin-scaffold/app/common/keys"
	"gin-scaffold/app/common/response"
	"gin-scaffold/app/services"
	"gin-scaffold/global"
	"gin-scaffold/locales"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"
)

func JWTAuthMiddleware(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.FailC(c, locales.CodeNeedLogin)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.FailC(c, locales.CodeInvalidToken)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secret), nil
		})

		if err != nil || services.JwtService.IsInBlackList(parts[1]) {
			response.FailC(c, locales.CodeInvalidToken)
			c.Abort()
			return
		}

		claims := token.Claims.(*services.CustomClaims)
		if claims.Issuer != GuardName {
			response.FailC(c, locales.CodeInvalidToken)
			c.Abort()
			return
		}

		if claims.ExpiresAt.Unix()-time.Now().Unix() < global.App.Config.Jwt.RefreshGracePeriod {
			lock := global.Lock(keys.LockRefreshToken, global.App.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				err, user := services.JwtService.GetUserInfo(GuardName, claims.UserID)
				if err != nil {
					global.App.Log.Error(err.Error())
					lock.Release()
				} else {
					tokenData, _ := services.JwtService.GetToken(GuardName, user)
					c.Header(keys.HeaderNewToken, tokenData.AccessToken)
					c.Header(keys.HeaderNewExpiresIn, strconv.Itoa(tokenData.ExpireIn))
					_ = services.JwtService.JoinBlackList(token)
				}
			}
		}

		c.Set(keys.CtxUserIDKey, claims.UserID)
		c.Set(keys.CtxUserNameKey, claims.UserName)
		c.Set(keys.CtxUserToken, token)
		c.Next()
	}
}
