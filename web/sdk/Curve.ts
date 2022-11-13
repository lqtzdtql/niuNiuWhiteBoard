import { FabricObject } from './FabricObject';

/** 曲线类 */
export class Curve extends FabricObject {
  public type: string = 'Curve';
  public ca: number = 0;
  public ch: number = 0;
  constructor(options: any) {
    super(options);
    this._initStateProperties();
    this._initCa(options);
  }

  _initStateProperties() {
    this.stateProperties = this.stateProperties.concat(['ca', 'ch']);
  }

  _initCa(options: any) {
    this.ca = options.ca || 0;
    this.ch = options.ch || 0;
  }

  _render(ctx: CanvasRenderingContext2D) {
    let x = -this.width / 2,
      y = -this.height / 2,
      w = this.width,
      h = this.height;

    ctx.beginPath();
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
    ctx.moveTo(x, y);
    ctx.quadraticCurveTo((w * this.ca) / 100, h, x + w, y + (h * this.ch) / 100);
    ctx.closePath();

    if (this.fill) ctx.fill();

    if (this.stroke) ctx.stroke();
  }
}
