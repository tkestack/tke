/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
import { Method, reduceNetworkRequest } from '.';
import { RequestParams } from '../src/modules/common/models';

const tips = seajs.require('tips');

export class RequestArgs {
  url: string;
  method?: string;
  bodyData?: any;
  clusterId?: string;
  projectId?: string;
  keyword?: string;
  headers?: Record<string, string>;
}

export class RequestResult {
  data: any;
  error: any;
}

export const SEND = async (args: RequestArgs) => {
  // 构建参数
  const params: RequestParams = {
    method: args.method,
    url: args.url,
    data: args.bodyData,
    userDefinedHeader: args?.headers ?? {}
  };
  const resp = new RequestResult();
  try {
    const response = await reduceNetworkRequest(params, args.clusterId, args.projectId, args.keyword);
    if (response.code !== 0) {
      tips.error(response.message, 2000);
      resp.data = args.bodyData;
      resp.error = response?.message;
    } else {
      if (args.method !== Method.get) {
        tips.success('操作成功', 2000);
      }
      resp.data = response.data;
      resp.error = null;
    }
    return resp;
  } catch (error) {
    tips.error(error?.response?.data?.ErrStatus?.message ?? error.response.data.message, 2000);
    resp.data = args.bodyData;
    resp.error = error.response.data.message;
    return resp;
  }
};

export const GET = async (args: RequestArgs) => {
  args.method = Method.get;
  args.bodyData = null;
  const response = await SEND(args);
  return response;
};
export const DELETE = async (args: RequestArgs) => {
  args.method = Method.delete;
  args.bodyData = null;
  const response = await SEND(args);
  return response;
};
export const POST = async (args: RequestArgs) => {
  args.method = Method.post;
  args.bodyData = JSON.stringify(args.bodyData);
  const response = await SEND(args);
  return response;
};
export const PUT = async (args: RequestArgs) => {
  args.method = Method.put;
  args.bodyData = JSON.stringify(args.bodyData);
  const response = await SEND(args);
  return response;
};
export const PATCH = async (args: RequestArgs) => {
  args.method = Method.patch;
  args.bodyData = JSON.stringify(args.bodyData);
  const response = await SEND(args);
  return response;
};
