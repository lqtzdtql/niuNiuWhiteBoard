# web侧SDK使用文档

# 新建房间/加入房间

**负责人：刘骐铜**

```JavaScript
import {joinRoom} from 'sdk/index'

const room = joinRoom(params: paramsType)
//type paramsType = {
//  token: string;
//  onlyRead: boolean;
//  el: HTMLCanvasElement;
  //elOptions?: {};
//};

room.on('leaveRoom',(options:{leaveUserId:string])=>{
    ...业务侧对有人离开房间的逻辑处理
})

room.on('joinRoom',(options:{joinUserId:string])=>{
    ...业务侧对有人进入房间的逻辑处理
})

room.on('new',（options:{canvasId:string，userId:string}）=>{
    ...业务侧对房间内出现新白板（可能是自己创建也可能是他人创建）的逻辑处理,
    通过userId可以知道这个新白板是谁创建的
})

room.on('customizeMessage',(options:{message:any})=>{
    ...业务侧拿到房间内其他用户发给自己的消息后的处理逻辑，此处的message等同于
    下面向房间内所有人发消息接口里的入参message
})
```

# 新建白板(新建房间不会默认自动新建白板，需要调用接口）

```JavaScript
room.createCanvas()；
```

# 切换白板

```JavaScript
room.switchCanvas(canvasId: string)
```

# 获取对应编号白板

```JavaScript
const canvas = room.getCanvas(canvasId: string)
```

# 切换房间的只读/协作模式

```JavaScript
room.modifyOnlyRead()
```

# 踢人/退出房间

```JavaScript
room.kickOutRoom(userId: string)
退出房间时userId传自己的即可
```

# 导入白板

```JavaScript
room.importCanvas(canvasId: string, canvasData: string)
```

# 导出白板

```JavaScript
const canvasData = room.exportCanvas(canvasId: string);
```

# 向房间内所有人发消息

```JavaScript
room.sendMessageToAll(message: any);
```

# 笔刷

```JavaScript
const modifyBrush = (options) => {
    canvas.current.modifyBrush(options);
}
return (
<div
    onClick={() => {
    modifyBrush(options);
}}>
   直线
</div>
)
```

![img](https://bytedance.feishu.cn/space/api/box/stream/download/asynccode/?code=MDlkOWNhYjFmYjdlNDMzMTQyMmMwYzhjZTY4NGMzNWJfVmJRUWw0M0FLMmhpQkJnc3BMSG1CWUNyNWhkT1NKYVJfVG9rZW46Ym94Y25XdXEycjFVNW9hbVQ3TFFPejZmQk9jXzE2NjgzOTU5MDU6MTY2ODM5OTUwNV9WNA)

新增类型11：涂鸦（类似自由线条，区别是涂鸦只有自己能看见，不是对象也不可部分擦除，可以一键清除）

options:

{

type:number,//笔刷类型

fill:颜色,   //填充颜色

stroke:颜色，//描边颜色

strokeWidth:number,//描边线宽

headlen:number, //箭头斜线长度

text:string,//要渲染的文字

size:number,//字号

font:size,//除字号外其余字体样式

fillText:boolean,//文字填充渲染

strokeText:boolean,//文字描边渲染

rx:number,//矩形圆角

ry:number,//矩形圆角

}

# 清除涂鸦

```JavaScript
canvas.clearGraffiti();
```

# 撤销与重做

```JavaScript
const revoke = () => {
    canvas.current.revoke();
  };
  const redo = () => {
    canvas.current.redo();
  };
  return (
      <div onClick={revoke}>撤销</div>
      <div onClick={redo}>重做</div>
  )
```

# 属性

白板（canvas类）

```JavaScript
 public canvasId: string;
  /** 画布宽度 */
  public width: number;
  /** 画布高度 */
  public height: number;
  /** 画布背景颜色 */
  public backgroundColor;
  /** 包围 canvas 的外层 div 容器 */
  public wrapperEl: HTMLElement;
  /** 下层 canvas 画布，主要用于绘制所有物体 */
  public lowerCanvasEl: HTMLCanvasElement;
  /** 上层 canvas，主要用于监听鼠标事件、涂鸦模式、左键点击拖蓝框选区域 */
  public upperCanvasEl: HTMLCanvasElement;
  /** 上层画布环境 */
  public contextTop: CanvasRenderingContext2D;
  /** 下层画布环境 */
  public contextContainer: CanvasRenderingContext2D;
  /** 缓冲层画布环境，方便某些情况方便计算用的，比如检测物体是否透明 */
  public cacheCanvasEl: HTMLCanvasElement;
  public contextCache: CanvasRenderingContext2D;
  public containerClass: string = 'canvas-container';

  /** 记录最近一个激活的物体，可以优化点选过程，也就是点选的时候先判断是否是当前激活物体 */
  // public lastRenderedObjectWithControlsAboveOverlay;
  /** 通过像素来检测物体而不是通过包围盒 */
  // public perPixelTargetFind: boolean = false;

  /** 一些鼠标样式 */
  public defaultCursor: string = 'default';
  public hoverCursor: string = 'move';
  public moveCursor: string = 'move';
  public rotationCursor: string = 'crosshair';
  /**笔刷： 0默认1直线2曲线3矩形4菱形5三角形6圆形7箭头8橡皮9文字10自由线条 */
  public brush: {} = { type: 0 };
  public start: Pos = {};
  public end: Pos = {};
  public temp: Pos = {};
  public penPath: Pos[] = [];

  public viewportTransform: number[] = [1, 0, 0, 1, 0, 0];
  public vptCoords: {};

  // public relatedTarget;
  /** 选择区域框的背景颜色 */
  public selectionColor: string = 'rgba(100, 100, 255, 0.3)';
  /** 选择区域框的边框颜色 */
  public selectionBorderColor: string = 'red';
  /** 选择区域的边框大小，拖蓝的线宽 */
  public selectionLineWidth: number = 1;
  /** 左键拖拽的产生的选择区域，拖蓝区域 */
  private _groupSelector: GroupSelector;
  /** 当前选中的组 */
  public _activeGroup: Group;
  public canvasId: string = '';
  public modifiedList: [] = [];
  public modifiedAgainList: [] = [];

  /** 画布中所有添加的物体 */
  private _objects: FabricObject[];
  /** 整个画布到上面和左边的偏移量 */
  private _offset: Offset;
  /** 当前物体的变换信息，src 目录下中有截图 */
  private _currentTransform: CurrentTransform;
  /** 当前激活物体 */
  private _activeObject;
  /** 变换之前的中心点方式 */
  // private _previousOriginX;
  private _previousPointer: Pos;
  
```

物体基类

```JavaScript
 /** 物体类型标识 */
  public type: string = 'object';
  /** 是否处于激活态，也就是是否被选中 */
  public active: boolean = false;
  /** 是否可见 */
  public visible: boolean = true;
  /** 默认水平变换中心 left | right | center */
  public originX: string = 'center';
  /** 默认垂直变换中心 top | bottom | center */
  public originY: string = 'center';
  /** 物体位置 top 值 */
  public top: number = 0;
  /** 物体位置 left 值 */
  public left: number = 0;
  /** 物体原始宽度 */
  public width: number = 0;
  /** 物体原始高度 */
  public height: number = 0;
  /** 物体当前的缩放倍数 x */
  public scaleX: number = 1;
  /** 物体当前的缩放倍数 y */
  public scaleY: number = 1;
  /** 物体当前的旋转角度 */
  public angle: number = 0;
  /** 左右镜像，比如反向拉伸控制点 */
  public flipX: boolean = false;
  /** 上下镜像，比如反向拉伸控制点 */
  public flipY: boolean = false;
  /** 选中态物体和边框之间的距离 */
  public padding: number = 0;
  /** 物体缩放后的宽度 */
  public currentWidth: number = 0;
  /** 物体缩放后的高度 */
  public currentHeight: number = 0;
  /** 激活态边框颜色 */
  public borderColor: string = 'red';
  /** 激活态控制点颜色 */
  public cornerColor: string = 'red';
  /** 物体默认填充颜色 */
  public fill: string = 'rgb(0,0,0)';
  /** 混合模式 globalCompositeOperation */
  // public fillRule: string = 'source-over';
  /** 物体默认描边颜色，默认无 */
  public stroke: string;
  /** 物体默认描边宽度 */
  public strokeWidth: number = 1;
  /** 矩阵变换 */
  // public transformMatrix: number[];
  /** 最小缩放值 */
  // public minScaleLimit: number = 0.01;
  /** 是否有控制点 */
  public hasControls: boolean = true;
  /** 是否有旋转控制点 */
  public hasRotatingPoint: boolean = true;
  /** 旋转控制点偏移量 */
  public rotatingPointOffset: number = 40;
  /** 移动的时候边框透明度 */
  public borderOpacityWhenMoving: number = 0.4;
  /** 物体是否在移动中 */
  public isMoving: boolean = false;
  /** 选中态的边框宽度 */
  public borderWidth: number = 1;
  /** 物体控制点用 stroke 还是 fill */
  public transparentCorners: boolean = false;
  /** 物体控制点大小，单位 px */
  public cornerSize: number = 12;
  /** 通过像素来检测物体而不是通过包围盒 */
  public perPixelTargetFind: boolean = false;
  /** 物体控制点位置，随时变化 */
  public oCoords: Coords;
  /** 物体所在的 canvas 画布 */
  public canvas;
  /** 物体执行变换之前的状态 */
  public originalState;
  /** 物体所属的组 */
  public group;
  /** 物体被拖蓝选区保存的时候需要临时保存下 hasControls 的值 */
  public orignHasControls: boolean = true;
  public stateProperties: string[] = (
    'top left width height scaleX scaleY ' +
    'flipX flipY angle cornerSize fill originX originY ' +
    'stroke strokeWidth ' +
    'borderWidth transformMatrix visible'
  ).split(' ');
  public timestamp: number = 0;
  public objectId: string = '';

  private _cacheCanvas: HTMLCanvasElement;
  private _cacheContext: CanvasRenderingContext2D;
  public cacheWidth: number;
  public cacheHeight: number;
  public dirty: boolean;
```

组

```JavaScript
public type: string = 'group';
  // 组中所有的物体
  public objects: FabricObject[];
  public originalState;
```

矩形

```JavaScript
public type: string = 'rect';
  /** 圆角 rx */
  public rx: number = 0;
  /** 圆角 ry */
  public ry: number = 0;
```

三角形

```JavaScript
public type: string = 'triangle';
```

圆形

```JavaScript
public type: string = 'round';
  public roundAngle: number = 360;
```

菱形

```JavaScript
 public type: string = 'diamond';
```

直线

```JavaScript
 public type: string = 'line';
```

箭头

```JavaScript
public type: string = 'array';
  public headlen: number = 15;
  public direction: number = 0;
```

文字

```JavaScript
public type: string = 'text';
  public text: string = '';
  public size: number = 10;
  public font: string = 'Microsoft YaHei';
  public fillText: boolean = true;
  public strokeText: boolean = false;
```

自由线条

```JavaScript
public type: string = 'text';
  public penPath: [] = [];
  public start;
```