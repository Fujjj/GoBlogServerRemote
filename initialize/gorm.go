package initialize

import (
	"os"
	"server/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() *gorm.DB {
	mysqlCfg := global.Configs.Mysql
	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlCfg.LogLevel()), //根据配置文件设置日志级别
	})
	if err != nil {
		global.Log.Error("Failed to connect mySQL", zap.Error(err))
		os.Exit(1)
	}
	global.Log.Info("Connect mySQL success")
	//获取底层的 SQL 数据库连接对象
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	return db
}
