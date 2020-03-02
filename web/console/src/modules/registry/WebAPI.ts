import { collectionPaging, RecordSet, uuid } from '@tencent/qcloud-lib';
import { QueryState } from '@tencent/qcloud-redux-query';
import { OperationResult } from '@tencent/qcloud-redux-workflow';

import { resourceConfig } from '../../../config/resourceConfig';
import { reduceK8sRestfulPath, reduceNetworkRequest, reduceNetworkWorkflow } from '../../../helpers';
import { Method } from '../../../helpers/reduceNetwork';
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
  ChartIns,
  ChartInsFilter,
  Image,
  ImageCreation,
  ImageFilter
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
  let apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
  let url = reduceK8sRestfulPath({
    resourceInfo: apiKeyResourceInfo
  });

  let params: RequestParams = {
    method: Method.get,
    url
  };

  let response = await reduceNetworkRequest(params);
  let apiKeyList = [];
  try {
    if (response.code === 0) {
      let listItems = response.data;
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
    let apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
    let url = reduceK8sRestfulPath({
      resourceInfo: apiKeyResourceInfo
    });

    let apiKey = apiKeys[0];
    /** 构建参数 */
    let requestParams = {
      description: apiKey.description,
      expire: apiKey.expire + apiKey.unit
    };
    let params: RequestParams = {
      method: Method.post,
      url: url + '/default/token',
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);
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
    let apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
    let url = reduceK8sRestfulPath({
      resourceInfo: apiKeyResourceInfo,
      specificName: apiKeys[0].metadata.name
    });

    let params: RequestParams = {
      method: Method.delete,
      url
    };

    let response = await reduceNetworkRequest(params);
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
    let apiKeyResourceInfo: ResourceInfo = resourceConfig()['apiKey'];
    let url = reduceK8sRestfulPath({
      resourceInfo: apiKeyResourceInfo,
      specificName: apiKeys[0].metadata.name
    });

    apiKeys[0].status = Object.assign({}, apiKeys[0].status, {
      disabled: !apiKeys[0].status.disabled
    });

    let requestParams = apiKeys[0];
    let params: RequestParams = {
      method: Method.put,
      url: url,
      data: requestParams
    };

    let response = await reduceNetworkRequest(params);
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
  let { search, paging } = query;

  let params: RequestParams = {
    method: Method.get,
    url: REPO_URL
  };

  let response = await reduceNetworkRequest(params);
  let repoList = [],
    total = 0;
  try {
    if (response.code === 0) {
      let listItems = response.data;
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
    let repo = repos[0];
    /** 构建参数 */
    let requestParams = {
      apiVersion: 'registry.tkestack.io/v1',
      kind: 'Namespace',
      spec: {
        displayName: repo.displayName,
        name: repo.name,
        visibility: repo.visibility || 'Public'
      }
    };
    let params: RequestParams = {
      method: Method.post,
      url: REPO_URL,
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);
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
    let params: RequestParams = {
      method: Method.delete,
      url: REPO_URL + repos[0].metadata.name
    };

    let response = await reduceNetworkRequest(params);
    if (response.code === 0) {
      return operationResult(repos);
    } else {
      return operationResult(repos, response);
    }
  } catch (error) {
    throw reduceNetworkWorkflow(error);
  }
}

/** Chart Group */
export async function fetchChartList(query: QueryState<ChartFilter>) {
  let { search, paging } = query;

  let params: RequestParams = {
    method: Method.get,
    url: CHART_URL
  };

  let response = await reduceNetworkRequest(params);
  let chartList = [],
    total = 0;
  try {
    if (response.code === 0) {
      let listItems = response.data;
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

  const result: RecordSet<Repo> = {
    recordCount: total,
    records: chartList
  };

  return result;
}

export async function createChart(charts: ChartCreation[]) {
  try {
    let chart = charts[0];
    /** 构建参数 */
    let requestParams = {
      apiVersion: 'registry.tkestack.io/v1',
      kind: 'ChartGroup',
      spec: {
        displayName: chart.displayName,
        name: chart.name,
        visibility: chart.visibility || 'Public'
      }
    };
    let params: RequestParams = {
      method: Method.post,
      url: CHART_URL,
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);
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
    let params: RequestParams = {
      method: Method.delete,
      url: CHART_URL + charts[0].metadata.name
    };

    let response = await reduceNetworkRequest(params);
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
  let { search, paging, filter } = query;

  let params: RequestParams = {
    method: Method.get,
    url: `${REPO_URL}${filter.chartgroup}/charts`
  };

  let response = await reduceNetworkRequest(params);
  let chartList = [],
    total = 0;
  try {
    if (response.code === 0) {
      let listItems = response.data;
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
  let { search, paging, filter } = query;

  let params: RequestParams = {
    method: Method.get,
    url: `${REPO_URL}${filter.namespace}/repositories`
  };

  let response = await reduceNetworkRequest(params);
  let imageList = [],
    total = 0;
  try {
    if (response.code === 0) {
      let listItems = response.data;
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
    let image = images[0];
    /** 构建参数 */
    let requestParams = {
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
    let params: RequestParams = {
      method: Method.post,
      url: `${REPO_URL}${image.namespace}/repositories`,
      data: requestParams
    };
    let response = await reduceNetworkRequest(params);
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
    let image = images[0];
    let params: RequestParams = {
      method: Method.delete,
      url: `${REPO_URL}${image.metadata.namespace}/repositories/${image.metadata.name}`
    };

    let response = await reduceNetworkRequest(params);
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
  let _localUrl = localStorage.getItem('_registry_url');
  if (_localUrl && !/^https:\/\/.*/.test(_localUrl)) {
    return _localUrl;
  } else {
    let paramsInfo: RequestParams = {
      method: Method.get,
      url: '/apis/gateway.tkestack.io/v1/tokens/info'
    };
    let paramsSysInfo: RequestParams = {
      method: Method.get,
      url: '/apis/gateway.tkestack.io/v1/sysinfo'
    };

    let info = await reduceNetworkRequest(paramsInfo);
    let sysInfo = await reduceNetworkRequest(paramsSysInfo);
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
