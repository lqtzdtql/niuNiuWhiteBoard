import { FabricObject } from './FabricObject';

/** 直线类 */
export class Line extends FabricObject {
  public type: string = 'Line';
  public zHeight: number = 0;
  public direction: number = 0;

  constructor(options: any) {
    super(options);
    this._initStateProperties();
    this._initZHeight(options);
  }

  _initStateProperties() {
    this.stateProperties = this.stateProperties.concat(['zHeight', 'direction']);
  }

  _initZHeight(options) {
    this.zHeight = options.zHeight;
    this.direction = options.direction;
  }

  _render(ctx: CanvasRenderingContext2D) {
    ctx.beginPath();
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
    if (this.direction >= 0) {
      ctx.moveTo(-this.width / 2, -this.zHeight / 2 + 1);
      ctx.lineTo(this.width / 2, this.zHeight / 2 - 1);
    } else {
      ctx.moveTo(this.width / 2, -this.zHeight / 2 + 1);
      ctx.lineTo(-this.width / 2, this.zHeight / 2 - 1);
    }
    ctx.closePath();
    ctx.stroke();
  }
}
