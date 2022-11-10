package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Host         string        `json:"host"` //域名+端口
	Routes       []string      `json:"routes"`
	OpenJwt      bool          `json:"openJwt"`
	LogConfig    *LogConfig    `json:"logConfig"`
	DbConfig     *DbConfig     `json:"dbConfig"`
	QiniuService *QiniuService `json:"qiniuService"`
	Whiteboard   *Whiteboard   `json:"whiteboard"`
}

type DbConfig struct {
	DriverName   string `json:"driverName"`
	Dsn          string `json:"dsn"`
	ShowSql      bool   `json:"showSql"`
	ShowExecTime bool   `json:"showExecTime"`
	MaxIdle      int    `json:"maxIdle"`
	MaxOpen      int    `json:"maxOpen"`
}

// 日志保存地址
type LogConfig struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

type QiniuService struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	RTCAppID  string `json:"rtcAppId"`
}

type Whiteboard struct {
	AK string `json:"ak"`
}

var Cfg Config

func init() {
	// 设置文件名
	viper.SetConfigName("config")
	// 设置文件类型
	viper.SetConfigType("yaml")
	// 设置文件路径，可以多个viper会根据设置顺序依次查找
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
