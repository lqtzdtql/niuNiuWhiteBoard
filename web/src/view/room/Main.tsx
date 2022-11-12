import triangle from '@Public/static/img/三角形.png';
import word from '@Public/static/img/文字.png';
import ellipse from '@Public/static/img/椭圆形.png';
import rubber from '@Public/static/img/橡皮擦.png';
import brush from '@Public/static/img/画笔.png';
import line from '@Public/static/img/直线.png';
import rectangle from '@Public/static/img/矩形.png';
import arrow from '@Public/static/img/箭头.png';
import { Canvas } from '@Sdk/Canvas';
import { Button, Form, Input, Popover } from 'antd';
import React, { useEffect, useRef } from 'react';

export const Main = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null as unknown as HTMLCanvasElement);
  const canvas = useRef<Canvas>(null as unknown as Canvas);

  useEffect(() => {
    canvas.current = new Canvas(canvasRef.current, {});
  }, []);

  const onFinish = (values: any) => {
    console.log('Success:', values);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('Failed:', errorInfo);
  };

  const modifyBrush = (type: number) => {
    canvas.current.modifyBrush({ type, stroke: 'green', fill: 'red', text: '哈哈哈' });
  };

  const revoke = () => {
    canvas.current.revoke();
    console.log('撤销');
  };

  const redo = () => {
    canvas.current.redo();
    console.log('重做');
  };

  const Content = () => {
    return (
      <Form
        name="basic"
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        initialValues={{ remember: true }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        autoComplete="off"
      >
        <Form.Item
          label="Username"
          name="username"
          rules={[{ required: true, message: 'Please input your username!' }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password"
          rules={[{ required: true, message: 'Please input your password!' }]}
        >
          <Input.Password />
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
          <Button type="primary" htmlType="submit">
            Submit
          </Button>
        </Form.Item>
      </Form>
    );
  };

  return (
    <div id="main" className="room-main">
      <div className="room-sider">
        <Popover content={<Content />} title="画笔" trigger="click" placement="right">
          <img className="room-sider-item" src={brush} />
        </Popover>
        <Popover content={<Content />} title="橡皮擦" trigger="click" placement="right">
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
      <canvas
        height={document.getElementById('main')?.clientHeight}
        width={((document.getElementById('main')?.clientHeight || 1) * 1920) / 1080}
        id="whiteboard"
        className="room-content"
        ref={canvasRef}
      />
    </div>
  );
};
