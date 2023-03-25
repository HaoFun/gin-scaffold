package bootstrap

import (
	"gin-scaffold/global"
	"gin-scaffold/utils"
)

func InitializeSnowflake() {
	utils.InitSnowflake(
		global.App.Config.App.StartTime,
		global.App.Config.App.MachineId,
	)
}
