import { IBaseResponse } from './../IBaseService';
export interface ICreate {
  name: string;
  type: string;
}
export interface ICreateResponse extends IBaseResponse {
  room: IRoom;
}

export interface IRoom {
  created_time?: string;
  host_name: string;
  host_uuid?: string;
  name: string;
  participants?: IParticipant[];
  type: string;
  updated_time?: string;
  uuid: string;
}

export interface IParticipant {
  name?: string;
  permission?: string;
  user_uuid?: string;
}

export interface IRoomListResponse extends IBaseResponse {
  roomlist: IRoom[];
}
