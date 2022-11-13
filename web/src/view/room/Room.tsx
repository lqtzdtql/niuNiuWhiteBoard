import { Canvas } from '@Sdk/Canvas';
import { joinRoom } from '@Sdk/index';
import { Room as SDKRoom } from '@Sdk/Room';
import { IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { getUserInfo } from '@Src/service/login/LoginService';
import { getWhiteBoardToken } from '@Src/service/room/RoomService';
import React, { useEffect, useRef, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { Context } from './context';
import { Main } from './Main';
import './style/Room.less';
import { Title } from './Title';

const Room = () => {
  const location = useLocation();
  const roomInfo = location.state as IRoom;
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const canvasClass = useRef<Canvas>(null);
  const [boards, setBoards] = useState<string[]>([]);
  const [partic, setPartic] = useState<IUserInfo[]>([]);
  const roomRef = useRef<SDKRoom | null>(null);
  const [checked, setChecked] = useState(true);

  useEffect(() => {
    joinWhiteBoard(roomInfo.uuid);
    if (userInfo.name !== roomInfo.host_name) {
      setPartic([userInfo]);
    }
  }, []);

  async function joinWhiteBoard(uuid: string) {
    const response = await getWhiteBoardToken(uuid);

    roomRef.current = await joinRoom({
      token: response.token,
      onlyRead: false,
      el: canvasRef.current as HTMLCanvasElement,
    });

    if (roomRef.current) {
      roomRef.current.on('leaveRoom', async (options: { leaveUserId: string }) => {
        console.log('abc', 'leaveRoom');
        setPartic((pre) => pre.filter((item) => item.uuid !== options.leaveUserId));
      });
      roomRef.current.on('joinRoom', async (options: { joinUserId: string }) => {
        const response = await getUserInfo(options.joinUserId);
        setPartic((pre) => [...pre, response]);
      });
      roomRef.current.on('new', (options: { canvasId: string; userId: string }) => {
        setBoards((pre) => [...pre, options.canvasId]);
      });
      roomRef.current.on('modifyOnlyRead', () => {
        setChecked((pre) => !pre);
      });
    }
  }

  const context = {
    canvasRef,
    canvasClass,
    roomRef,
    boards,
    partic,
    checked,
    setChecked,
  };

  return (
    <Context.Provider value={context}>
      <div className="room-box">
        <Title roomInfo={roomInfo} userInfo={userInfo} />
        <Main roomInfo={roomInfo} userInfo={userInfo} />
      </div>
    </Context.Provider>
  );
};
export default Room;
