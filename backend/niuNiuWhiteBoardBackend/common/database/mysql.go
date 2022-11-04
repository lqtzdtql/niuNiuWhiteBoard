package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var MEngine *xorm.Engine

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

type DbConfig struct {
	DriverName   string
	Dsn          string
	ShowSql      bool
	ShowExecTime bool
	MaxIdle      int
	MaxOpen      int
}

var Db = map[string]DbConfig{
	"db1": {
		DriverName:   "mysql",
		Dsn:          "root:password@tcp(127.0.0.1:3306)/niuNiuWhiteBoard?charset=utf8mb4&parseTime=true&loc=Local",
		ShowSql:      true,
		ShowExecTime: false,
		MaxIdle:      10,
		MaxOpen:      200,
	},
}
