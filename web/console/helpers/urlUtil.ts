import { remove, isEmpty } from '../src/modules/common/utils';
import { ResourceInfo } from '../src/modules/common';
import { resourceConfig } from '../config';

export function parseQueryString(str: string = '') {
  let result = {};
  str
    .replace(/^\?*/, '')
    .split('&')
    .forEach(item => {
      let keyVal = item.split('=');
      if (keyVal.length > 0) {
        let key = decodeURIComponent(keyVal[0]);
        result[key] = keyVal[1] ? decodeURIComponent(keyVal[1]) : '';
      }
    });
  return result;
}

export function buildQueryString(obj: any = {}) {
  let keys = remove(Object.keys(obj), value => value === '');
  let queryStr = keys.map(key => `${encodeURIComponent(key)}=${encodeURIComponent(obj[key])}`).join('&');

  if (queryStr) {
    return '?' + queryStr;
  } else {
    return '';
  }
}

/**
 * 用于获取queryString
 * @param k8sQueryObj
 *  eg: ?fieldSelector=involvedObject.name=*,involvedObject.kind=*&limit=1
 *  传入进来的结构
 *  {
 *      fieldSelector: {
 *          involvedObject.name: *,
 *          involvedObject.kind: *
 *      },
 *      limit: 1
 *  }
 * @param options: K8sRestfulPathOptions
 */
export const reduceK8sQueryString = ({
  k8sQueryObj = {},
  restfulPath = ''
}: {
  k8sQueryObj: any;
  restfulPath?: string;
}) => {
  let operator = '?';
  let queryString = '';
  if (!isEmpty(k8sQueryObj)) {
    let queryKeys = Object.keys(k8sQueryObj);
    queryKeys.forEach((queryKey, index) => {
      if (index !== 0) {
        queryString += '&';
      }

      // 这里去判断每种资源的query，eg：fieldSelector、limit等
      let specificQuery = k8sQueryObj[queryKey];

      if (typeof specificQuery === 'object') {
        // 这里是对于 query的字段里面，还有多种过滤条件，比如fieldSelector支持 involvedObject.name=*,involvedObject.kind=*
        let specificKeys = Object.keys(specificQuery),
          specificString = '';
        specificKeys.forEach((speKey, index) => {
          if (index !== 0) {
            specificString += ',';
          }
          specificString += speKey + '=' + specificQuery[speKey];
        });
        if (specificString) {
          queryString += queryKey + '=' + specificString;
        }
      } else {
        queryString += queryKey + '=' + k8sQueryObj[queryKey];
      }
    });
  }

  /** 如果原本的url里面已经有 ? 了，则我们这里的query的内容，必须是拼接在后面，而不能直接加多一个 ? */
  if (restfulPath.includes('?')) {
    operator = '&';
  }

  return queryString ? `${operator}${queryString}` : '';
};

interface K8sRestfulPathOptions {
  /** 资源的配置 */
  resourceInfo: ResourceInfo;

  /** 命名空间，具体的ns */
  namespace?: string;

  isSpetialNamespace?: boolean;

  /** 不在路径最后的变量，比如projectId*/
  middleKey?: string;

  /** 具体的资源名称 */
  specificName?: string;

  /** 某个具体资源下的子资源，eg: deployment/---/pos */
  extraResource?: string;

  /** 集群id，适用于addon 请求平台转发的场景 */
  clusterId?: string;

  /** 集群logAgentName */
  logAgentName?: string;

  meshId?: string;
}

/**
 * 获取k8s 的restful 风格的path
 * @param resourceInfo: ResourceInfo  资源的配置
 * @param namespace: string 具体的命名空间
 * @param specificName: string  具体的资源名称
 * @param extraResource: string 某个具体资源下的子资源
 * @param clusterId: string 集群id，适用于addon 请求平台转发的场景
 */
export const reduceK8sRestfulPath = (options: K8sRestfulPathOptions) => {
  let {
    resourceInfo,
    namespace = '',
    isSpetialNamespace = false,
    specificName = '',
    extraResource = '',
    clusterId = '',
    logAgentName = '',
    meshId
  } = options;

  namespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
  let url: string = '';
  let isAddon = resourceInfo.requestType && resourceInfo.requestType.addon ? resourceInfo.requestType.addon : false;

  /**
   * addon 和 非 addon的资源，请求的url 不太一样
   * addon:
   *  1. 如果包含有extraResource（目前仅支持 events 和 pods）
   *     => /apis/platfor.tke/v1/clusters/${cluster}/${addonNames}${extraResource}?message={namespace}&reason={specificName}
   *  2. 如果不包含extraResource
   *     => /apis/platfor.tke/v1/clusters/${cluster}/${addonNames}?namespace={namespace}&name={specificName}
   *
   * 非addon（以deployment为例):  /apis/apps/v1beta2/namespaces/${namespace}/deployments/${deployment}/${extraResource}
   */
  if (isAddon) {
    // 兼容新旧日志组件
    let baseInfo: ResourceInfo = resourceConfig()[logAgentName ? 'logagent' : 'cluster'];
    let baseValue = logAgentName || clusterId;
    url = `/${baseInfo.basicEntry}/${baseInfo.group}/${baseInfo.version}/${baseInfo.requestType['list']}/${baseValue}/${resourceInfo.requestType['list']}`;

    if (extraResource || resourceInfo['namespaces'] || specificName) {
      let queryArr: string[] = [];
      resourceInfo.namespaces && namespace && queryArr.push(`namespace=${namespace}`);
      specificName && queryArr.push(`name=${specificName}`);
      extraResource && queryArr.push(`action=${extraResource}`);
      url += `?${queryArr.join('&')}`;
    }
  } else {
    url =
      `/${resourceInfo.basicEntry}/` +
      (resourceInfo.group ? `${resourceInfo.group}/` : '') +
      `${resourceInfo.version}/` +
      (resourceInfo.namespaces ? `${resourceInfo.namespaces}/${namespace}/` : '') +
      `${resourceInfo.requestType.list}` +
      (specificName ? `/${specificName}` : '') +
      (extraResource ? `/${extraResource}` : '');
  }
  return url;
};

export function reduceNs(namesapce) {
  let newNs = namesapce;
  /// #if project
  //业务侧ns eg: cls-xxx-ns 需要去除前缀
  if (newNs) {
    newNs = newNs.startsWith('global') ? newNs.split('-').splice(1).join('-') : newNs.split('-').splice(2).join('-');
  }
  /// #endif
  return newNs;
}

export function reverseReduceNs(clusterId: string, namespace: string) {
  let newNs = namespace;
  /// #if project
  //业务侧ns eg: cls-xxx-ns 需要去除前缀
  if (newNs) {
    newNs = `${clusterId}-${newNs}`;
  }
  /// #endif
  return newNs;
}
