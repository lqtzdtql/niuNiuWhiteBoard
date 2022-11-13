# 安卓端设计文档

**负责人：邓未央**

### 登录界面（HomeActivity）

进入页面先读取 SP，看看是否已经有保存好的账号密码，如果有就填上去。

为页面下方的用户协议视图添加 SpannableString，当用户点击的时候跳转到用户协议界面。

为账号和密码视图添加改变监听，当每次字符串改变结束后检查是否符合规范，显示提示信息。

重写时间分发方法，判断当前用户点按位置进行隐藏软键盘操作。

当用户点击登录时先检查是否勾选同意用户协议和手机号格式，没有的话就给出对应的提示，否则的话执行登录操作，登录失败则再进行注册和登录操作，登录成功之后保存当前账号密码到 SP，带着用户信息跳转到大厅页面并 finish 当前页面。

### 用户协议界面（UserAgreementActivity）

逻辑非常简单，使用 WebView 加载在 assets 中已写好的一个 html 格式的用户协议书。 

### 大厅页面（HomeActivity）

初始化视图后立刻起一个线程请求所有房间的情况，拿到数据就显示到视图上。房间列表视图使用的 Adapter 是 RoomListAdapter，它会持有房间信息数据，为其设置一个选择房间后的监听，这样当点击一个条目后可以由它拿到房间信息然后告知 HomeActivity，由 HomeActivity 完成跳转操作，并携带上房间信息和当前用户信息。

当用户点击右下角的创建房间按钮就直接创建房间页面。

### 创建房间页面（CreateRoomActivity）

初始化视图，给 RadioButton 加上监听，当选择不需要密码时就隐藏填入密码的视图，默认是需要密码的视图。

当点击创建房间的按钮时起一个线程去服务端请求创建房间，如果创建成功的话就直接跳转到房间页面并 finish 当前页面，否则显示对应提示。

### 绘画房间页面（RoomActivity）

进入页面先初始化画板和设置相关的视图，并设置对应的监听，如选择画笔大小事件监听、选择画笔颜色事件监听、选择形状事件监听等等，拿到 HomeActivity 传过来的 roomUUID 和当前用户信息，起一个线程去服务端请求加入房间，建立连接。

这个页面负责真实绘画的是 ArtBoard，这个是一个完全的自定义视图，直接继承于 View，可以调用 setCanvasWidth 和 setCanvasHeight 来设置画布宽高，可以使用 loadFile(String filename) 指定文件名来从文件加载一幅之前保存的绘画。设置 ShapeChangeListener 可以在有新的图像新增、改变、移除的时候接受到通知，对应的回调方法为 void onAddShape(Shape shape) 、void onMoveShape(Shape shape)、void onDeleteShape(Shape shape)。可以调用 getSaveShapeList() 来获取当前绘画的所有图形。ArtBoard 内部持有一个 Canvas，在其上绑定了一个 Bitmap，所有的绘画操作也都会保存在 Bitmap 上，可以调用 getmBitmap() 来获取。ArtBoard 重写了触摸事件，来处理各种绘画操作。

ArtBoard 中保存的图像皆为 Shape 类型，这是一个抽象类，所有的图形都继承于这个类并实现自己的特性，Shape 具有 getKind() 方法，可以获取到它本身是什么类型的，内部都保存了这个图型的所有信息，包括路径、画笔、宽高、颜色等。Shape 具有抽象方法 isInterSect，所有非抽象子类必须实现，该方法用户判断路径是否与自己相交。

页面中负责维持长链接和处理信息解析、分发、包装的是 WhiteBoardConn ，它继承于 WebSocketClient ，当连接上后会直接起一个定时线程，每隔 3s 发送一个心跳包，当 onMessage 接受到消息后根据其中的消息类型进行判断和解析，分发给应该处理它的对象。它对外界暴露多个方法如 addShape、changeShape，在内部进行信息抽取和包装并发送。

页面底下各个按钮点开后的面板分别由 ChooseUIManager、FootUIManager、SaveMenuManager 控制。

应用中所有的网络请求都由 AppServer 负责，使用单例控制，内部使用 OkHttp 实现，对外提供登录、注册、加入房间等等方法，返回数据实体类，数据实体类由 GsonFormat 生成，对外提供 parse 方法由 gson 构造。

