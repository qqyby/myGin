package inits

import (
	"fmt"
	"go.uber.org/zap"
	"myGin/dao/mysql"
	"myGin/dao/redis"
	"myGin/logger"
	"myGin/pkg/snowflake"
	"myGin/settings"
	"sync"

)

var (
	settingOnce sync.Once
	clusterOnce sync.Once
)

func InitCfgAndLog(cfgPath string) {
	settingOnce.Do(func() {
		// 记载配置
		if err := settings.Init(cfgPath); err != nil {
			panic(fmt.Sprintf("init settings error:%v ", err))
		}

		// 初始化日志
		if err := logger.Init(); err != nil {
			panic(fmt.Sprintf("init logger error:%v ", err))
		}

		// 初始化 snowflake
		if err := snowflake.Init(settings.AppCfg.SnowflakeStartTime, settings.AppCfg.SnowflakeMachineID); err != nil {
			panic(fmt.Sprintf("init snowflake error:%v ", err))
		}
	})
}

func InitDao() {
	clusterOnce.Do(func() {
		// 初始化mysql
		if err := mysql.Init(settings.MysqlCfg.User, settings.MysqlCfg.Password,
			settings.MysqlCfg.DbName, settings.MysqlCfg.Host, settings.MysqlCfg.Port,
			settings.MysqlCfg.MaxOpenConn, settings.MysqlCfg.MaxIdleConn); err != nil {
			panic(fmt.Sprintf("init mysql error:%v ", err))
		}

		// 初始化redis  host string, port int, pwd string, db, poolSize int
		if err := redis.Init(settings.RedisCfg.MasterName, settings.RedisCfg.Host, settings.RedisCfg.Password,
			settings.RedisCfg.Db, settings.RedisCfg.PoolSize); err != nil {
			panic(fmt.Sprintf("init redis error:%v ", err))
		}
	})
}

func Close() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	_ = mysql.Close()
	_ = zap.L().Sync()
	_ = redis.Close()
}
