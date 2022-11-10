import { IUserInfo } from '@Src/service/login/ILoginService';
import { fetchRes } from '@Src/utils/NetworkUtils';
import { ILogin, ILoginResponse } from './ILoginService';

async function login(params: ILogin): Promise<ILoginResponse> {
  return fetchRes('post', '/login', params);
}

async function getUserInfo(uuid: string): Promise<IUserInfo> {
  return fetchRes('get', '/v1/userinfo/' + uuid, {});
}

export { login, getUserInfo };
