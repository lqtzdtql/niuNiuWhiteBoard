// @ts-nocheck
import React, { useEffect, useRef } from 'react';
import { FabricObjects } from '../sdk/index';
import { Canvas } from '../sdk/Canvas';

function App() {
  const canvasRef = useRef(null);
  const canvas = useRef(null);
  useEffect(() => {
    canvas.current = new Canvas(canvasRef.current, {});

    // tri1.on('added', () => {
    //   console.log('tri1 被添加了');
    // });
    // rect.on('rotating', () => {
    //   console.log('rect 被旋转了');
    // });
    // rect.on('modified', () => {
    //   console.log('rect 被修改了');
    // });
  }, []);
  const modifyBrush = (type) => {
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

  return (
    <div className="app">
      <canvas id="canvas" width="800" height="600" ref={canvasRef}></canvas>
      <div
        onClick={() => {
          modifyBrush(0);
        }}
      >
        默认
      </div>
      <div
        onClick={() => {
          modifyBrush(1);
        }}
      >
        直线
      </div>
      <div
        onClick={() => {
          modifyBrush(3);
        }}
      >
        矩形
      </div>
      <div
        onClick={() => {
          modifyBrush(4);
        }}
      >
        菱形
      </div>
      <div
        onClick={() => {
          modifyBrush(5);
        }}
      >
        三角形
      </div>
      <div
        onClick={() => {
          modifyBrush(6);
        }}
      >
        圆形
      </div>
      <div
        onClick={() => {
          modifyBrush(7);
        }}
      >
        箭头
      </div>
      <div
        onClick={() => {
          modifyBrush(8);
        }}
      >
        橡皮
      </div>
      <div
        onClick={() => {
          modifyBrush(9);
        }}
      >
        文字
      </div>
      <div
        onClick={() => {
          modifyBrush(10);
        }}
      >
        自由线条
      </div>
      <div onClick={revoke}>撤销</div>
      <div onClick={redo}>重做</div>
    </div>
  );
}

export default App;
