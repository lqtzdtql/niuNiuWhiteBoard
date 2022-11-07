import { FabricObject } from './FabricObject';

/** 直线类 */
export class Line extends FabricObject {
  public type: string = 'line';

  constructor(options: any) {
    super(options);
    this.height = 30;
  }

  _render(ctx: CanvasRenderingContext2D) {
    ctx.beginPath();
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
    ctx.moveTo(-this.width / 2, 0);
    ctx.lineTo(this.width / 2, 0);
    ctx.closePath();

    ctx.stroke();
  }
}
