import { IBaseResponse } from './../IBaseService';
export interface ILogin {
  mobile: string;
  passwd: string;
}

export interface ILoginResponse extends IBaseResponse {
  user_info: IUserInfo;
}

export interface IUserInfo {
  id: number;
  uuid: string;
  name: string;
  mobile: string;
}

export interface IGetUserInfoResponse extends IBaseResponse {}
