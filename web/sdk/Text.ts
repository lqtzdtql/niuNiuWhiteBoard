import { FabricObject } from './FabricObject';
import { Util } from './Util';

/** 文字类 */
export class Text extends FabricObject {
  public type: string = 'text';
  public text: string = '';
  public size: number = 10;
  public font: string = 'Microsoft YaHei';
  public fillText: boolean = true;
  public strokeText: boolean = false;
  constructor(options: any) {
    super(options);
    this._initStateProperties();
    this._initText(options);
  }
  _initStateProperties() {
    this.stateProperties = this.stateProperties.concat(['text', 'size', 'font', 'fillText', 'strokeText']);
  }

  _initText(options: any) {
    this.text = options.text || '';
    this.font = options.font || 'Microsoft YaHei';
    this.size = options.size || 10;
    this.height = Util.getTextHeight(this.font, this.size);
    if (options.strokeText) {
      this.fillText = false;
      this.strokeText = true;
    }
    let ctx = document.createElement('canvas').getContext('2d');
    ctx!.font = this.font ? this.size + 'px ' + this.font : this.size + 'px';
    this.width = ctx!.measureText(this.text).width || 10;
    ctx = null;
  }
  _render(ctx: CanvasRenderingContext2D) {
    ctx.font = this.font ? this.size + 'px ' + this.font : this.size + 'px';
    if (this.group) ctx.translate(-this.group.width / 2 + this.width / 2, -this.group.height / 2 + this.height / 2);

    if (this.fillText) ctx.fillText(this.text, -this.width / 2, this.height / 2);

    if (this.strokeText) ctx.strokeText(this.text, 0, 0);
  }
}
