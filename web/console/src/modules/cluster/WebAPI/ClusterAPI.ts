import { CreateResource } from './../../common/models/CreateResource';
import { QueryState } from '@tencent/qcloud-redux-query';
import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { RequestParams, ResourceInfo, ClusterFilter, Cluster } from '../../common/models';
import {
  reduceNetworkRequest,
  Method,
  reduceK8sRestfulPath,
  operationResult,
  reduceNetworkWorkflow,
  requestMethodForAction
} from '../../../../helpers';
import { resourceConfig } from '../../../../config';
import { CreateIC } from '../models';
import { authTypeMapping } from '../constants/Config';

/**
 * 集群列表的查询
 * @param query 集群列表查询的一些过滤条件
 */
export async function fetchClusterList(query: QueryState<ClusterFilter>, regionId: number = 1) {
  let { search } = query;
  let clusters: Cluster[] = [];

  let clusterResourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = reduceK8sRestfulPath({ resourceInfo: clusterResourceInfo });

  if (search) {
    url += '/' + search;
  }

  let params: RequestParams = {
    method: Method.get,
    url
  };
  try {
    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      if (response.data.items) {
        clusters = response.data.items.map(item => {
          return Object.assign({}, item, { id: uuid() });
        });
      } else {
        // 这里是拉取某个具体的resource的时候，没有items属性
        clusters.push({
          id: uuid(),
          metadata: response.data.metadata,
          spec: response.data.spec,
          status: response.data.status
        });
      }
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  const result: RecordSet<Cluster> = {
    recordCount: clusters.length,
    records: clusters
  };

  return result;
}

export async function fetchPrometheuses() {
  let resourceInfo: ResourceInfo = resourceConfig().prometheus;
  let url = reduceK8sRestfulPath({
    resourceInfo
  });
  let params: RequestParams = {
    method: Method.get,
    url
  };
  let records = [];
  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      records = response.data.items.map(item => {
        return Object.assign({}, item, { id: uuid() });
      });
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (error.code !== 'ResourceNotFound') {
      throw error;
    }
  }
  return {
    records,
    recordCount: records.length
  };
}

/**
 * 创建独立集群
 * @param resource: CreateIC   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function createIC(clusters: CreateIC[]) {
  try {
    let {
      name,
      k8sVersion,
      cidr,
      computerList,
      networkDevice,
      maxClusterServiceNum,
      maxNodePodNum,
      vipAddress,
      vipPort,
      vip,
      gpu,
      gpuType
    } = clusters[0];

    let resourceInfo = resourceConfig()['cluster'];
    let url = reduceK8sRestfulPath({ resourceInfo });
    // 获取具体的请求方法，create为POST，modify为PUT
    let method = 'POST';

    let machines = [];

    computerList.forEach(computer => {
      computer.ipList.split(';').forEach(ip => {
        let labels = {};
        computer.labels.forEach(kv => {
          labels[kv.key] = kv.value;
        });
        if (computer.isGpu) {
          labels['nvidia-device-enable'] = 'enable';
        }
        machines.push({
          ip: ip,
          port: computer.ssh ? +computer.ssh : 22,
          username: computer.username ? computer.username : undefined,
          password:
            computer.authType === authTypeMapping.password && computer.password
              ? window.btoa(computer.password)
              : undefined,
          // role: computer.role,
          privateKey:
            computer.authType === authTypeMapping.cert && computer.privateKey
              ? window.btoa(computer.privateKey)
              : undefined,
          passPhrase:
            computer.authType === authTypeMapping.cert && computer.passPhrase
              ? window.btoa(computer.passPhrase)
              : undefined,
          labels: labels
        });
      });
    });

    let jsonData = {
      apiVersion: `${resourceInfo.group}/${resourceInfo.version}`,
      kind: resourceInfo.headTitle,
      metadata: {
        generateName: 'cls'
      },
      spec: {
        displayName: name,
        clusterCIDR: cidr,
        networkDevice: networkDevice,
        features: gpu
          ? {
              gpuType: gpuType
            }
          : undefined,
        properties: {
          maxClusterServiceNum: maxClusterServiceNum,
          maxNodePodNum: maxNodePodNum
        },
        type: 'Baremetal',
        version: k8sVersion,
        machines: machines
      },
      status: vip
        ? {
            addresses: [
              {
                host: vipAddress,
                type: 'Advertise',
                port: vipPort ? +vipPort : 6443
              }
            ]
          }
        : undefined
    };

    jsonData = JSON.parse(JSON.stringify(jsonData));
    // 构建参数
    let params: RequestParams = {
      method,
      url,
      data: JSON.stringify(jsonData)
    };
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(clusters);
    } else {
      return operationResult(clusters, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(clusters, reduceNetworkWorkflow(error));
  }
}

/**
 * 创建独立集群
 * @param resource: CreateIC   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyClusterName(clusters: CreateResource[]) {
  try {
    let { jsonData, resourceInfo, clusterId } = clusters[0];

    let url = reduceK8sRestfulPath({ resourceInfo, specificName: clusterId, clusterId: clusterId });
    // 构建参数
    let params: RequestParams = {
      method: Method.patch,
      url,
      userDefinedHeader: {
        'Content-Type': 'application/strategic-merge-patch+json'
      },
      data: jsonData
    };

    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(clusters);
    } else {
      return operationResult(clusters, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(clusters, reduceNetworkWorkflow(error));
  }
}

/**
 * 创建独立集群
 * @param resource: CreateIC   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function fetchCreateICK8sVersion() {
  return [{ text: '1.14.6', value: '1.14.6' }];
}

/**
 * 创建导入集群
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function createImportClsutter(resource: CreateResource[], regionId: number) {
  try {
    let { mode, resourceIns, clusterId, resourceInfo, namespace, jsonData } = resource[0];

    let clustercredentialResourceInfo = resourceConfig().clustercredential;
    let clusterUrl = reduceK8sRestfulPath({ resourceInfo, clusterId }),
      clustercredentialUrl = reduceK8sRestfulPath({ resourceInfo: clustercredentialResourceInfo, clusterId });
    let method = requestMethodForAction(mode);

    let clusterData = JSON.parse(jsonData);

    let clustercredentialData = {
      clusterName: '',
      metadata: {
        generateName: 'clustercredential'
      },
      caCert: clusterData.status.credential.caCert,
      token: clusterData.status.credential.token ? clusterData.status.credential.token : undefined
    };
    clusterData.status.credential = undefined;
    clusterData = JSON.stringify(clusterData);
    // 构建参数
    let params: RequestParams = {
      method,
      url: clusterUrl,
      data: clusterData
    };

    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      let clustercredentialParams: RequestParams = {
        method,
        url: clustercredentialUrl,
        data: clustercredentialData
      };
      clustercredentialData.clusterName = response.data.metadata.name;
      clustercredentialData = JSON.parse(JSON.stringify(clustercredentialData));
      let clustercredentialResponce = await reduceNetworkRequest(clustercredentialParams, clusterId);
      if (clustercredentialResponce.code === 0) {
        return operationResult(resource);
      } else {
        return operationResult(resource, reduceNetworkWorkflow(clustercredentialResponce));
      }
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}
