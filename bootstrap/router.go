package bootstrap

import (
	"context"
	"gin-scaffold/app/middlewares"
	"gin-scaffold/global"
	"gin-scaffold/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.GinLogger())
	router.Use(middlewares.GinRecovery(true))
	router.Use(middlewares.LanguageMiddleware())
	router.Use(middlewares.Cors())

	router.Static("/public", "./static")
	router.Static("/storage", "./storage/app/public")

	apiGroup := router.Group("/api")
	routes.SetApiGroupRoutes(apiGroup)

	return router
}

func RunServer() {
	r := setRouter()

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting.")
}
