import axios from 'axios';
import { OperationResult } from '@tencent/ff-redux';
import { RequestParams, ResourceInfo } from '../src/modules/common/models';
import { changeForbiddentConfig } from '../index';

/** 是否展示没有权限的弹窗 */
export let Init_Forbiddent_Config = {
  isShow: false,
  message: ''
};
/** 获取当前的uuid */

export const uuid = () => {
  let s = [];
  let hexDigits = '0123456789abcdef';
  for (let i = 0; i < 36; i++) {
    s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
  }
  s[14] = '4'; // bits 12-15 of the time_hi_and_version field to 0010
  s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1); // bits 6-7 of the clock_seq_hi_and_reserved to 01
  s[8] = s[13] = s[18] = s[23] = '-';

  let uuid = s.join('');
  return uuid;
};

/** 获取当前控制台modules 的域名匹配项 */
const GET_CONSOLE_MODULE_BASE_URL = location.origin || '';

interface CommonModuleProps {
  groupName: string;
  version: string;
}

export interface ConsoleModuleMapProps {
  [props: string]: CommonModuleProps;
}

/** 设置ConsoleAPIAddress的配置 */
export const setConsoleAPIAddress = (apiData: ConsoleModuleMapProps) => {
  window['modules'] = apiData;
};

/** RESTFUL风格的请求方法 */
export const Method = {
  get: 'GET',
  post: 'POST',
  patch: 'PATCH',
  delete: 'DELETE',
  put: 'PUT'
};

/**
 * 用于获取当前操作的正确请求方法
 */
export const requestMethodForAction = (type: string) => {
  const mapMethod = {
    create: Method.post,
    modify: Method.put,
    delete: Method.delete,
    list: Method.get,
    update: Method.patch
  };

  let method = mapMethod[type] ? mapMethod[type] : 'get';

  return method;
};

/**
 * 统一的请求处理
 * @param userParams: RequestParams
 */
export const reduceNetworkRequest = async (userParams: RequestParams, clusterId?: string, projectId?: string, keyword?: string) => {
  let {
    method,
    url,
    userDefinedHeader = {},
    data = {},
    apiParams,
    // baseURL = getConsoleAPIAddress(ConsoleModuleAddressEnum.PLATFORM)
    baseURL = GET_CONSOLE_MODULE_BASE_URL
  } = userParams;

  let rsp;
  // 请求tke-apiserver的 cluster的header
  if (clusterId) {
    userDefinedHeader = Object.assign({}, userDefinedHeader, {
      'X-TKE-ClusterName': clusterId
    });
  }
  if (projectId) {
    userDefinedHeader = Object.assign({}, userDefinedHeader, {
      'X-TKE-ProjectName': projectId
    });
  }
  if (keyword) {
    userDefinedHeader = Object.assign({}, userDefinedHeader, {
      'X-TKE-FuzzyResourceName': keyword
    });
  }

  let params = {
    method,
    baseURL,
    url,
    withCredentials: true,
    headers: Object.assign(
      {},
      {
        'X-Remote-Extra-RequestID': uuid(),
        'Content-Type': 'application/json'
      },
      userDefinedHeader
    )
  };

  if (data) {
    params = Object.assign({}, params, { data: data });
  }

  // 发送请求
  try {
    rsp = await axios(params as any);
  } catch (error) {
    // 如果返回是 401的话，自动登出，此时是鉴权不过，cookies失效了
    if (error.response && error.response.status === 401) {
      // location.reload();
    } else if (error.response && error.response.status === 403) {
      changeForbiddentConfig({
        isShow: true,
        message: error.response.data.message
      });
      throw error;
    } else if (error.response === undefined) {
      let uuid =
        error.config && error.config.headers && error.config.headers['X-Remote-Extra-RequestID']
          ? error.config.headers['X-Remote-Extra-RequestID']
          : '';
      error.response = {
        data: {
          message: `系统内部服务错误（${uuid}）`
        }
      };
      throw error;
    } else {
      throw error;
    }
  }

  // 处理回报请求
  let response = reduceNetworkResponse(rsp);
  return response;
};

/**
 * 处理返回的数据
 * @param type  判断当前控制台的类型
 * @param response  请求返回的数据
 */
const reduceNetworkResponse = (response: any = {}) => {
  let result;
  result = {
    code: response.status >= 200 && response.status < 300 ? 0 : response.status,
    data: response.data,
    message: response.statusText
  };

  return result;
};

/**
 * 处理workflow发生的错误
 * @param error workflow失败的错误信息
 */
export const reduceNetworkWorkflow = (error: any) => {
  return error.response.data && error.response.data.ErrStatus ? error.response.data.ErrStatus : error.response.data;
};

/**
 * 处理workflow的返回结果
 * @param target T[]
 * @param error any
 */
export const operationResult = function<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
};
