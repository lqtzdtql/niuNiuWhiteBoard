package conf

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var MEngine *xorm.Engine

func init() {
	if MEngine == nil {
		var err error
		MEngine, err = xorm.NewEngine(Db["db1"].DriverName, Db["db1"].Dsn)
		if err != nil {
			log.Fatal(err)
		}
		MEngine.SetMaxIdleConns(Db["db1"].MaxIdle) //空闲连接
		MEngine.SetMaxOpenConns(Db["db1"].MaxOpen) //最大连接数
		MEngine.ShowSQL(Db["db1"].ShowSql)
		MEngine.ShowExecTime(Db["db1"].ShowExecTime)
	}

}
