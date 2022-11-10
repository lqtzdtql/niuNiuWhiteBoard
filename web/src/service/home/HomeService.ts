import { fetchRes } from '@Src/utils/NetworkUtils';
import { IBaseResponse } from './../IBaseService';
import { ICreate, ICreateResponse, IRoomListResponse } from './IHomeService';

async function create(params: ICreate): Promise<ICreateResponse> {
  return fetchRes('post', '/v1/rooms', params);
}
async function roomList(): Promise<IRoomListResponse> {
  return fetchRes('get', '/v1/roomlist', {});
}

async function logout(): Promise<IBaseResponse> {
  return fetchRes('get', '/v1/logout', {});
}

async function getChatRoomToken(uuid: string) {
  return fetchRes('get', `/v1/rooms/${uuid}/rtc`, {});
}

export { create, roomList, logout, getChatRoomToken };
