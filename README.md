# 小画家


## 分工：

**刘骐铜--Web端SDK+Web前端**

**邓未央--安卓端+Web前端**

**巩宸旭--SDK服务端+小画家业务服务端**
## 体验链接
http://81.68.68.216:10000/

## 简介

随着在线教育，视频会议的发展，人们对在线协作创作的需求逐渐旺盛。在在线课堂中，老师可能会随时勾画重点，而学生可以在本地进行涂鸦，演算。在视频会议中，参会人员可以利用白板书写灵感，头脑风暴。而在在线创作领域，则需要“小画家”们共同协作，一起创作出属于大家的作品。

基于此，贝极星都是很有用的人团队共同研发了小画家，支持web端，安卓端。有绘图，文字，导入导出，撤销重做等功能。

但是，我们不希望自己的白板产品，仅仅是一个白板。因此，我们模仿如七牛等Paas层厂商，创作了自己的白板SDK——牛牛白板。牛牛白板接口简单，易于接入。

根据调研，我们并没有发现开源的白板SDK的产品，绝大多数的白板产品都是在Paas层之上封装了自己的服务。我们希望牛牛SDK，成为第一款开源白板SDK产品。

## 设计架构

小画家的白板能力是基于牛牛白板SDK的，但是为了处理业务需求的相关的事件，我们设计了单点登录（包括手机号注册，登录），用户信息管理（如账号密码管理），房间管理。

使用牛牛白板需要进行鉴权，流程如下：

1. 用户通过牛牛白板的网址获取SK，在业务服务器配置时使用。

2. 客户端 APP 将 room_uuid 和 user_uuid 作为参数请求客户的业务服务器，客户业务服务器通过 SK 签算出 RoomToken；

3. 客户端 SDK 通过客户业务服务器签发的 RoomToken 请求 SDK server,  SDK server验证请求合法性后创建房间或者加入房间，返回sdk端的user_uuid和room_uuid。

4. 建立websocket持久连接，SDK和SDK server端进行消息交互。

![image](https://user-images.githubusercontent.com/84149464/203592496-c633d1c2-add3-47bd-8b05-a1c13ad7318e.png)

后端交互如下
![image](https://user-images.githubusercontent.com/84149464/203592582-fc84f4fe-bfb4-4b83-885a-defa0fb1ec98.png)

前端设计如下
![image](https://user-images.githubusercontent.com/84149464/203592684-197aeee6-8223-4a60-9319-957627f61424.png)


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
go build -o littlepainter cmd/main.go
```

在niuNiuSDKBackend和niuNiuWhiteBoardBackend之下，各有一个models.sql。在mysql中建立两个数据库，database niuNiuSDK和database niuNiuWhiteBoard，分别导入对应目录下的models.sql。

两个文件下各有一个yaml配置文件，将里面的数据库配置更改（尤其是dsn）。七牛服务的AK和SK不进行上传。牛牛白板SDK的AK可以继续使用。

## Android端运行

使用Android Studio构建运行

##后续计划

1. UI重新设计。
2. 整体代码重构，提升代码质量。
3. 白板计划添加动画功能，修复组框选存在的问题。
4. 业务端将引入长连接，避免目前用户异常退出时房间不销毁的问题。
5. SDK使用接口优化，SDK文档优化。
6. 实现自己的音视频服务（从SDK到SFU的构建）
我们更期待明年11月，一个全新的，更加好用的牛牛白板(或者说白板语聊房)展现在用户面前~



