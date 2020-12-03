import { OperationResult, QueryState, RecordSet } from '@tencent/ff-redux';

import {
  Method,
  // operationResult,
  reduceK8sQueryString,
  reduceK8sRestfulPath,
  reduceNetworkRequest,
  reduceNetworkWorkflow,
  RequestResult,
  GET,
  POST,
  PUT,
  PATCH,
  DELETE
} from '../../../helpers';
import { resourceConfig } from '../../../config/resourceConfig';
import { RequestParams, ResourceInfo } from '../common/models';
import { CHART_URL, REPO_URL, Default_D_URL } from './constants/Config';
import {
  ApiKey,
  ApiKeyCreation,
  ApiKeyFilter,
  Repo,
  RepoCreation,
  RepoFilter,
  Chart,
  ChartCreation,
  ChartFilter,
  ChartDetailFilter,
  ChartVersionFilter,
  ChartIns,
  ChartInsFilter,
  Image,
  ImageCreation,
  ImageFilter,
  ChartGroup,
  ChartGroupFilter,
  ChartGroupDetailFilter,
  Project,
  ChartVersion,
  ChartInfoFilter,
  Cluster,
  ClusterFilter,
  Namespace,
  ProjectNamespace,
  NamespaceFilter,
  ProjectNamespaceFilter,
  ProjectFilter,
  UserPlain
} from './models';

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
}

/** 访问凭证相关 */
export async function fetchApiKeyList(query: QueryState<ApiKeyFilter>) {
  const { search, paging } = query;
  const apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
  const url = reduceK8sRestfulPath({
    resourceInfo: apiKeyResourceInfo
  });

  const params: RequestParams = {
    method: Method.get,
    url
  };

  const response = await reduceNetworkRequest(params);
  let apiKeyList = [];
  try {
    if (response.code === 0) {
      const listItems = response.data;
      if (listItems.items) {
        apiKeyList = listItems.items.map((item, index) => {
          return Object.assign({}, item, { id: index });
        });
      }
    }
  } catch (error) {
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  if (search) {
    apiKeyList = apiKeyList.filter(x => x.description.includes(query.search));
  }

  const result: RecordSet<ApiKey> = {
    recordCount: apiKeyList.length,
    records: apiKeyList
  };

  return result;
}

export async function createApiKey(apiKeys: ApiKeyCreation[]) {
  try {
    const apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
    const url = reduceK8sRestfulPath({
      resourceInfo: apiKeyResourceInfo
    });

    const apiKey = apiKeys[0];
    /** 构建参数 */
    const requestParams = {
      description: apiKey.description,
      expire: apiKey.expire + apiKey.unit
    };
    const params: RequestParams = {
      method: Method.post,
      url: url + '/default/token',
      data: requestParams
    };
    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(apiKey);
    } else {
      return operationResult(apiKey, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function deleteApiKey(apiKeys: ApiKey[]) {
  try {
    const apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
    const url = reduceK8sRestfulPath({
      resourceInfo: apiKeyResourceInfo,
      specificName: apiKeys[0].metadata.name
    });

    const params: RequestParams = {
      method: Method.delete,
      url
    };

    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(apiKeys);
    } else {
      return operationResult(apiKeys, response);
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function toggleKeyStatus(apiKeys: ApiKey[]) {
  try {
    const apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
    const url = reduceK8sRestfulPath({
      resourceInfo: apiKeyResourceInfo,
      specificName: apiKeys[0].metadata.name
    });

    apiKeys[0].status = Object.assign({}, apiKeys[0].status, {
      disabled: !apiKeys[0].status.disabled
    });

    const requestParams = apiKeys[0];
    const params: RequestParams = {
      method: Method.put,
      url: url,
      data: requestParams
    };

    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(apiKeys);
    } else {
      return operationResult(apiKeys, response);
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

/** 镜像仓库相关 */
export async function fetchRepoList(query: QueryState<RepoFilter>) {
  const { search, paging } = query;

  const params: RequestParams = {
    method: Method.get,
    url: REPO_URL
  };

  const response = await reduceNetworkRequest(params);
  let repoList = [],
    total = 0;
  try {
    if (response.code === 0) {
      const listItems = response.data;
      if (listItems.items) {
        repoList = listItems.items.map((item, index) => {
          return Object.assign({}, item, { id: index });
        });
      }
    }
  } catch (error) {
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  if (search) {
    repoList = repoList.filter(x => x.spec.displayName.includes(query.search));
  }
  total = repoList.length;

  const result: RecordSet<Repo> = {
    recordCount: total,
    records: repoList
  };

  return result;
}

export async function createRepo(repos: RepoCreation[]) {
  try {
    const repo = repos[0];
    /** 构建参数 */
    const requestParams = {
      apiVersion: 'registry.tkestack.io/v1',
      kind: 'Namespace',
      spec: {
        displayName: repo.displayName,
        name: repo.name,
        visibility: repo.visibility || 'Public'
      }
    };
    const params: RequestParams = {
      method: Method.post,
      url: REPO_URL,
      data: requestParams
    };
    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(repo);
    } else {
      return operationResult(repo, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function deleteRepo(repos: Repo[]) {
  try {
    const params: RequestParams = {
      method: Method.delete,
      url: REPO_URL + repos[0].metadata.name
    };

    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(repos);
    } else {
      return operationResult(repos, response);
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function createChart(charts: ChartCreation[]) {
  try {
    const chart = charts[0];
    /** 构建参数 */
    const requestParams = {
      apiVersion: 'registry.tkestack.io/v1',
      kind: 'ChartGroup',
      spec: {
        displayName: chart.displayName,
        name: chart.name,
        visibility: chart.visibility || 'Public'
      }
    };
    const params: RequestParams = {
      method: Method.post,
      url: CHART_URL,
      data: requestParams
    };
    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(chart);
    } else {
      return operationResult(chart, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function deleteChart(charts: Chart[]) {
  try {
    const params: RequestParams = {
      method: Method.delete,
      url: CHART_URL + charts[0].metadata.name
    };

    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(charts);
    } else {
      return operationResult(charts, response);
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function fetchChartInsList(query: QueryState<ChartInsFilter>) {
  const { search, paging, filter } = query;

  const params: RequestParams = {
    method: Method.get,
    url: `${REPO_URL}${filter.chartgroup}/charts`
  };

  const response = await reduceNetworkRequest(params);
  let chartList = [],
    total = 0;
  try {
    if (response.code === 0) {
      const listItems = response.data;
      if (listItems.items) {
        chartList = listItems.items.map((item, index) => {
          return Object.assign({}, item, { id: index });
        });
      }
    }
  } catch (error) {
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  if (search) {
    chartList = chartList.filter(x => x.spec.displayName.includes(query.search));
  }
  total = chartList.length;

  const result: RecordSet<ChartIns> = {
    recordCount: total,
    records: chartList
  };

  return result;
}

/** 镜像相关 */
export async function fetchImageList(query: QueryState<ImageFilter>) {
  const { search, paging, filter } = query;

  const params: RequestParams = {
    method: Method.get,
    url: `${REPO_URL}${filter.namespace}/repositories`
  };

  const response = await reduceNetworkRequest(params);
  let imageList = [],
    total = 0;
  try {
    if (response.code === 0) {
      const listItems = response.data;
      if (listItems.items) {
        imageList = listItems.items.map((item, index) => {
          return Object.assign({}, item, { id: index });
        });
      }
    }
  } catch (error) {
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  if (search) {
    imageList = imageList.filter(x => x.spec.displayName.includes(query.search));
  }
  total = imageList.length;

  const result: RecordSet<Image> = {
    recordCount: total,
    records: imageList
  };

  return result;
}

export async function createImage(images: ImageCreation[]) {
  try {
    const image = images[0];
    /** 构建参数 */
    const requestParams = {
      apiVersion: 'registry.tkestack.io/v1',
      kind: 'Repository',
      metadata: {
        namespace: image.namespace
      },
      spec: {
        displayName: image.displayName,
        name: image.name,
        namespaceName: image.namespaceName,
        visibility: image.visibility || 'Public'
      }
    };
    const params: RequestParams = {
      method: Method.post,
      url: `${REPO_URL}${image.namespace}/repositories`,
      data: requestParams
    };
    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(image);
    } else {
      return operationResult(image, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function deleteImage(images: Image[]) {
  try {
    const image = images[0];
    const params: RequestParams = {
      method: Method.delete,
      url: `${REPO_URL}${image.metadata.namespace}/repositories/${image.metadata.name}`
    };

    const response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(images);
    } else {
      return operationResult(images, response);
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

export async function fetchDockerRegUrl() {
  const _localUrl = localStorage.getItem('_registry_url');
  if (_localUrl && !/^https:\/\/.*/.test(_localUrl)) {
    return _localUrl;
  } else {
    const paramsInfo: RequestParams = {
      method: Method.get,
      url: '/apis/gateway.tkestack.io/v1/tokens/info'
    };
    const paramsSysInfo: RequestParams = {
      method: Method.get,
      url: '/apis/gateway.tkestack.io/v1/sysinfo'
    };

    const info = await reduceNetworkRequest(paramsInfo);
    const sysInfo = await reduceNetworkRequest(paramsSysInfo);
    try {
      let url = '';
      if (info.code === 0) {
        url += info.data.extra.tenantid[0];
      } else if (sysInfo.code === 0) {
        url += sysInfo.data.registry.defaultTenant;
      } else {
        url += 'default';
      }

      url += '.';

      if (sysInfo.code === 0) {
        url += sysInfo.data.registry.domainSuffix;
      } else {
        url += 'registry.com';
      }
      if (info.code === 0 && sysInfo.code === 0) {
        try {
          localStorage.setItem('_registry_url', url);
        } catch (e) {}
      }
      return url;
    } catch (e) {
      return Default_D_URL;
    }
  }
}

/** 以下代码为新版本代码 */
/** 模板仓库 */
/**
 * 仓库列表
 * @param query
 */
export async function fetchChartGroupList(query: QueryState<ChartGroupFilter>) {
  const { keyword, filter } = query;
  // TODO FIX
  // 根据业务过滤
  const queryObj = filter.repoType
    ? {
        'fieldSelector=repoType': filter.repoType
      }
    : {};
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const rr: RequestResult = await GET({ url: url + queryString, keyword: keyword });
  const objs: ChartGroup[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<ChartGroup> = {
    recordCount: objs.length,
    records: objs
  };
  return result;
}

/**
 * 查询仓库
 * @param filter 查询条件参数
 */
export async function fetchChartGroup(filter: ChartGroupDetailFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.name });
  const rr: RequestResult = await GET({ url });
  return rr.data;
}

/**
 * 修改仓库
 * @param chartGroupInfo
 */
export async function updateChartGroup([chartGroupInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: chartGroupInfo.metadata.name });
  const rr: RequestResult = await PUT({
    url,
    bodyData: chartGroupInfo
  });
  return operationResult(rr.data, rr.error);
}

/**
 * 增加仓库
 * @param chartGroupInfo
 */
export async function addChartGroup([chartGroupInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  chartGroupInfo.spec.importedInfo.password = chartGroupInfo.spec.importedInfo.password
    ? btoa(chartGroupInfo.spec.importedInfo.password)
    : '';
  const rr: RequestResult = await POST({
    url,
    bodyData: chartGroupInfo
  });
  return operationResult(rr.data, rr.error);
}

/**
 * 删除仓库
 * @param group
 */
export async function deleteChartGroup([chartGroup]: ChartGroup[]) {
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: chartGroup.metadata.name });
  const rr: RequestResult = await DELETE({
    url
  });
  return operationResult(rr.data, rr.error);
}

/**
 * 同步仓库
 * @param group
 */
export async function repoUpdateChartGroup([chartGroup]: ChartGroup[]) {
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    extraResource: 'repoupdating',
    specificName: chartGroup.metadata.name
  });
  const rr: RequestResult = await POST({
    url
  });
  return operationResult(rr.data, rr.error);
}

/**
 * 有权限的业务列表
 * @param query
 */
// export async function fetchManagedProjectList(query: QueryState<void>) {
//   const empty: RecordSet<Project> = {
//     recordCount: 0,
//     records: []
//   };
//   const resourceInfo: ResourceInfo = resourceConfig()['info'];
//   const url = reduceK8sRestfulPath({ resourceInfo });
//   let rr: RequestResult = await GET(url, true);
//   if (rr.error) {
//     return empty;
//   }
//   let uid: string = rr.data.uid;
//   const projectResourceInfo: ResourceInfo = resourceConfig()['projects'];
//   const projectUrl = reduceK8sRestfulPath({ resourceInfo: projectResourceInfo });
//   let prr: RequestResult = await GET(projectUrl, true);
//   const projectBelongResourceInfo: ResourceInfo = resourceConfig()['user'];
//   const projectBelongUrl = reduceK8sRestfulPath({
//     resourceInfo: projectBelongResourceInfo,
//     specificName: uid,
//     extraResource: 'projects'
//   });
//   let pbrr: RequestResult = await GET(projectBelongUrl, true);
//   if (prr.error || pbrr.error) {
//     return empty;
//   }
//   let managedProjects = !pbrr.error && pbrr.data.managedProjects ? Object.keys(pbrr.data.managedProjects) : [];
//   let items = [];
//   if (!prr.error && prr.data.items) {
//     prr.data.items.forEach(i => {
//       if (managedProjects.indexOf(i.metadata.name) > -1) {
//         items.push({
//           id: i.metadata && i.metadata.name,
//           metadata: {
//             name: i.metadata.name
//           },
//           spec: {
//             displayName: i.spec.displayName
//           }
//         });
//       }
//     });
//   }
//   const result: RecordSet<Project> = {
//     recordCount: items.length,
//     records: items
//   };
//   return result;
// }

/**
 * 有权限的业务列表
 * @param query
 */
export async function fetchPortalProjectList(query: QueryState<ProjectFilter>) {
  const empty: RecordSet<Project> = {
    recordCount: 0,
    records: []
  };
  const resourceInfo: ResourceInfo = resourceConfig()['portal'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const rr: RequestResult = await GET({ url });
  if (rr.error) {
    return empty;
  }
  const items = Object.keys(rr.data.projects).map(key => {
    return {
      id: key,
      metadata: {
        name: key
      },
      spec: {
        displayName: rr.data.projects[key]
      }
    };
  });
  const result: RecordSet<Project, ChartInfoFilter> = {
    recordCount: items.length,
    records: items,
    data: query.filter.chartInfoFilter
  };
  return result;
}

/**
 * 获取人员信息
 * @param query
 */
export async function fetchUserInfo() {
  const resourceInfo: ResourceInfo = resourceConfig()['info'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const rr: RequestResult = await GET({ url });
  return rr.data;
}

/**
 * 仓库列表
 * @param query
 */
export async function fetchChartList(query: QueryState<ChartFilter>) {
  const { keyword, filter } = query;
  let fieldSelector = '';
  if (filter.repoType) {
    fieldSelector = 'repoType=' + filter.repoType;
  }
  if (filter.projectID) {
    fieldSelector = fieldSelector + ',projectID=' + filter.projectID;
  }
  const queryObj = fieldSelector
    ? {
        fieldSelector: fieldSelector
      }
    : {};
  const resourceInfo: ResourceInfo = resourceConfig()['chart'];
  const opts = { resourceInfo: resourceInfo };
  // if (filter.namespace) {
  //   opts['namespace'] = filter.namespace;
  //   opts['isSpecialNamespace'] = true;
  // }
  const url = reduceK8sRestfulPath(opts);
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  // reduceNetworkRequest设置了未传业务id则从cookie中读取的逻辑，传空的业务id是为了适配这段逻辑
  const rr: RequestResult = await GET({
    url: url + queryString,
    keyword
  });
  const objs: Chart[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<Chart> = {
    recordCount: objs.length,
    records: objs
  };
  return result;
}

/**
 * 查询Chart
 * @param filter 查询条件参数
 */
export async function fetchChart(filter: ChartDetailFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['chart'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: filter.namespace,
    specificName: filter.name,
    isSpecialNamespace: true
  });
  const rr: RequestResult = await GET({ url });
  return rr.data;
}

/**
 * 修改Chart
 * @param chartInfo
 */
export async function updateChart([chartInfo], filter: ChartDetailFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['chart'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: chartInfo.metadata.namespace,
    specificName: chartInfo.metadata.name,
    isSpecialNamespace: true
  });
  const rr: RequestResult = await PUT({ url, bodyData: chartInfo });
  return operationResult(rr.data, rr.error);
}

/**
 * 删除模板
 * @param group
 */
export async function deleteChartVersion([chartVersion]: ChartVersion[], filter: ChartVersionFilter) {
  // const url = `/chart/api/${filter.chartGroupName}/charts/${filter.chartName}/${filter.chartVersion}`;
  // let rr: RequestResult = await DELETE({ url });
  // return operationResult(rr.data, rr.error);
  const resourceInfo: ResourceInfo = resourceConfig()['chart'];
  // const queryObj = {
  //   version: filter.chartVersion
  // };
  // const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: filter.chartDetailFilter.namespace,
    specificName: filter.chartDetailFilter.name,
    extraResource: 'version',
    isSpecialNamespace: true
  });
  const rr: RequestResult = await DELETE({
    url: url + '/' + filter.chartVersion
  });
  return operationResult(rr.data, rr.error);
}

/**
 * 获取模板
 * @param group
 */
export async function fetchChartVersionFile(filter: ChartVersionFilter) {
  const url = `/chart/${filter.chartGroupName}/charts/${filter.chartName}-${filter.chartVersion}.tgz`;
  const rr: RequestResult = await GET({ url });
  return rr.data;
}

/**
 * 查询仓库
 * @param filter 查询条件参数
 */
export async function fetchChartInfo(filter: ChartInfoFilter) {
  const queryObj = {
    // version: filter.chartVersion,
    cluster: filter.cluster
  };
  const resourceInfo: ResourceInfo = resourceConfig()['chart'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: filter.metadata.namespace,
    specificName: filter.metadata.name,
    extraResource: 'version',
    isSpecialNamespace: true
  });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const rr: RequestResult = await GET({
    url: url + '/' + filter.chartVersion + queryString
  });
  return rr.data;
}

/**
 * 增加应用
 * @param appInfo
 */
export async function addApp([appInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: appInfo.metadata.namespace,
    isSpecialNamespace: true
  });
  const rr: RequestResult = await POST({ url, bodyData: appInfo });
  return operationResult(rr.data, rr.error);
}

/**
 * 集群列表
 * @param query
 */
export async function fetchClusterList(query: QueryState<ClusterFilter>) {
  const resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const rr: RequestResult = await GET({ url });
  const objs: Cluster[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<Cluster, ChartInfoFilter> = {
    recordCount: objs.length,
    records: objs,
    data: query.filter.chartInfoFilter
  };
  return result;
}

/**
 * 命名空间列表
 * @param query
 */
export async function fetchNamespaceList(query: QueryState<NamespaceFilter>) {
  const { keyword, filter } = query;
  const resourceInfo: ResourceInfo = resourceConfig()['ns'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const rr: RequestResult = await GET({ url, clusterId: filter.cluster });
  const objs: Namespace[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<Namespace, ChartInfoFilter> = {
    recordCount: objs.length,
    records: objs,
    data: query.filter.chartInfoFilter
  };
  return result;
}

/**
 * 业务Namespace查询
 * @param query Namespace查询的一些过滤条件
 */
export async function fetchProjectNamespaceList(query: QueryState<ProjectNamespaceFilter>) {
  const { keyword, filter } = query;
  const resourceInfo: ResourceInfo = resourceConfig()['namespaces'];
  const url = reduceK8sRestfulPath({ resourceInfo, specificName: filter.projectId, extraResource: 'namespaces' });
  const rr: RequestResult = await GET({ url });
  const objs: ProjectNamespace[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<ProjectNamespace, ChartInfoFilter> = {
    recordCount: objs.length,
    records: objs,
    data: query.filter.chartInfoFilter
  };
  return result;
}

/**
 * 用户列表的查询，不跟localidentities混用，不参杂其他场景参数，如策略、角色
 * @param query 列表查询条件参数
 */
export async function fetchCommonUserList(query: QueryState<void>) {
  const { search } = query;
  const queryObj = {
    'fieldSelector=keyword': search || ''
  };

  const resourceInfo: ResourceInfo = resourceConfig()['user'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const rr: RequestResult = await GET({ url: url + queryString });
  const users: UserPlain[] =
    !rr.error && rr.data.items
      ? rr.data.items.map(i => {
          return {
            id: i.metadata && i.metadata.name,
            name: i.spec && (i.spec.name ? i.spec.name : i.spec.username),
            displayName: i.spec && i.spec.displayName
          };
        })
      : [];
  const result: RecordSet<UserPlain> = {
    recordCount: users.length,
    records: users
  };
  return result;
}
