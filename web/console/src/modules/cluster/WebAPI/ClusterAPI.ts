import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { resourceConfig } from '../../../../config';
import {
  Method,
  operationResult,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  reduceNetworkWorkflow,
  requestMethodForAction
} from '../../../../helpers';
import { Cluster, ClusterFilter, RequestParams, ResourceInfo } from '../../common/models';
import { CreateResource } from '../../common/models/CreateResource';
import { authTypeMapping, CreateICVipType } from '../constants/Config';
import { CreateIC } from '../models';
import { deleteResourceIns } from './K8sResourceAPI';
import {getK8sValidVersions} from '@src/webApi/cluster'

/**
 * 集群列表的查询
 * @param query 集群列表查询的一些过滤条件
 */
export async function fetchClusterList(query: QueryState<ClusterFilter>, regionId = 1) {
  const { search } = query;
  let clusters: Cluster[] = [];

  const clusterResourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = reduceK8sRestfulPath({ resourceInfo: clusterResourceInfo });

  if (search) {
    url += '/' + search;
  }

  const params: RequestParams = {
    method: Method.get,
    url
  };
  try {
    const response = await reduceNetworkRequest(params);

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
  const resourceInfo: ResourceInfo = resourceConfig().prometheus;
  const url = reduceK8sRestfulPath({
    resourceInfo
  });
  const params: RequestParams = {
    method: Method.get,
    url
  };

  // 兼容新的monitor版本peomethus
  const monitorParams: RequestParams = {
    method: Method.get,
    url: '/apis/monitor.tkestack.io/v1/prometheuses'
  };

  let records = [];
  try {
    const [response, monitorResponse] = await Promise.all([
      reduceNetworkRequest(params),
      reduceNetworkRequest(monitorParams)
    ]);
    if (response.code === 0) {
      records = response.data.items.map(item => {
        return Object.assign({}, item, { id: uuid() });
      });
    }

    if (monitorResponse.code === 0) {
      records = records.concat(
        monitorResponse.data.items.map(item => {
          return Object.assign({}, item, { id: uuid() });
        })
      );
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
    const {
      name,
      k8sVersion,
      cidr,
      computerList,
      networkDevice,
      maxClusterServiceNum,
      maxNodePodNum,
      vipAddress,
      vipPort,
      vipType,
      gpu,
      gpuType,
      merticsServer
    } = clusters[0];

    const resourceInfo = resourceConfig()['cluster'];
    const url = reduceK8sRestfulPath({ resourceInfo });
    // 获取具体的请求方法，create为POST，modify为PUT
    const method = 'POST';

    const machines = [];

    computerList.forEach(computer => {
      computer.ipList.split(';').forEach(ip => {
        const labels = {};
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
        features: {
          gpuType: gpu ? gpuType : undefined,
          ha:
            vipType !== CreateICVipType.unuse
              ? {
                  tke:
                    vipType === CreateICVipType.tke
                      ? {
                          vip: vipAddress
                        }
                      : undefined,
                  thirdParty:
                    vipType === CreateICVipType.existed
                      ? {
                          vip: vipAddress,
                          vport: +vipPort
                        }
                      : undefined
                }
              : undefined,

          enableMetricsServer: merticsServer
        },
        properties: {
          maxClusterServiceNum: maxClusterServiceNum,
          maxNodePodNum: maxNodePodNum
        },
        type: 'Baremetal',
        version: k8sVersion,
        machines: machines
      }
    };

    jsonData = JSON.parse(JSON.stringify(jsonData));
    // 构建参数
    const params: RequestParams = {
      method,
      url,
      data: JSON.stringify(jsonData)
    };
    const response = await reduceNetworkRequest(params);
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
    const { jsonData, resourceInfo, clusterId } = clusters[0];

    const url = reduceK8sRestfulPath({ resourceInfo, specificName: clusterId, clusterId: clusterId });
    // 构建参数
    const params: RequestParams = {
      method: Method.patch,
      url,
      userDefinedHeader: {
        'Content-Type': 'application/strategic-merge-patch+json'
      },
      data: jsonData
    };

    const response = await reduceNetworkRequest(params);
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
  const list = await getK8sValidVersions()
  return list
    .map(_ => ({text: _, value: _}))
  
}

/**
 * 创建导入集群
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function createImportClsutter(resource: CreateResource[], regionId: number) {
  try {
    const { mode, resourceIns, clusterId, resourceInfo, namespace, jsonData } = resource[0];

    const clustercredentialResourceInfo = resourceConfig().clustercredential;
    const clusterUrl = reduceK8sRestfulPath({ resourceInfo, clusterId }),
      clustercredentialUrl = reduceK8sRestfulPath({ resourceInfo: clustercredentialResourceInfo, clusterId });
    const method = requestMethodForAction(mode);

    let clusterData = JSON.parse(jsonData);

    const clustercredentialData = {
      metadata: {
        generateName: 'cc'
      },
      caCert: clusterData.status.credential.caCert,
      token: clusterData.status.credential.token ? clusterData.status.credential.token : undefined,
      clientKey: clusterData.status.credential.clientKey || undefined,
      clientCert: clusterData.status.credential.clientCert || undefined
    };
    // 构建参数
    const clustercredentialParams: RequestParams = {
      method,
      url: clustercredentialUrl,
      data: clustercredentialData
    };

    const clustercredentialResponce = await reduceNetworkRequest(clustercredentialParams, clusterId);
    if (clustercredentialResponce.code === 0) {
      const clusterParams: RequestParams = {
        method,
        url: clusterUrl,
        data: clusterData
      };
      clusterData.spec.clusterCredentialRef = { name: clustercredentialResponce.data.metadata.name };
      clusterData.status.credential = undefined;
      clusterData = JSON.parse(JSON.stringify(clusterData));
      try {
        const clusterResponce = await reduceNetworkRequest(clusterParams, clusterId);
        if (clusterResponce.code === 0) {
          return operationResult(resource);
        } else {
          return operationResult(resource, reduceNetworkWorkflow(clusterResponce));
        }
      } catch (error) {
        deleteResourceIns(
          [
            {
              id: uuid(),
              resourceIns: clustercredentialResponce.data.metadata.name,
              resourceInfo: clustercredentialResourceInfo
            }
          ],
          1
        );
        throw error;
      }
    } else {
      return operationResult(resource, reduceNetworkWorkflow(clustercredentialResponce));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}
