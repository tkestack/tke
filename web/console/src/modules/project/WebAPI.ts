/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

import { collectionPaging, OperationResult, QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { resourceConfig } from '../../../config/resourceConfig';
import { reduceK8sRestfulPath, reduceNetworkRequest, reduceNetworkWorkflow } from '../../../helpers';
import { Method } from '../../../helpers/reduceNetwork';
import { reduceK8sQueryString } from '../../../helpers/urlUtil';
import { Region, RegionFilter, RequestParams, ResourceInfo } from '../common/models';
import { resourceTypeToUnit } from './constants/Config';
import {
  Cluster,
  ClusterFilter,
  Manager,
  ManagerFilter,
  Namespace,
  NamespaceEdition,
  NamespaceFilter,
  NamespaceOperator,
  Project,
  ProjectEdition,
  ProjectFilter
} from './models';
import { ProjectResourceLimit, UserManagedProject, UserManagedProjectFilter } from './models/Project';

// 返回标准操作结果
function operationResult<T>(target: T[] | T, error?: any): OperationResult<T>[] {
  if (target instanceof Array) {
    return target.map(x => ({ success: !error, target: x, error }));
  }
  return [{ success: !error, target: target as T, error }];
}

/**
 * 业务查询
 * @param query 地域查询的一些过滤条件
 */
export async function fetchProjectList(query: QueryState<ProjectFilter>) {
  let {
    search,
    filter: { parentProject },
    searchFilter: { projectId }
  } = query;
  let projectList = [],
    total = 0;
  let projectResourceInfo: ResourceInfo = resourceConfig().projects;

  let userResourceInfo: ResourceInfo = resourceConfig().portal;
  let portalUrl = reduceK8sRestfulPath({ resourceInfo: userResourceInfo });
  let portalparams: RequestParams = {
    method: Method.get,
    url: portalUrl
  };
  let response = await reduceNetworkRequest(portalparams);
  let isAdmin = true,
    userProjects = [];
  if (response.code === 0) {
    userProjects = Object.keys(response.data.projects).map(key => key);
    isAdmin = response.data.administrator;
  }

  if (isAdmin || parentProject) {
    let k8sQueryObj = {
      fieldSelector: {
        'spec.parentProjectName': parentProject ? parentProject : undefined
      }
    };
    k8sQueryObj = JSON.parse(JSON.stringify(k8sQueryObj));

    let k8sUrl = reduceK8sRestfulPath({
      resourceInfo: projectResourceInfo,
      specificName: projectId ? projectId : null
    });

    let queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: k8sUrl });

    let url = k8sUrl + queryString;

    let params: RequestParams = {
      method: Method.get,
      url
    };

    try {
      let response = await reduceNetworkRequest(params);

      if (response.code === 0) {
        let listItems = response.data;
        if (listItems.items) {
          projectList = listItems.items.map(item => {
            return Object.assign({}, item, { id: item.metadata.name });
          });
        } else {
          // 这里是拉取某个具体的resource的时候，没有items属性
          projectList.push({
            id: listItems.metadata.name,
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
  } else {
    let requests = userProjects.map(projectId => {
      let url = reduceK8sRestfulPath({
        resourceInfo: projectResourceInfo,
        specificName: projectId
      });

      let params: RequestParams = {
        method: Method.get,
        url,
        userDefinedHeader: {
          'X-TKE-ProjectName': projectId
        }
      };
      return reduceNetworkRequest(params).catch(e => null);
    });
    let response = await Promise.all(requests);
    response.forEach(r => {
      if (r && r.code === 0) {
        let listItems = r.data;
        if (listItems.items) {
          projectList = listItems.items.map(item => {
            return Object.assign({}, item, { id: item.metadata.name });
          });
        } else {
          // 这里是拉取某个具体的resource的时候，没有items属性
          projectList.push({
            id: listItems.metadata.name,
            metadata: listItems.metadata,
            spec: listItems.spec,
            status: listItems.status
          });
        }
      }
    });
    total = projectList.length;
  }

  if (search) {
    projectList = projectList.filter(x => x.spec.displayName.includes(query.search));
  }
  total = projectList.length;
  const result: RecordSet<Project> = {
    recordCount: total,
    records: projectList
  };
  return result;
}

/**
 * 业务查询
 * @param query 地域查询的一些过滤条件
 */
export async function fetchProjectDetail(projectId?: string) {
  let projectResourceInfo: ResourceInfo = resourceConfig()['projects'];
  let url = reduceK8sRestfulPath({ resourceInfo: projectResourceInfo, specificName: projectId });
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  if (response.code === 0) {
    return Object.assign({}, response.data, { id: response.data.metadata.name });
  }
}

function _reduceProjectLimit(projectResourceLimit: ProjectResourceLimit[]) {
  let hardInfo = {};
  projectResourceLimit.forEach(item => {
    let value;
    if (resourceTypeToUnit[item.type] === '个' || resourceTypeToUnit[item.type] === '核') {
      value = +item.value;
    } else {
      value = item.value + 'Mi';
    }
    hardInfo[item.type] = value;
  });
  return hardInfo;
}

/**
 * 业务编辑
 */
export async function editProject(projects: ProjectEdition[]) {
  try {
    let projectResourceInfo: ResourceInfo = resourceConfig()['projects'];
    let url = reduceK8sRestfulPath({ resourceInfo: projectResourceInfo });

    /** 构建参数 */

    let clusterObject = {};
    projects[0].clusters.forEach(cluster => {
      let resourceLimitObject = {};
      cluster.resourceLimits.forEach(resourceLimit => {
        resourceLimitObject[resourceLimit.type] = resourceLimit.value;
        if (resourceTypeToUnit[resourceLimit.type] === 'MiB') {
          resourceLimitObject[resourceLimit.type] += 'Mi';
        }
      });
      clusterObject[cluster.name] = { hard: resourceLimitObject };
    });

    let requestParams = {
        kind: projectResourceInfo.headTitle,
        apiVersion: `${projectResourceInfo.group}/${projectResourceInfo.version}`,
        spec: {
          displayName: projects[0].displayName,
          members: projects[0].members.map(m => m.name),
          clusters: clusterObject,
          parentProjectName: projects[0].parentProject ? projects[0].parentProject : undefined
        }
      },
      method = 'POST';

    if (projects[0].id) {
      //修改
      method = 'PUT';
      url += '/' + projects[0].id;
      requestParams = JSON.parse(
        JSON.stringify({
          kind: projectResourceInfo.headTitle,
          apiVersion: `${projectResourceInfo.group}/${projectResourceInfo.version}`,
          metadata: {
            name: projects[0].id,
            resourceVersion: projects[0].resourceVersion
          },
          spec: {
            displayName: projects[0].displayName ? projects[0].displayName : null,
            members: projects[0].members.length ? projects[0].members.map(m => m.name) : null,
            clusters: clusterObject,
            parentProjectName: projects[0].parentProject ? projects[0].parentProject : null
          },
          status: projects[0].status
        })
      );
    }
    let params: RequestParams = {
      method,
      url,
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(projects);
    } else {
      return operationResult(projects, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(projects, reduceNetworkWorkflow(error));
  }
}

/**
 * 业务删除
 */
export async function deleteProject(projects: Project[]) {
  try {
    let projectResourceInfo: ResourceInfo = resourceConfig()['projects'];
    let url = reduceK8sRestfulPath({ resourceInfo: projectResourceInfo, specificName: projects[0].id + '' });
    let params: RequestParams = {
      method: Method.delete,
      url,
      userDefinedHeader: {
        'X-TKE-ProjectName': projects[0].id + ''
      }
    };

    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(projects);
    } else {
      return operationResult(projects, response);
    }
  } catch (error) {
    return operationResult(projects, reduceNetworkWorkflow(error));
  }
}
/**
 * Namespace查询
 * @param query Namespace查询的一些过滤条件
 */
export async function fetchNamespaceList(query: QueryState<NamespaceFilter>) {
  let { filter, search } = query;
  let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
  let url = reduceK8sRestfulPath({
    resourceInfo: NamespaceResourceInfo,
    specificName: filter.projectId,
    extraResource: 'namespaces'
  });
  let namespaceList = [];
  if (search) {
    url = url + '/' + search;
  }

  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };
  try {
    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      let listItems = response.data;
      if (listItems.items) {
        namespaceList = listItems.items.map(item => {
          return Object.assign({}, item, { id: item.metadata.name });
        });
      } else {
        // 这里是拉取某个具体的resource的时候，没有items属性
        namespaceList.push({
          id: listItems.metadata.name,
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

  const result: RecordSet<Namespace> = {
    recordCount: namespaceList.length,
    records: namespaceList
  };

  return result;
}

/**
 * 地域的查询
 * @param query 地域查询的一些过滤条件
 */
export async function fetchRegionList(query?: QueryState<RegionFilter>) {
  let regionList = [];

  regionList = [
    {
      Alias: 'gz',
      CreatedAt: '2018-01-24T19:58:09+08:00',
      Id: 1,
      RegionId: 1,
      RegionName: 'ap-guangzhou',
      Status: 'alluser',
      UpdatedAt: '2018-01-24T19:58:09+08:00'
    }
  ];

  const result: RecordSet<Region> = {
    recordCount: regionList.length,
    records: regionList
  };

  return result;
}

/**
 * 集群列表的查询
 * @param query 集群列表查询的一些过滤条件
 */
export async function fetchClusterList(query: QueryState<ClusterFilter>) {
  let clsuterResource: ResourceInfo = resourceConfig().cluster;
  let url = reduceK8sRestfulPath({ resourceInfo: clsuterResource });

  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };

  let response = await reduceNetworkRequest(params);
  let clusterList = [];
  if (response.code === 0) {
    let list = response.data;
    clusterList = list.items.map(item => {
      return { id: uuid(), clusterId: item.metadata.name, clusterName: item.spec.displayName };
    });
  }

  const result: RecordSet<Cluster> = {
    recordCount: clusterList.length,
    records: clusterList
  };

  return result;
}

/**
 * Namespace编辑
 */
export async function editNamespace(namespaces: NamespaceEdition[], op: NamespaceOperator) {
  try {
    let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
    let url = reduceK8sRestfulPath({
      resourceInfo: NamespaceResourceInfo,
      specificName: op.projectId,
      extraResource: 'namespaces'
    });
    /** 构建参数 */
    let requestParams = {
        kind: NamespaceResourceInfo.headTitle,
        apiVersion: `${NamespaceResourceInfo.group}/${NamespaceResourceInfo.version}`,
        spec: {
          clusterName: namespaces[0].clusterName,
          namespace: namespaces[0].namespaceName,
          projectName: op.projectId,
          hard: _reduceProjectLimit(namespaces[0].resourceLimits)
        }
      },
      method = 'POST';

    if (namespaces[0].id) {
      //修改
      method = 'PUT';
      url += '/' + namespaces[0].id;
      requestParams['metadata'] = {
        name: namespaces[0].id,
        resourceVersion: namespaces[0].resourceVersion,
        projectName: op.projectId
      };
      requestParams['status'] = namespaces[0].status;
    }

    let params: RequestParams = {
      method,
      url,
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      return operationResult(namespaces);
    } else {
      return operationResult(namespaces, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(namespaces, reduceNetworkWorkflow(error));
  }
}

/**
 * Namespace删除
 */
export async function deleteNamespace(namespaces: Namespace[], op: NamespaceOperator) {
  try {
    let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
    let url = reduceK8sRestfulPath({
      resourceInfo: NamespaceResourceInfo,
      specificName: op.projectId,
      extraResource: `namespaces/${namespaces[0].metadata.name}`
    });
    // 是用于后台去异步的删除resource当中的pod
    let extraParamsForDelete = {
      propagationPolicy: 'Background'
    };
    extraParamsForDelete['gracePeriodSeconds'] = 0;
    /** 构建参数 */
    let method = 'DELETE';
    let params: RequestParams = {
      method,
      url,
      data: extraParamsForDelete
    };

    let response = await reduceNetworkRequest(params);

    if (response.code === 0) {
      return operationResult(namespaces);
    } else {
      return operationResult(namespaces, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(namespaces, reduceNetworkWorkflow(error));
  }
}

/**
 * user列表的查询
 * @param query 集群列表查询的一些过滤条件
 */
export async function fetchUser(query: QueryState<ManagerFilter>) {
  let userInfo: ResourceInfo = resourceConfig()['user'];
  let url = reduceK8sRestfulPath({ resourceInfo: userInfo });
  let { filter, search } = query;
  /** 构建参数 */
  if (search) {
    url += `?fieldSelector=keyword=${search}`;
  }
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };

  let response = await reduceNetworkRequest(params);
  let userList = [];
  if (response.code === 0) {
    let list = response.data;
    userList = list.items
      ? list.items.map(item => {
          return {
            id: item.metadata.name,
            displayName: item.spec && item.spec.displayName,
            name: item.spec && item.spec.name
          };
        })
      : [];
  }

  const result: RecordSet<Manager> = {
    recordCount: userList.length,
    records: userList
  };

  return result;
}

/**
 * 查询登陆用户信息
 * @param query
 */
// export async function fetchLoginUserInfo() {
//   let userInfo: ResourceInfo = resourceConfig()['info'];
//   let url = reduceK8sRestfulPath({ resourceInfo: userInfo });
//   let method = 'GET';
//   let params: RequestParams = {
//     method,
//     url
//   };
//   let response = await reduceNetworkRequest(params);
//
//   return response;
// }

/**
 *
 * @param query 集群列表查询的一些过滤条件
 */
export async function fetchAdminstratorInfo() {
  let userResourceInfo: ResourceInfo = resourceConfig().platforms;
  let url = reduceK8sRestfulPath({ resourceInfo: userResourceInfo });
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  let info = {};
  if (response.code === 0) {
    let list = response.data;
    if (list.items) {
      info = list.items.length ? list.items[0] : {};
    }
  }

  return info;
}

/**
 * 业务编辑
 */
export async function modifyAdminstrator(projects: ProjectEdition[]) {
  try {
    let platformsResourceInfo: ResourceInfo = resourceConfig().platforms;
    let url = reduceK8sRestfulPath({ resourceInfo: platformsResourceInfo });

    /** 构建参数 */
    let requestParams = {
        kind: platformsResourceInfo.headTitle,
        apiVersion: `${platformsResourceInfo.group}/${platformsResourceInfo.version}`,
        spec: {
          administrators: projects[0].members.map(m => m.name)
        }
      },
      method = 'POST';

    if (projects[0].id) {
      //修改
      method = 'PUT';
      url += '/' + projects[0].id;
      requestParams = JSON.parse(
        JSON.stringify({
          kind: platformsResourceInfo.headTitle,
          apiVersion: `${platformsResourceInfo.group}/${platformsResourceInfo.version}`,
          metadata: {
            name: projects[0].id,
            resourceVersion: projects[0].resourceVersion
          },
          spec: {
            administrators: projects[0].members.length ? projects[0].members.map(m => m.name) : null
          }
        })
      );
    }
    let params: RequestParams = {
      method,
      url,
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(projects);
    } else {
      return operationResult(projects, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(projects, reduceNetworkWorkflow(error));
  }
}

/**
 * 业务查询
 * @param query 地域查询的一些过滤条件
 */
export async function fetchProjectUserInfo(query: QueryState<ProjectFilter>) {
  let projectResourceInfo: ResourceInfo = resourceConfig().auth_project;
  let url = reduceK8sRestfulPath({ resourceInfo: projectResourceInfo });
  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  let projectUserMap = {};
  try {
    if (response.code === 0) {
      let listItems = response.data;
      if (listItems.items) {
        listItems.items.forEach(item => {
          let userInfo = item.members
            ? Object.keys(item.members).map(key => ({
                id: key,
                username: item.members[key]
              }))
            : [];
          projectUserMap[item.metadata.name] = userInfo;
        });
      }
    }
  } catch (error) {
    // 这里是搜索的时候，如果搜索不到的话，会报404的错误，只有在 resourceNotFound的时候，不把错误抛出去
    if (+error.response.status !== 404) {
      throw error;
    }
  }

  return projectUserMap;
}
/**
 * 添加已有业务为子业务
 * @param query
 */
export async function addExistMultiProject(projects: Project[], parentProjectName: string) {
  let resourceInfo = resourceConfig().projects;
  try {
    let requests = projects.map(async item => {
      let url = reduceK8sRestfulPath({ resourceInfo, specificName: item.metadata.name });
      let method = Method.patch;
      let param = {
        method,
        url,
        userDefinedHeader: {
          'Content-Type': 'application/strategic-merge-patch+json'
        },
        data: {
          spec: {
            parentProjectName
          }
        }
      };
      let response = reduceNetworkRequest(param);
      return response;
    });
    // 构建参数
    let response = await Promise.all(requests);
    if (response.every(r => r.code === 0)) {
      return operationResult(projects);
    } else {
      return operationResult(projects, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(projects, reduceNetworkWorkflow(error));
  }
}

export async function deleteParentProject(projects: Project[]) {
  let resourceInfo = resourceConfig().projects;
  try {
    let url = reduceK8sRestfulPath({ resourceInfo, specificName: projects[0].metadata.name });
    let method = Method.patch;
    let param = {
      method,
      url,
      userDefinedHeader: {
        'Content-Type': 'application/strategic-merge-patch+json'
      },
      data: {
        spec: {
          parentProjectName: null
        }
      }
    };
    // 构建参数s
    let response = await reduceNetworkRequest(param);
    if (response.code === 0) {
      return operationResult(projects);
    } else {
      return operationResult(projects, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(projects, reduceNetworkWorkflow(error));
  }
}

export async function fetchNamespaceKubectlConfig(query: QueryState<NamespaceFilter>) {
  let {
    filter: { projectId, np }
  } = query;
  let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
  let url = reduceK8sRestfulPath({
    resourceInfo: NamespaceResourceInfo,
    specificName: projectId,
    extraResource: `namespaces/${np}/certificate`
  });

  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };
  let result = {
    certPem: '',
    keyPem: '',
    caCertPem: '',
    apiServer: ''
  };
  try {
    let response = await reduceNetworkRequest(params);
    if (response.code === 0 && response.data.status.certificate) {
      result = {
        certPem: response.data.status.certificate.certPem,
        keyPem: response.data.status.certificate.keyPem,
        caCertPem: response.data.status.certificate.caCertPem,
        apiServer: response.data.status.certificate.apiServer
      };
    }
  } catch (error) {}

  return result;
}

export async function migrateNamesapce(namespaces: Namespace[], options: NamespaceOperator) {
  try {
    let { projectId, desProjectId } = options;
    let NamespaceResourceInfo: ResourceInfo = resourceConfig().namespaces;
    let url = reduceK8sRestfulPath({
      resourceInfo: NamespaceResourceInfo,
      specificName: projectId,
      extraResource: `nsemigrations`
    });

    let method = Method.post;
    let param = {
      method,
      url,
      data: {
        kind: 'NsEmigration',
        apiVersion: 'business.tkestack.io/v1',
        metadata: {
          namespace: projectId
        },
        spec: {
          namespace: namespaces[0].metadata.name,
          nsShowName: namespaces[0].metadata.namespace,
          destination: desProjectId
        }
      }
    };
    // 构建参数
    let response = await reduceNetworkRequest(param);
    if (response.code === 0) {
      return operationResult(namespaces);
    } else {
      return operationResult(namespaces, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(namespaces, reduceNetworkWorkflow(error));
  }
}
/**
 * 集群列表的查询
 * @param query 集群列表查询的一些过滤条件
 */
export async function fetchUserManagedProjects(query: QueryState<UserManagedProjectFilter>) {
  let { userId } = query.filter;
  let userResourceInfo: ResourceInfo = resourceConfig().user;
  let url = reduceK8sRestfulPath({
    resourceInfo: userResourceInfo,
    specificName: userId,
    extraResource: 'projects'
  });

  /** 构建参数 */
  let method = 'GET';
  let params: RequestParams = {
    method,
    url
  };

  let response = await reduceNetworkRequest(params);
  let managedProjects = [];
  if (response.code === 0) {
    let list = response.data;
    managedProjects = Object.keys(list.managedProjects).map(item => ({ id: item, name: item }));
  }

  const result: RecordSet<UserManagedProject> = {
    recordCount: managedProjects.length,
    records: managedProjects
  };

  return result;
}
export async function fetchUserId(query: QueryState<string>) {
  let infoResourceInfo: ResourceInfo = resourceConfig()['info'];
  let url = reduceK8sRestfulPath({ resourceInfo: infoResourceInfo });
  let params: RequestParams = {
    method: Method.get,
    url
  };
  let result;
  try {
    let response = await reduceNetworkRequest(params);
    result = response.data;
  } catch (error) {}

  return result;
}
