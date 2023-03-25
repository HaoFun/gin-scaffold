package bootstrap

import (
	"gin-scaffold/app/models"
	"gin-scaffold/global"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func InitializeDB() *gorm.DB {
	switch global.App.Config.Database.Driver {
	case "mysql":
		return initMySqlGorm()
	default:
		return initMySqlGorm()
	}
}

func initMySqlGorm() *gorm.DB {
	dbConfig := global.App.Config.Database

	if dbConfig.Database == "" {
		return nil
	}

	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host +
		":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.Database + "?charset=" +
		dbConfig.Charset + "&parseTime=True&loc=Local"

	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   getGormLogger(),
	}); err != nil {
		global.App.Log.Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		initMysqlTables(db)
		return db
	}
}

func initMysqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		models.User{},
	)

	if err != nil {
		global.App.Log.Error("mysql db migrate table failed, err:", zap.Any("err", err))
		os.Exit(0)
	}
}

func getGormLogWriter() logger.Writer {
	var writer io.Writer

	if global.App.Config.Database.EnableFileLogWriter {
		writer = &lumberjack.Logger{
			Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Database.LogFilename,
			MaxSize:    global.App.Config.Log.MaxSize,
			MaxBackups: global.App.Config.Log.MaxBackups,
			MaxAge:     global.App.Config.Log.MaxAge,
			Compress:   global.App.Config.Log.Compress,
		}
	} else {
		writer = os.Stdout
	}

	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch global.App.Config.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logMode,
		IgnoreRecordNotFoundError: false, // ignore ErrRecordNotFound Error
		Colorful:                  !global.App.Config.Database.EnableFileLogWriter,
	})
}
