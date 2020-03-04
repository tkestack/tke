import { Resource, ResourceFilter } from '@src/modules/common';
import { OperationResult, QueryState } from '@tencent/ff-redux';
import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { t } from '@tencent/tea-app/lib/i18n';
import * as JsYAML from 'js-yaml';
import { resourceConfig } from '../../../config';
import { reduceK8sRestfulPath, reduceNetworkRequest } from '../../../helpers';
import { RequestParams, ResourceInfo } from '../common/models';
import { ClusterHelmStatus } from './constants/Config';
import {
  Helm,
  HelmFilter,
  HelmHistory,
  HelmKeyValue,
  InstallingHelm,
  TencenthubChart,
  TencenthubChartVersion,
  TencenthubNamespace
} from './models';

// 提示
const tips = seajs.require('tips');

// 获取cvm配置化平台的配置文件
const cvmConfig = seajs.require('config_cvm');

/** RESTFUL风格的请求方法 */
const Method = {
  get: 'GET',
  post: 'POST',
  patch: 'PATCH',
  delete: 'DELETE',
  put: 'PUT'
};

const SEND = async (
  url: string,
  method: string,
  bodyData: any,
  regionId: number,
  clusterId: string,
  tipErr: boolean = true
) => {
  // 构建参数
  let params: RequestParams = {
    method: method,
    url,
    data: bodyData
  };
  let response = await reduceNetworkRequest(params, clusterId);
  return response.data;
};

const GET = async (url: string, regionId: number, clusterId: string, tipErr: boolean = true) => {
  let response = await SEND(url, Method.get, null, regionId, clusterId, tipErr);
  return response;
};
const DELETE = async (url: string, regionId: number, clusterId: string, tipErr: boolean = true) => {
  let response = await SEND(url, Method.delete, null, regionId, clusterId, tipErr);
  return response;
};
const POST = async (url: string, bodyData: any, regionId: number, clusterId: string, tipErr: boolean = true) => {
  let response = await SEND(url, Method.post, JSON.stringify(bodyData), regionId, clusterId, tipErr);
  return response;
};

const PUT = async (url: string, bodyData: any, regionId: number, clusterId: string, tipErr: boolean = true) => {
  let response = await SEND(url, Method.put, JSON.stringify(bodyData), regionId, clusterId, tipErr);
  return response;
};

const PATCH = async (url: string, bodyData: any, regionId: number, clusterId: string, tipErr: boolean = true) => {
  let response = await SEND(url, Method.patch, JSON.stringify(bodyData), regionId, clusterId, tipErr);
  return response;
};

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
}

/**
 * 查询集群是否开通helm
 */
export async function checkClusterHelmStatus(regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['helm'];
  let url = `${reduceK8sRestfulPath({ resourceInfo })}?fieldSelector=spec.clusterName=${clusterId}`;
  let response = await GET(url, regionId, clusterId);

  let ret = response;
  if (
    ret.items &&
    ret.items.length &&
    ret.items[0].status.phase &&
    (ret.items[0].status.phase as string).toLowerCase() === 'running'
  ) {
    return { code: ClusterHelmStatus.RUNNING, reason: '' };
  } else if (
    ret.items &&
    ret.items.length &&
    ret.items[0].status.phase &&
    (ret.items[0].status.phase as string).toLowerCase() === 'checking'
  ) {
    return { code: ClusterHelmStatus.CHECKING, reason: '' };
  } else if (
    ret.items &&
    ret.items.length &&
    ret.items[0].status.phase &&
    (ret.items[0].status.phase as string).toLowerCase() === 'initializing'
  ) {
    return { code: ClusterHelmStatus.INIT, reason: '' };
  } else if (
    ret.items &&
    ret.items.length &&
    ret.items[0].status.phase &&
    (ret.items[0].status.phase as string).toLowerCase() === 'failed'
  ) {
    return {
      code: ClusterHelmStatus.ERROR,
      reason: ret.items[0].status.reason
    };
  } else if (
    ret.items &&
    ret.items.length &&
    ret.items[0].status.phase &&
    (ret.items[0].status.phase as string).toLowerCase() === 'reinitializing'
  ) {
    return {
      code: ClusterHelmStatus.REINIT,
      reason: ret.items[0].status.reason
    };
  } else {
    return { code: ClusterHelmStatus.NONE, reason: '' };
  }
}

/**
 * 开通集群helm
 */
export async function setupHelm(regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['helm'];
  let url = reduceK8sRestfulPath({ resourceInfo });
  await POST(
    url,
    {
      kind: resourceInfo.headTitle,
      apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
      metadata: {
        generateName: 'hm'
      },
      spec: {
        clusterName: clusterId
      }
    },
    regionId,
    clusterId
  );
}

const formatKeyValue = (kvs: HelmKeyValue[]) => {
  let values = {
    raw_original: '',
    values_type: 'kv'
  };
  let vs = [];
  kvs.forEach(item => {
    vs.push(item.key + '=' + item.value);
  });
  values.raw_original = vs.join(',');

  return values;
};

/**
 * 创建Helm应用
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function createHelm(
  params: {
    helmName: string;
    chart_url: string;
    resource: string;
    namespace: string;
    username?: string;
    password?: string;
    kvs?: HelmKeyValue[];
  },
  regionId: number = 1,
  clusterId: string
) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({ resourceInfo, specificName: clusterId })}/helm/tiller/v2/releases/${
    params.helmName
  }/json`;
  let data = {
    chart_url: params.chart_url,
    repo: params.resource,
    chart_ns: params.namespace,
    // token: params.token
    username: params.username || seajs.require('util').getUin() + '',
    password: params.password || ''
  };
  if (params.kvs && params.kvs.length) {
    data['values'] = formatKeyValue(params.kvs);
  }
  await POST(url, data, regionId, clusterId);
}
export async function createHelmByOther(
  params: {
    helmName: string;
    chart_url: string;
    resource: string;
    username?: string;
    password?: string;
    kvs?: HelmKeyValue[];
  },
  regionId: number = 1,
  clusterId: string
) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({ resourceInfo, specificName: clusterId })}/helm/tiller/v2/releases/${
    params.helmName
  }/json`;
  let data = {
    chart_url: params.chart_url,
    repo: params.resource
  };
  if (params.username) {
    data['username'] = params.username;
    data['password'] = params.password;
  }
  if (params.kvs && params.kvs.length) {
    data['values'] = formatKeyValue(params.kvs);
  }
  await POST(url, data, regionId, clusterId);
}

/**
 * 更新Helm应用
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function updateHelm(
  params: {
    helmName: string;
    chart_url: string;
    token: string;
    kvs?: HelmKeyValue[];
  },
  regionId: number = 1,
  clusterId: string
) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({ resourceInfo, specificName: clusterId })}/helm/tiller/v2/releases/${
    params.helmName
  }/json`;
  let data = {
    chart_url: params.chart_url,
    username: seajs.require('util').getUin() + '',
    password: params.token
  };
  if (params.kvs && params.kvs.length) {
    data['values'] = formatKeyValue(params.kvs);
  }
  await PUT(url, data, regionId, clusterId);
}

/**
 * 更新Helm应用
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function updateHelmByOther(
  params: {
    helmName: string;
    chart_url: string;
    username?: string;
    password?: string;
    kvs?: HelmKeyValue[];
  },
  regionId: number = 1,
  clusterId: string
) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({ resourceInfo, specificName: clusterId })}/helm/tiller/v2/releases/${
    params.helmName
  }/json`;
  let data = {
    chart_url: params.chart_url
  };
  if (params.username) {
    data['username'] = params.username;
    data['password'] = params.password;
  }
  if (params.kvs && params.kvs.length) {
    data['values'] = formatKeyValue(params.kvs);
  }
  await PUT(url, data, regionId, clusterId);
}

/**
 * 回滚到指定版本
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function rollbackVersion(
  params: {
    helmName: string;
    version: number;
  },
  regionId: number = 1,
  clusterId: string
) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({ resourceInfo, specificName: clusterId })}/helm/tiller/v2/releases/${
    params.helmName
  }/rollback/json?version=${params.version}`;
  await GET(url, regionId, clusterId);
}

/**
 * 获取Helm应用安装历史
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function fetchHistory(params: { helmName: string }, regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({ resourceInfo, specificName: clusterId })}/helm/tiller/v2/releases/${
    params.helmName
  }/json?max=10`;
  let response = await GET(url, regionId, clusterId);

  let history: HelmHistory[] = response.releases;

  const result: RecordSet<HelmHistory> = {
    recordCount: history.length,
    records: history
  };

  return result;
}

/**
 * 获取安装中的Helm应用列表
 * @param regionId
 * @param clusterId
 */
export async function fetchInstallingHelmList(regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/installing/json`;

  let response = await GET(url, regionId, clusterId);

  let helms: InstallingHelm[] = [];
  for (let key in response) {
    helms.push({
      id: key,
      name: key,
      status: +response[key]
    });
  }

  const result: RecordSet<InstallingHelm> = {
    recordCount: helms.length,
    records: helms
  };

  return result;
}

/**
 * 获取安装中的Helm应用详情
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function fetchInstallingHelm(params: { helmName: string }, regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/installing/${params.helmName}/content/json`;
  let response = await GET(url, regionId, clusterId);
  return response;
}

/**
 * 忽略安装中的Helm应用
 * @param params
 * @param regionId
 * @param clusterId
 */
export async function ignoreInstallingHelm(params: { helmName: string }, regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/installing/${params.helmName}/json`;
  let response = await DELETE(url, regionId, clusterId);
  return response.data;
}

/**
 * helm列表的查询
 * @param query helm列表查询的一些过滤条件
 */
export async function fetchHelmList(query: QueryState<HelmFilter>, regionId: number = 1, clusterId: string) {
  let { paging } = query;

  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/json?status_codes=DEPLOYED&&status_codes=FAILED&&status_codes=DELETING&&status_codes=DELETED&&status_codes=UNKNOWN&&sort_by=LAST_RELEASED&&sort_order=DESC&&limit=${
    paging.pageSize
  }`;

  try {
    let response = await GET(url, regionId, clusterId, false);
    let data = response;
    let helms: Helm[] = [];
    helms =
      data && data.releases
        ? data.releases.map(item => {
            let configs = [];
            if (item.config && item.config.raw) {
              let configValue = {};
              try {
                configValue = JSON.parse(item.config.raw);
              } catch (e) {
                try {
                  configValue = JsYAML.safeLoad(item.config.raw);
                } catch (ee) {}
              }
              for (let key in configValue) {
                configs.push({
                  key,
                  value: configValue[key]
                });
              }
            }
            return { ...item, id: item.name, config: configs };
          })
        : [];
    const result: RecordSet<Helm> = {
      recordCount: helms.length,
      records: helms
    };

    return result;
  } catch (error) {
    if (error.message === 'EOF') {
      return {
        recordCount: 0,
        records: []
      };
    } else {
      try {
        let errorInfo = JSON.parse(error.data.Response.Error.Message);
        error.message = errorInfo.message;
        tips.success(error.message, 2000);
        throw error;
      } catch (e) {
        throw error;
      }
    }
  }
}

/**
 * 拉取Helm详情
 */
export async function fetchHelm(params: { helmName: string }, regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/${params.helmName}/status/json`;
  let response = await GET(url, regionId, clusterId);
  return response;
}

/**
 * 拉取Helm详情
 */
export async function fetchHelmResourceList(params: { helmName: string }, regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/${params.helmName}/content/json`;
  let response = await GET(url, regionId, clusterId);
  return response;
}

/**
 * 删除Helm
 */
export async function deleteHelm(params: { helmName: string }, regionId: number = 1, clusterId: string) {
  let resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  let url = `${reduceK8sRestfulPath({
    resourceInfo,
    specificName: clusterId
  })}/helm/tiller/v2/releases/${params.helmName}/json?purge=true`;
  await DELETE(url, regionId, clusterId);
}

export async function getTencentHubToken() {
  let params: RequestParams = {
    apiParams: {
      module: 'thub',
      interfaceName: 'GetUserToken',
      regionId: 1,
      restParams: {
        Version: '2018-04-18'
      }
    }
  };

  let response = await reduceNetworkRequest(params);
  if (response.code === 0) {
    return JSON.parse(response.data.Data).token;
  }
}

export async function fetchTencenthubNamespaceList() {
  let params: RequestParams = {
    apiParams: {
      module: 'thub',
      interfaceName: 'DescribeMasterNamespaces',
      regionId: 1,
      restParams: {
        Version: '2018-04-18'
      }
    }
  };

  let response = await reduceNetworkRequest(params);
  if (response.code === 0) {
    let data = JSON.parse(response.data.Data);
    data = data.map(item => {
      return {
        name: item
      };
    });
    const result: RecordSet<TencenthubNamespace> = {
      recordCount: data.length,
      records: data
    };

    return result;
  }
}
export async function fetchTencenthubChartList(namespace: string) {
  let params: RequestParams = {
    apiParams: {
      module: 'thub',
      interfaceName: 'ListChart',
      regionId: 1,
      restParams: {
        Version: '2018-04-18',
        Namespace: namespace
      },
      opts: {
        tipErr: false
      }
    }
  };

  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      let data = JSON.parse(response.data.Data);
      const result: RecordSet<TencenthubChart> = {
        recordCount: data.list.length,
        records: data.list
      };

      return result;
    } else {
      return {
        recordCount: 0,
        records: []
      };
    }
  } catch (e) {
    if (e.message === 'authentication required') {
      setTimeout(() => {
        tips.error(t('您还没有开通TencentHub帐户，请先到TencentHub开通帐户'), 5000);
      }, 1000);
    } else {
      tips.error(e.message, 2000);
    }
    return {
      recordCount: 0,
      records: []
    };
    // return { content: t('请先去TencentHub开通帐号') };
  }
}
export async function fetchTencenthubChartVersionList(namespace: string, helmName: string) {
  let params: RequestParams = {
    apiParams: {
      module: 'thub',
      interfaceName: 'ListChartVersion',
      regionId: 1,
      restParams: {
        Version: '2018-04-18',
        Namespace: namespace,
        ChartName: helmName
      }
    }
  };
  let response = await reduceNetworkRequest(params);
  if (response.code === 0) {
    let data = JSON.parse(response.data.Data);

    const result: RecordSet<TencenthubChartVersion> = {
      recordCount: data.list.length,
      records: data.list
    };
    return result;
  }
}
export async function fetchTencenthubChartReadMe(namespace: string, helmName: string, version: string) {
  let params: RequestParams = {
    apiParams: {
      module: 'thub',
      interfaceName: 'DescribeChartFile',
      regionId: 1,
      restParams: {
        Version: '2018-04-18',
        Namespace: namespace,
        ChartName: helmName,
        ChartVersion: version,
        File: helmName + '/README.md'
      },
      opts: {
        tipErr: false
      }
    }
  };

  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      let data = JSON.parse(response.data.Data);
      return data;
    }
  } catch (e) {
    return { content: t('该Chart暂未提供ReadMe') };
  }
}

/**
 * 获取资源的具体的 yaml文件
 * @param resourceIns: Resource[]   当前需要请求的具体资源数据
 * @param resourceInfo: ResouceInfo 当前请求数据url的基本配置
 */
export async function fetchUserPortal(resourceInfo: ResourceInfo) {
  let url = reduceK8sRestfulPath({ resourceInfo });

  // 构建参数
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  return response.data;
}

/**
 * Namespace查询
 * @param query Namespace查询的一些过滤条件
 */
export async function fetchProjectNamespaceList(query: QueryState<ResourceFilter>) {
  let { filter } = query;
  let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
  let url = reduceK8sRestfulPath({
    resourceInfo: NamespaceResourceInfo,
    specificName: filter.specificName,
    extraResource: 'namespaces'
  });
  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };

  let response = await reduceNetworkRequest(params);
  let namespaceList = [],
    total = 0;
  if (response.code === 0) {
    let list = response.data;
    total = list.items.length;
    namespaceList = list.items.map(item => {
      return Object.assign({}, item, { id: uuid(), name: item.metadata.name });
    });
  }

  const result: RecordSet<Resource> = {
    recordCount: total,
    records: namespaceList
  };

  return result;
}
