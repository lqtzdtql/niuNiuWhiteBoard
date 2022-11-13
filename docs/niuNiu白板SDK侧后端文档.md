# niuNiuSDK后端文档

**负责人：巩宸旭**

**niuNiuSDK后端主要负责消息管理，鉴权管理，连接管理。**

**SDK server使用gin框架+mysql+redis。**

## **消息管理**

为了使得客户端之间能够协同，必须有一套信令系统。

我们对消息格式做如下规定，无论是客户端或是业务端，对收到的json消息都用以下结构进行解析。其中最重要的字段是ContentType，为消息类型。

```Go
type Message struct {
   From         string `json:"from,omitempty"`
   ToRoom       string `json:"to,omitempty"`
   ToWhiteBoard string `json:"toWhiteBoard,omitempty"`
   ToUser       string `json:"toUser,omitempty"`
   ObjectId     string `json:"objectId,omitempty"`
   ContentType  int32  `json:"contentType"`
   Content      string `json:"content,omitempty"`
   Timestamp    int64  `json:"timestamp,omitempty"`
   IsLock       bool   `json:"isLock,omitempty"`
   ReadOnly     bool   `json:"readOnly,omitempty"`
   UserName    string  `json:"userName,omitempty"`
}
```

信令类型目前有以下几种：

```Go
HEAT_BEAT     = 1 
UPDATE_BOARD  = 2 
OBJECT_NEW    = 3
OBJECT_MODIFY = 4
OBJECT_DELETE = 5
SWITCH_BOARD  = 6
DRAWING_LOCK  = 7
CREATE_BOARD  = 8
CAN_LOCK      = 9
LEAVE_ROOM    = 10
CANVAS_LIST   = 11
CUSTOMIZE_MESSAGE = 12
ENTER_ROOM    = 13
```

**注意：发送或者接收消息的字段第一个字母都是小写！**

### *HEAT_BEAT*     = 1 （心跳）

- 客户端发送的结构为，字段中包含以下两个部分。本类型消息由客户端定时发送。

```Go
 type HeatBeatReq struct{
     from //来自于哪个用户，为用户的UUID
     contentType //信令类型，为HEAT_BEAT
 }
```

- 对于服务端，心跳不需要保存，发送给指定参与者。服务端回复的字段中包括contentType

```JavaScript
 type HeatBeatRes struct{
     contentType //信令类型为HEAT_BEAT
 }
```

### *UPDATE_BOARD*  = 2 （更新白板）

- 客户端发送的结构为，字段中包含以下两个部分。本类型消息由客户端每隔一段时间主动发送。

```JavaScript
from
toRoom
toWhiteBoard //注意：此时的toWhiteBoard为要更新哪个白板，而不是发送到哪个白板
contentType
```

- 服务端收到消息后，查找数据库，找出该白板内的所有图形，并且回复给发送者。其中，content中为{objects:string[]}的JSON。

```JavaScript
contentType  //类型为UPDATE_BOARD
toWhiteBoard
content //{objects:string[]}的JSON
```

### *OBJECT_NEW*    = 3（创建对象）

- 客户端发送的请求带以下几个字段

```JavaScript
from //发送方uuid
toRoom //发送给哪个房间
toWhiteBoard //发送给哪个白板
objectId //对象id
content  //内容为该对象的json字符串，服务端不关心其内容
contentType
timestamp //该对象产生时的时间戳
```

- 服务端回复带以下几个字段，广播给除了发送者以外的所有人。

```JavaScript
contentType  //类型为OBJECT_NEW
toWhiteBoard // 该对象的白板号，用户收到后跟自己当前的白板号对比，决定是否渲染
objectId //对象id
content  //该对象的json字符串，服务端不关心其内容
```

### *OBJECT_MODIFY* = 4（修改对象）

客户端发来的格式类型和新建对象一样

```JavaScript
from   
toRoom
toWhiteBoard
objectId
content
contentType //OBJECT_MODIFY
timestamp   //修改时的时间戳
```

服务端收到之后，回复中带有以下几种字段，需要对时间戳做对比，决定是否修改，如果可以修改，则将修改后的对象广播给除发送者以外的所有人，否则丢弃请求。

```JavaScript
contentType  //回复的类型为OBJECT_MODIFY
toWhiteBoard //该对象的白板号，用户收到后跟自己当前的白板号对比，决定是否渲染
objectId
content
```

### *OBJECT_DELETE* = 5（删除对象）

客户端发来的格式类型和新建对象一样，带有以下几种字段

```JavaScript
from
toRoom
toWhiteBoard
objectId
contentType //OBJECT_DELETE
timestamp   //删除时的时间戳
```

服务端收到后，先去找出该对象的图元锁，该对象如果有图元锁，先删掉其图元锁，再删掉对象。回复中带有以下几种字段。

```JavaScript
contentType 
toWhiteBoard
objectId
content
```

### *SWITCH_BOARD*  = 6（切换白板）

客户端切换白板有两种情况

1. 在只读模式下，只能由主持人发送切换白板的消息，所有的人会跟着主持人一起切换到指定白板。

1. 在协作模式下，每个人都可以切换，每个人切换的是自己的白板。

```JavaScript
from
toRoom
toWhiteBoard //注意：此时的toWhiteBoard为要切换到哪个白板的白板号，而不是发送到哪个白板
contentType //SWITCH_BOARD
readOnly //bool类型，true表示此时是只读模式，false代表是协作模式
```

服务端收到后，查找数据库，找出该白板内的所有图形

1. 如果是只读模式

- 发送者以外的其他人，会收到*SWITCH_BOARD，切换到新白板并且收到新白板的所有对象内容。*

- 如果是发送者，收到的类型为UPDATA_BOARD，更新当前白板内容。

1. 如果是协作模式

- 发送者收到的类型为UPDATA_BOARD，更新当前白板内容。其他人不会收到*。*

```JavaScript
contentType 
toWhiteBoard //此时的toWhiteBoard为要切换到哪个白板的白板号，而不是发送到哪个白板
content //content中的内容为{objects:string[]}的JSON
```

### *DRAWING_LOCK*  = 7（图元锁）

图元锁是当用户触发某个图形时，先会拥有该图形的所有权。在该用户对这个图形操作的过程中，即解锁之前，其他人无法对该图形进行操作。这样可以避免两个人抢占一个图形，发生意想不到的情况。

客户端发来的消息中，相比于普通对象，会多一个isLock。

- 如果用户要对某个图形加锁，isLock设置为true。

- 如果要解锁，isLock设置为false。

- 其他字段含义与普通对象一致。

```JavaScript
from
toRoom
toWhiteBoard
objectId //图元锁锁的是哪个对象，该字段为指定对象的id
isLock   //区分是加锁请求还是解锁请求
contentType
timestamp 
```

服务端对收到的图元锁类型的消息做如下处理：

如果是解锁请求，如果存在这个锁，就删除这个锁，并且向发送者以外的所有人广播解锁消息（**contentType为*****DRAWING_LOCK，isLock为false***），如果不存在，就忽略。

如果是加锁请求，服务端会先判断此时是否已经有该锁。

- 如果已经上锁，给发送者回复消息中， **contentType为 CAN_LOCK ，islock 为true**。表示该图形已经上锁，上锁失败。不会给其他人发送消息。

- 如果此时该对象还没有上锁，给发送者回复的消息中， **contentType为 CAN_LOCK ，islock 为false，**表示上锁成功。同时，给发送者以外的其他人广播的消息中，**contentType为*****DRAWING_LOCK*** **，islock 为true**。表示该图形已经被上锁

**为什么要引入CAN_LOCK信令？**因为如果只有*DRAWING_LOCK，当*有两个人同时申请加锁时，就会造成死锁，该图形无法被任何人操作的情况。这种现象是不可容忍的，因此引入CAN_LOCK信令，由服务端向发送者来发送，搭配isLock字段，可以让发送者知道自己是否加锁成功，非常重要。

无论是CANL_LOCK还是DRAWING_LOCK，携带的都是以下四个字段：

```JavaScript
contentType
toWhiteBoard
objectId
isLock
```

### *CREATE_BOARD*  = 8（创建白板）

客户端，发送的contentType为*CREATE_BOARD。*

```JavaScript
from
toRoom
contentType
```

服务端回复的contentType为CREATE_BOARD，content的json格式为{"canvasId" : "xxxxx"}，客户端收到后要对这个content进行解析，获取白板号。

```JavaScript
content    
contentType
```

### *LEAVE_ROOM*    = 10（离开房间）

离开房间的contentType为*LEAVE_ROOM*

- 房主踢人，会将 leaveUser设置为被踢者的UserUuid。

- 用户主动退出， leaveUser和from的内容是一致的。

```JavaScript
 from
 toRoom
 userName  //sdk端的user_uuid
 contentType
```

对于服务端，当有人发送*LEAVE_ROOM* 消息或者有人掉线时*，*会广播给其他人LEAVE_ROOM给其他人。用户退出后，会从用户列表中删除该用户，并且判断该房间人数是否为0，如果为0，则销毁房间。

注意：业务端的房主收到sdk端的user_name，也就是业务端的user_uuid后，会发布踢人的请求。业务服务器也会从参与者列表中删除对应的参与者。

```Go
 contentType
 user  //注意：这里回复的是sdk端的user_name，收到后要传给业务端。（不是sdk端的user_uuid！！）
```

### CANVAS_LIST   = 11（获取该房间的白板列表）

当用户进入房间后，给用户展示的不是一个白板，而是白板列表。白板列表需要实时更新，因此需要客户端每隔一段时间进行请求。

```JavaScript
 from
 toRoom  
 contentType //CANVAS_LIST
```

服务端找出该房间内的所有白板，构造content，content为字符串数组json后的结果。数组里的内容是所有的白板id。

```JavaScript
contentType
content
```

### CUSTOMIZE_MESSAGE = 12(用户侧自定义广播)

客户端发送的消息中包含以下字段

```JavaScript
from
toRoom
contentType
content
```

服务端返回字段，将content广播给所有人。

```JavaScript
contentType
content
```

### ENTER_ROOM    = 13

客户端：进房

服务端广播有人进房的消息

```Go
 contentType
 user  //注意：这里回复的是sdk端的user_name，收到后要传给业务端。（不是sdk端的user_uuid！！）
```

## **鉴权管理**

鉴权管理用来验证能否能够合法的使用我们的牛牛白板SDK，流程如下：

1. 用户通过牛牛白板的网址获取SK，将SK放在业务服务器的配置中，SK需要保存好。（考虑到工作量，我们暂时放弃使用诸如七牛等公司正在使用AK/SK机制，仅仅以SK表达鉴权的思想）

1. 客户端 APP 将 room_uuid 和 user_uuid 作为参数请求客户的业务服务器，客户业务服务器通过 AK 签算出 RoomToken。

1. 客户端 SDK 通过客户业务服务器签发的 RoomToken 请求 SDK server,  SDK server验证请求合法性后创建房间或者加入房间，返回sdk端的user_uuid和room_uuid。

### **GET auth**

GET /auth

### **请求参数**

| 名称  | 位置  | 类型   | 必选 |
| ----- | ----- | ------ | ---- |
| token | query | string | 是   |

> 返回示例
>
> 成功

```Go
{
  "message": "enter room success",
  "user_uuid": "participantuuid",
  "user_name": "username",
  "room_uuid": "roomuuid",
  "code": 200
}
```

### **返回结果**

| 状态码 | 状态码含义                                              | 说明 |
| ------ | ------------------------------------------------------- | ---- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | 成功 |

### **返回数据结构**

状态码 **200**

| 名称        | 类型    | 必选 |
| ----------- | ------- | ---- |
| » message   | string  | TRUE |
| » user_uuid | string  | TRUE |
| » user_name | string  | TRUE |
| » room_uuid | string  | TRUE |
| » code      | integer | TRUE |

## **连接管理**

服务端为了向客户端主动发送消息，需要在连接管理上具有以下功能：

1. 建立连接(保持连接)

1. 断开连接(删除连接)

1. 维护连接(心跳检测)

1. 接收消息

1. 发送消息

白板和音视频不同，白板需要保障可靠传输。在连接管理上，我们提出了两种协议保障可靠传输。一种是WebSocket，另一种是QUIC。考虑到使用的广泛度，我们选择使用WebSocket。实际上通过数据显示，相同带宽下，QUIC的建连时间和传输时间要更短。未来websocket也有将底层迁移到QUIC的趋势。未来我们也可能将协议迁移到QUIC上。

请求websocket时继续使用token鉴权，连接成功后会给房间内其他的人广播有用户进入。

 **路径 /websocket**

### **请求参数**

| 名称  | 位置  | 类型   | 必选 |
| ----- | ----- | ------ | ---- |
| token | query | string | 是   |