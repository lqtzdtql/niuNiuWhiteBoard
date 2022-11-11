import triangle from '@Public/static/img/三角形.png';
import word from '@Public/static/img/文字.png';
import ellipse from '@Public/static/img/椭圆形.png';
import rubber from '@Public/static/img/橡皮擦.png';
import brush from '@Public/static/img/画笔.png';
import line from '@Public/static/img/直线.png';
import rectangle from '@Public/static/img/矩形.png';
import arrow from '@Public/static/img/箭头.png';
import { joinRoom } from '@Sdk/index';
import { Room as SDKRoom } from '@Sdk/Room';
import { IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { getWhiteBoardToken } from '@Src/service/room/RoomService';
import { Popover } from 'antd';
import React, { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import './style/Room.less';
import { Title } from './Title';

const Room = () => {
  const location = useLocation();
  const roomInfo = location.state as IRoom;
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');

  let room: SDKRoom;

  // const BoardList = () => {
  //   return <List size="small" dataSource={data} renderItem={(item: number) => <List.Item>页面{item}</List.Item>} />;
  // };

  useEffect(() => {
    // joinWhiteBoard(roomInfo.uuid);
  }, []);

  async function joinWhiteBoard(uuid: string) {
    const response = await getWhiteBoardToken(uuid);
    room = await joinRoom({
      token: response.token,
      onlyRead: false,
      el: document.getElementById('whiteboard') as HTMLCanvasElement,
    });
  }

  const Sider: React.FC = () => {
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

  return (
    <div className="room-box">
      <Title roomInfo={roomInfo} userInfo={userInfo} />
      <div className="room-main">
        <Sider />
        <canvas className="room-content" />
      </div>
    </div>
  );
};
export default Room;
