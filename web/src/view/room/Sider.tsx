import triangle from '@Public/static/img/三角形.png';
import word from '@Public/static/img/文字.png';
import ellipse from '@Public/static/img/椭圆形.png';
import rubber from '@Public/static/img/橡皮擦.png';
import brush from '@Public/static/img/画笔.png';
import line from '@Public/static/img/直线.png';
import rectangle from '@Public/static/img/矩形.png';
import arrow from '@Public/static/img/箭头.png';
import { Popover } from 'antd';
import React from 'react';

export const Sider: React.FC = () => {
  return (
    <div className="room-sider">
      <Popover content={<div>Content</div>} title="画笔" trigger="click" placement="right">
        <img className="room-sider-item" src={brush} />
      </Popover>
      <Popover content={<div>Content</div>} title="橡皮擦" trigger="click" placement="right">
        <img className="room-sider-item" src={rubber} />
      </Popover>
      <Popover content={<div>Content</div>} title="箭头" trigger="click" placement="right">
        <img className="room-sider-item" src={arrow} />
      </Popover>
      <Popover content={<div>Content</div>} title="直线" trigger="click" placement="right">
        <img className="room-sider-item" src={line} />
      </Popover>
      <Popover content={<div>Content</div>} title="文字" trigger="click" placement="right">
        <img className="room-sider-item" src={word} />
      </Popover>
      <Popover content={<div>Content</div>} title="矩形" trigger="click" placement="right">
        <img className="room-sider-item" src={rectangle} />
      </Popover>{' '}
      <Popover content={<div>Content</div>} title="三角形" trigger="click" placement="right">
        <img className="room-sider-item" src={triangle} />
      </Popover>
      <Popover content={<div>Content</div>} title="椭圆" trigger="click" placement="right">
        <img className="room-sider-item" src={ellipse} />
      </Popover>
    </div>
  );
};
