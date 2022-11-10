import triangle from '@Public/static/img/三角形.png';
import download from '@Public/static/img/下载.png';
import share from '@Public/static/img/分享.png';
import friends from '@Public/static/img/好友.png';
import word from '@Public/static/img/文字.png';
import newPage from '@Public/static/img/新建页面.png';
import curve from '@Public/static/img/曲线.png';
import ellipse from '@Public/static/img/椭圆形.png';
import rubber from '@Public/static/img/橡皮擦.png';
import brush from '@Public/static/img/画笔.png';
import line from '@Public/static/img/直线.png';
import rectangle from '@Public/static/img/矩形.png';
import arrow from '@Public/static/img/箭头.png';
import page from '@Public/static/img/纸张.png';
import quit from '@Public/static/img/退出.png';
import { typeMap } from '@Src/constants/Constants';
import { IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { exitRoom } from '@Src/service/room/RoomService';
import { List, Popover, Switch } from 'antd';
import React, { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import './style/Room.less';
import WhiteBoard from './WhiteBoard';

const Room = () => {
  const location = useLocation();
  const roomInfo = location.state as IRoom;
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');

  console.log(roomInfo);

  const navigate = useNavigate();
  const [index, setIndex] = useState(1);
  const data = [1, 2, 3, 4];
  const BoardList = () => {
    return <List size="small" dataSource={data} renderItem={(item: number) => <List.Item>页面{item}</List.Item>} />;
  };

  const quitRoom = async () => {
    console.log(roomInfo);

    const response = await exitRoom(roomInfo.uuid);
    if (response.code === 200) {
      navigate('/home');
    }
  };

  const Title: React.FC = () => {
    return (
      <div className="title">
        <img src={quit} onClick={quitRoom} />
        <h3 className="room-name">
          {typeMap.get(roomInfo.type)}:{roomInfo.name}的页面{index}
        </h3>
        <Switch checkedChildren="协作模式" unCheckedChildren="只读模式" defaultChecked />

        <Popover content={<div>Content</div>} title={`主持人:${roomInfo.host_name}`} trigger="click" placement="bottom">
          <img className="room-title-item" src={friends} />
        </Popover>
        <img className="room-title-item" src={download} />
        <img className="room-title-item" src={share} />
        <img className="room-title-item" src={newPage} />
        <Popover content={<BoardList />} trigger="click" placement="bottomRight">
          <img className="room-title-item" src={page} />
        </Popover>
      </div>
    );
  };

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
        <Popover content={<div>Content</div>} title="曲线" trigger="click" placement="right">
          <img className="room-sider-item" src={curve} />
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
      <Title />
      <div className="room-main">
        <Sider />
        <div className="room-content">
          <WhiteBoard />
        </div>
      </div>
    </div>
  );
};
export default Room;
