import { Room } from './Room';

type paramsType = {
  token: string;
  onlyRead: boolean;
  el: HTMLCanvasElement;
  elOptions?: {};
};

export async function joinRoom(params: paramsType) {
  const res = await fetch(`/auth?token=${params.token}`);
  if (res.status >= 200 && res.status < 300) {
    const resData = await res.json();
    const { user_uuid, room_uuid, code, message } = resData;
    if (code !== 200) {
      return message;
    } else {
      const roomData = {
        roomId: room_uuid,
        userId: user_uuid,
        onlyRead: params.onlyRead,
        el: params.el,
        elOptions: params.elOptions,
      };
      const room = new Room(roomData);
      return room;
    }
  } else {
    return res.statusText;
  }
}
