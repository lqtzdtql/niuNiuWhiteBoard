package database

import (
	"log"

	conf "niuNiuWhiteBoardBackend/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var MEngine *xorm.Engine

var Db = map[string]conf.DbConfig{
	"db1": {
		DriverName:   conf.Cfg.DbConfig.DriverName,
		Dsn:          conf.Cfg.DbConfig.Dsn,
		ShowSql:      conf.Cfg.DbConfig.ShowSql,
		ShowExecTime: conf.Cfg.DbConfig.ShowExecTime,
		MaxIdle:      conf.Cfg.DbConfig.MaxIdle,
		MaxOpen:      conf.Cfg.DbConfig.MaxOpen,
	},
}

func InitDatabase() (*xorm.Engine, error) {
	if MEngine == nil {
		MEngine, err := xorm.NewEngine(Db["db1"].DriverName, Db["db1"].Dsn)
		if err != nil {
			log.Fatalln("Database NewEngine fatal")
		}
		MEngine.SetMaxIdleConns(Db["db1"].MaxIdle) //空闲连接
		MEngine.SetMaxOpenConns(Db["db1"].MaxOpen) //最大连接数
		MEngine.ShowSQL(Db["db1"].ShowSql)
		MEngine.ShowExecTime(Db["db1"].ShowExecTime)
		return MEngine, nil
	}
	return MEngine, nil
}
