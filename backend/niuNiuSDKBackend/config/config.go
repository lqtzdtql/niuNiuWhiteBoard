package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbConfig    DbConfig    `json:"dbConfig"`
	Log         LogConfig   `json:"log"`
	RedisConfig RedisConfig `json:"redisConfig"`
}

type DbConfig struct {
	DriverName   string `json:"driverName"`
	Dsn          string `json:"dsn"`
	ShowSql      bool   `json:"showSql"`
	ShowExecTime bool   `json:"showExecTime"`
	MaxIdle      int    `json:"maxIdle"`
	MaxOpen      int    `json:"maxOpen"`
}

type RedisConfig struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	Auth    string `json:"auth"`
}

// 日志保存地址
type LogConfig struct {
	Path  string
	Level string
}

var cfg Config

func init() {
	// 设置文件名
	viper.SetConfigName("config")
	// 设置文件类型
	viper.SetConfigType("yaml")
	// 设置文件路径，可以多个viper会根据设置顺序依次查找
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	viper.Unmarshal(&cfg)

}
func GetConfig() Config {
	return cfg
}
