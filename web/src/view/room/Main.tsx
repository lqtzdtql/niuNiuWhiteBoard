import triangle from '@Public/static/img/三角形.png';
import revokeImg from '@Public/static/img/撤销.png';
import word from '@Public/static/img/文字.png';
import ellipse from '@Public/static/img/椭圆形.png';
import rubber from '@Public/static/img/橡皮擦.png';
import brush from '@Public/static/img/画笔.png';
import line from '@Public/static/img/直线.png';
import rectangle from '@Public/static/img/矩形.png';
import arrow from '@Public/static/img/箭头.png';
import diamond from '@Public/static/img/菱形.png';
import pointer from '@Public/static/img/选择.png';
import redoImg from '@Public/static/img/重做.png';
import { Button, Form, Input, Popover, Radio, Select } from 'antd';
import { Colorpicker } from 'antd-colorpicker';
import React, { useContext, useState } from 'react';
import { Context } from './context';
import { TitleProps } from './Title';

interface brushOptions {
  type: number; //笔刷类型
  fill?: string; //填充颜色
  stroke?: string; //描边颜色
  strokeWidth?: number; //描边线宽
  text?: string; //要渲染的文字
  size?: number; //字号
  font?: string; //除字号外其余字体样式
  fillText?: boolean; //文字填充渲染
  strokeText?: boolean; //文字描边渲染
}

export const Main = (props: TitleProps) => {
  const [openPen, setOpenPen] = useState(false);
  const [openArrow, setOpenArrow] = useState(false);
  const [openLine, setOpenLine] = useState(false);
  const [openWord, setOpenWord] = useState(false);
  const [openRect, setOpenRect] = useState(false);
  const [openTri, setOpenTri] = useState(false);
  const [openEl, setOpenEl] = useState(false);
  const [openDiamond, setOpenDiamond] = useState(false);
  const [openNote, setOpenNote] = useState(false);

  const { canvasClass, canvasRef } = useContext(Context);

  const whiteboardHeight = document.body.clientHeight - 52 || 1080;

  // useEffect(() => {
  //   if (canvasRef.current) {
  //     canvasClass.current = new Canvas(canvasRef.current, {});
  //   }
  // }, []);

  const handleOpenPenChange = (newOpen: boolean) => {
    setOpenPen(newOpen);
  };

  const handleOpenArrowChange = (newOpen: boolean) => {
    setOpenArrow(newOpen);
  };
  const handleOpenLineChange = (newOpen: boolean) => {
    setOpenLine(newOpen);
  };
  const handleOpenWordChange = (newOpen: boolean) => {
    setOpenWord(newOpen);
  };
  const handleOpenRectChange = (newOpen: boolean) => {
    setOpenRect(newOpen);
  };
  const handleOpenTriChange = (newOpen: boolean) => {
    setOpenTri(newOpen);
  };
  const handleOpenElChange = (newOpen: boolean) => {
    setOpenEl(newOpen);
  };
  const handleOpenDiamondChange = (newOpen: boolean) => {
    setOpenDiamond(newOpen);
  };
  const handleOpenNoteChange = (newOpen: boolean) => {
    setOpenNote(newOpen);
  };

  const onFinish = (values: any) => {
    setOpenPen(false);
    setOpenArrow(false);
    setOpenLine(false);
    setOpenWord(false);
    setOpenRect(false);
    setOpenTri(false);
    setOpenEl(false);
    setOpenDiamond(false);
    setOpenNote(false);

    values.fillText = values.fillText === 'fillText';
    values.strokeText = !values.fillText;

    values.stroke =
      `rgba(${values.stroke.rgb.r},${values.stroke.rgb.g},${values.stroke.rgb.b},${values.stroke.rgb.a})` || 'black';
    values.fill =
      `rgba(${values.fill.rgb.r},${values.fill.rgb.g},${values.fill.rgb.b},${values.fill.rgb.a})` || 'black';

    modifyBrush(values);
  };

  const modifyBrush = (options: brushOptions) => {
    if (canvasClass.current) {
      canvasClass.current.modifyBrush(options);
      console.log(canvasClass.current);
    }
  };

  const Content = ({ type }: { type: number }) => {
    return (
      <Form
        name="basic"
        onFinish={onFinish}
        initialValues={{
          type: type, //笔刷类型
          fill: {
            hsl: {
              h: 0,
              s: 0,
              l: 0,
              a: 1,
            },
            hex: '#000000',
            rgb: {
              r: 0,
              g: 0,
              b: 0,
              a: 1,
            },
            hsv: {
              h: 0,
              s: 0,
              v: 0,
              a: 1,
            },
            oldHue: 0,
            source: 'hsv',
          }, //填充颜色
          stroke: {
            hsl: {
              h: 0,
              s: 0,
              l: 0,
              a: 1,
            },
            hex: '#000000',
            rgb: {
              r: 0,
              g: 0,
              b: 0,
              a: 1,
            },
            hsv: {
              h: 0,
              s: 0,
              v: 0,
              a: 1,
            },
            oldHue: 0,
            source: 'hsv',
          }, //描边颜色
          strokeWidth: 2, //描边线宽
          text: '', //要渲染的文字
          size: 18, //字号
          font: 'serif', //除字号外其余字体样式
          fillText: 'fillText', //文字填充渲染
          strokeText: false, //文字描边渲染
        }}
        autoComplete="off"
      >
        <Form.Item hidden name="type" />
        <Form.Item hidden={type === 1} label={'颜色'} name="stroke">
          <Colorpicker popup={true} />
        </Form.Item>
        <Form.Item hidden={type === 10 || type === 11 || type === 7 || type === 1} label={'填充颜色'} name="fill">
          <Colorpicker popup={true} />
        </Form.Item>
        <Form.Item hidden={type === 9} label="线条宽度(px)" name="strokeWidth">
          <Input />
        </Form.Item>
        <Form.Item
          hidden={
            type === 10 ||
            type === 11 ||
            type === 7 ||
            type === 1 ||
            type === 3 ||
            type === 5 ||
            type === 6 ||
            type === 4
          }
          label="内容"
          name="text"
        >
          <Input />
        </Form.Item>
        <Form.Item
          hidden={
            type === 10 ||
            type === 11 ||
            type === 7 ||
            type === 1 ||
            type === 3 ||
            type === 5 ||
            type === 6 ||
            type === 4
          }
          label="字号(px)"
          name="size"
        >
          <Input />
        </Form.Item>
        <Form.Item
          hidden={
            type === 10 ||
            type === 11 ||
            type === 7 ||
            type === 1 ||
            type === 3 ||
            type === 5 ||
            type === 6 ||
            type === 4
          }
          label="字体"
          name="font"
        >
          <Select>
            <Select.Option value="serif">serif</Select.Option>
            <Select.Option value="Mircosoft Yahei">Mircosoft Yahei</Select.Option>
            <Select.Option value="Times New Roman">Times New Roman</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          hidden={
            type === 10 ||
            type === 11 ||
            type === 7 ||
            type === 1 ||
            type === 3 ||
            type === 5 ||
            type === 6 ||
            type === 4
          }
          name="fillText"
        >
          <Radio.Group>
            <Radio value="fillText"> 文字填充渲染 </Radio>
            <Radio value="strokeText"> 文字描边渲染 </Radio>
          </Radio.Group>
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
          <Button type="primary" htmlType="submit">
            确定
          </Button>
        </Form.Item>
      </Form>
    );
  };

  return (
    <div id="main" className="room-main">
      <div className="room-sider">
        <img className="room-sider-item" src={pointer} onClick={() => modifyBrush({ type: 0 })} />
        <Popover
          content={<Content type={10} />}
          onOpenChange={handleOpenPenChange}
          open={openPen}
          title={<b>画笔</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={brush} />
        </Popover>
        <img className="room-sider-item" src={rubber} onClick={() => modifyBrush({ type: 8 })} />
        <Popover
          content={<Content type={7} />}
          onOpenChange={handleOpenArrowChange}
          open={openArrow}
          title={<b>箭头</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={arrow} />
        </Popover>
        <Popover
          content={<Content type={1} />}
          onOpenChange={handleOpenLineChange}
          open={openLine}
          title={<b>直线</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={line} />
        </Popover>
        <Popover
          content={<Content type={9} />}
          onOpenChange={handleOpenWordChange}
          open={openWord}
          title={<b>文本</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={word} />
        </Popover>
        <Popover
          content={<Content type={3} />}
          onOpenChange={handleOpenRectChange}
          open={openRect}
          title={<b>矩形</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={rectangle} />
        </Popover>
        <Popover
          content={<Content type={5} />}
          onOpenChange={handleOpenTriChange}
          open={openTri}
          title={<b>三角形</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={triangle} />
        </Popover>
        <Popover
          content={<Content type={6} />}
          onOpenChange={handleOpenElChange}
          open={openEl}
          title={<b>椭圆形</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={ellipse} />
        </Popover>
        <Popover
          content={<Content type={4} />}
          onOpenChange={handleOpenDiamondChange}
          open={openDiamond}
          title={<b>菱形</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={diamond} />
        </Popover>
        <img
          className="room-sider-item"
          src={revokeImg}
          onClick={() => {
            if (canvasClass.current) {
              canvasClass.current.revoke();
              console.log('撤销');
            }
          }}
        />
        <img
          className="room-sider-item"
          src={redoImg}
          onClick={() => {
            if (canvasClass.current) {
              canvasClass.current.redo();
              console.log('重做');
            }
          }}
        />
        {/* <Popover
          content={<Content type={11} />}
          onOpenChange={handleOpenNoteChange}
          open={openNote}
          title={<b>笔记</b>}
          trigger="click"
          placement="rightTop"
        >
          <img className="room-sider-item" src={note} />
        </Popover>
        <img
          className="room-sider-item"
          src={noteClear}
          onClick={() => {
            if (canvasClass.current != null) {
              canvasClass.current.clearGraffiti();
            }
          }}
        // */}
      </div>
      <canvas
        height={whiteboardHeight}
        width={(whiteboardHeight * 1920) / 1080}
        // width={1036}
        // height={583}
        id="whiteboard"
        className="room-content"
        ref={canvasRef}
      />
    </div>
  );
};
