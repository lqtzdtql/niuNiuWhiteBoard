package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	conf "niuNiuSDKBackend/config"
)

var MEngine *xorm.Engine

var Db = map[string]conf.DbConfig{
	"db1": {
		DriverName:   conf.GetConfig().DbConfig.DriverName,
		Dsn:          conf.GetConfig().DbConfig.Dsn,
		ShowSql:      conf.GetConfig().DbConfig.ShowSql,
		ShowExecTime: conf.GetConfig().DbConfig.ShowExecTime,
		MaxIdle:      conf.GetConfig().DbConfig.MaxIdle,
		MaxOpen:      conf.GetConfig().DbConfig.MaxOpen,
	},
}

func init() {
	if MEngine == nil {
		var err error
		MEngine, err = xorm.NewEngine(Db["db1"].DriverName, Db["db1"].Dsn)
		if err != nil {
			log.Fatalln("Database NewEngine fatal")
		}
		MEngine.SetMaxIdleConns(Db["db1"].MaxIdle) //空闲连接
		MEngine.SetMaxOpenConns(Db["db1"].MaxOpen) //最大连接数
		MEngine.ShowSQL(Db["db1"].ShowSql)
		MEngine.ShowExecTime(Db["db1"].ShowExecTime)
		return
	}
	return
}
