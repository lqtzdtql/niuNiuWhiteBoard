export class Util {
  /** 给元素添加类名 */
  static addClass(element: HTMLElement, className: string) {
    if ((' ' + element.className + ' ').indexOf(' ' + className + ' ') === -1) {
      element.className += (element.className ? ' ' : '') + className;
    }
  }

  /** 给元素设置样式 */
  static setStyle(element: HTMLElement, styles) {
    let elementStyle = element.style;

    if (typeof styles === 'string') {
      element.style.cssText += ';' + styles;
      return styles.indexOf('opacity') > -1
        ? Util.setOpacity(element, styles.match(/opacity:\s*(\d?\.?\d*)/)[1])
        : element;
    }
    for (let property in styles) {
      if (property === 'opacity') {
        Util.setOpacity(element, styles[property]);
      } else {
        elementStyle[property] = styles[property];
      }
    }
    return element;
  }

  /** 设置元素透明度 */
  static setOpacity(element: HTMLElement, value: string) {
    element.style.opacity = value;
    return element;
  }

  /** 设置 css 的 userSelect 样式为 none，也就是不可选中的状态 */
  static makeElementUnselectable(element: HTMLElement): HTMLElement {
    element.style.userSelect = 'none';
    return element;
  }

  /** 计算元素偏移值 */
  static getElementOffset(element): Offset {
    let valueT = 0;
    let valueL = 0;
    do {
      valueT += element.offsetTop || 0;
      valueL += element.offsetLeft || 0;
      element = element.offsetParent;
    } while (element);
    return { left: valueL, top: valueT };
  }
}
