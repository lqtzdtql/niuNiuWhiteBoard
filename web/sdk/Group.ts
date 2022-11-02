import { FabricObject } from './FabricObject';
import { Util } from './Util';

/**
 * 组类，也就是拖蓝框选区域包围的那些物体构成了一个组
 * Group 虽然继承至 FabricObject，但是要注意获取某些属性有时是没有的
 */

export class Group extends FabricObject {
  public type: string = 'group';
  /** 组中所有的物体 */
  public objects: FabricObject[];

  /** 物体是否在 group 中 */
  contains(object: FabricObject) {
    return this.objects.indexOf(object) > -1;
  }

  destroy() {
    return this._restoreObjectsState();
  }

  /** 还原创建 group 之前的状态 */
  _restoreObjectsState(): Group {
    this.objects.forEach(this._restoreObjectState, this);
    return this;
  }

  /** 还原 group 中某个物体的初始状态 */
  _restoreObjectState(object: FabricObject): Group {
    let groupLeft = this.get('left'),
      groupTop = this.get('top'),
      groupAngle = this.getAngle() * (Math.PI / 180),
      rotatedTop = Math.cos(groupAngle) * object.get('top') + Math.sin(groupAngle) * object.get('left'),
      rotatedLeft = -Math.sin(groupAngle) * object.get('top') + Math.cos(groupAngle) * object.get('left');

    object.setAngle(object.getAngle() + this.getAngle());

    object.set('left', groupLeft + rotatedLeft * this.get('scaleX'));
    object.set('top', groupTop + rotatedTop * this.get('scaleY'));

    object.set('scaleX', object.get('scaleX') * this.get('scaleX'));
    object.set('scaleY', object.get('scaleY') * this.get('scaleY'));

    object.setCoords();
    object.hasControls = object.orignHasControls;
    // delete object.__origHasControls;
    object.setActive(false);
    object.setCoords();

    return this;
  }

  get(prop: string) {
    // 组里面有很多元素，所以虽然继承至 Fabric，但是有很多属性读取是无效的，设置同理
    return this[prop];
  }
}
