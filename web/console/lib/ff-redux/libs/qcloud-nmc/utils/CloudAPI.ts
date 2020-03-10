import { OperationResult } from '../../../';
import { extend } from '../../qcloud-lib';

interface Window {
  Promise?: PromiseConstructorLike;
}

const DEV_ENV = process.env.NODE_ENV === 'development';

const tips = seajs.require('tips');
const login = seajs.require('login');

const errorLog = (apiName: string, request: any, response: any) => `
## ${apiName}
Request: ${JSON.stringify(request, null, 4) || '(empty)'}
Response: ${JSON.stringify(response, null, 4)}
`;

let net = seajs.require('net');

export class CloudAPIError extends Error {
  /**
   * 用户需要登录
   */
  public static CODE_NEED_LOGIN = 7;

  apiName: string;
  errorCode: number;
  request: any;
  response: any;
  constructor(apiName: string, errorCode: number, extraMessage: string, request: any, response: any) {
    super(extraMessage || '未知错误');
    this.apiName = apiName;
    this.errorCode = errorCode;
    this.request = request;
    this.response = response;
  }
}
/* eslint-disable */
export class CloudAPI {
  /**
   * 创建一个云 API 访问实例
   * @param {string} serviceType 请求的云 API 命名空间，如 `cvm`，默认的云 API 地址为 `${serviceType}.api.qcloud.com`
   * @param {string} endpoint 自定义请求的 API 地址，如 "http://example.api.tencentyun.com"，注意，该地址只在开发模式下可用
   */
  constructor(public serviceType: string, public endpoint?: string, public options?: any) {}

  /**
   * 请求云 API
   * @param  {string} apiName  API 名称
   * @param  {any}    params   API 参数
   * @param  {number} regionId 区域 ID，默认为广州（1）
   */
  public async request<TResult>(apiName: string, params: any = {}, regionId = 1) {
    return new Promise<TResult>((resolve, reject) => {
      let url = '/cgi/capi?action=delegate';
      if (this.options && this.options.secure) {
        url += '&secure=1';
      }
      if (this.options && this.options.version) {
        url += '&version=' + this.options.version;
      }
      const config = { method: 'POST', url: url };
      let options: any = {
        serviceType: this.serviceType,
        cmd: apiName,
        region: regionId,
        regionId,
        data: params
      };

      if (DEV_ENV && this.endpoint) {
        options.endpoint = this.endpoint;
      }

      const callback = response => {
        let error: CloudAPIError = null;

        // 非法返回
        if (response.code) {
          // 未收到云 API 数据
          if (!response.data) {
            const request = extend({}, config, options);
            error = new CloudAPIError(apiName, response.code, '数据未返回', request, response);
          }
          // 云 API 返回错误码
          else if (response.data.code) {
            const request = extend({}, config, options);
            error = new CloudAPIError(apiName, response.data.code, response.data.message, request, response);
          } else {
            const request = extend({}, config, options);
            error = new CloudAPIError(apiName, response.code, response.data, request, response);
          }
        }

        if (error) {
          // 没有登录，进行登录检查并弹出登录对话框
          if (error.errorCode === CloudAPIError.CODE_NEED_LOGIN) {
            login.checkLogin();
            tips.error('登录会话超时，请重新登陆');
            reject(error);
            return;
          }

          console.warn(errorLog(apiName, params, response.data));
          if (!DEV_ENV) {
            error.message = '服务器开了个小差，请稍后重试';
          } else {
            console.dir(error);
          }
          reject(error);
        } else {
          resolve(response.data as TResult);
        }
      };
      net.send(config, {
        data: options,
        cb: callback,
        global: false
      });
    });
  }

  // tips: 数组的重载写在前面会优先匹配
  /**
   * 请求执行某个 API 操作，操作的目标是指定类型的集合。
   * */
  public async operation<TTarget>(
    apiName: string,
    params: any,
    targets: TTarget[],
    regiondId?: number
  ): Promise<OperationResult<TTarget>[]>;

  /**
   * 请求执行某个 API 操作，操作的目标是指定的类型
   * */
  public async operation<TTarget>(
    apiName: string,
    params: any,
    target: TTarget,
    regiondId?: number
  ): Promise<OperationResult<TTarget>>;

  public async operation<TTarget>(
    apiName: string,
    params: any,
    target: TTarget | TTarget[],
    regionId?: number
  ): Promise<OperationResult<TTarget> | OperationResult<TTarget>[]> {
    function makeResult(target: TTarget, error?: any): OperationResult<TTarget> {
      return error ? { success: false, target, error } : { success: true, target };
    }

    return this.request(apiName, params, regionId).then(
      () => {
        if (target instanceof Array) {
          return target.map(x => makeResult(x));
        } else {
          return makeResult(target as TTarget);
        }
      },
      error => {
        if (target instanceof Array) {
          return target.map(x => makeResult(x, error));
        } else {
          return makeResult(target as TTarget, error);
        }
      }
    );
  }
}
