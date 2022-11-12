import { joinRoom } from '@Sdk/index';
import { Room as SDKRoom } from '@Sdk/Room';
import { IRoom } from '@Src/service/home/IHomeService';
import { IUserInfo } from '@Src/service/login/ILoginService';
import { getWhiteBoardToken } from '@Src/service/room/RoomService';
import React, { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { Main } from './Main';
import './style/Room.less';
import { Title } from './Title';

const Room = () => {
  const location = useLocation();
  const roomInfo = location.state as IRoom;
  const userInfo: IUserInfo = JSON.parse(localStorage.getItem('userInfo') || '');

  let room: SDKRoom;

  useEffect(() => {
    joinWhiteBoard(roomInfo.uuid);
  }, []);

  async function joinWhiteBoard(uuid: string) {
    const response = await getWhiteBoardToken(uuid);
    room = await joinRoom({
      token: response.token,
      onlyRead: false,
      el: document.getElementById('whiteboard') as HTMLCanvasElement,
    });
  }

  return (
    <div className="room-box">
      <Title roomInfo={roomInfo} userInfo={userInfo} />
      <Main />
    </div>
  );
};
export default Room;
