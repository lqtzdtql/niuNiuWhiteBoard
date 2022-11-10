import { IBaseResponse } from './../IBaseService';
import { IUserInfo } from './../login/ILoginService';
export interface IRegister {
  mobile: string;
  passwd: string;
  name: string;
}

export interface IRegisterResponse extends IBaseResponse {
  user_info: IUserInfo;
}
