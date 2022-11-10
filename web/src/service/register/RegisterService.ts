import { fetchRes } from '@Src/utils/NetworkUtils';
import { IRegister, IRegisterResponse } from './IRegisterService';

async function register(params: IRegister): Promise<IRegisterResponse> {
  return fetchRes('post', '/signup', params);
}
export { register };
