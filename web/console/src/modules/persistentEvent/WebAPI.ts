import { RequestParams, UserDefinedHeader } from '../common/models';
import {
  reduceNetworkRequest,
  reduceNetworkWorkflow,
  requestMethodForAction,
  reduceK8sRestfulPath,
  operationResult
} from '../../../helpers';
import { CreateResource } from './models';

/**
 * 设置集群持久化事件
 * @param resource: CreateResource 创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyPeConfig(resource: CreateResource[], regionId: number) {
  try {
    let { mode, resourceInfo, clusterId, jsonData, namespace, resourceIns } = resource[0];
    let url = reduceK8sRestfulPath({ resourceInfo, namespace });
    let userDefinedHeader: UserDefinedHeader = {};

    if (mode === 'update') {
      url += `/${resourceIns}`;
      userDefinedHeader = {
        'Content-Type': 'application/strategic-merge-patch+json'
      };
    }

    let method = requestMethodForAction(mode);

    // 构建参数
    let params: RequestParams = {
      method,
      url,
      userDefinedHeader,
      data: jsonData
    };

    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}
