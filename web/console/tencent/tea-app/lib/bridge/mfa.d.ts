import { RequestBody, RequestOptions } from "./capi";
/**
 * MFA 二次验证业务接入指南
 * http://tapd.oa.com/10103951/markdown_wikis/view/#1010103951008387339
 */
export declare const mfa: {
    /**
   * 对指定的云 API 进行 MFA 验证
   *
   * @param api 要校验的云 API，需要包括业务和接口名两部分，如 "cvm:DestroyInstance"
   *
   * @returns 返回值为 boolean 的 Promise，为 true 则表示校验通过，可以调用云 API
   *
   * @example
  ```js
  // 发起 MFA 校验
  const mfaPassed = await app.mfa.verify('cvm:DestroyInstance');
  
  if (!mfaPassed) {
    // 校验取消，跳过后续业务
    return;
  }
  
  // 校验完成，调用云 API
  const result = await app.capi.request({
    serviceType: 'cvm',
    cmd: 'DestroyInstance',
    // ...
  });
  
  ```
  */
    verify(api: string): Promise<boolean>;
    /**
     * 跟 verify 类似，不过为指定的 ownerUin 校验
     */
    verifyForOwner(api: string, ownerUin: string): Promise<boolean>;
    /**
   * 校验 MFA 后调用云 API，使用该方法具备失败重新校验的能力
   *
   * @example
  ```js
  const result = await app.mfa.request({
    regionId: 1,
    serviceType: 'cvm',
    cmd: 'DestroyInstance',
    data: {
      instanceId: 'ins-a5d3ccw8c'
    }
  }, {
    onMFAError: error => {
      // 碰到 MFA 的错误请进行重试逻辑，业务可以自己限制重试次数
      return error.retry();
    }
  });
  ```
    *
    * > 注意：*如果已经使用了 `app.mfa.verify()` 方法进行 MFA 校验，则无需再使用该方法发起 API 请求，直接使用 `app.capi.request()` 模块发起即可*
    */
    request(body: RequestBody, options: MFARequestOptions): Promise<any>;
};
export interface MFARequestOptions extends RequestOptions {
    onMFAError(error: MFAError): Promise<any>;
}
export interface MFAError extends Error {
    retry(): Promise<any>;
}
