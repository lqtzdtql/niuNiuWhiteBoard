import { EventCenter } from './EventCenter';
import { Util } from './Util';

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

  /** 整个画布到上面和左边的偏移量 */
  private _offset: Offset;
  /** 画布中所有添加的物体 */
  private _objects: FabricObject[];

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
}
