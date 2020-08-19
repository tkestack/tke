import { resourceConfig } from '@config';
import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { apiServerVersion } from '../../../../config';
import {
  Method,
  operationResult,
  reduceK8sQueryString,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  reduceNetworkWorkflow,
  requestMethodForAction
} from '../../../../helpers';
import { isEmpty } from '../../common';
import { CreateResource, MergeType, RequestParams, ResourceInfo, UserDefinedHeader } from '../../common/models';
import {
  DifferentInterfaceResourceOperation,
  LogContentQuery,
  LogHierarchyQuery,
  Namespace,
  Resource,
  ResourceFilter
} from '../models';

// 提示框
const tips = seajs.require('tips');

/**
 * namespace列表的查询
 * @param query: QueryState<ResourceFilter> namespace列表的查询
 * @param resourceInfo: ResourceInfo 当前namespace查询api的配置
 */
export async function fetchNamespaceList(query: QueryState<ResourceFilter>, resourceInfo: ResourceInfo) {
  let { filter, search } = query;
  let { clusterId, regionId } = filter;

  let namespaceList = [];

  // 获取k8s的url
  let url = reduceK8sRestfulPath({ resourceInfo });

  if (search) {
    url = url + '/' + search;
  }

  /** 构建参数 */
  let params: RequestParams = {
    method: Method.get,
    url
  };

  try {
    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      let list = response.data;
      if (list.items) {
        namespaceList = list.items.map(item => {
          return {
            id: uuid(),
            name: item.metadata.name,
            displayName: item.metadata.name
          };
        });
      } else {
        namespaceList.push({
          id: uuid(),
          name: list.metadata.name,
          displayName: list.metadata.name
        });
      }
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  const result: RecordSet<Namespace> = {
    recordCount: namespaceList.length,
    records: namespaceList
  };

  return result;
}

/**
 * 下发、同步secret
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyNamespaceSecret(resource: CreateResource[], regionId: number) {
  try {
    let { resourceIns, mode, clusterId, resourceInfo, namespace, jsonData } = resource[0];

    let isCreate = mode === 'create';
    let userDefinedHeader: UserDefinedHeader = {};
    let method = Method.post;
    // 获取k8s restfulpath
    let url = reduceK8sRestfulPath({ resourceInfo, namespace });

    // 如果为 同步秘钥，则需要改为patch的方式
    if (!isCreate) {
      url += `/${resourceIns}`;
      userDefinedHeader = {
        'Content-Type': 'application/strategic-merge-patch+json'
      };
      method = Method.patch;
    }

    // 构建参数 requestBody 当中
    let params: RequestParams = {
      method,
      url,
      userDefinedHeader,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId,
        restParams: {
          Method: method,
          Path: url,
          Version: '2018-05-25',
          RequestBody: jsonData
        },
        opts: {
          tipErr: false
        }
      }
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      tips.success(t('下发成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, response);
    }
  } catch (error) {
    tips.error(t('下发失败'), 2000);
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

/**
 * Resource列表的查询
 * @param query:    Resource 的查询过滤条件
 * @param resourceInfo:ResourceInfo 资源的相关配置
 * @param isClearData:  是否清空数据
 * @param k8sQueryObj: any  是否有queryString
 * @param isNeedDes: boolean    是否需要降序展示
 */
export async function fetchResourceList(
  query: QueryState<ResourceFilter>,
  options: {
    resourceInfo: ResourceInfo;
    isClearData?: boolean;
    k8sQueryObj?: any;
    isNeedDes?: boolean;
    isNeedSpecific?: boolean;
    isContinue?: boolean;
  }
) {
  let { filter, search, paging, continueToken } = query,
    { namespace, clusterId, regionId, specificName, meshId } = filter;

  let {
    resourceInfo,
    isClearData = false,
    k8sQueryObj = {},
    isNeedDes = false,
    isNeedSpecific = true,
    isContinue = false
  } = options;
  let resourceList = [];
  let nextContinueToken: string;

  // 如果是主动清空 或者 resourceInfo 为空，都不需要发请求
  if (!isClearData && !isEmpty(resourceInfo)) {
    let k8sUrl = '';
    // 如果有搜索字段的话
    if (search) {
      k8sUrl = reduceK8sRestfulPath({
        resourceInfo,
        namespace,
        specificName: isNeedSpecific ? search : '',
        clusterId,
        meshId
      });
    } else {
      k8sUrl = reduceK8sRestfulPath({
        resourceInfo,
        namespace,
        specificName: isNeedSpecific ? specificName : '',
        clusterId,
        meshId
      });
    }

    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    if (isContinue && !search) {
      let { pageSize } = paging;
      k8sQueryObj = JSON.parse(
        JSON.stringify(
          Object.assign({}, k8sQueryObj, { limit: pageSize, continue: continueToken ? continueToken : undefined })
        )
      );
    }

    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    let queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: k8sUrl });
    let url = k8sUrl + queryString;

    // 构建参数
    let params: RequestParams = {
      method: Method.get,
      url
    };

    try {
      let response = await reduceNetworkRequest(params, clusterId);

      if (response.code === 0) {
        let listItems = response.data;

        // 这里将继续拉取数据的token传递下去
        if (isContinue && listItems.metadata && listItems.metadata.continue) {
          nextContinueToken = listItems.metadata.continue;
        }

        if (listItems.items) {
          resourceList = listItems.items.map(item => {
            return Object.assign({}, item, { id: uuid() });
          });
        } else {
          // 这里是拉取某个具体的resource的时候，没有items属性
          resourceList.push({
            metadata: listItems.metadata,
            spec: listItems.spec,
            status: listItems.status
          });
        }
      }
    } catch (error) {
      // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
      if (+error.response.status !== 404) {
        throw error;
      }
    }
  }

  const result: RecordSet<Resource> = {
    recordCount: resourceList.length,
    records: isNeedDes && resourceList.length > 1 ? resourceList.reverse() : resourceList,
    continueToken: nextContinueToken,
    continue: nextContinueToken ? true : false
  };

  return result;
}

/**
 * 获取具体的某个资源，用于在某个资源下，获取其他资源的途径,
 * @param query:    ResourceFilter 的查询过滤条件
 * @param resourceInfo: ResourcrInfo    资源的具体配置
 * @param isClearData: boolean  是否清除数据
 * @param isRecordSet: boolean  返回的数据是否为recordset类型
 * @param k8sQueryObj: any  是否有queryString
 */
export async function fetchSpecificResourceList(
  query: QueryState<ResourceFilter>,
  resourceInfo: ResourceInfo,
  isClearData: boolean = false,
  isRecordSet: boolean = false,
  k8sQueryObj: any = {}
) {
  let { filter } = query,
    { namespace, clusterId, regionId, specificName } = filter;

  let result: any;
  let resourceList = [];

  if (!isClearData) {
    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName, clusterId });
    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    let queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: k8sUrl });
    // 这里是拼接查询的 queryString
    let url = k8sUrl + queryString;

    // 构建参数
    let params: RequestParams = {
      method: Method.get,
      url
    };

    try {
      let response = await reduceNetworkRequest(params, clusterId);

      if (response.code === 0) {
        let listItems = response.data;
        if (listItems.items) {
          resourceList = listItems.items.map(item => {
            return Object.assign({}, item, { id: uuid() });
          });
        } else {
          // 这里是拉取某个具体的resource的时候，没有items属性
          resourceList.push({
            metadata: listItems.metadata,
            spec: listItems.spec,
            status: listItems.status
          });
        }
      }
    } catch (error) {
      if (+error.response.status !== 404) {
        throw error;
      }
    }

    // 这里主要是根据需要返回的类型，如只需要纯数组 还是需要recordSet这种格式的返回
    if (isRecordSet) {
      result = {
        recordCount: resourceList.length,
        records: resourceList
      };
    } else {
      result = resourceList;
    }

    return result;
  }
}

/**
 * 拉取某个具体资源下的额外资源，如 event、pod等
 * @param query: ResourceFilter 的查询过滤条件
 * @param resourceInfo: ResourceInfo 资源的具体配置
 * @param isClearData: boolean 是否清除数据
 * @param extraResource: string 额外的资源，如event
 * @param k8sQueryObj: any  是否有queryString
 * @param isNeedDes: boolean    是否需要降序展示
 */
export async function fetchExtraResourceList(
  query: QueryState<ResourceFilter>,
  resourceInfo: ResourceInfo,
  isClearData: boolean = false,
  extraResource: string = '',
  k8sQueryObj: any = {},
  isNeedDes: boolean = false
) {
  let { filter } = query,
    { namespace, clusterId, regionId, specificName } = filter;

  let extraResourceList = [];

  if (!isClearData) {
    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName, extraResource, clusterId });
    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    let queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: k8sUrl });
    // 这里是拼接查询的 queryString
    let url = k8sUrl + queryString;

    // 构建参数
    let params: RequestParams = {
      method: Method.get,
      url
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      let listItems = response.data;
      if (listItems.items) {
        extraResourceList = listItems.items.map(item => {
          return Object.assign({}, item, { id: uuid() });
        });
      }
    }
  }

  const result: RecordSet<any> = {
    recordCount: extraResourceList.length,
    records: isNeedDes && extraResourceList.length ? extraResourceList.reverse() : extraResourceList
  };

  return result;
}

/**
 * 拉取资源下的日志
 * @param query: ResourceFilter 的查询过滤条件
 * @param resourceInfo: ResourceInfo 资源的具体配置
 * @param isClearData: boolean 是否清除数据
 * @param k8sQueryObj: any  是否有queryString
 */
export async function fetchResourceLogList(
  query: QueryState<ResourceFilter>,
  resourceInfo: ResourceInfo,
  isClearData: boolean = false,
  k8sQueryObj: any = {}
) {
  let { filter } = query,
    { namespace, clusterId, regionId, specificName } = filter;

  let logList = [];

  if (!isClearData) {
    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName, extraResource: 'log', clusterId });
    // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
    let queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: k8sUrl });
    let url = k8sUrl + queryString;

    // 构建参数
    let params: RequestParams = {
      method: Method.get,
      url,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId: regionId,
        restParams: {
          Method: Method.get,
          Path: url,
          Version: '2018-05-25'
        },
        opts: {
          tipErr: false
        }
      }
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      let content = response.data;
      content !== '' && logList.push(content);
    }
  }

  const result: RecordSet<any> = {
    recordCount: logList.length,
    records: logList
  };

  return result;
}

/**
 * 获取日志组件的组件名称
 */
export async function fetchLogagentName(resourceInfo: ResourceInfo, clusterId: string, k8sQueryObj: any = {}) {
  let logAgent = {};
  let k8sUrl = reduceK8sRestfulPath({ resourceInfo });
  let queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: k8sUrl });
  let url = k8sUrl + queryString;
  // 构建参数
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params, clusterId);

  if (response.code === 0) {
    const { items } = response.data;
    if (!isEmpty(items)) {
      // 返回的是数组形式，理论上可以有多个 logAgent，实际上默认取第一个即可
      logAgent = items[0];
    }
  }

  return logAgent;
}

/**
 * 获取日志目录结构
 */
export async function fetchResourceLogHierarchy(query: LogHierarchyQuery) {
  let { agentName, clusterId, namespace, pod, container } = query;
  let logList = [];

  let url = `/apis/logagent.tkestack.io/v1/logagents/${agentName}/filetree`;
  const payload = {
    kind: 'LogFileTree',
    apiVersion: 'logagent.tkestack.io/v1',
    spec: {
      clusterId,
      namespace: namespace.replace(new RegExp(`^${clusterId}-`), ''),
      container,
      pod
    }
  };
  let params: RequestParams = {
    method: Method.post,
    url,
    userDefinedHeader: {},
    data: payload
  };

  let response = await reduceNetworkRequest(params, clusterId);

  const traverse = (hierarchyData, path = '') => {
    let { path: subPath, isDir, children } = hierarchyData;
    // 如果是日志文件的话，构造完整路径，附加到日志文件列表，返回
    if (!isDir) {
      logList.push(path ? path + '/' + subPath : subPath);
      return;
    }
    for (let i = 0; i < children.length; i++) {
      const item = children[i];
      traverse(item, path ? path + '/' + subPath : subPath);
    }
  };

  if (response.code === 0) {
    // 接口成功的话 response.data 为日志内容，失败的话为 { Code: '', Message: '' } 格式的错误
    if (response.data && !response.data.Code) {
      let content = response.data;
      !isEmpty(content) && traverse(content);
    }
  }

  return logList;
}

/**
 * 获取日志内容
 */
export async function fetchResourceLogContent(query: LogContentQuery) {
  let content = '';
  let { agentName, clusterId, namespace, pod, container, start, length, filepath } = query;

  let url = `/apis/logagent.tkestack.io/v1/logagents/${agentName}/filecontent`;
  const payload = {
    kind: 'LogFileContent',
    apiVersion: 'logagent.tkestack.io/v1',
    spec: {
      clusterId,
      namespace: namespace.replace(new RegExp(`^${clusterId}-`), ''),
      container,
      pod,
      start,
      length,
      filepath
    }
  };
  let params: RequestParams = {
    method: Method.post,
    url,
    userDefinedHeader: {},
    data: payload
  };

  let response = await reduceNetworkRequest(params, clusterId);

  if (response.code === 0) {
    let { data } = response;
    if (data && data.content) {
      content = data.content;
    }
  }

  return content;
}

/**
 * 下载日志文件
 * @param query
 */
export async function downloadLogFile(query) {
  let content = '';
  let { agentName, clusterId, namespace, pod, container, filepath } = query;

  let url = `/apis/logagent.tkestack.io/v1/logagents/${agentName}/filedownload`;
  const payload = {
    pod,
    namespace: namespace.replace(new RegExp(`^${clusterId}-`), ''),
    container,
    path: filepath
  };
  // 构建参数
  let params: RequestParams = {
    method: Method.post,
    url,
    data: payload
  };

  let response = await reduceNetworkRequest(params, clusterId);

  if (response.code === 0) {
    const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/x-tar' }));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filepath);
    document.body.appendChild(link);
    link.click();
  }

  return content;
}

/**
 * 同时创建多种资源
 * @param resource: CreateResource 创建resourceIns的相关信息
 * @param regionId: number 地域的ID
 */
export async function applyResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { clusterId, yamlData, jsonData } = resource[0];

    let url = `/${apiServerVersion.basicUrl}/${apiServerVersion.group}/${apiServerVersion.version}/clusters/${clusterId}/apply`;

    // 这里是独立部署版 和 控制台共用的参数，只有是yamlData的时候才需要userdefinedHeader，如果是jaonData的话，就不需要了
    let userDefinedHeader: UserDefinedHeader = yamlData
      ? {
          Accept: 'application/json',
          'Content-Type': 'application/yaml'
        }
      : {};

    // 构建参数
    let params: RequestParams = {
      method: Method.post,
      url,
      userDefinedHeader,
      data: yamlData ? yamlData : jsonData
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

/**创建多种资源 跟applyResourceIns不同的是  每个资源调用的是不同的接口
 *  operations :resources[index] 对应 operations[index]
 */
export async function applyDifferentInterfaceResource(
  resources: CreateResource[],
  operations: DifferentInterfaceResourceOperation[] = []
) {
  let allResponses = []; //收集所有资源的返回结果
  try {
    for (let index = 0; index < resources.length; index++) {
      let { mode, resourceIns, clusterId, yamlData, resourceInfo, namespace, jsonData } = resources[index];
      let extraResource = operations[index] && operations[index].extraResource ? operations[index].extraResource : '';
      let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, extraResource, clusterId });
      //拼接字符串查询参数
      let queryUrl =
        operations[index] && operations[index].query
          ? reduceK8sQueryString({ k8sQueryObj: operations[index].query })
          : '';
      url = url + queryUrl;
      let method = requestMethodForAction(mode);
      // 这里是独立部署版 和 控制台共用的参数，只有是yamlData的时候才需要userdefinedHeader，如果是jaonData的话，就不需要了
      let userDefinedHeader: UserDefinedHeader = yamlData
        ? {
            Accept: 'application/json',
            'Content-Type': 'application/yaml'
          }
        : {};
      // 构建参数
      let params: RequestParams = {
        method,
        url,
        userDefinedHeader,
        data: yamlData ? yamlData : jsonData
      };
      let response = await reduceNetworkRequest(params, clusterId);
      allResponses.push(response);
    }

    //统一处理相应结果
    allResponses.forEach(response => {
      //有一个响应出错
      if (response.code !== 0) {
        return operationResult(resources, reduceNetworkWorkflow(response));
      }
    });
    //所有的响应都OK的话
    return operationResult(resources);
  } catch (error) {
    return operationResult(resources, reduceNetworkWorkflow(error));
  }
}

/**
 * 创建ResourceIns
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { mode, resourceIns, clusterId, yamlData, resourceInfo, namespace, jsonData, meshId } = resource[0];

    let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId, meshId });
    // 获取具体的请求方法，create为POST，modify为PUT
    let method = requestMethodForAction(mode);
    // 这里是独立部署版 和 控制台共用的参数，只有是yamlData的时候才需要userdefinedHeader，如果是jaonData的话，就不需要了
    let userDefinedHeader: UserDefinedHeader = yamlData
      ? {
          Accept: 'application/json',
          'Content-Type': 'application/yaml'
        }
      : {};

    // 构建参数
    let params: RequestParams = {
      method,
      url,
      userDefinedHeader,
      data: yamlData ? yamlData : jsonData
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

/**
 * 创建ResourceIns
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function modifyMultiResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let requests = resource.map(async item => {
      let { mode, resourceIns, clusterId, yamlData, resourceInfo, namespace, jsonData } = item;
      let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId });
      // 获取具体的请求方法，create为POST，modify为PUT
      let method = requestMethodForAction(mode);
      // 这里是独立部署版 和 控制台共用的参数，只有是yamlData的时候才需要userdefinedHeader，如果是jaonData的话，就不需要了
      let userDefinedHeader: UserDefinedHeader = yamlData
        ? {
            Accept: 'application/json',
            'Content-Type': 'application/yaml'
          }
        : {};
      let param = {
        method,
        url,
        userDefinedHeader,
        data: yamlData ? yamlData : jsonData,
        apiParams: {
          module: 'tke',
          interfaceName: 'ForwardRequest',
          regionId,
          restParams: {
            Method: method,
            Path: url,
            Version: '2018-05-25',
            RequestBody: yamlData ? yamlData : jsonData
          },
          opts: {
            tipErr: false
          }
        }
      };
      let response = reduceNetworkRequest(param, clusterId);
      return response;
    });
    // 构建参数
    let response = await Promise.all(requests);
    if (response.every(r => r.code === 0)) {
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

/**
 * 删除ResourceIns
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function deleteResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { resourceIns, clusterId, resourceInfo, namespace, meshId, isSpetialNamespace = true } = resource[0];

    let k8sUrl = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId, meshId, isSpetialNamespace });
    let url = k8sUrl;

    // 是用于后台去异步的删除resource当中的pod
    let extraParamsForDelete = {
      propagationPolicy: 'Background'
    };
    if (resourceInfo.headTitle === 'Namespace') {
      extraParamsForDelete['gracePeriodSeconds'] = 0;
    }

    // 构建参数 requestBody 当中
    let params: RequestParams = {
      method: Method.delete,
      url,
      data: JSON.stringify(extraParamsForDelete)
    };
    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      tips.success(t('删除成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}
/**
 * 回滚ResourceIns
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function rollbackResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { resourceIns, clusterId, resourceInfo, namespace, jsonData } = resource[0];

    let rsResourceInfo = resourceConfig(resourceInfo.k8sVersion).rs;
    /// #if project
    //业务侧ns eg: cls-xxx-ns 需要去除前缀
    if (resourceInfo.namespaces) {
      namespace = namespace.split('-').splice(2).join('-');
    }
    /// #endif
    // 因为回滚需要使用特定的apiVersion，故不用reduceK8sRestful
    let k8sUrl =
      `/${resourceInfo.basicEntry}/apps/v1beta1/` +
      (resourceInfo.namespaces ? `${resourceInfo.namespaces}/${namespace}/` : '') +
      `${resourceInfo.requestType['list']}/${resourceIns}/rollback`;
    let url = k8sUrl;

    // 构建参数 requestBody 当中
    let params: RequestParams = {
      method: Method.post,
      url,
      data: jsonData,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId,
        restParams: {
          Method: Method.post,
          Path: url,
          Version: '2018-05-25',
          RequestBody: jsonData
        }
      }
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      tips.success(t('回滚成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

/**
 * 更新某个具体的deployment的资源的yaml文件
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function updateResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let { resourceIns, clusterId, resourceInfo, namespace, jsonData, isStrategic = true } = resource[0];

    let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId });
    let params: RequestParams = {
      method: Method.patch,
      url,
      userDefinedHeader: {
        'Content-Type': isStrategic ? 'application/strategic-merge-patch+json' : 'application/merge-patch+json'
      },
      data: jsonData,
      apiParams: {
        module: 'tke',
        interfaceName: 'ForwardRequest',
        regionId,
        restParams: {
          Method: Method.patch,
          Path: url,
          Version: '2018-05-25',
          RequestBody: jsonData
        }
      }
    };

    let response = await reduceNetworkRequest(params, clusterId);

    if (response.code === 0) {
      tips.success(t('更新成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

/**
 * 更新某个具体的deployment的资源的yaml文件
 * @param resource: CreateResource   创建resourceIns的相关信息
 * @param regionId: number 地域的id
 */
export async function updateMultiResourceIns(resource: CreateResource[], regionId: number) {
  try {
    let requests = resource.map(async item => {
      let { resourceIns, clusterId, resourceInfo, namespace, jsonData, mergeType } = item;

      let url = reduceK8sRestfulPath({ resourceInfo, namespace, specificName: resourceIns, clusterId });
      let params: RequestParams = {
        method: Method.patch,
        url,
        userDefinedHeader: {
          'Content-Type': mergeType ? mergeType : MergeType.StrategicMerge
        },
        data: jsonData,
        apiParams: {
          module: 'tke',
          interfaceName: 'ForwardRequest',
          regionId,
          restParams: {
            Method: Method.patch,
            Path: url,
            Version: '2018-05-25',
            RequestBody: jsonData
          }
        }
      };

      let response = reduceNetworkRequest(params, clusterId);
      return response;
    });
    let response = await Promise.all(requests);
    if (response.every(r => r.code === 0)) {
      tips.success(t('更新成功'), 2000);
      return operationResult(resource);
    } else {
      return operationResult(resource, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(resource, reduceNetworkWorkflow(error));
  }
}

/**
 * 获取资源的具体的 yaml文件
 * @param resourceIns: Resource[]   当前需要请求的具体资源数据
 * @param resourceInfo: ResouceInfo 当前请求数据url的基本配置
 */
export async function fetchResourceYaml(
  resourceIns: Resource[] | string,
  resourceInfo: ResourceInfo,
  namespace: string,
  clusterId: string,
  regionId: number
) {
  let url = reduceK8sRestfulPath({
    resourceInfo,
    namespace,
    specificName: Array.isArray(resourceIns) ? resourceIns[0].metadata.name : resourceIns,
    clusterId
  });

  let userDefinedHeader = {
    Accept: 'application/yaml'
  };

  // 构建参数
  let params: RequestParams = {
    method: Method.get,
    url,
    userDefinedHeader
  };

  let response = await reduceNetworkRequest(params, clusterId);
  let yamlList = response.code === 0 ? [response.data] : [];

  const result: RecordSet<Resource> = {
    recordCount: yamlList.length,
    records: yamlList
  };

  return result;
}
