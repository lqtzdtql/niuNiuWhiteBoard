import { fetchRes } from '@Src/utils/NetworkUtils';
import { IBaseResponse } from './../IBaseService';

// /v1/rooms/{uuid}/exit
async function exitRoom(uuid: string): Promise<IBaseResponse> {
  return fetchRes('get', `/v1/rooms/${uuid}/exit`, {});
}
export { exitRoom };
