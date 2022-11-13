# 小画家



## Web端演示视频

见目录下视频文件

## web端项目运行

yarn start

## 后端使用文档

后端需要安装go环境，Mysql环境和Redis。

```
cd niuNiuWhiteBoard/backend/niuNiuSDKBackend
go build -o SDKServer cmd/main.go

cd niuNiuWhiteBoard/backend/niuNiuWhiteBoardBackend
go build -o SDKServer cmd/main.go
```

在niuNiuSDKBackend和niuNiuWhiteBoardBackend之下，各有一个models.sql。在mysql中建立两个数据库，database niuNiuSDK和database niuNiuWhiteBoard，分别导入对应目录下的models.sql。

两个文件下各有一个yaml配置文件，将里面的数据库配置更改（尤其是dsn）。七牛服务的AK和SK不进行上传。牛牛白板SDK的AK可以继续使用。

## Android端运行

使用Android Studio构建运行
