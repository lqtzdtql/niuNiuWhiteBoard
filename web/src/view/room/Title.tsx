import download from '@Public/static/img/下载.png';
import friends from '@Public/static/img/好友.png';
import upload from '@Public/static/img/导入.png';
import newPage from '@Public/static/img/新建页面.png';
import page from '@Public/static/img/纸张.png';
import quit from '@Public/static/img/退出.png';
import muteMicrophone from '@Public/static/img/静音.png';
import microphone from '@Public/static/img/麦克风.png';
import { typeMap } from '@Src/constants/Constants';
import { IParticipant, IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { exitRoom, getChatRoomToken } from '@Src/service/room/RoomService';
import { Button, Form, Input, List, Popover, Switch } from 'antd';
import { runInAction } from 'mobx';
import QNRTC, {
  QNConnectionDisconnectedInfo,
  QNConnectionDisconnectedReason,
  QNConnectionState,
  QNMicrophoneAudioTrack,
  QNRemoteAudioTrack,
  QNRemoteTrack,
  QNRTCClient,
} from 'qnweb-rtc';
import React, { useContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Context } from './context';
import { Stream } from './stream';

export interface TitleProps {
  roomInfo: IRoom;
  userInfo: IUserInfo;
}

let client: QNRTCClient = QNRTC.createClient();
let localTracks: QNMicrophoneAudioTrack;
let remoteTracks: QNRemoteAudioTrack[];
const streams: Stream[] = [];
const localStream: Stream = new Stream();

export const Title = (props: TitleProps) => {
  const roomInfo = props.roomInfo;
  const userInfo = props.userInfo;

  const navigate = useNavigate();
  const [listIndex, setListIndex] = useState(0);
  const [openMicrophone, setOpenMicrophone] = useState(true);
  const { roomRef } = useContext(Context);
  const { partic } = useContext(Context);
  const { boards } = useContext(Context);
  const { checked } = useContext(Context);
  const { setChecked } = useContext(Context);
  const { canvasClass } = useContext(Context);
  // const [openHear, setOpenHear] = useState(true);
  const [open, setOpen] = useState(false);
  const [listOpen, setListOpen] = useState(false);

  useEffect(() => {
    joinRTCRoom();
  }, []);

  const ParticipantsList = () => {
    return (
      <List
        size="small"
        dataSource={partic?.filter((item) => item.permission !== 'host')}
        renderItem={(Item: IParticipant) => <List.Item>{Item.name}</List.Item>}
      ></List>
    );
  };

  const BoardList = () => {
    return (
      <List
        size="small"
        dataSource={boards}
        renderItem={(item: string, index: number) => (
          <List.Item onClick={() => changeBoard(item, index)}>页面{index + 1}</List.Item>
        )}
      />
    );
  };

  const changeBoard = (canvasId: string, index: number) => {
    setListIndex(index + 1);
    canvasClass.current = roomRef.current.getCanvas(canvasId);
    roomRef.current.switchCanvas(canvasId);
    setListOpen(false);
  };

  const joinRTCRoom = async () => {
    const response = await getChatRoomToken(roomInfo.uuid);
    client.on('user-left', (remoteUserID: string) => {
      runInAction(() => {
        const index = streams.findIndex((item) => item.user_id === remoteUserID);
        streams.splice(index, 1);
      });
    });

    // 远端音频取消发布
    client.on('user-unpublished', async (userID: string, qntrack: QNRemoteAudioTrack[]) => {
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

    client.on('connection-state-changed', function (connectionState: string, info: QNConnectionDisconnectedInfo) {
      console.log('connection-state-changed', connectionState);
      // 当进入连接断开状态
      if (connectionState === QNConnectionState.DISCONNECTED) {
        // 监控断开原因
        switch (info.reason) {
          // 当异常断开时
          case QNConnectionDisconnectedReason.ERROR:
            break;
          // 当被踢出房间时
          case QNConnectionDisconnectedReason.KICKED_OUT:
            break;
          // 当调用接口，主动离开房间时
          case QNConnectionDisconnectedReason.LEAVE:
            break;
        }
      }
    });

    // 订阅远端音视频
    client.on('user-published', async (userID: string, qntrack: QNRemoteAudioTrack[]) => {
      runInAction(async () => {
        const { audioTracks } = await client.subscribe(qntrack);
        remoteTracks = audioTracks;
        remoteTracks.forEach((track) => {
          let stream = streams.find((item) => item.user_id === userID && item.tag === track.tag);
          if (stream === undefined) {
            stream = new Stream();
            stream.user_id = userID;
            stream.tag = track.tag || 'mc';
            streams.push(stream);
          }
          stream.audioTrack = track;
        });

        muteStateChanged([...remoteTracks], streams);
        const remoteElement = document.getElementById('body') as HTMLElement;
        // 遍历返回的远端 Track，调用 play 方法完成在页面上的播放
        for (const remoteTrack of [...remoteTracks]) {
          remoteTrack.play(remoteElement);
        }
      });
    });

    const muteStateChanged = (tracks: QNRemoteTrack[], streams: Stream[]) => {
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

        muteStateChanged([...audioTracks], streams);
      });
    };

    await client.join(response.token);

    const audioConfig = { tag: 'mc' };
    localTracks = await QNRTC.createMicrophoneAudioTrack(audioConfig);
    console.log('my local tracks', localTracks);

    localStream.user_id = client.userID;
    localStream.isLocal = true;
    localStream.tag = 'mc';
    localStream.audioTrack = localTracks;

    await client.publish(localTracks);

    streams.push(localStream);

    subscribeRemoteUser();
  };

  const quitRTCRoom = async () => {
    await client.leave();
    sessionStorage.removeItem('boards');
    const response = await exitRoom(roomInfo.uuid);

    roomRef.current.kickOutRoom(userInfo.uuid);
    roomRef.current.closeWs();
    if (response.code === 200) {
      navigate('/home', { replace: true });
    }
  };

  // const changeHear = async () => {
  //   setOpenHear(!openHear);
  // };

  const changeMicrophone = async () => {
    if (openMicrophone) {
      localStream.muteTrack('audio', true);
      setOpenMicrophone(false);
    } else {
      localStream.muteTrack('audio', false);
      setOpenMicrophone(true);
    }
  };

  function downloadTXT(filename: string, text: string, format: string, isImport: boolean) {
    const link = document.createElement('a');
    if (isImport) {
      link.setAttribute('href', `data:text/${format};charset=utf-8,${encodeURIComponent(text)}`);
    } else {
      // "\ ufeff" to solve the problem of CSV Chinese garbled code
      const blob = new Blob([`\ufeff${text}`], {
        type: `data:text/${format};charset=utf-8`,
      });
      const url = URL.createObjectURL(blob);
      link.href = url;
    }
    link.setAttribute('download', filename);
    link.style.display = 'none';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  const onFinish = (values: any) => {
    roomRef.current.importCanvas(canvasClass.current.canvasId, values.data);
    setOpen(false);
  };
  const handleOpen = (newOpen: boolean) => {
    setOpen(newOpen);
  };
  const handleListOpen = (newOpen: boolean) => {
    setListOpen(newOpen);
  };

  return (
    <div className="title">
      <img src={quit} onClick={quitRTCRoom} />
      <h3 className="room-name">
        {typeMap.get(roomInfo.type)}:{roomInfo.name} {listIndex > 0 ? `的页面${listIndex}` : ''}
      </h3>
      <Switch
        checked={checked}
        checkedChildren="协作模式"
        unCheckedChildren="只读模式"
        defaultChecked
        onChange={() => {
          roomRef.current.modifyOnlyRead();
        }}
        disabled={userInfo.name !== roomInfo.host_name}
      />
      <Popover content={<ParticipantsList />} title={`主持人:${roomInfo.host_name}`} trigger="click" placement="bottom">
        <img className="room-title-item" src={friends} />
      </Popover>
      {/* <img id="player" className="room-title-item" src={openHear ? hear : muteHear} onClick={changeHear} /> */}
      <img className="room-title-item" src={openMicrophone ? microphone : muteMicrophone} onClick={changeMicrophone} />
      <img
        className="room-title-item"
        src={download}
        onClick={() => {
          const canvasData = roomRef.current.exportCanvas(canvasClass.current.canvasId);
          downloadTXT(`${canvasClass.current.canvasId}.txt`, canvasData, 'plain', true);
        }}
      />

      <Popover
        content={
          <Form onFinish={onFinish}>
            <Form.Item name="data">
              <Input />
            </Form.Item>
            <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
              <Button type="primary" htmlType="submit">
                确定
              </Button>
            </Form.Item>
          </Form>
        }
        onOpenChange={handleOpen}
        open={open}
        trigger="click"
      >
        <img className="room-title-item" src={upload} />
      </Popover>
      {/* <img className="room-title-item" src={share} /> */}
      <img
        className="room-title-item"
        src={newPage}
        onClick={() => {
          roomRef.current.createCanvas();
        }}
      />
      <Popover
        open={listOpen}
        onOpenChange={handleListOpen}
        content={<BoardList />}
        trigger="click"
        placement="bottomRight"
      >
        <img className="room-title-item" src={page} />
      </Popover>
    </div>
  );
};
