package routes

import (
	"gin-scaffold/app/controllers/api/v1"
	"gin-scaffold/app/middlewares"
	"gin-scaffold/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/auth/register", v1.Register)
	router.POST("/auth/login", v1.Login)

	authRouter := router.Group("").Use(middlewares.JWTAuthMiddleware(services.AppGuardName))
	{
		authRouter.GET("/auth/info", v1.Info)
		authRouter.POST("/auth/logout", v1.Logout)
	}
}
