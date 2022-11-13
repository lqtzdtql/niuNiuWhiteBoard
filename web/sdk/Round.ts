import { FabricObject } from './FabricObject';
import { Util } from './Util';

/** 圆形类 */
export class Round extends FabricObject {
  public type: string = 'Round';
  public roundAngle: number = 360;
  constructor(options: any) {
    super(options);
    this._initStateProperties();
    this._initAngle(options);
    this.height = this.width;
  }

  _initStateProperties() {
    this.stateProperties = this.stateProperties.concat(['roundAngle']);
  }

  _initAngle(options: any) {
    this.roundAngle = options.roundAngle || 360;
  }
  _render(ctx: CanvasRenderingContext2D) {
    let r = this.width / 2;
    ctx.beginPath();
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);
    ctx.arc(0, 0, r, 0, Util.degreesToRadians(this.roundAngle), false);
    ctx.closePath();

    if (this.fill) ctx.fill();

    if (this.stroke) ctx.stroke();
  }
}
