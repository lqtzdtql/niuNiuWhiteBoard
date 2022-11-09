import { FabricObject } from './FabricObject';
import { Rect } from './Rect';
import { Triangle } from './Triangle';
import { Round } from './Round';
import { Group } from './Group';
import { Curve } from './Curve';
import { Diamond } from './Diamond';
import { Line } from './Line';
import { Arrow } from './Arrow';
import { Text } from './Text';
import { Pen } from './Pen';
import { FabricImage } from './FabricImage';

// 最终导出的东西都挂载到 fabric 上面

export class FabricObjects {
  static FabricObject = FabricObject;
  static Rect = Rect;
  static Triangle = Triangle;
  static Round = Round;
  static Curve = Curve;
  static Diamond = Diamond;
  static Line = Line;
  static Arrow = Arrow;
  static Text = Text;
  static Pen = Pen;
  static Group = Group;
  static FabricImage = FabricImage;
}
