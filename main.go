package main

import (
	"gin-scaffold/bootstrap"
	"gin-scaffold/global"
)

func main() {
	bootstrap.InitializeConfig()

	global.App.Log = bootstrap.InitializeLog()
	global.App.Log.Info("log init success!")
	global.App.DB = bootstrap.InitializeDB()
	global.App.Redis = bootstrap.InitializeRedis()
	if global.App.Redis == nil {
		global.App.Log.Error("redis init failed!")
		return
	}
	bootstrap.InitializeTranslator()
	bootstrap.InitializeSnowflake()

	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()

	bootstrap.RunServer()
}
