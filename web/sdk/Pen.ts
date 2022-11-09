// @ts-nocheck
import { FabricObject } from './FabricObject';
import { Util } from './Util';

/** 自由线条类 */
export class Pen extends FabricObject {
  public type: string = 'text';
  public penPath: [] = [];
  public start;
  constructor(options: any) {
    super(options);
    this._initStateProperties();
    this._initText(options);
  }
  _initStateProperties() {
    this.stateProperties = this.stateProperties.concat(['penPath', 'start']);
  }

  _initText(options: any) {
    this.penPath = options.penPath;
    this.start = [this.left, this.top];
  }

  _render(ctx: CanvasRenderingContext2D) {
    for (let i = 1; i < this.penPath.length; i++) {
      ctx.beginPath();
      if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
      ctx.moveTo(this.penPath[i - 1].x - this.start[0], this.penPath[i - 1].y - this.start[1]);
      ctx.lineTo(this.penPath[i].x - this.start[0], this.penPath[i].y - this.start[1]);
      ctx.closePath();
      ctx.stroke();
    }
  }
}
