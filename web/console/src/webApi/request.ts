import { changeForbiddentConfig } from '@/index';
import Axios from 'axios';
import { v4 as uuidv4 } from 'uuid';

const instance = Axios.create({
  baseURL: '/apis',
  timeout: 3000
});

instance.interceptors.request.use(
  config => {
    config.headers['X-Remote-Extra-RequestID'] = uuidv4();
    return config;
  },
  error => {
    console.log('request error:', error);
    return Promise.reject(error);
  }
);

instance.interceptors.response.use(
  ({ data }) => data,
  error => {
    if (!error.response) {
      error.response = {
        data: {
          message: `系统内部服务错误（${error?.config?.heraders?.['X-Remote-Extra-RequestID'] || ''}）`
        }
      };
    }

    if (error.response.status === 401) {
      location.reload();
    }

    if (error.response.status === 403) {
      changeForbiddentConfig({
        isShow: true,
        message: error.response.data.message
      });
    }

    return Promise.reject(error);
  }
);

export default instance;
