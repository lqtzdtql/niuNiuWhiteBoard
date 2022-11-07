// @ts-nocheck
import React, { useEffect, useRef } from 'react';
import { fabric } from '../sdk/index';
import { Canvas } from '../sdk/Canvas';

function App() {
  const canvasRef = useRef(null);
  const canvas = useRef(null);
  useEffect(() => {
    canvas.current = new Canvas(canvasRef.current, {});
    const rect = new fabric.Rect({
      top: 200,
      left: 100,
      width: 100,
      height: 100,
      fill: 'green',
      rx: 20,
      ry: 20,
    });
    const rect2 = new fabric.Rect({
      top: 300,
      left: 100,
      width: 100,
      height: 100,
      fill: 'blue',
      angle: 45,
    });
    const rect3 = new fabric.Rect({
      top: 200,
      left: 200,
      width: 50,
      height: 50,
      fill: 'pink',
    });
    const rect4 = new fabric.Rect({
      top: 200,
      left: 300,
      width: 50,
      height: 50,
      fill: 'pink',
    });
    const tri1 = new fabric.Triangle({
      top: 200,
      left: 200,
      width: 50,
      height: 50,
      stroke: 'green',
      fill: 'white',
    });
    const round = new fabric.Round({
      top: 200,
      left: 200,
      width: 100,
      height: 200,
      fill: 'black',
    });
    const curve = new fabric.Curve({
      top: 200,
      left: 200,
      width: 50,
      height: 50,
      ca: 100,
      ch: 50,
      stroke: 'red',
      fill: 'white',
    });
    const diamond = new fabric.Diamond({
      top: 200,
      left: 200,
      width: 150,
      height: 50,
      fill: 'green',
    });
    const line = new fabric.Line({
      start: { x: 100, y: 100 },
      end: { x: 300, y: 300 },
      top: 200,
      left: 200,
      width: 200,
      height: 200,
      stroke: 'green',
    });
    // const line = new fabric.Line({
    //   top: 300,
    //   left: 250,
    //   width: 100,
    //   height: 50,
    //   stroke: 'black',
    // });
    const arrow = new fabric.Arrow({
      top: 300,
      left: 250,
      width: 100,
      height: 50,
      stroke: 'red',
    });
    const g = new fabric.Group([rect3, rect4], {
      top: 200,
      left: 200,
    });

    tri1.on('added', () => {
      console.log('tri1 被添加了');
    });
    rect.on('rotating', () => {
      console.log('rect 被旋转了');
    });
    rect.on('modified', () => {
      console.log('rect 被修改了');
    });
  }, []);
  const modifyBrush = (type) => {
    canvas.current.modifyBrush({ type, stroke: 'green', fill: 'red' });
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
    </div>
  );
}

export default App;
