package bootstrap

import (
	"fmt"
	"gin-scaffold/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func InitializeConfig() *viper.Viper {
	config := "config.yaml"

	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		config = configEnv
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed:%s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)

		if err := v.Unmarshal(&global.App.Config); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&global.App.Config); err != nil {
		fmt.Println(err)
	}

	return v
}
