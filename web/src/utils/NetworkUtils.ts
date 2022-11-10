import { baseUrl } from '@Src/constants/Constants';

const getUrl = (suffix: string) => {
  return baseUrl + suffix;
};
export function isValidKey(key: string | number | symbol, object: object): key is keyof typeof object {
  return key in object;
}

const fetchRes = async (type: string, url: string, params: {}) => {
  url = getUrl(url);
  let headers = new Headers();
  const token = localStorage.getItem('token');
  if (token && token !== '') {
    headers.set('Access-Token', token);
  }
  let requestConfig: RequestInit = {
    method: type,
    headers,
  };

  if (type == 'get') {
    let dataStr = '';
    Object.keys(params).forEach((key) => {
      if (isValidKey(key, params)) {
        dataStr += key + '=' + params[key] + '&';
      }
    });
    if (dataStr !== '') {
      dataStr = dataStr.substr(0, dataStr.lastIndexOf('&'));
      url = url + '?' + dataStr;
    }
  }
  if (type == 'post') {
    Object.defineProperty(requestConfig, 'body', {
      value: JSON.stringify(params),
    });
    Object.defineProperty(requestConfig.headers, 'Content-Length', {
      value: JSON.stringify(params).length,
    });
  }
  try {
    const response = await fetch(url, requestConfig);
    // console.log(response);

    const responseJson = response.json();

    const token = response.headers.get('Refresh-Token');
    if (token && token !== '') {
      localStorage.setItem('token', response.headers.get('Refresh-Token') || '');
    }

    return responseJson;
  } catch (error) {
    throw error as Error;
  }
};

export { fetchRes };
