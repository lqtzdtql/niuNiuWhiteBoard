import { FabricObject } from './FabricObject';

/** 箭头类 */
export class Arrow extends FabricObject {
  public type: string = 'Arrow';
  public headlen: number = 15;
  public direction: number = 0;
  constructor(options: any) {
    super(options);
    this._initStateProperties();
    this._initArrow(options);
  }

  _initStateProperties() {
    this.stateProperties = this.stateProperties.concat(['theta', 'headlen', 'direction']);
  }

  _initArrow(options: any) {
    this.direction = options.direction || 0;
    this.headlen = options.headlen || 15;
    this.height = 30;
  }

  _render(ctx: CanvasRenderingContext2D) {
    let x = -this.width / 2,
      w = this.width;
    ctx.beginPath();
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
    if (this.direction >= 0) {
      ctx.moveTo(x, 0);
      ctx.lineTo(x + w - this.headlen, 0);
      ctx.lineTo(x + w - this.headlen, this.headlen);
      ctx.lineTo(x + w, 0);
      ctx.lineTo(x + w - this.headlen, -this.headlen);
      ctx.lineTo(x + w - this.headlen, 0);
    } else {
      ctx.moveTo(x + w, 0);
      ctx.lineTo(x + this.headlen, 0);
      ctx.lineTo(x + this.headlen, this.headlen);
      ctx.lineTo(x, 0);
      ctx.lineTo(x + this.headlen, -this.headlen);
      ctx.lineTo(x + this.headlen, 0);
    }
    ctx.closePath();

    ctx.stroke();
  }
}
