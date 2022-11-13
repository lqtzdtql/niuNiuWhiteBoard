import { FabricObject } from './FabricObject';

/** 三角形类 */
export class Triangle extends FabricObject {
  public type: string = 'Triangle';
  constructor(options: any) {
    super(options);
    this.height = this.width;
  }
  _render(ctx: CanvasRenderingContext2D) {
    let x = -this.width / 2,
      y = -this.height / 2,
      w = this.width,
      h = this.height;

    ctx.beginPath();
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
    ctx.moveTo(x, y + h);
    ctx.lineTo(0, y);
    ctx.lineTo(x + w, y + h);
    ctx.lineTo(x, y + h);
    ctx.closePath();

    if (this.fill) ctx.fill();

    if (this.stroke) ctx.stroke();
  }
}
