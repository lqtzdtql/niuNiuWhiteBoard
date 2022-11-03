import { Util } from './Util';
import { Point } from './Point';
import { Intersection } from './Intersection';
import { Offset, Coords, Corner, IAnimationOption } from './interface';
import { EventCenter } from './EventCenter';

/** 物体基类，有一些共同属性和方法 */
export class FabricObject extends EventCenter {
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
  public originalState: {} = {};
  /** 物体所属的组 */
  public group;
  /** 物体被拖蓝选区保存的时候需要临时保存下 hasControls 的值 */
  public originHasControls: boolean = true;
  public stateProperties: string[] = (
    'top left width height scaleX scaleY ' +
    'flipX flipY angle cornerSize fill originX originY ' +
    'stroke strokeWidth ' +
    'borderWidth transformMatrix visible'
  ).split(' ');

  private _cacheCanvas: HTMLCanvasElement;
  private _cacheContext: CanvasRenderingContext2D;
  public cacheWidth: number;
  public cacheHeight: number;
  public dirty: boolean;

  constructor(options) {
    super();
    this.initialize(options);
  }

  initialize(options) {
    options && this.setOptions(options);
  }

  setOptions(options) {
    for (let prop in options) {
      this[prop] = options[prop];
    }
  }

  /** 获取包围盒的四条边 */
  _getImageLines(corner: Corner) {
    return {
      topline: {
        o: corner.tl,
        d: corner.tr,
      },
      rightline: {
        o: corner.tr,
        d: corner.br,
      },
      bottomline: {
        o: corner.br,
        d: corner.bl,
      },
      leftline: {
        o: corner.bl,
        d: corner.tl,
      },
    };
  }

  /**
   * 射线检测法：以鼠标坐标点为参照，水平向右做一条射线，求坐标点与多条边的交点个数
   * 如果和物体相交的个数为偶数点则点在物体外部；如果为奇数点则点在内部
   * 不过 fabric 的点选多边形都是用于包围盒，也就是矩形，所以该方法是专门针对矩形的，并且针对矩形做了一些优化
   */
  _findCrossPoints(ex: number, ey: number, lines): number {
    let b1, // 射线的斜率
      b2, // 边的斜率
      a1,
      a2,
      xi, // 射线与边的交点
      // yi, // 射线与边的交点
      xcount = 0,
      iLine; // 当前边

    // 遍历包围盒的四条边
    for (let lineKey in lines) {
      iLine = lines[lineKey];

      // 优化1：如果边的两个端点的 y 值都小于鼠标点的 y 值，则跳过
      if (iLine.o.y < ey && iLine.d.y < ey) continue;
      // 优化2：如果边的两个端点的 y 值都大于鼠标点的 y 值，则跳过
      if (iLine.o.y >= ey && iLine.d.y >= ey) continue;

      // 优化3：如果边是一条垂线
      if (iLine.o.x === iLine.d.x && iLine.o.x >= ex) {
        xi = iLine.o.x;
        // yi = ey;
      } else {
        // 简单计算下射线与边的交点，看式子容易晕，建议自己手动算一下
        b1 = 0;
        b2 = (iLine.d.y - iLine.o.y) / (iLine.d.x - iLine.o.x);
        a1 = ey - b1 * ex;
        a2 = iLine.o.y - b2 * iLine.o.x;

        xi = -(a1 - a2) / (b1 - b2);
        // yi = a1 + b1 * xi;
      }
      // 只需要计数 xi >= ex 的情况
      if (xi >= ex) {
        xcount += 1;
      }
      // 优化4：因为 fabric 中的多边形只需要用到矩形，所以根据矩形的特质，顶多只有两个交点，所以我们可以提前结束循环
      if (xcount === 2) {
        break;
      }
    }
    return xcount;
  }

  /** 检测哪个控制点被点击了 */
  _findTargetCorner(e: MouseEvent, offset: Offset): boolean | string {
    if (!this.hasControls || !this.active) return false;

    let pointer = Util.getPointer(e, this.canvas.upperCanvasEl);
    let ex = pointer.x - offset.left;
    let ey = pointer.y - offset.top;
    let xpoints;
    let lines;

    for (let i in this.oCoords) {
      if (i === 'mtr' && !this.hasRotatingPoint) {
        continue;
      }

      lines = this._getImageLines(this.oCoords[i].corner);

      // debugger 绘制物体控制点的四个顶点
      // this.canvas.contextTop.fillRect(lines.bottomline.d.x, lines.bottomline.d.y, 2, 2);
      // this.canvas.contextTop.fillRect(lines.bottomline.o.x, lines.bottomline.o.y, 2, 2);

      // this.canvas.contextTop.fillRect(lines.leftline.d.x, lines.leftline.d.y, 2, 2);
      // this.canvas.contextTop.fillRect(lines.leftline.o.x, lines.leftline.o.y, 2, 2);

      // this.canvas.contextTop.fillRect(lines.topline.d.x, lines.topline.d.y, 2, 2);
      // this.canvas.contextTop.fillRect(lines.topline.o.x, lines.topline.o.y, 2, 2);

      // this.canvas.contextTop.fillRect(lines.rightline.d.x, lines.rightline.d.y, 2, 2);
      // this.canvas.contextTop.fillRect(lines.rightline.o.x, lines.rightline.o.y, 2, 2);

      xpoints = this._findCrossPoints(ex, ey, lines);
      if (xpoints % 2 === 1 && xpoints !== 0) {
        return i;
      }
    }
    return false;
  }

  setActive(active: boolean = false): FabricObject {
    this.active = !!active;
    return this;
  }

  getAngle(): number {
    return this.angle;
  }

  get(key: string) {
    return this[key];
  }

  setAngle(angle: number) {
    this.angle = angle;
  }

  set(key: string, value): FabricObject {
    // if (typeof value === 'function') value = value(this.get(key));
    // if (key === 'scaleX' || key === 'scaleY') {
    //     value = this._constrainScale(value);
    // }
    // if (key === 'width' || key === 'height') {
    //     this.minScaleLimit = Util.toFixed(Math.min(0.1, 1 / Math.max(this.width, this.height)), 2);
    // }
    if (key === 'scaleX' && value < 0) {
      this.flipX = !this.flipX;
      value *= -1;
    } else if (key === 'scaleY' && value < 0) {
      this.flipY = !this.flipY;
      value *= -1;
    }
    this[key] = value;
    return this;
  }

  /** 重新设置物体包围盒的边框和各个控制点，包括位置和大小 */
  setCoords(): FabricObject {
    let strokeWidth = this.strokeWidth > 1 ? this.strokeWidth : 0,
      padding = this.padding,
      radian = Util.degreesToRadians(this.angle);

    this.currentWidth = (this.width + strokeWidth) * this.scaleX + padding * 2;
    this.currentHeight = (this.height + strokeWidth) * this.scaleY + padding * 2;

    // If width is negative, make postive. Fixes path selection issue
    // if (this.currentWidth < 0) {
    //     this.currentWidth = Math.abs(this.currentWidth);
    // }

    // 物体中心点到顶点的斜边长度
    let _hypotenuse = Math.sqrt(Math.pow(this.currentWidth / 2, 2) + Math.pow(this.currentHeight / 2, 2));
    let _angle = Math.atan(this.currentHeight / this.currentWidth);
    // let _angle = Math.atan2(this.currentHeight, this.currentWidth);

    // offset added for rotate and scale actions
    let offsetX = Math.cos(_angle + radian) * _hypotenuse,
      offsetY = Math.sin(_angle + radian) * _hypotenuse,
      sinTh = Math.sin(radian),
      cosTh = Math.cos(radian);

    let coords = this.getCenterPoint();
    let tl = {
      x: coords.x - offsetX,
      y: coords.y - offsetY,
    };
    let tr = {
      x: tl.x + this.currentWidth * cosTh,
      y: tl.y + this.currentWidth * sinTh,
    };
    let br = {
      x: tr.x - this.currentHeight * sinTh,
      y: tr.y + this.currentHeight * cosTh,
    };
    let bl = {
      x: tl.x - this.currentHeight * sinTh,
      y: tl.y + this.currentHeight * cosTh,
    };
    let ml = {
      x: tl.x - (this.currentHeight / 2) * sinTh,
      y: tl.y + (this.currentHeight / 2) * cosTh,
    };
    let mt = {
      x: tl.x + (this.currentWidth / 2) * cosTh,
      y: tl.y + (this.currentWidth / 2) * sinTh,
    };
    let mr = {
      x: tr.x - (this.currentHeight / 2) * sinTh,
      y: tr.y + (this.currentHeight / 2) * cosTh,
    };
    let mb = {
      x: bl.x + (this.currentWidth / 2) * cosTh,
      y: bl.y + (this.currentWidth / 2) * sinTh,
    };
    let mtr = {
      x: tl.x + (this.currentWidth / 2) * cosTh,
      y: tl.y + (this.currentWidth / 2) * sinTh,
    };

    // clockwise
    this.oCoords = { tl, tr, br, bl, ml, mt, mr, mb, mtr };

    // set coordinates of the draggable boxes in the corners used to scale/rotate the image
    this._setCornerCoords();

    return this;
  }

  /** 获取物体中心点 */
  getCenterPoint() {
    return this.translateToCenterPoint(new Point(this.left, this.top), this.originX, this.originY);
  }

  /** 将中心点移到变换基点 */
  translateToCenterPoint(point: Point, originX: string, originY: string): Point {
    let cx = point.x,
      cy = point.y;

    if (originX === 'left') {
      cx = point.x + this.getWidth() / 2;
    } else if (originX === 'right') {
      cx = point.x - this.getWidth() / 2;
    }

    if (originY === 'top') {
      cy = point.y + this.getHeight() / 2;
    } else if (originY === 'bottom') {
      cy = point.y - this.getHeight() / 2;
    }
    const p = new Point(cx, cy);
    if (this.angle) {
      return Util.rotatePoint(p, point, Util.degreesToRadians(this.angle));
    } else {
      return p;
    }
  }

  /** 获取当前大小，包含缩放效果 */
  getWidth(): number {
    return this.width * this.scaleX;
  }
  /** 获取当前大小，包含缩放效果 */
  getHeight(): number {
    return this.height * this.scaleY;
  }

  /** 重新设置物体的每个控制点，包括位置和大小 */
  _setCornerCoords() {
    let coords = this.oCoords,
      radian = Util.degreesToRadians(this.angle),
      newTheta = Util.degreesToRadians(45 - this.angle),
      cornerHypotenuse = Math.sqrt(2 * Math.pow(this.cornerSize, 2)) / 2,
      cosHalfOffset = cornerHypotenuse * Math.cos(newTheta),
      sinHalfOffset = cornerHypotenuse * Math.sin(newTheta),
      sinTh = Math.sin(radian),
      cosTh = Math.cos(radian);

    coords.tl.corner = {
      tl: {
        x: coords.tl.x - sinHalfOffset,
        y: coords.tl.y - cosHalfOffset,
      },
      tr: {
        x: coords.tl.x + cosHalfOffset,
        y: coords.tl.y - sinHalfOffset,
      },
      bl: {
        x: coords.tl.x - cosHalfOffset,
        y: coords.tl.y + sinHalfOffset,
      },
      br: {
        x: coords.tl.x + sinHalfOffset,
        y: coords.tl.y + cosHalfOffset,
      },
    };

    coords.tr.corner = {
      tl: {
        x: coords.tr.x - sinHalfOffset,
        y: coords.tr.y - cosHalfOffset,
      },
      tr: {
        x: coords.tr.x + cosHalfOffset,
        y: coords.tr.y - sinHalfOffset,
      },
      br: {
        x: coords.tr.x + sinHalfOffset,
        y: coords.tr.y + cosHalfOffset,
      },
      bl: {
        x: coords.tr.x - cosHalfOffset,
        y: coords.tr.y + sinHalfOffset,
      },
    };

    coords.bl.corner = {
      tl: {
        x: coords.bl.x - sinHalfOffset,
        y: coords.bl.y - cosHalfOffset,
      },
      bl: {
        x: coords.bl.x - cosHalfOffset,
        y: coords.bl.y + sinHalfOffset,
      },
      br: {
        x: coords.bl.x + sinHalfOffset,
        y: coords.bl.y + cosHalfOffset,
      },
      tr: {
        x: coords.bl.x + cosHalfOffset,
        y: coords.bl.y - sinHalfOffset,
      },
    };

    coords.br.corner = {
      tr: {
        x: coords.br.x + cosHalfOffset,
        y: coords.br.y - sinHalfOffset,
      },
      bl: {
        x: coords.br.x - cosHalfOffset,
        y: coords.br.y + sinHalfOffset,
      },
      br: {
        x: coords.br.x + sinHalfOffset,
        y: coords.br.y + cosHalfOffset,
      },
      tl: {
        x: coords.br.x - sinHalfOffset,
        y: coords.br.y - cosHalfOffset,
      },
    };

    coords.ml.corner = {
      tl: {
        x: coords.ml.x - sinHalfOffset,
        y: coords.ml.y - cosHalfOffset,
      },
      tr: {
        x: coords.ml.x + cosHalfOffset,
        y: coords.ml.y - sinHalfOffset,
      },
      bl: {
        x: coords.ml.x - cosHalfOffset,
        y: coords.ml.y + sinHalfOffset,
      },
      br: {
        x: coords.ml.x + sinHalfOffset,
        y: coords.ml.y + cosHalfOffset,
      },
    };

    coords.mt.corner = {
      tl: {
        x: coords.mt.x - sinHalfOffset,
        y: coords.mt.y - cosHalfOffset,
      },
      tr: {
        x: coords.mt.x + cosHalfOffset,
        y: coords.mt.y - sinHalfOffset,
      },
      bl: {
        x: coords.mt.x - cosHalfOffset,
        y: coords.mt.y + sinHalfOffset,
      },
      br: {
        x: coords.mt.x + sinHalfOffset,
        y: coords.mt.y + cosHalfOffset,
      },
    };

    coords.mr.corner = {
      tl: {
        x: coords.mr.x - sinHalfOffset,
        y: coords.mr.y - cosHalfOffset,
      },
      tr: {
        x: coords.mr.x + cosHalfOffset,
        y: coords.mr.y - sinHalfOffset,
      },
      bl: {
        x: coords.mr.x - cosHalfOffset,
        y: coords.mr.y + sinHalfOffset,
      },
      br: {
        x: coords.mr.x + sinHalfOffset,
        y: coords.mr.y + cosHalfOffset,
      },
    };

    coords.mb.corner = {
      tl: {
        x: coords.mb.x - sinHalfOffset,
        y: coords.mb.y - cosHalfOffset,
      },
      tr: {
        x: coords.mb.x + cosHalfOffset,
        y: coords.mb.y - sinHalfOffset,
      },
      bl: {
        x: coords.mb.x - cosHalfOffset,
        y: coords.mb.y + sinHalfOffset,
      },
      br: {
        x: coords.mb.x + sinHalfOffset,
        y: coords.mb.y + cosHalfOffset,
      },
    };

    coords.mtr.corner = {
      tl: {
        x: coords.mtr.x - sinHalfOffset + sinTh * this.rotatingPointOffset,
        y: coords.mtr.y - cosHalfOffset - cosTh * this.rotatingPointOffset,
      },
      tr: {
        x: coords.mtr.x + cosHalfOffset + sinTh * this.rotatingPointOffset,
        y: coords.mtr.y - sinHalfOffset - cosTh * this.rotatingPointOffset,
      },
      bl: {
        x: coords.mtr.x - cosHalfOffset + sinTh * this.rotatingPointOffset,
        y: coords.mtr.y + sinHalfOffset - cosTh * this.rotatingPointOffset,
      },
      br: {
        x: coords.mtr.x + sinHalfOffset + sinTh * this.rotatingPointOffset,
        y: coords.mtr.y + cosHalfOffset - cosTh * this.rotatingPointOffset,
      },
    };
  }

  /** 保存物体当前的状态到 originalState 中 */
  saveState(): FabricObject {
    this.originalState = {};
    this.stateProperties.forEach((prop) => {
      this.originalState[prop] = this[prop];
    });
    return this;
  }

  /**
   * 平移坐标系到中心点
   * @param center
   * @param {string} originX  left | center | right
   * @param {string} originY  top | center | bottom
   * @returns
   */
  translateToOriginPoint(center: Point, originX: string, originY: string): Point {
    let x = center.x,
      y = center.y;

    // Get the point coordinates
    if (originX === 'left') {
      x = center.x - this.getWidth() / 2;
    } else if (originX === 'right') {
      x = center.x + this.getWidth() / 2;
    }
    if (originY === 'top') {
      y = center.y - this.getHeight() / 2;
    } else if (originY === 'bottom') {
      y = center.y + this.getHeight() / 2;
    }

    // Apply the rotation to the point (it's already scaled properly)
    return Util.rotatePoint(new Point(x, y), center, Util.degreesToRadians(this.angle));
  }

  /** 转换成本地坐标 */
  toLocalPoint(point: Point, originX: string, originY: string): Point {
    let center = this.getCenterPoint();

    let x, y;
    if (originX !== undefined && originY !== undefined) {
      if (originX === 'left') {
        x = center.x - this.getWidth() / 2;
      } else if (originX === 'right') {
        x = center.x + this.getWidth() / 2;
      } else {
        x = center.x;
      }

      if (originY === 'top') {
        y = center.y - this.getHeight() / 2;
      } else if (originY === 'bottom') {
        y = center.y + this.getHeight() / 2;
      } else {
        y = center.y;
      }
    } else {
      x = this.left;
      y = this.top;
    }

    return Util.rotatePoint(new Point(point.x, point.y), center, -Util.degreesToRadians(this.angle)).subtractEquals(
      new Point(x, y),
    );
  }

  /**
   * 根据物体的 origin 来设置物体的位置
   * @method setPositionByOrigin
   * @param {Point} pos
   * @param {string} originX left | center | right
   * @param {string} originY top | center | bottom
   */
  setPositionByOrigin(pos: Point, originX: string, originY: string) {
    let center = this.translateToCenterPoint(pos, originX, originY);
    let position = this.translateToOriginPoint(center, this.originX, this.originY);
    // console.log(`更新缩放的物体位置:[${position.x}，${position.y}]`);
    this.set('left', position.x);
    this.set('top', position.y);
  }

  /**
   * 物体与框选区域是否相交，用框选区域的四条边分别与物体的四条边求交
   * @param {Point} selectionTL 拖蓝框选区域左上角的点
   * @param {Point} selectionBR 拖蓝框选区域右下角的点
   * @returns {boolean}
   */
  intersectsWithRect(selectionTL: Point, selectionBR: Point): boolean {
    let oCoords = this.oCoords,
      tl = new Point(oCoords.tl.x, oCoords.tl.y),
      tr = new Point(oCoords.tr.x, oCoords.tr.y),
      bl = new Point(oCoords.bl.x, oCoords.bl.y),
      br = new Point(oCoords.br.x, oCoords.br.y);

    let intersection = Intersection.intersectPolygonRectangle([tl, tr, br, bl], selectionTL, selectionBR);
    return intersection.status === 'Intersection';
  }

  /**
   * 物体是否被框选区域包含
   * @param {Point} selectionTL 拖蓝框选区域左上角的点
   * @param {Point} selectionBR 拖蓝框选区域右下角的点
   * @returns {boolean}
   */
  isContainedWithinRect(selectionTL: Point, selectionBR: Point): boolean {
    let oCoords = this.oCoords,
      tl = new Point(oCoords.tl.x, oCoords.tl.y),
      tr = new Point(oCoords.tr.x, oCoords.tr.y),
      bl = new Point(oCoords.bl.x, oCoords.bl.y);

    return tl.x > selectionTL.x && tr.x < selectionBR.x && tl.y > selectionTL.y && bl.y < selectionBR.y;
  }
}
