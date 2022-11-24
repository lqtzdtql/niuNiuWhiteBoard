import { Room } from './Room';

type paramsType = {
  token: string;
  onlyRead: boolean;
  el: HTMLCanvasElement;
  elOptions?: {};
};

export async function joinRoom(params: paramsType) {
  const res = await fetch(`http://81.68.68.216:8888/auth?token=${params.token}`);
  if (res.status >= 200 && res.status < 300) {
    const resData = await res.json();
    const { user_uuid, room_uuid, code, message, user_name } = resData;
    if (code !== 200) {
      return message;
    } else {
      const roomData = {
        roomId: room_uuid,
        userId: user_uuid,
        onlyRead: params.onlyRead,
        el: params.el,
        userName: user_name,
        token: params.token,
      };
      const room = new Room(roomData);
      return room;
    }
  } else {
    return res.statusText;
  }
}
