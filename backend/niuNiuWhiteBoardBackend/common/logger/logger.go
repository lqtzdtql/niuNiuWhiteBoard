package logger

import (
	"log"
	"os"
)

func CreateDir(dir string) (bool, error) {
	_, err := os.Stat(dir)

	if err == nil {
		return true, nil
	}

	err2 := os.MkdirAll(dir, 0755)
	if err2 != nil {
		return false, err2
	}

	return true, nil
}

func init() {
	res, err := CreateDir("../../LOG") //创建文件夹
	if res == false {
		panic(err)
	}
	file, _ := os.OpenFile("./LOG/run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //打开日志文件，不存在则创建

	log.SetOutput(file)                                 //设置输出流
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime) //日志输出样式
}
