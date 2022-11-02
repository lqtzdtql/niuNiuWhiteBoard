import { EventCenter } from './EventCenter';
import { Util } from './Util';
import { FabricObject } from './FabricObject';
import { Group } from './Group';
import { Offset, Pos, GroupSelector } from './interface';

export class Canvas extends EventCenter {
  /** 画布宽度 */
  public width: number;
  /** 画布高度 */
  public height: number;
  /** 画布背景颜色 */
  public backgroundColor: any;
  /** 包围 canvas 的外层 div 容器 */
  public wrapperEl: HTMLElement;
  /** 下层 canvas 画布，主要用于绘制所有物体 */
  public lowerCanvasEl: HTMLCanvasElement;
  /** 上层 canvas，主要用于监听鼠标事件、涂鸦模式、左键点击拖蓝框选区域 */
  public upperCanvasEl: HTMLCanvasElement;
  /** 缓冲层画布 */
  public cacheCanvasEl: HTMLCanvasElement;
  /** 上层画布环境 */
  public contextTop: CanvasRenderingContext2D;
  /** 下层画布环境 */
  public contextContainer: CanvasRenderingContext2D;
  /** 缓冲层画布环境，方便某些情况方便计算用的，比如检测物体是否透明 */
  public contextCache: CanvasRenderingContext2D;
  public containerClass: string = 'canvas-container';

  /** 左键拖拽的产生的选择区域，拖蓝区域 */
  private _groupSelector: GroupSelector | null;
  /** 当前选中的组 */
  public _activeGroup: Group | null;

  /** 整个画布到上面和左边的偏移量 */
  private _offset: Offset;
  /** 画布中所有添加的物体 */
  private _objects: FabricObject[];
  /** 当前物体的变换信息 */
  private _currentTransform: CurrentTransform;
  /** 当前激活物体 */
  private _activeObject;
  /** 变换之前的中心点方式 */
  // private _previousOriginX;
  private _previousPointer: Pos;

  constructor(el: HTMLCanvasElement, options) {
    super();
    // 初始化下层画布 lower-canvas
    this._initStatic(el, options);
    // 初始化上层画布 upper-canvas
    this._initInteractive();
    // 初始化缓冲层画布
    this._createCacheCanvas();
    // 处理模糊问题
    this._initRetinaScaling();
  }

  /** 初始化 _objects、lower-canvas 宽高、options 赋值 */
  _initStatic(el: HTMLCanvasElement, options) {
    this._objects = [];

    this._createLowerCanvas(el);
    this._initOptions(options);

    this.calcOffset();
  }

  _createLowerCanvas(el: HTMLCanvasElement) {
    this.lowerCanvasEl = el;
    Util.addClass(this.lowerCanvasEl, 'lower-canvas');
    this._applyCanvasStyle(this.lowerCanvasEl);
    this.contextContainer = this.lowerCanvasEl.getContext('2d') as CanvasRenderingContext2D;
  }

  _applyCanvasStyle(el: HTMLCanvasElement) {
    let width = this.width || el.width;
    let height = this.height || el.height;
    Util.setStyle(el, {
      position: 'absolute',
      width: width + 'px',
      height: height + 'px',
      left: 0,
      top: 0,
    });
    el.width = width;
    el.height = height;
    Util.makeElementUnselectable(el);
  }

  _initOptions(options) {
    for (let prop in options) {
      this[prop] = options[prop];
    }

    this.width = +this.lowerCanvasEl.width || 0;
    this.height = +this.lowerCanvasEl.height || 0;

    this.lowerCanvasEl.style.width = this.width + 'px';
    this.lowerCanvasEl.style.height = this.height + 'px';
  }

  /** 获取画布的偏移量，到时计算鼠标点击位置需要用到 */
  calcOffset(): Canvas {
    this._offset = Util.getElementOffset(this.lowerCanvasEl);
    return this;
  }

  /** 初始化交互层，也就是 upper-canvas */
  _initInteractive() {
    this._currentTransform = null;
    this._groupSelector = null;
    this._initWrapperElement();
    this._createUpperCanvas();
    this._initEvents();
    this.calcOffset();
  }

  /** 因为我们用了两个 canvas，所以在 canvas 的外面再多包一个 div 容器 */
  _initWrapperElement() {
    this.wrapperEl = Util.wrapElement(this.lowerCanvasEl, 'div', {
      class: this.containerClass,
    });
    Util.setStyle(this.wrapperEl, {
      width: this.width + 'px',
      height: this.height + 'px',
      position: 'relative',
    });
    Util.makeElementUnselectable(this.wrapperEl);
  }

  /** 创建上层画布，主要用于鼠标交互和涂鸦模式 */
  _createUpperCanvas() {
    this.upperCanvasEl = Util.createCanvasElement();
    this.upperCanvasEl.className = 'upper-canvas';
    this.wrapperEl.appendChild(this.upperCanvasEl);
    this._applyCanvasStyle(this.upperCanvasEl);
    this.contextTop = this.upperCanvasEl.getContext('2d') as CanvasRenderingContext2D;
  }

  /** 给上层画布增加鼠标事件 */
  _initEvents() {
    this._onMouseDown = this._onMouseDown.bind(this);
    this._onMouseMove = this._onMouseMove.bind(this);
    this._onMouseUp = this._onMouseUp.bind(this);
    this._onResize = this._onResize.bind(this);

    Util.addListener(window, 'resize', this._onResize);
    Util.addListener(this.upperCanvasEl, 'mousedown', this._onMouseDown);
    Util.addListener(this.upperCanvasEl, 'mousemove', this._onMouseMove);
  }

  _onMouseDown(e: MouseEvent) {
    this.__onMouseDown(e);
    Util.addListener(document, 'mouseup', this._onMouseUp);
    Util.addListener(document, 'mousemove', this._onMouseMove);
    Util.removeListener(this.upperCanvasEl, 'mousemove', this._onMouseMove);
  }

  __onMouseDown(e: MouseEvent) {
    // 只处理左键点击，要么是拖蓝事件、要么是点选事件
    let isLeftClick = 'which' in e ? e.which === 1 : e.button === 0;
    if (!isLeftClick) return;

    // 这个是为了保险起见，ignore if some object is being transformed at this moment
    if (this._currentTransform) return;

    let target = this.findTarget(e);
    let pointer = this.getPointer(e);
    let corner;
    this._previousPointer = pointer;

    if (this._shouldClearSelection(e)) {
      // 如果是拖蓝选区事件
      this._groupSelector = {
        // 重置选区状态
        ex: pointer.x,
        ey: pointer.y,
        top: 0,
        left: 0,
      };
      // 让所有元素失去激活状态
      this.deactivateAllWithDispatch();
      // this.renderAll();
    } else {
      // 如果是点选操作，接下来就要为各种变换做准备
      target?.saveState();

      // 判断点击的是不是控制点
      corner = target?._findTargetCorner(e, this._offset);
      // if ((corner = target._findTargetCorner(e, this._offset))) {
      //     this.onBeforeScaleRotate(target);
      // }
      if (this._shouldHandleGroupLogic(e, target)) {
        // 如果是选中组
        this._handleGroupLogic(e, target);
        target = this.getActiveGroup();
      } else {
        // 如果是选中单个物体
        if (target !== this.getActiveGroup()) {
          this.deactivateAll();
        }
        this.setActiveObject(target, e);
      }
      this._setupCurrentTransform(e, target);

      // if (target) this.renderAll();
    }
    // 不论是拖蓝选区事件还是点选事件，都需要重新绘制
    // 拖蓝选区：需要把之前激活的物体取消选中态
    // 点选事件：需要把当前激活的物体置顶
    this.renderAll();

    this.emit('mouse:down', { target, e });
    target && target.emit('mousedown', { e });
    // if (corner === 'mtr') {
    //     // 如果点击的是上方的控制点，也就是旋转操作，我们需要临时改一下变换中心，因为我们一直就是以 center 为中心，所以可以先不管
    //     this._previousOriginX = this._currentTransform.target.originX;
    //     this._currentTransform.target.adjustPosition('center');
    //     this._currentTransform.left = this._currentTransform.target.left;
    //     this._currentTransform.top = this._currentTransform.target.top;
    // }
  }

  /** 检测是否有物体在鼠标位置 */
  findTarget(e: MouseEvent, skipGroup: boolean = false): FabricObject | undefined {
    let target;
    // let pointer = this.getPointer(e);

    // 优先考虑当前组中的物体，因为激活的物体被选中的概率大
    let activeGroup = this.getActiveGroup();
    if (activeGroup && !skipGroup && this.containsPoint(e, activeGroup)) {
      target = activeGroup;
      return target;
    }

    // 遍历所有物体，判断鼠标点是否在物体包围盒内
    for (let i = this._objects.length; i--; ) {
      if (this._objects[i] && this.containsPoint(e, this._objects[i])) {
        target = this._objects[i];
        break;
      }
    }

    // 如果不根据包围盒来判断，而是根据透明度的话，可以用下面的代码
    // 先通过包围盒找出可能点选的物体，再通过透明度具体判断，具体思路可参考 _isTargetTransparent 方法
    // let possibleTargets = [];
    // for (let i = this._objects.length; i--; ) {
    //     if (this._objects[i] && this.containsPoint(e, this._objects[i])) {
    //         if (this.perPixelTargetFind || this._objects[i].perPixelTargetFind) {
    //             possibleTargets[possibleTargets.length] = this._objects[i];
    //         } else {
    //             target = this._objects[i];
    //             this.relatedTarget = target;
    //             break;
    //         }
    //         break;
    //     }
    // }
    // for (let j = 0, len = possibleTargets.length; j < len; j++) {
    //     pointer = this.getPointer(e);
    //     let isTransparent = this._isTargetTransparent(possibleTargets[j], pointer.x, pointer.y);
    //     if (!isTransparent) {
    //         target = possibleTargets[j];
    //         this.relatedTarget = target;
    //         break;
    //     }
    // }

    if (target) return target;
  }

  getActiveGroup(): Group | null {
    return this._activeGroup;
  }

  containsPoint(e: MouseEvent, target: FabricObject): boolean {
    let pointer = this.getPointer(e);
    let xy = this._normalizePointer(target, pointer);
    let x = xy.x;
    let y = xy.y;

    // we iterate through each object. If target found, return it.
    let iLines = target._getImageLines(target.oCoords);
    let xpoints = target._findCrossPoints(x, y, iLines);

    // if xcount is odd then we clicked inside the object
    // For the specific case of square images xcount === 1 in all true cases
    if ((xpoints && xpoints % 2 === 1) || target._findTargetCorner(e, this._offset)) {
      return true;
    }
    return false;
  }

  /** 获取相对于画布左上角的坐标 */
  getPointer(e: MouseEvent): Pos {
    let pointer = Util.getPointer(e, this.upperCanvasEl);
    return {
      x: pointer.x - this._offset.left,
      y: pointer.y - this._offset.top,
    };
  }

  /** 如果当前的物体在当前的组内，则要考虑扣去组的 top、left 值 */
  _normalizePointer(object: FabricObject, pointer: Pos) {
    let activeGroup = this.getActiveGroup(),
      x = pointer.x,
      y = pointer.y;

    let isObjectInGroup = activeGroup && object.type !== 'group' && activeGroup.contains(object);

    if (isObjectInGroup) {
      x -= activeGroup.left;
      y -= activeGroup.top;
    }
    return { x, y };
  }

  _shouldClearSelection(e: MouseEvent) {
    let target = this.findTarget(e),
      activeGroup = this.getActiveGroup();
    return !target || (target && activeGroup && !activeGroup.contains(target) && activeGroup !== target && !e.shiftKey);
  }

  /** 使所有元素失活，并触发相应事件 */
  deactivateAllWithDispatch(): Canvas {
    // let activeObject = this.getActiveGroup() || this.getActiveObject();
    // if (activeObject) {
    //     this.emit('before:selection:cleared', { target: activeObject });
    // }
    this.deactivateAll();
    // if (activeObject) {
    //     this.emit('selection:cleared');
    // }
    return this;
  }

  /** 将所有物体设置成未激活态 */
  deactivateAll() {
    let allObjects = this._objects;
    for (let i = 0; i < allObjects.length; i++) {
      allObjects[i].setActive(false);
    }
    this.discardActiveGroup();
    this.discardActiveObject();
    return this;
  }

  /** 将当前选中组失活 */
  discardActiveGroup(): Canvas {
    let g = this.getActiveGroup();
    if (g) g.destroy();
    return this.setActiveGroup(null);
  }

  setActiveGroup(group: Group | null): Canvas {
    this._activeGroup = group;
    if (group) {
      group.canvas = this;
      group.setActive(true);
    }
    return this;
  }

  /** 清空所有激活物体 */
  discardActiveObject() {
    if (this._activeObject) {
      this._activeObject.setActive(false);
    }
    this._activeObject = null;
    return this;
  }

  /** 是否要处理组的逻辑 */
  _shouldHandleGroupLogic(e: MouseEvent, target: FabricObject) {
    let activeObject = this._activeObject;
    return e.shiftKey && (this.getActiveGroup() || (activeObject && activeObject !== target));
  }

  _handleGroupLogic(e, target) {
    if (target === this.getActiveGroup()) {
      // if it's a group, find target again, this time skipping group
      target = this.findTarget(e, true);
      // if even object is not found, bail out
      if (!target || target.isType('group')) {
        return;
      }
    }
    let activeGroup = this.getActiveGroup();
    if (activeGroup) {
      if (activeGroup.contains(target)) {
        activeGroup.removeWithUpdate(target);
        this._resetObjectTransform(activeGroup);
        target.setActive(false);
        if (activeGroup.size() === 1) {
          // remove group alltogether if after removal it only contains 1 object
          this.discardActiveGroup();
        }
      } else {
        activeGroup.addWithUpdate(target);
        this._resetObjectTransform(activeGroup);
      }
      // this.emit('selection:created', { target: activeGroup, e: e });
      activeGroup.setActive(true);
    } else {
      // group does not exist
      if (this._activeObject) {
        // only if there's an active object
        if (target !== this._activeObject) {
          // and that object is not the actual target
          let group = new Group([this._activeObject, target]);
          this.setActiveGroup(group);
          activeGroup = this.getActiveGroup();
        }
      }
      // activate target object in any case
      target.setActive(true);
    }
    // if (activeGroup) {
    //     activeGroup.saveCoords();
    // }
  }
}
