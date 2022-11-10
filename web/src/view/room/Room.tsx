import triangle from '@Public/static/img/三角形.png';
import download from '@Public/static/img/下载.png';
import share from '@Public/static/img/分享.png';
import hear from '@Public/static/img/声音.png';
import muteHear from '@Public/static/img/声音关闭.png';
import friends from '@Public/static/img/好友.png';
import word from '@Public/static/img/文字.png';
import newPage from '@Public/static/img/新建页面.png';
import ellipse from '@Public/static/img/椭圆形.png';
import rubber from '@Public/static/img/橡皮擦.png';
import brush from '@Public/static/img/画笔.png';
import line from '@Public/static/img/直线.png';
import rectangle from '@Public/static/img/矩形.png';
import arrow from '@Public/static/img/箭头.png';
import page from '@Public/static/img/纸张.png';
import quit from '@Public/static/img/退出.png';
import muteMicrophone from '@Public/static/img/静音.png';
import microphone from '@Public/static/img/麦克风.png';
import { typeMap } from '@Src/constants/Constants';
import { getChatRoomToken } from '@Src/service/home/HomeService';
import { IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { exitRoom } from '@Src/service/room/RoomService';
import { List, Popover, Switch } from 'antd';
import QNRTC, { QNMicrophoneAudioTrack, QNRemoteTrack, QNRTCClient } from 'qnweb-rtc';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import './style/Room.less';
import WhiteBoard from './WhiteBoard';

const Room = () => {
  const location = useLocation();
  const roomInfo = location.state as IRoom;
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');
  // 创建QNRTCClient对象
  const client = QNRTC.createClient();
  let localTracks: QNMicrophoneAudioTrack[] = [];

  const [openMicrophone, setOpenMicrophone] = useState(true);
  const [openHear, setOpenHear] = useState(true);

  const navigate = useNavigate();
  const [index, setIndex] = useState(1);
  const data = [1, 2, 3, 4];
  const BoardList = () => {
    return <List size="small" dataSource={data} renderItem={(item: number) => <List.Item>页面{item}</List.Item>} />;
  };

  useEffect(() => {
    joinChat();
  }, []);

  async function joinChat() {
    const response = await getChatRoomToken(roomInfo.uuid);
    joinChatRoom(response.token);
  }

  async function joinChatRoom(roomToken: string) {
    // 需要先监听对应事件再加入房间
    autoSubscribe(client);
    // 这里替换成刚刚生成的 RoomToken
    await client.join(roomToken);
    console.log('joinRoom success!');
    await publish(client);
  }

  const quitRoom = async () => {
    console.log(roomInfo);
    client.leave().then(() => {
      for (const track of localTracks) {
        track.destroy();
      }
      localTracks = [];
    });
    const response = await exitRoom(roomInfo.uuid);
    if (response.code === 200) {
      navigate('/home', { replace: true });
    }
  };

  const changeHear = async () => {
    setOpenHear(!openHear);
  };

  const changeMicrophone = async () => {
    setOpenMicrophone(!openMicrophone);
  };

  // 这里的参数 client 是指刚刚初始化的 QNRTCClient 对象
  async function subscribe(client: QNRTCClient, tracks: QNRemoteTrack | QNRemoteTrack[]) {
    // 传入 Track 对象数组调用订阅方法发起订阅，异步返回成功订阅的 Track 对象。
    const remoteTracks = await client.subscribe(tracks);
    // 选择页面上的一个元素作为父元素，播放远端的音视频轨
    const remoteElement = document.getElementById('player');
    // 遍历返回的远端 Track，调用 play 方法完成在页面上的播放
    if (remoteElement && openHear) {
      for (const remoteTrack of [...remoteTracks.audioTracks]) {
        remoteTrack.play(remoteElement);
      }
    }
  }

  function autoSubscribe(client: QNRTCClient) {
    // 添加事件监听，当房间中出现新的 Track 时就会触发，参数是 trackInfo 列表
    client.on('user-published', (userId: any, tracks: QNRemoteTrack | QNRemoteTrack[]) => {
      console.log('user-published!', userId, tracks);
      subscribe(client, tracks)
        .then(() => console.log('subscribe success!'))
        .catch((e) => console.error('subscribe error', e));
    });
    // 就是这样，就像监听 DOM 事件一样通过 on 方法监听相应的事件并给出处理函数即可
  }
  async function publish(client: QNRTCClient) {
    // 同时采集麦克风音频和摄像头视频轨道。
    // 这个函数会返回一组audio track 与 video track
    const localTracks = await QNRTC.createMicrophoneAudioTrack();
    console.log('my local tracks', localTracks);
    // 将刚刚的 Track 列表发布到房间中
    if (openMicrophone) {
      await client.publish(localTracks);
      console.log('publish success!');
    }
  }

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
        <img id="player" className="room-title-item" src={openHear ? hear : muteHear} onClick={changeHear} />
        <img
          className="room-title-item"
          src={openMicrophone ? microphone : muteMicrophone}
          onClick={changeMicrophone}
        />
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
