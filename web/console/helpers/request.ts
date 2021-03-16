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
    data: args.bodyData
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
