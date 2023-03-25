package middlewares

import (
	"gin-scaffold/app/common/keys"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept-Language", "User-Agent"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{keys.HeaderNewExpiresIn, keys.HeaderNewToken, "Content-Disposition"}

	return cors.New(config)
}
