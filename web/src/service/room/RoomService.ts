import { fetchRes } from '@Src/utils/NetworkUtils';
import { IBaseResponse } from './../IBaseService';

// /v1/rooms/{uuid}/exit
async function exitRoom(uuid: string): Promise<IBaseResponse> {
  return fetchRes('get', `/v1/rooms/${uuid}/exit`, {});
}

async function getChatRoomToken(uuid: string) {
  return fetchRes('get', `/v1/rooms/${uuid}/rtc`, {});
}

async function getWhiteBoardToken(uuid: string) {
  return fetchRes('get', `/v1/rooms/${uuid}/whiteboard`, {});
}

export { exitRoom, getChatRoomToken, getWhiteBoardToken };
