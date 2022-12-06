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
import { OperationResult, QueryState, RecordSet } from '@tencent/ff-redux';

import {
  Method,
  operationResult,
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
import {
  App,
  AppFilter,
  AppDetailFilter,
  AppResource,
  AppResourceFilter,
  AppHistory,
  AppHistoryFilter,
  History,
  Cluster,
  Namespace,
  ProjectNamespace,
  NamespaceFilter,
  Chart,
  ChartFilter,
  ChartInfo,
  ChartInfoFilter,
  ChartGroup,
  ChartGroupFilter,
  Project,
  ProjectNamespaceFilter
} from './models';

/**
 * 集群列表
 * @param query
 */
export async function fetchClusterList(query: QueryState<void>) {
  const resourceInfo: ResourceInfo = resourceConfig()['cluster'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const rr: RequestResult = await GET({ url });
  const objs: Cluster[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<Cluster> = {
    recordCount: objs.length,
    records: objs
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
  const result: RecordSet<Namespace> = {
    recordCount: objs.length,
    records: objs
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
  const result: RecordSet<ProjectNamespace> = {
    recordCount: objs.length,
    records: objs
  };
  return result;
}

/** 应用 */
/**
 * 应用列表
 * @param query
 */
export async function fetchAppList(query: QueryState<AppFilter>) {
  const { keyword, filter } = query;
  const queryObj = {
    fieldSelector: {
      'spec.targetNamespace': filter.namespace,
      'spec.targetCluster': filter.cluster
    }
  };
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({ resourceInfo: { ...resourceInfo, namespaces: undefined } });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });

  const rr: RequestResult = await GET({ url: url + queryString, clusterId: filter.cluster, keyword });
  const objs: App[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<App> = {
    recordCount: objs.length,
    records: objs
  };
  return result;
}

/**
 * 查询应用
 * @param filter 查询条件参数
 */
export async function fetchApp(filter: AppDetailFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: filter.namespace,
    specificName: filter.name,
    isSpecialNamespace: true
  });
  const rr: RequestResult = await GET({ url, clusterId: filter.cluster });
  return rr.data;
}

/**
 * 修改应用
 * @param appInfo
 */
export async function updateApp([appInfo]) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: appInfo.metadata.namespace,
    specificName: appInfo.metadata.name,
    isSpecialNamespace: true
  });
  // const rr: RequestResult = await PUT({ url, bodyData: appInfo });
  const rr: RequestResult = await PATCH({
    url,
    bodyData: {
      spec: {
        chart: {
          ...appInfo?.spec?.chart
        },

        values: {
          rawValues: appInfo?.spec?.values?.rawValues
        }
      }
    },
    headers: {
      'Content-Type': 'application/strategic-merge-patch+json'
    }
  });
  return operationResult(rr.data, rr.error);
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
 * 删除应用
 * @param group
 */
export async function deleteApp([app]: App[]) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: app.metadata.namespace,
    specificName: app.metadata.name,
    isSpecialNamespace: true
  });
  const rr: RequestResult = await DELETE({ url });
  return operationResult(rr.data, rr.error);
}

/**
 * 查询资源
 * @param filter 查询条件参数
 */
export async function fetchAppResource(filter: AppResourceFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: filter.namespace,
    specificName: filter.name,
    extraResource: 'resources',
    isSpecialNamespace: true
  });
  const rr: RequestResult = await GET({ url, clusterId: filter.cluster });
  return rr.data;
}

/**
 * 查询历史
 * @param filter 查询条件参数
 */
export async function fetchAppHistory(filter: AppHistoryFilter) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: filter.namespace,
    specificName: filter.name,
    extraResource: 'histories',
    isSpecialNamespace: true
  });
  const rr: RequestResult = await GET({ url, clusterId: filter.cluster });
  return rr.data;
}

/**
 * 回滚应用
 * @param group
 */
export async function rollbackApp([app]: History[]) {
  const resourceInfo: ResourceInfo = resourceConfig()['app'];
  const namespace = (app.involvedObject && app.involvedObject.metadata && app.involvedObject.metadata.namespace) || '';
  const name = (app.involvedObject && app.involvedObject.metadata && app.involvedObject.metadata.name) || '';
  const cluster = (app.involvedObject && app.involvedObject.spec && app.involvedObject.spec.targetCluster) || '';
  const queryObj = {
    revision: app.revision,
    cluster: cluster
  };
  const url = reduceK8sRestfulPath({
    resourceInfo,
    namespace: namespace,
    specificName: name,
    extraResource: 'rollback',
    isSpecialNamespace: true
  });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const rr: RequestResult = await POST({ url: url + queryString, bodyData: {}, clusterId: cluster });
  return operationResult(rr.data, rr.error);
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
  const opts = { resourceInfo: { ...resourceInfo, namespaces: undefined } };

  const url = reduceK8sRestfulPath(opts);
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
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
 * 查询仓库
 * @param filter 查询条件参数
 */
export async function fetchChartInfo(filter: ChartInfoFilter) {
  const queryObj = {
    // version: filter.chartVersion,
    cluster: filter.cluster,
    namespace: filter.namespace
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
 * 仓库列表
 * @param query
 */
export async function fetchChartGroupList(query: QueryState<ChartGroupFilter>) {
  const { keyword, filter } = query;
  const queryObj = {};
  const resourceInfo: ResourceInfo = resourceConfig()['chartgroup'];
  const url = reduceK8sRestfulPath({ resourceInfo });
  const queryString = reduceK8sQueryString({ k8sQueryObj: queryObj });
  const rr: RequestResult = await GET({ url: url + queryString, keyword });
  const objs: ChartGroup[] = !rr.error && rr.data.items ? rr.data.items : [];
  const result: RecordSet<ChartGroup> = {
    recordCount: objs.length,
    records: objs
  };
  return result;
}

/**
 * 有管理权限的业务列表
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
export async function fetchPortalProjectList(query: QueryState<void>) {
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
  const result: RecordSet<Project> = {
    recordCount: items.length,
    records: items
  };
  return result;
}
