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
import { runInAction } from 'mobx';
import QNRTC, { QNMicrophoneAudioTrack, QNRemoteAudioTrack, QNRemoteTrack, QNRemoteVideoTrack } from 'qnweb-rtc';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Stream } from './stream';
import './style/Room.less';
import WhiteBoard from './WhiteBoard';

const Room = () => {
  const location = useLocation();
  const roomInfo = location.state as IRoom;
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');
  // 创建QNRTCClient对象
  const client = QNRTC.createClient();
  const streams: Stream[] = [];
  const localStream: Stream = new Stream();

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

  client.on('user-left', (remoteUserID: string) => {
    runInAction(() => {
      const index = streams.findIndex((item) => item.user_id === remoteUserID);
      streams.splice(index, 1);
    });
  });

  // 订阅远端音视频
  client.on('user-published', async (userID: string, qntrack: (QNRemoteAudioTrack | QNRemoteVideoTrack)[]) => {
    runInAction(async () => {
      const { audioTracks } = await client.subscribe(qntrack);

      audioTracks.forEach((track) => {
        let stream = streams.find((item) => item.user_id === userID && item.tag === track.tag);
        if (stream === undefined) {
          stream = new Stream();
          stream.user_id = userID;
          stream.tag = track.tag || 'mc';
          streams.push(stream);
        }
        stream.audioTrack = track;
      });

      muteStateChanged([...audioTracks]);
    });
  });

  // ----------------------------------------------------------------

  // 远端音视频取消发布
  client.on('user-unpublished', async (userID: string, qntrack: (QNRemoteAudioTrack | QNRemoteVideoTrack)[]) => {
    runInAction(async () => {
      qntrack.forEach((track) => {
        const index = streams.findIndex((item) => item.user_id === userID && item.tag === track.tag);
        if (index >= 0) {
          const stream = streams[index];
          if (track.isAudio()) {
            stream.audioTrack = undefined;
          }

          if (stream.audioTrack === undefined) {
            streams.splice(index, 1);
          }
        }
      });
    });
  });

  async function joinChat() {
    const response = await getChatRoomToken(roomInfo.uuid);
    joinChatRoom(response.token);
  }

  async function joinChatRoom(roomToken: string) {
    // // 需要先监听对应事件再加入房间
    // autoSubscribe(client);
    // // 这里替换成刚刚生成的 RoomToken
    // await client.join(roomToken);
    // console.log('joinRoom success!');
    // await publish(client);

    await client.join(roomToken);

    const audioConfig = { tag: 'mc' };
    const localTracks = await QNRTC.createMicrophoneAudioTrack(audioConfig);
    console.log('my local tracks', localTracks);

    localStream.user_id = client.userID;
    localStream.isLocal = true;
    localStream.tag = 'mc';
    localStream.audioTrack = localTracks;

    await client.publish(localTracks);
    console.log('publish success! client.userID: ', client.userID);

    streams.push(localStream);

    // ----------------------------------------------------------------

    subscribeRemoteUser();
  }
  // 远程用户更新静音状态
  const muteStateChanged = (tracks: QNRemoteTrack[]) => {
    tracks.forEach((track) => {
      (function (track, streams) {
        track.on('mute-state-changed', (isMuted: boolean) => {
          runInAction(async () => {
            const stream = streams.find((item: Stream) => item.user_id === track.userID && item.tag === track.tag);
            if (stream === undefined) {
              return;
            }
            if (track.isAudio()) stream.audioMuted = isMuted;
          });
        });
      })(track, streams);
    });
  };

  const quitRoom = async () => {
    console.log(roomInfo);
    await client.leave().then(() => {
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

  const subscribeRemoteUser = () => {
    client.remoteUsers.forEach(async (user) => {
      const { audioTracks } = await client.subscribe([...user.getAudioTracks()]);

      const mcStream = new Stream();
      const screenStream = new Stream();

      audioTracks.forEach((track) => {
        if (track.tag === 'mc') mcStream.audioTrack = track;
        if (track.tag === 'screen') screenStream.audioTrack = track;
      });

      if (mcStream.audioTrack !== undefined) {
        mcStream.user_id = user.userID;
        streams.push(mcStream);
      }

      if (screenStream.audioTrack !== undefined) {
        screenStream.user_id = user.userID;
        streams.push(screenStream);
      }

      muteStateChanged([...audioTracks]);
    });
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
