'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

function _interopDefault (ex) { return (ex && (typeof ex === 'object') && 'default' in ex) ? ex['default'] : ex; }

var teaApp = require('@tencent/tea-app');
require('@tencent/tea-component/lib/tea.css');
var tslib = require('tslib');
var React = require('react');
var React__default = _interopDefault(React);
var i18n = require('@tencent/tea-app/lib/i18n');
var ffValidator = require('@tencent/ff-validator');
var teaComponent = require('@tencent/tea-component');
var reactRedux = require('react-redux');
var ffRedux = require('@tencent/ff-redux');
var ffComponent = require('@tencent/ff-component');
var ModalMain = require('@tencent/tea-component/lib/modal/ModalMain');
var redux = require('redux');
var Cookies = _interopDefault(require('js-cookie'));
var SparkMD5 = _interopDefault(require('spark-md5'));
var axios = _interopDefault(require('axios'));
var jsBase64 = require('js-base64');
var bridge = require('@tencent/tea-app/lib/bridge');
var teaComponent$1 = require('tea-component');
var InputPassword = require('@tencent/tea-component/lib/input/InputPassword');
var reduxLogger = require('redux-logger');
var thunk = _interopDefault(require('redux-thunk'));

/**
 * @fileoverview
 *
 * 本文件词条由 `tea scan` 命令扫描生成，可提交给翻译团队进行翻译。
 * 具体使用方法，请参考 http://tapd.oa.com/QCloud_2015/markdown_wikis/view/#1010103951008365523
 */

/**
 * @type {import('@tencent/tea-app').I18NTranslation}
 */
var translation = {
  // 使用 `tea scan` 命令扫描词条
};
var zh = {
  translation: translation
};
var zh_1 = zh.translation;

// 导入依赖
// 初始化国际化词条
teaApp.i18n.init({
  translation: zh_1
});

var UrlParams;
(function (UrlParams) {
  var Sub;
  (function (Sub) {
    Sub["redirect"] = "redirect";
    Sub["startUp"] = "startUp";
    Sub["sub"] = "sub";
    Sub["create"] = "create";
    Sub["addExist"] = "addExist";
    Sub["upgrade"] = "upgrade";
    Sub["upgradeMaster"] = "upgradeMaster";
  })(Sub = UrlParams.Sub || (UrlParams.Sub = {}));
  var Mode;
  (function (Mode) {
    Mode["list"] = "list";
    Mode["detail"] = "detail";
    Mode["create"] = "create";
    Mode["addnode"] = "addnode";
    Mode["modify"] = "modify";
    Mode["apply"] = "apply";
    Mode["createnode"] = "createnode";
    Mode["update"] = "update";
  })(Mode = UrlParams.Mode || (UrlParams.Mode = {}));
})(UrlParams || (UrlParams = {}));

function _typeof(obj) {
  "@babel/helpers - typeof";

  return _typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function (obj) {
    return typeof obj;
  } : function (obj) {
    return obj && "function" == typeof Symbol && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;
  }, _typeof(obj);
}

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
var isEmpty = function isEmpty(value) {
  if (Array.isArray(value)) {
    //value为数组
    return !value.length;
  } else if (_typeof(value) === 'object') {
    //value为对象
    if (value === null) {
      //value为null
      return true;
    } else {
      //value是否没有key
      return !Object.keys(value).length;
    }
  } else if (typeof value === 'undefined') {
    //value为undefinded
    return true;
  } else if (Number.isFinite(value)) {
    //value为数值
    return false;
  } else {
    //value为默认值
    return !value;
  }
};

/*
 * @File: 这是文件的描述
 * @Description: 这是文件的描述
 * @Version: 1.0
 * @Autor: brycewwang
 * @Date: 2022-06-14 22:27:26
 * @LastEditors: brycewwang
 * @LastEditTime: 2022-06-14 22:27:26
 */
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
/**
 * @returns 删除数据后的新数组，改变原数组
 */
var remove = function remove(arr, func) {
  var targets = arr.filter(func);
  targets.forEach(function (elem) {
    arr.forEach(function (value, index) {
      if (JSON.stringify(value) === JSON.stringify(elem)) {
        arr.splice(index, 1);
      }
    });
  });
  return arr;
};

function parseQueryString(str) {
  if (str === void 0) {
    str = '';
  }
  var result = {};
  str.replace(/^\?*/, '').split('&').forEach(function (item) {
    var keyVal = item.split('=');
    if (keyVal.length > 0) {
      var key = decodeURIComponent(keyVal[0]);
      result[key] = keyVal[1] ? decodeURIComponent(keyVal[1]) : '';
    }
  });
  return result;
}
function buildQueryString(obj) {
  if (obj === void 0) {
    obj = {};
  }
  var keys = remove(Object.keys(obj), function (value) {
    return value === '';
  });
  var queryStr = keys.map(function (key) {
    return "".concat(encodeURIComponent(key), "=").concat(encodeURIComponent(obj[key]));
  }).join('&');
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
var reduceK8sQueryString = function reduceK8sQueryString(_a) {
  var _b;
  var _c = _a.k8sQueryObj,
    k8sQueryObj = _c === void 0 ? {} : _c,
    _d = _a.restfulPath;
  var operator = '?';
  var queryString = '';
  if (!isEmpty(k8sQueryObj)) {
    var queryKeys = (_b = Object.keys(k8sQueryObj)) === null || _b === void 0 ? void 0 : _b.filter(function (key) {
      return key === 'labelSelector';
    });
    queryKeys.forEach(function (queryKey, index) {
      if (index !== 0) {
        queryString += '&';
      }
      // 这里去判断每种资源的query，eg：fieldSelector、limit等
      var specificQuery = k8sQueryObj[queryKey];
      if (_typeof(specificQuery) === 'object') {
        // 这里是对于 query的字段里面，还有多种过滤条件，比如fieldSelector支持 involvedObject.name=*,involvedObject.kind=*
        var specificKeys = Object.keys(specificQuery),
          specificString_1 = '';
        specificKeys.forEach(function (speKey, index) {
          if (index !== 0) {
            specificString_1 += ',';
          }
          specificString_1 += speKey + '=' + specificQuery[speKey];
        });
        if (specificString_1) {
          queryString += queryKey + '=' + specificString_1;
        }
      } else {
        queryString += queryKey + '=' + k8sQueryObj[queryKey];
      }
    });
  }
  /** 如果原本的url里面已经有 ? 了，则我们这里的query的内容，必须是拼接在后面，而不能直接加多一个 ? */
  // if (restfulPath.includes('?')) {
  //   operator = '&';
  // }
  return queryString ? "".concat(operator).concat(queryString) : '';
};
// interface K8sRestfulPathOptions {
//   /** 资源的配置 */
//   // resourceInfo: ResourceInfo;
//   resourceInfo: any;
//   /** 命名空间，具体的ns */
//   namespace?: string;
//   /** 业务视图是否切分namespace */
//   isSpecialNamespace?: boolean;
//   /** 不在路径最后的变量，比如projectId*/
//   middleKey?: string;
//   /** 具体的资源名称 */
//   specificName?: string;
//   /** 某个具体资源下的子资源，eg: deployment/---/pos */
//   extraResource?: string;
//   /** 集群id，适用于addon 请求平台转发的场景 */
//   clusterId?: string;
//   /** 集群logAgentName */
//   logAgentName?: string;
//   meshId?: string;
// }
// /**
//  * 获取k8s 的restful 风格的path
//  * @param resourceInfo: ResourceInfo  资源的配置
//  * @param namespace: string 具体的命名空间
//  * @param specificName: string  具体的资源名称
//  * @param extraResource: string 某个具体资源下的子资源
//  * @param clusterId: string 集群id，适用于addon 请求平台转发的场景
//  */
// export const reduceK8sRestfulPath = (options: K8sRestfulPathOptions) => {
//   let {
//     resourceInfo,
//     namespace = '',
//     isSpecialNamespace = false,
//     specificName = '',
//     extraResource = '',
//     clusterId = '',
//     meshId = '',
//     logAgentName = ''
//   } = options;
//   namespace = namespace.replace(new RegExp(`^${clusterId}-`), '');
//   let url = '';
//   const isAddon = resourceInfo.requestType && resourceInfo.requestType.addon ? resourceInfo.requestType.addon : false;
//   /**
//    * addon 和 非 addon的资源，请求的url 不太一样
//    * addon:
//    *  1. 如果包含有extraResource（目前仅支持 events 和 pods）
//    *     => /apis/platfor.tke/v1/clusters/${cluster}/${addonNames}${extraResource}?message={namespace}&reason={specificName}
//    *  2. 如果不包含extraResource
//    *     => /apis/platfor.tke/v1/clusters/${cluster}/${addonNames}?namespace={namespace}&name={specificName}
//    *
//    * 非addon（以deployment为例):  /apis/apps/v1beta2/namespaces/${namespace}/deployments/${deployment}/${extraResource}
//    */
//   if (isAddon) {
//     // 兼容新旧日志组件
//     // const baseInfo: ResourceInfo = resourceConfig()[logAgentName ? 'logagent' : 'cluster'];
//     const baseInfo = '';
//     const baseValue = logAgentName || clusterId;
//     url = `/${baseInfo.basicEntry}/${baseInfo.group}/${baseInfo.version}/${baseInfo.requestType['list']}/${baseValue}/${resourceInfo.requestType['list']}`;
//     if (extraResource || resourceInfo['namespaces'] || specificName) {
//       const queryArr: string[] = [];
//       resourceInfo.namespaces && namespace && queryArr.push(`namespace=${namespace}`);
//       specificName && queryArr.push(`name=${specificName}`);
//       extraResource && queryArr.push(`action=${extraResource}`);
//       url += `?${queryArr.join('&')}`;
//     }
//   } else {
//     url =
//       `/${resourceInfo.basicEntry}/` +
//       (resourceInfo.group ? `${resourceInfo.group}/` : '') +
//       `${resourceInfo.version}/` +
//       (resourceInfo.namespaces ? `${resourceInfo.namespaces}/${namespace}/` : '') +
//       `${resourceInfo.requestType.list}` +
//       (specificName ? `/${specificName}` : '') +
//       (extraResource ? `/${extraResource}` : '');
//   }
//   return url;
// };
// export function cutNsStartClusterId({ namespace, clusterId }) {
//   return namespace.replace(new RegExp(`^${clusterId}-`), '');
// }
// export function reduceNs(namesapce) {
//   let newNs = namesapce;
//   /// #if project
//   //业务侧ns eg: cls-xxx-ns 需要去除前缀
//   if (newNs) {
//     newNs = newNs.startsWith('global') ? newNs.split('-').splice(1).join('-') : newNs.split('-').splice(2).join('-');
//   }
//   /// #endif
//   return newNs;
// }
// export function reverseReduceNs(clusterId: string, namespace: string) {
//   let newNs = namespace;
//   /// #if project
//   //业务侧ns eg: cls-xxx-ns 需要去除前缀
//   if (newNs) {
//     newNs = `${clusterId}-${newNs}`;
//   }
//   /// #endif
//   return newNs;
// }

var Cookie;
(function (Cookie) {
  Cookie.getCookie = function (name) {
    var _a, _b;
    var result = '';
    if (document === null || document === void 0 ? void 0 : document.cookie) {
      (_b = (_a = document === null || document === void 0 ? void 0 : document.cookie) === null || _a === void 0 ? void 0 : _a.split(';')) === null || _b === void 0 ? void 0 : _b.forEach(function (item) {
        var _a = item === null || item === void 0 ? void 0 : item.split('='),
          key = _a[0],
          value = _a[1];
        if ((key === null || key === void 0 ? void 0 : key.trim()) === (name === null || name === void 0 ? void 0 : name.trim())) {
          result = value === null || value === void 0 ? void 0 : value.trim();
        }
      });
    } else {
      result = '';
    }
    return result;
  };
})(Cookie || (Cookie = {}));

var Util;
(function (Util) {
  Util.TkeStackDefaultClusterId = 'global';
  Util.getCreator = function (platform, resource) {
    var _a, _b, _c, _d, _e, _f;
    var name = (_f = (_c = (_b = (_a = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _a === void 0 ? void 0 : _a.annotations) === null || _b === void 0 ? void 0 : _b['ssm.infra.tce.io/creator']) !== null && _c !== void 0 ? _c : (_e = (_d = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _d === void 0 ? void 0 : _d.labels) === null || _e === void 0 ? void 0 : _e['ssm.infra.tce.io/creator']) !== null && _f !== void 0 ? _f : '-';
    return name;
  };
  Util.getUserName = function () {
    var _a, _b, _c;
    var name = '';
    if (document.cookie) {
      try {
        name = (_c = (_b = (_a = decodeURIComponent(Cookie === null || Cookie === void 0 ? void 0 : Cookie.getCookie('nick'))) === null || _a === void 0 ? void 0 : _a.split('@')) === null || _b === void 0 ? void 0 : _b[0]) === null || _c === void 0 ? void 0 : _c.toString();
      } catch (error) {
        name = '';
      }
    }
    return name;
  };
  /**
   * 查询云上平台uin
   * @param platform
   * @returns
   */
  Util.getUin = function () {
    var _a, _b;
    var result = '';
    if (document.cookie) {
      try {
        result = (_b = (_a = decodeURIComponent(Cookie === null || Cookie === void 0 ? void 0 : Cookie.getCookie('uin'))) === null || _a === void 0 ? void 0 : _a.slice(1)) === null || _b === void 0 ? void 0 : _b.toString();
      } catch (error) {
        result = '';
      }
    }
    return result;
  };
  Util.getReadableFileSizeString = function (fileSizeInBytes) {
    if (fileSizeInBytes === undefined) {
      return '-';
    }
    var i = -1;
    var byteUnits = [' kB', ' MB', ' GB', ' TB', 'PB', 'EB', 'ZB', 'YB'];
    var size = fileSizeInBytes;
    do {
      size = size / 1024;
      i = i + 1;
    } while (size > 1024);
    return Math.max(size, 0.1).toFixed(1) + byteUnits[i];
  };
  Util.getCOSClusterId = function (platform, hubCluster, cluster) {
    var clusterId = '';
    switch (platform) {
      case PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TKESTACK:
        clusterId = (cluster === null || cluster === void 0 ? void 0 : cluster.clusterId) || Util.TkeStackDefaultClusterId;
        break;
      case PlatformType.TDCC:
      default:
        clusterId = hubCluster === null || hubCluster === void 0 ? void 0 : hubCluster.clusterId;
        break;
    }
    return clusterId;
  };
  Util.getDefaultRegion = function (platform, route) {
    var _a, _b;
    var regionId;
    switch (platform) {
      case PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TKESTACK:
        regionId = (_b = +((_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.rid)) !== null && _b !== void 0 ? _b : HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
        break;
      case PlatformType.TDCC:
      default:
        regionId = HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
        break;
    }
    return regionId;
  };
  Util.getClusterId = function (platform, resource, route) {
    var _a, _b, _c, _d, _e;
    var clusterId;
    switch (platform) {
      case PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TKESTACK:
        // clusterId = resource?.metadata?.labels?.['tdcc.cloud.tencent.com/cluster-id'] || resource?.metadata?.labels?.['ssm.infra.tce.io/cluster-id'] || route?.queries?.clusterid || TkeStackDefaultClusterId || '-';
        clusterId = Util.TkeStackDefaultClusterId || '-';
        break;
      case PlatformType.TDCC:
      default:
        clusterId = ((_b = (_a = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['tdcc.cloud.tencent.com/cluster-id']) || ((_d = (_c = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _c === void 0 ? void 0 : _c.labels) === null || _d === void 0 ? void 0 : _d['ssm.infra.tce.io/cluster-id']) || ((_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.clusterid) || '-';
        break;
    }
    return clusterId;
  };
  Util.getClusterName = function (platform, resource, route) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
    var clusterName;
    switch (platform) {
      case PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TKESTACK:
        clusterName = ((_b = (_a = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['tdcc.cloud.tencent.com/cluster-name']) || ((_d = (_c = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _c === void 0 ? void 0 : _c.labels) === null || _d === void 0 ? void 0 : _d['ssm.infra.tce.io/cluster-name']) || ((_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.clusterid) || '-';
        break;
      case PlatformType.TDCC:
      default:
        clusterName = ((_g = (_f = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _f === void 0 ? void 0 : _f.labels) === null || _g === void 0 ? void 0 : _g['tdcc.cloud.tencent.com/cluster-name']) || ((_j = (_h = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _h === void 0 ? void 0 : _h.labels) === null || _j === void 0 ? void 0 : _j['ssm.infra.tce.io/cluster-name']) || ((_k = route === null || route === void 0 ? void 0 : route.queries) === null || _k === void 0 ? void 0 : _k.clusterid) || '-';
        break;
    }
    return clusterName;
  };
  Util.getInstanceId = function (platform, instanceName, serviceResource) {
    var _a, _b;
    var result = !instanceName ? ((_a = serviceResource === null || serviceResource === void 0 ? void 0 : serviceResource.spec) === null || _a === void 0 ? void 0 : _a.externalID) ? (_b = serviceResource === null || serviceResource === void 0 ? void 0 : serviceResource.spec) === null || _b === void 0 ? void 0 : _b.externalID : '' : DefaultNamespace + '-' + instanceName;
    return result;
  };
  Util.getRouterPath = function (pathname) {
    var _a;
    var path;
    if (pathname === null || pathname === void 0 ? void 0 : pathname.includes((_a = PlatformType.TKESTACK) === null || _a === void 0 ? void 0 : _a.toLowerCase())) {
      path = "/tkestack/middleware";
    } else {
      path = "/tdcc/middleware";
    }
    return path;
  };
})(Util || (Util = {}));

var tips = seajs.require('tips');
function downloadText(crtText, filename, contentType) {
  if (contentType === void 0) {
    contentType = 'text/plain;charset=utf-8;';
  }
  var userAgent = navigator.userAgent;
  if (navigator === null || navigator === void 0 ? void 0 : navigator['msSaveBlob']) {
    var blob = new Blob([crtText], {
      type: contentType
    });
    navigator === null || navigator === void 0 ? void 0 : navigator['msSaveBlob'](blob, filename);
  } else if (userAgent.indexOf('MSIE 9.0') > 0) {
    tips.error(i18n.t('该浏览器暂不支持导出功能'));
  } else {
    var blob = new Blob([crtText], {
      type: contentType
    });
    var link = document.createElement('a');
    if (link.download !== undefined) {
      var url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', filename);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  }
}
function downloadKubeconfig(crtText, filename) {
  if (filename === void 0) {
    filename = 'kubeconfig';
  }
  downloadText(crtText, filename, 'applicatoin/octet-stream;charset=utf-8;');
}

var _a, _b;
var ResourceTypeEnum;
(function (ResourceTypeEnum) {
  ResourceTypeEnum["ServiceBinding"] = "ServiceBinding";
  ResourceTypeEnum["Backup"] = "ServiceOpsPlan";
  ResourceTypeEnum["ServiceResource"] = "ServiceInstance";
  ResourceTypeEnum["ServicePlan"] = "ServicePlan";
  ResourceTypeEnum["ServiceOpsBackup"] = "ServiceOpsBackup";
  ResourceTypeEnum["Secret"] = "Secret";
})(ResourceTypeEnum || (ResourceTypeEnum = {}));
var ResourceTypeMap = (_a = {}, _a[ResourceTypeEnum.ServiceBinding] = {
  path: 'servicebindings',
  title: '服务绑定',
  resourceKind: 'ServiceBinding',
  schemaKey: 'BindingCreateParamsSchema'
}, _a[ResourceTypeEnum.ServicePlan] = {
  path: 'serviceplans',
  title: '规格',
  resourceKind: 'ServicePlan',
  schemaKey: 'planSchema'
}, _a[ResourceTypeEnum.ServiceResource] = {
  path: 'serviceinstances',
  title: '实例',
  resourceKind: 'ServiceInstance',
  schemaKey: 'instanceCreateParameterSchema'
}, _a[ResourceTypeEnum.Backup] = {
  path: 'serviceopsplans',
  title: '备份',
  resourceKind: 'ServiceOpsPlan',
  schemaKey: 'instanceCreateParameterSchema'
}, _a[ResourceTypeEnum.ServiceOpsBackup] = {
  path: 'serviceopsbackups',
  title: '备份',
  resourceKind: 'ServiceOpsBackup',
  schemaKey: ''
}, _a[ResourceTypeEnum.Secret] = {
  path: 'secrets',
  title: 'Secret',
  resourceKind: 'Secret',
  schemaKey: ''
}, _a);
var serviceMngTabs = [{
  id: ResourceTypeEnum.ServiceResource,
  label: '实例管理'
}, {
  id: ResourceTypeEnum.ServicePlan,
  label: '规格管理'
}];
var CreateSpecificOperatorEnum;
(function (CreateSpecificOperatorEnum) {
  CreateSpecificOperatorEnum["CreateResource"] = "CreateResource";
  CreateSpecificOperatorEnum["BackupStrategy"] = "BackupStrategy";
  CreateSpecificOperatorEnum["BackupNow"] = "BackupNow";
  CreateSpecificOperatorEnum["CreateServiceBinding"] = "CreateServiceBinding";
})(CreateSpecificOperatorEnum || (CreateSpecificOperatorEnum = {}));
var CreateSpecificOperatorMap = (_b = {}, _b[CreateSpecificOperatorEnum.CreateResource] = {
  msg: i18n.t('新建资源成功')
}, _b[CreateSpecificOperatorEnum.BackupStrategy] = {
  msg: i18n.t('备份策略成功')
}, _b[CreateSpecificOperatorEnum.BackupNow] = {
  msg: i18n.t('立即备份成功')
}, _b[CreateSpecificOperatorEnum.CreateServiceBinding] = {
  msg: i18n.t('新建资源成功')
}, _b);
var ErrorMsgEnum = {
  COS_Resource_Not_Found: i18n.t('您尚未配置云上备份地址或者地址配置不正确，请进入服务概览-设置页面下检查配置是否正确')
};
var DefaultNamespace = 'sso';
var SystemNamespace = 'ssm';
var showResourceDeleteLoading = function showResourceDeleteLoading(resource, originResources, platform) {
  var loading = false;
  var resourceType = resource === null || resource === void 0 ? void 0 : resource.kind;
  switch (resourceType) {
    case ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource:
      loading = originResources === null || originResources === void 0 ? void 0 : originResources.some(function (item) {
        var _a, _b, _c;
        return ((_a = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name) && Util.getClusterId(platform, resource) === Util.getClusterId(platform, item) && ((_c = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _c === void 0 ? void 0 : _c.deletionTimestamp);
      });
      break;
    case ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceOpsBackup:
    case ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding:
      loading = originResources === null || originResources === void 0 ? void 0 : originResources.some(function (item) {
        var _a, _b, _c;
        return ((_a = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name) && ((_c = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _c === void 0 ? void 0 : _c.deletionTimestamp);
      });
  }
  return loading;
};
var ExcludeNamespaces = ['clusternet-system', 'ssm'];

var _a$1, _b$1, _c, _d, _e, _f;
var PlatformType;
(function (PlatformType) {
  PlatformType["TDCC"] = "tdcc";
  PlatformType["TKESTACK"] = "tkeStack";
})(PlatformType || (PlatformType = {}));
var YunApiName = {
  DescribeServiceVendors: 'DescribeServiceVendors',
  DescribeServiceResources: 'DescribeServiceInstances',
  DescribeExternalClusters: 'DescribeExternalClusters',
  DescribeHubClusters: 'DescribeHubClusters'
};
var ServiceMngYunApiName = (_a$1 = {}, _a$1[ResourceTypeEnum.ServiceResource] = 'DescribeServiceInstances', _a$1[ResourceTypeEnum.ServicePlan] = 'DescribeServicePlans', _a$1[ResourceTypeEnum.Backup] = 'ForwardRequestTDCC', _a$1[ResourceTypeEnum.ServiceBinding] = 'ForwardRequestTDCC', _a$1);
var ResourceModuleName = (_b$1 = {}, _b$1[PlatformType.TDCC] = 'tdcc', _b$1);
var ResourceVersionName = (_c = {}, _c[PlatformType.TDCC] = '2022-01-25', _c);
var ResourceApiName = (_d = {}, _d[PlatformType.TDCC] = 'ForwardRequestTDCC', _d);
var middlewareRouteRegStr = (_e = {}, _e[PlatformType.TDCC] = {
  exactMatch: 'tdcc\/middleware$',
  dimMatch: 'tdcc\/middleware'
}, _e[PlatformType.TKESTACK] = {
  exactMatch: 'middleware\/middleware$',
  dimMatch: 'middleware\/middleware'
}, _e);
var rootNodeSelector = (_f = {}, _f[PlatformType.TDCC] = '#apaas-midleware', _f[PlatformType.TKESTACK] = '#apaas-midleware', _f);

var requestResourceType;
(function (requestResourceType) {
  requestResourceType["MC"] = "MC";
})(requestResourceType || (requestResourceType = {}));
var ServiceDetailType;
(function (ServiceDetailType) {
  ServiceDetailType["ServiceClasses"] = "serviceclasses";
  ServiceDetailType["ServiceBinding"] = "servicebinding";
})(ServiceDetailType || (ServiceDetailType = {}));
var ServiceNameType;
(function (ServiceNameType) {
  ServiceNameType["ETCD"] = "etcd";
  ServiceNameType["MONGO"] = "mongo";
  ServiceNameType["Redis"] = "redis";
  ServiceNameType["Mariadb"] = "mariadb";
})(ServiceNameType || (ServiceNameType = {}));

var _a$2, _b$2, _c$1, _d$1, _e$1, _f$1;
/**
 * 类型
 */
var SchemaType;
(function (SchemaType) {
  SchemaType["Integer"] = "Integer";
  SchemaType["String"] = "String";
  SchemaType["Select"] = "Select";
  SchemaType["Boolean"] = "Boolean";
  SchemaType["CPU"] = "CPU";
  SchemaType["Storage"] = "Storage";
  SchemaType["List"] = "List";
  SchemaType["Map"] = "Map";
  SchemaType["Custom"] = "Custom";
})(SchemaType || (SchemaType = {}));
var DashboardStatusEnum;
(function (DashboardStatusEnum) {
  DashboardStatusEnum["Creating"] = "Creating";
  DashboardStatusEnum["Created"] = "Created";
  DashboardStatusEnum["Ready"] = "Ready";
  DashboardStatusEnum["Failed"] = "Failed";
})(DashboardStatusEnum || (DashboardStatusEnum = {}));
var SupportedOperationsEnum;
(function (SupportedOperationsEnum) {
  SupportedOperationsEnum["Backup"] = "Backup";
  SupportedOperationsEnum["Restore"] = "Restore";
})(SupportedOperationsEnum || (SupportedOperationsEnum = {}));
var DashboardStatusMap = (_a$2 = {}, _a$2[DashboardStatusEnum.Created] = {
  text: 'Creating',
  className: 'text-warning'
}, _a$2[DashboardStatusEnum.Created] = {
  text: 'Created',
  className: 'text-success'
}, _a$2[DashboardStatusEnum.Ready] = {
  text: 'Ready',
  className: 'text-success'
}, _a$2[DashboardStatusEnum.Failed] = {
  text: i18n.t('Failed'),
  className: 'text-danger'
}, _a$2);
//查询RSIP信息
var getRsIps = function getRsIps(NodeInfo) {
  if (!NodeInfo) {
    return [];
  }
  var mode = NodeInfo.mode,
    nodes = NodeInfo.nodes;
  switch (mode) {
    case 'VM':
      return nodes === null || nodes === void 0 ? void 0 : nodes.map(function (x) {
        var _a, _b;
        return "".concat(x === null || x === void 0 ? void 0 : x.hostNode.host, ":").concat((_b = (_a = x === null || x === void 0 ? void 0 : x.ports) === null || _a === void 0 ? void 0 : _a[0]) !== null && _b !== void 0 ? _b : '');
      });
    case 'Physical':
      return nodes === null || nodes === void 0 ? void 0 : nodes.map(function (x) {
        var _a, _b;
        return "".concat(x === null || x === void 0 ? void 0 : x.hostNode.host, ":").concat((_b = (_a = x === null || x === void 0 ? void 0 : x.ports) === null || _a === void 0 ? void 0 : _a[0]) !== null && _b !== void 0 ? _b : '');
      });
    case 'Container':
      return nodes === null || nodes === void 0 ? void 0 : nodes.map(function (x) {
        var _a, _b;
        return "".concat(x === null || x === void 0 ? void 0 : x.containerNode.host, ":").concat((_b = (_a = x === null || x === void 0 ? void 0 : x.ports) === null || _a === void 0 ? void 0 : _a[0]) !== null && _b !== void 0 ? _b : '');
      });
    default:
      return [];
  }
};
var initServiceInstanceEdit = function initServiceInstanceEdit(data) {
  return {
    formData: data === null || data === void 0 ? void 0 : data.reduce(function (pre, cur) {
      var _a, _b;
      return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = getFormItemDefaultValue(cur), _a.unitMap = tslib.__assign(tslib.__assign({}, pre === null || pre === void 0 ? void 0 : pre.unitMap), (_b = {}, _b[cur === null || cur === void 0 ? void 0 : cur.name] = cur === null || cur === void 0 ? void 0 : cur.unit, _b)), _a));
    }, {
      instanceName: '',
      timeBackup: false,
      isSetParams: false,
      unitMap: {}
    }),
    validator: data === null || data === void 0 ? void 0 : data.reduce(function (pre, cur) {
      var _a;
      return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = {
        status: 1,
        message: ''
      }, _a));
    }, {
      instanceName: {
        status: 1,
        message: ''
      },
      timeBackup: {
        status: 1,
        message: ''
      },
      isSetParams: {
        status: 1,
        message: ''
      }
    })
  };
};
var InstanceMustPropMap;
(function (InstanceMustPropMap) {
  InstanceMustPropMap["InstanceName"] = "instanceName";
  InstanceMustPropMap["ClusterId"] = "clusterId";
  InstanceMustPropMap["Version"] = "version";
  InstanceMustPropMap["Plan"] = "plan";
})(InstanceMustPropMap || (InstanceMustPropMap = {}));
var NumberAndSymbolReg = /^[a-z][0-9a-z-]{1,62}$/;
// 校验的相关配置
var InstanceBaseValidateSchema = (_b$2 = {}, _b$2[InstanceMustPropMap === null || InstanceMustPropMap === void 0 ? void 0 : InstanceMustPropMap.InstanceName] = {
  rules: [{
    type: ffValidator.RuleTypeEnum.custom,
    customFunc: function customFunc(value, store) {
      var instanceName = (store === null || store === void 0 ? void 0 : store.formData).instanceName;
      if (!instanceName || !NumberAndSymbolReg.test(instanceName)) {
        return {
          status: ffValidator.ValidatorStatusEnum.Failed,
          message: i18n.t('名称为2-63个字符，可包含数字、小写英文字以及短划线（-），且不能以短划线（-）开头')
        };
      }
      return {
        status: 1,
        message: ''
      };
    }
  }]
}, _b$2[InstanceMustPropMap === null || InstanceMustPropMap === void 0 ? void 0 : InstanceMustPropMap.ClusterId] = {
  rules: [{
    type: ffValidator.RuleTypeEnum.custom,
    customFunc: function customFunc(value, store) {
      var clusterId = (store === null || store === void 0 ? void 0 : store.formData).clusterId;
      if (!clusterId) {
        return {
          status: ffValidator.ValidatorStatusEnum.Failed,
          message: i18n.t('必须选择目标集群')
        };
      }
      return {
        status: 1,
        message: ''
      };
    }
  }]
}, _b$2[InstanceMustPropMap === null || InstanceMustPropMap === void 0 ? void 0 : InstanceMustPropMap.Version] = {
  rules: [{
    type: ffValidator.RuleTypeEnum.custom,
    customFunc: function customFunc(value, store) {
      var version = (store === null || store === void 0 ? void 0 : store.formData).version;
      if (!version) {
        return {
          status: ffValidator.ValidatorStatusEnum.Failed,
          message: i18n.t('版本不能为空')
        };
      }
      return {
        status: 1,
        message: ''
      };
    }
  }]
}, _b$2[InstanceMustPropMap === null || InstanceMustPropMap === void 0 ? void 0 : InstanceMustPropMap.Plan] = {
  rules: [{
    type: ffValidator.RuleTypeEnum.custom,
    customFunc: function customFunc(value, store) {
      var plan = (store === null || store === void 0 ? void 0 : store.formData).plan;
      if (!plan) {
        return {
          status: ffValidator.ValidatorStatusEnum.Failed,
          message: i18n.t('规格不能为空')
        };
      }
      return {
        status: 1,
        message: ''
      };
    }
  }]
}, _b$2);
var isRequireBaseProps = function isRequireBaseProps(key) {
  return ['instanceName', 'clusterId', 'version', 'plan'].includes(key);
};
var isInSchemaProps = function isInSchemaProps(key, data) {
  return data === null || data === void 0 ? void 0 : data.some(function (item) {
    return key === (item === null || item === void 0 ? void 0 : item.name);
  });
};
var validateAllProps = function validateAllProps(serviceInstanceEdit, instanceSchemas) {
  var _a;
  var vm;
  var formData = serviceInstanceEdit.formData;
  (_a = Object.keys(formData)) === null || _a === void 0 ? void 0 : _a.forEach(function (key) {
    var _a;
    var _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s;
    var validator;
    if (isRequireBaseProps(key)) {
      validator = (_d = (_c = (_b = InstanceBaseValidateSchema === null || InstanceBaseValidateSchema === void 0 ? void 0 : InstanceBaseValidateSchema[key]) === null || _b === void 0 ? void 0 : _b.rules) === null || _c === void 0 ? void 0 : _c[0]) === null || _d === void 0 ? void 0 : _d.customFunc(formData[key], serviceInstanceEdit);
    } else if (isInSchemaProps(key, instanceSchemas)) {
      var schema = instanceSchemas === null || instanceSchemas === void 0 ? void 0 : instanceSchemas.find(function (item) {
        return (item === null || item === void 0 ? void 0 : item.name) === key;
      });
      if ((_e = serviceInstanceEdit === null || serviceInstanceEdit === void 0 ? void 0 : serviceInstanceEdit.formData) === null || _e === void 0 ? void 0 : _e.isSetParams) {
        validator = validatePlanSchema(schema, formData[key] + ((_f = formData === null || formData === void 0 ? void 0 : formData['unitMap']) === null || _f === void 0 ? void 0 : _f[key]), OperatorNumberUnitReg, formData);
      } else {
        validator = {
          status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success,
          message: i18n.t('')
        };
      }
    } else if (key === null || key === void 0 ? void 0 : key.startsWith('backup')) {
      if (formData === null || formData === void 0 ? void 0 : formData['timeBackup']) {
        if ('backupDate' === key || 'backupTime' === key) {
          validator = {
            message: !!((_g = formData === null || formData === void 0 ? void 0 : formData['backupDate']) === null || _g === void 0 ? void 0 : _g.length) || !!((_h = formData === null || formData === void 0 ? void 0 : formData['backupTime']) === null || _h === void 0 ? void 0 : _h.length) ? i18n.t('') : i18n.t('备份日期和备份时间点至少选择一个'),
            status: !!((_j = formData === null || formData === void 0 ? void 0 : formData['backupDate']) === null || _j === void 0 ? void 0 : _j.length) || !!((_k = formData === null || formData === void 0 ? void 0 : formData['backupTime']) === null || _k === void 0 ? void 0 : _k.length) ? ffValidator.ValidatorStatusEnum.Success : ffValidator.ValidatorStatusEnum.Failed
          };
        } else if ('backupReserveTime' === key) {
          validator = {
            message: !!(formData === null || formData === void 0 ? void 0 : formData[key]) ? i18n.t('') : i18n.t('备份保留时间不能为空'),
            status: !!(formData === null || formData === void 0 ? void 0 : formData[key]) ? ffValidator.ValidatorStatusEnum.Success : ffValidator.ValidatorStatusEnum.Failed
          };
        } else {
          validator = {
            message: '',
            status: ffValidator.ValidatorStatusEnum.Success
          };
        }
      } else {
        validator = {
          status: ffValidator.ValidatorStatusEnum.Success,
          message: ''
        };
      }
    } else if (key === 'nodeSchedule') {
      validator = {
        status: ((_l = formData === null || formData === void 0 ? void 0 : formData['nodeSchedule']) === null || _l === void 0 ? void 0 : _l.enable) && (((_m = formData === null || formData === void 0 ? void 0 : formData['nodeSchedule']) === null || _m === void 0 ? void 0 : _m['nodeSelector']) && !((_p = (_o = formData === null || formData === void 0 ? void 0 : formData['nodeSchedule']) === null || _o === void 0 ? void 0 : _o['nodeSelector']) === null || _p === void 0 ? void 0 : _p.isValid) || ((_q = formData === null || formData === void 0 ? void 0 : formData['nodeSchedule']) === null || _q === void 0 ? void 0 : _q['nodeAffinity']) && !((_s = (_r = formData === null || formData === void 0 ? void 0 : formData['nodeSchedule']) === null || _r === void 0 ? void 0 : _r['nodeAffinity']) === null || _s === void 0 ? void 0 : _s.isValid)) ? ffValidator.ValidatorStatusEnum.Failed : ffValidator.ValidatorStatusEnum.Success,
        message: ''
      };
    } else {
      validator = {
        status: ffValidator.ValidatorStatusEnum.Success,
        message: ''
      };
    }
    vm = tslib.__assign(tslib.__assign({}, vm), (_a = {}, _a[key] = validator, _a));
  });
  return vm;
};
var _validateInstance = function _validateInstance(data, instanceSchemas) {
  var validator = validateAllProps(data, instanceSchemas);
  return Object.keys(validator).every(function (key) {
    var _a;
    return ((_a = validator === null || validator === void 0 ? void 0 : validator[key]) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success;
  });
};
var ServiceInstanceStatusEnum;
(function (ServiceInstanceStatusEnum) {
  ServiceInstanceStatusEnum["InProgress"] = "InProgress";
  ServiceInstanceStatusEnum["Ready"] = "Ready";
  ServiceInstanceStatusEnum["Failed"] = "Failed";
  ServiceInstanceStatusEnum["Unknown"] = "Unknown";
})(ServiceInstanceStatusEnum || (ServiceInstanceStatusEnum = {}));
var ServiceInstanceMap = (_c$1 = {}, _c$1[ServiceInstanceStatusEnum.InProgress] = {
  text: i18n.t('创建中'),
  className: 'text-warning'
}, _c$1[ServiceInstanceStatusEnum.Ready] = {
  text: i18n.t('运行中'),
  className: 'text-success'
}, _c$1[ServiceInstanceStatusEnum.Failed] = {
  text: i18n.t('创建失败'),
  className: 'text-danger'
}, _c$1[ServiceInstanceStatusEnum.Unknown] = {
  text: i18n.t('-'),
  className: 'text-weak'
}, _c$1);
var ServicePlanTypeEnum;
(function (ServicePlanTypeEnum) {
  ServicePlanTypeEnum["System"] = "System";
  ServicePlanTypeEnum["Custom"] = "Custom";
})(ServicePlanTypeEnum || (ServicePlanTypeEnum = {}));
var ServicePlanTypeMap = (_d$1 = {}, _d$1[ServicePlanTypeEnum.System] = {
  text: i18n.t('预设'),
  className: 'text-success'
}, _d$1[ServicePlanTypeEnum.Custom] = {
  text: i18n.t('自定义'),
  className: 'text-warning'
}, _d$1);
// export const getAuthorName =  (platform:string) => {
//   if(platform === PlatformType?.TDCC){
//     return localStorage?.getItem
//   }else{
//   }
// }
// schema相关 start
var FormItemType;
(function (FormItemType) {
  FormItemType["Select"] = "select";
  FormItemType["Input"] = "input";
  FormItemType["Switch"] = "switch";
  FormItemType["Paasword"] = "password";
  FormItemType["InputNumber"] = "inputNumber";
  FormItemType["MapField"] = "MapField";
})(FormItemType || (FormItemType = {}));
var getFormItemType = function getFormItemType(item) {
  var _a;
  var type;
  if ((_a = item === null || item === void 0 ? void 0 : item.candidates) === null || _a === void 0 ? void 0 : _a.length) {
    type = FormItemType.Select;
  } else if ((item === null || item === void 0 ? void 0 : item.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Boolean)) {
    type = FormItemType === null || FormItemType === void 0 ? void 0 : FormItemType.Switch;
  } else if ((item === null || item === void 0 ? void 0 : item.name) === FormItemType.Paasword) {
    type = FormItemType.Paasword;
  } else if ((item === null || item === void 0 ? void 0 : item.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Map)) {
    type = FormItemType.MapField;
  } else {
    type = FormItemType.Input;
  }
  return type;
};
var getFormItemDefaultValue = function getFormItemDefaultValue(item) {
  var value = '';
  if (item === null || item === void 0 ? void 0 : item.value) {
    if ((item === null || item === void 0 ? void 0 : item.type) === 'Boolean') {
      value = (item === null || item === void 0 ? void 0 : item.value) === 'true';
    } else {
      value = item === null || item === void 0 ? void 0 : item.value;
    }
  } else {
    if ((item === null || item === void 0 ? void 0 : item.type) === 'Integer') {
      value = '0';
    } else if ((item === null || item === void 0 ? void 0 : item.type) === 'Boolean') {
      value = false;
    } else if ((item === null || item === void 0 ? void 0 : item.type) === 'String') {
      value = '';
    } else {
      value = '';
    }
  }
  return value;
};
var prefixForSchema = function prefixForSchema(item, serviceName) {
  var _a;
  return serviceName && ((_a = ['size']) === null || _a === void 0 ? void 0 : _a.includes(item === null || item === void 0 ? void 0 : item.name)) ? "".concat(serviceName) : i18n.t('');
};
var hideSchema = function hideSchema(item) {
  var _a, _b, _c, _d;
  return ((_a = ['dashboardEnabled']) === null || _a === void 0 ? void 0 : _a.includes(item === null || item === void 0 ? void 0 : item.name)) || ((_b = ['spread']) === null || _b === void 0 ? void 0 : _b.includes(item === null || item === void 0 ? void 0 : item.name)) && ((_c = item === null || item === void 0 ? void 0 : item.candidates) === null || _c === void 0 ? void 0 : _c.length) === 1 && ((_d = item === null || item === void 0 ? void 0 : item.candidates) === null || _d === void 0 ? void 0 : _d[0]) === 'zone';
};
var SchemaInputNumOption = {
  min: 0
};
var OperatorNumberUnitReg = /(<?>?=?)(-?[0-9]{1,}[.]?[0-9]*e?[0-9]*)(\w*)/;
var NumberUnitReg = /(-?[0-9]{1,}[.]?[0-9]*e?[0-9]*)(\w*)/;
var compareOperator = {
  '<': {
    validate: function validate(a, b) {
      return a < b;
    },
    errMsg: i18n.t('小于')
  },
  '>': {
    validate: function validate(a, b) {
      return a > b;
    },
    errMsg: i18n.t('大于')
  },
  '<=': {
    validate: function validate(a, b) {
      return a <= b;
    },
    errMsg: i18n.t('小于等于')
  },
  '>=': {
    validate: function validate(a, b) {
      return a >= b;
    },
    errMsg: i18n.t('大于等于')
  },
  '==': {
    validate: function validate(a, b) {
      return a === b;
    },
    errMsg: i18n.t('等于')
  }
};
var storageUnitKeyMap = {
  G: 'G',
  M: 'M',
  "default": 'G'
};
var cpuUnitKeyMap = {
  C: '',
  M: 'm',
  "default": ''
};
var unitOptions = {
  CPU: [{
    text: i18n.t('核'),
    value: ''
  }, {
    text: i18n.t('毫核'),
    value: 'm'
  }],
  Storage: [{
    text: 'GiB',
    value: 'G'
  }, {
    text: 'MiB',
    value: 'M'
  }]
};
var cpuUnitMapping = {
  m: i18n.t('毫核'),
  mi: i18n.t('毫核'),
  defaultValue: i18n.t('核')
};
var storageUnitMapping = {
  g: 'GiB',
  gi: 'GiB',
  m: 'MiB',
  mi: 'MiB',
  defaultValue: 'GiB'
};
var DefaultUnitMap = (_e$1 = {}, _e$1[SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU] = cpuUnitMapping === null || cpuUnitMapping === void 0 ? void 0 : cpuUnitMapping.defaultValue, _e$1[SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage] = storageUnitMapping === null || storageUnitMapping === void 0 ? void 0 : storageUnitMapping.defaultValue, _e$1);
var suffixUnitForSchema = function suffixUnitForSchema(item) {
  var _a;
  var unit = '';
  return i18n.t('{{unit}}', {
    unit: ((_a = item === null || item === void 0 ? void 0 : item.label) === null || _a === void 0 ? void 0 : _a.includes(i18n.t('(核)'))) ? i18n.t('') : unit
  });
};
var UnitMappingList = (_f$1 = {}, _f$1[SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU] = cpuUnitMapping, _f$1[SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage] = storageUnitMapping, _f$1);
var getFormattedValue = function getFormattedValue(_a) {
  var _b, _c, _d, _e;
  var value = _a.value,
    unit = _a.unit,
    field = _a.field;
  if (((_b = [SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU]) === null || _b === void 0 ? void 0 : _b.includes(field === null || field === void 0 ? void 0 : field.type)) || ((_c = ['cpu']) === null || _c === void 0 ? void 0 : _c.some(function (item) {
    var _a, _b;
    return (_b = (_a = field === null || field === void 0 ? void 0 : field.name) === null || _a === void 0 ? void 0 : _a.toLocaleLowerCase()) === null || _b === void 0 ? void 0 : _b.includes(item);
  }))) {
    return ['m', 'mi'].some(function (x) {
      return x === (unit === null || unit === void 0 ? void 0 : unit.toLowerCase());
    }) ? Number(value) : Number(value) * 1000;
  } else if (((_d = [SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage]) === null || _d === void 0 ? void 0 : _d.includes(field === null || field === void 0 ? void 0 : field.type)) || ((_e = ['memory', 'storage']) === null || _e === void 0 ? void 0 : _e.some(function (item) {
    var _a, _b;
    return (_b = (_a = field === null || field === void 0 ? void 0 : field.name) === null || _a === void 0 ? void 0 : _a.toLocaleLowerCase()) === null || _b === void 0 ? void 0 : _b.includes(item);
  }))) {
    return ['m', 'mi'].some(function (x) {
      return x === (unit === null || unit === void 0 ? void 0 : unit.toLowerCase());
    }) ? Number(value) : Number(value) * 1000;
  } else {
    return Number(value);
  }
};
var validateFormItemEmpty = function validateFormItemEmpty(field, values) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v;
  var validator;
  var fieldValue = values === null || values === void 0 ? void 0 : values[field === null || field === void 0 ? void 0 : field.name];
  // 依赖项未启用，校验直接通过
  if (field.enabledCondition && values) {
    var _w = field === null || field === void 0 ? void 0 : field.enabledCondition.split('=='),
      conditionKey = _w[0],
      conditionValue = _w[1];
    var value = values[conditionKey];
    if (String(value) !== String(conditionValue)) {
      return {
        status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success,
        message: i18n.t('')
      };
    }
  }
  if (!field.optional && (field === null || field === void 0 ? void 0 : field.type) !== (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Boolean)) {
    validator = {
      status: isEmpty(fieldValue) || ([SchemaType.Integer, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage].includes(field === null || field === void 0 ? void 0 : field.type) || ((_a = field === null || field === void 0 ? void 0 : field.name) === null || _a === void 0 ? void 0 : _a.includes((_b = SchemaType.Storage) === null || _b === void 0 ? void 0 : _b.toLocaleLowerCase())) || ((_c = field === null || field === void 0 ? void 0 : field.name) === null || _c === void 0 ? void 0 : _c.includes((_d = SchemaType.CPU) === null || _d === void 0 ? void 0 : _d.toLocaleLowerCase())) || ((_e = field === null || field === void 0 ? void 0 : field.name) === null || _e === void 0 ? void 0 : _e.includes('memory'))) && +fieldValue < 0 ? ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed : ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success,
      message: isEmpty(fieldValue) || ([SchemaType.Integer, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage].includes(field === null || field === void 0 ? void 0 : field.type) || ((_f = field === null || field === void 0 ? void 0 : field.name) === null || _f === void 0 ? void 0 : _f.includes((_g = SchemaType.Storage) === null || _g === void 0 ? void 0 : _g.toLocaleLowerCase())) || ((_h = field === null || field === void 0 ? void 0 : field.name) === null || _h === void 0 ? void 0 : _h.includes((_j = SchemaType.CPU) === null || _j === void 0 ? void 0 : _j.toLocaleLowerCase())) || ((_k = field === null || field === void 0 ? void 0 : field.name) === null || _k === void 0 ? void 0 : _k.includes('memory'))) && +fieldValue < 0 ? getFormItemEmptyMsg(field) : i18n.t('')
    };
  } else {
    validator = {
      status: ([SchemaType.Integer, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage].includes(field === null || field === void 0 ? void 0 : field.type) || ((_l = field === null || field === void 0 ? void 0 : field.name) === null || _l === void 0 ? void 0 : _l.includes((_m = SchemaType.Storage) === null || _m === void 0 ? void 0 : _m.toLocaleLowerCase())) || ((_o = field === null || field === void 0 ? void 0 : field.name) === null || _o === void 0 ? void 0 : _o.includes((_p = SchemaType.CPU) === null || _p === void 0 ? void 0 : _p.toLocaleLowerCase())) || ((_q = field === null || field === void 0 ? void 0 : field.name) === null || _q === void 0 ? void 0 : _q.includes('memory'))) && +fieldValue < 0 ? ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed : ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success,
      message: ([SchemaType.Integer, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage].includes(field === null || field === void 0 ? void 0 : field.type) || ((_r = field === null || field === void 0 ? void 0 : field.name) === null || _r === void 0 ? void 0 : _r.includes((_s = SchemaType.Storage) === null || _s === void 0 ? void 0 : _s.toLocaleLowerCase())) || ((_t = field === null || field === void 0 ? void 0 : field.name) === null || _t === void 0 ? void 0 : _t.includes((_u = SchemaType.CPU) === null || _u === void 0 ? void 0 : _u.toLocaleLowerCase())) || ((_v = field === null || field === void 0 ? void 0 : field.name) === null || _v === void 0 ? void 0 : _v.includes('memory'))) && +fieldValue < 0 ? getFormItemEmptyMsg(field) : i18n.t('')
    };
  }
  return validator;
};
var getFormItemEmptyMsg = function getFormItemEmptyMsg(field) {
  var _a, _b, _c, _d, _e;
  var msg = '';
  if ([SchemaType.Integer, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage].includes(field === null || field === void 0 ? void 0 : field.type) || ((_a = field === null || field === void 0 ? void 0 : field.name) === null || _a === void 0 ? void 0 : _a.includes((_b = SchemaType.Storage) === null || _b === void 0 ? void 0 : _b.toLocaleLowerCase())) || ((_c = field === null || field === void 0 ? void 0 : field.name) === null || _c === void 0 ? void 0 : _c.includes((_d = SchemaType.CPU) === null || _d === void 0 ? void 0 : _d.toLocaleLowerCase())) || ((_e = field === null || field === void 0 ? void 0 : field.name) === null || _e === void 0 ? void 0 : _e.includes('memory'))) {
    msg = "".concat(field === null || field === void 0 ? void 0 : field.label, "\u9700\u8981\u662F\u5927\u4E8E\u7B49\u4E8E0\u7684\u6570\u5B57");
  } else if ((field === null || field === void 0 ? void 0 : field.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Boolean)) {
    msg = "".concat(field === null || field === void 0 ? void 0 : field.label, "\u9700\u8981\u5F00\u542F");
  } else {
    msg = "".concat(field === null || field === void 0 ? void 0 : field.label, "\u4E0D\u80FD\u4E3A\u7A7A");
  }
  return i18n.t('{{msg}}', {
    msg: msg
  });
};
var validatePlanSchema = function validatePlanSchema(field, fieldValue, validator, values, serviceName, instance) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o;
  if (validator === void 0) {
    validator = OperatorNumberUnitReg;
  }
  // 优先进行空校验
  var emptyValidator = validateFormItemEmpty(field, values);
  if ((emptyValidator === null || emptyValidator === void 0 ? void 0 : emptyValidator.status) === (ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed)) {
    return emptyValidator;
  }
  // 依赖项未启用，校验直接通过
  if (field.enabledCondition && values) {
    var _p = field === null || field === void 0 ? void 0 : field.enabledCondition.split('=='),
      conditionKey = _p[0],
      conditionValue = _p[1];
    var value_1 = values[conditionKey];
    if (String(value_1) !== String(conditionValue)) {
      return {
        status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success,
        message: i18n.t('')
      };
    }
  }
  //校验rabbitmq的实例服务用户名不能和实例的管理用户相同
  if (serviceName === 'rabbitmq' && instance && (field === null || field === void 0 ? void 0 : field.name) === 'user' && ((_b = (_a = instance === null || instance === void 0 ? void 0 : instance.spec) === null || _a === void 0 ? void 0 : _a.parameters) === null || _b === void 0 ? void 0 : _b.adminUser) === fieldValue) {
    return {
      status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed,
      message: i18n.t('{{message}}', {
        message: "\u7528\u6237\u540D\u4E0D\u80FD\u548C\u5B9E\u4F8B\u7BA1\u7406\u7528\u6237\u540D\u3010".concat((_d = (_c = instance === null || instance === void 0 ? void 0 : instance.spec) === null || _c === void 0 ? void 0 : _c.parameters) === null || _d === void 0 ? void 0 : _d.adminUser, "\u3011\u76F8\u540C")
      })
    };
  }
  if ([SchemaType.Integer, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.CPU, SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Storage].includes(field === null || field === void 0 ? void 0 : field.type) || ((_e = field === null || field === void 0 ? void 0 : field.name) === null || _e === void 0 ? void 0 : _e.includes((_f = SchemaType.Storage) === null || _f === void 0 ? void 0 : _f.toLocaleLowerCase())) || ((_g = field === null || field === void 0 ? void 0 : field.name) === null || _g === void 0 ? void 0 : _g.includes((_h = SchemaType.CPU) === null || _h === void 0 ? void 0 : _h.toLocaleLowerCase())) || ((_j = field === null || field === void 0 ? void 0 : field.name) === null || _j === void 0 ? void 0 : _j.includes('memory'))) {
    if (isEmpty(values[field === null || field === void 0 ? void 0 : field.name]) || isNaN(Number(values[field === null || field === void 0 ? void 0 : field.name])) || Number(values[field === null || field === void 0 ? void 0 : field.name]) < 0) {
      return {
        status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed,
        message: i18n.t('{{message}}', {
          message: "".concat(field === null || field === void 0 ? void 0 : field.label, "\u9700\u8981\u662F\u5927\u4E8E\u7B49\u4E8E0\u7684\u6570\u5B57")
        })
      };
    }
  }
  var unitMapping = (_k = UnitMappingList === null || UnitMappingList === void 0 ? void 0 : UnitMappingList[field === null || field === void 0 ? void 0 : field.type]) !== null && _k !== void 0 ? _k : {
    defaultValue: ''
  };
  var _q = (_l = NumberUnitReg === null || NumberUnitReg === void 0 ? void 0 : NumberUnitReg.exec(fieldValue)) !== null && _l !== void 0 ? _l : ['', ''],
    value = _q[1],
    valueUnit = _q[2];
  var isOr = (_m = field === null || field === void 0 ? void 0 : field.validator) === null || _m === void 0 ? void 0 : _m.includes('|');
  var validations = isOr ? field === null || field === void 0 ? void 0 : field.validator.split('|') : (_o = field === null || field === void 0 ? void 0 : field.validator) === null || _o === void 0 ? void 0 : _o.split('&');
  var errorMsgs = [];
  var inValids = validations === null || validations === void 0 ? void 0 : validations.map(function (validation) {
    var _a;
    var fromValidation = validator === null || validator === void 0 ? void 0 : validator.exec(validation === null || validation === void 0 ? void 0 : validation.trim());
    if (!fromValidation) {
      return true;
    }
    var operator = fromValidation[1],
      number = fromValidation[2],
      unit = fromValidation[3];
    var formattedValue = getFormattedValue({
      value: value,
      unit: valueUnit,
      field: field
    });
    var formattedValidateValue = getFormattedValue({
      value: number,
      unit: unit,
      field: field
    });
    var competitor = compareOperator[operator || '=='];
    var v = competitor === null || competitor === void 0 ? void 0 : competitor.validate(formattedValue, formattedValidateValue);
    errorMsgs.push("".concat(competitor.errMsg).concat(number).concat((_a = unitMapping === null || unitMapping === void 0 ? void 0 : unitMapping[(unit === null || unit === void 0 ? void 0 : unit.toLowerCase()) || 'defaultValue']) !== null && _a !== void 0 ? _a : unit));
    return !v;
  }).filter(function (x) {
    return x;
  });
  var errorMsg = errorMsgs.join(isOr ? i18n.t('或者') : i18n.t('并且'));
  var result = (inValids === null || inValids === void 0 ? void 0 : inValids.length) ? isOr ? inValids === null || inValids === void 0 ? void 0 : inValids.some(function (item) {
    return item;
  }) : inValids === null || inValids === void 0 ? void 0 : inValids.every(function (item) {
    return item;
  }) : false;
  var message = result ? "".concat(field.label).concat(i18n.t('需要')).concat(errorMsg) : i18n.t('');
  var status = result ? ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed : ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success;
  return {
    status: status,
    message: message
  };
};
var formatPlanSchemaData = function formatPlanSchemaData(key, value) {
  var result = '';
  if ((key === null || key === void 0 ? void 0 : key.toLowerCase()) === 'cpu') {
    result = value === null || value === void 0 ? void 0 : value.replace(/m$|mi$/, '');
  } else if (['memory', 'storage'].includes(key === null || key === void 0 ? void 0 : key.toLowerCase())) {
    result = value === null || value === void 0 ? void 0 : value.replace(/m$|mi$|mib$|M$|Mi$|MiB$|g$|G$|Gi$/, '');
  } else {
    result = value;
  }
  return result;
};
var showUnitOptions = function showUnitOptions(field) {
  var _a, _b, _c, _d, _e, _f, _g;
  var result = false;
  if (((_a = field === null || field === void 0 ? void 0 : field.type) === null || _a === void 0 ? void 0 : _a.includes(SchemaType.CPU)) || ((_b = field === null || field === void 0 ? void 0 : field.name) === null || _b === void 0 ? void 0 : _b.includes((_c = SchemaType.CPU) === null || _c === void 0 ? void 0 : _c.toLocaleLowerCase()))) {
    result = true;
  } else if (((_d = field === null || field === void 0 ? void 0 : field.type) === null || _d === void 0 ? void 0 : _d.includes(SchemaType.Storage)) || ((_e = field === null || field === void 0 ? void 0 : field.name) === null || _e === void 0 ? void 0 : _e.includes((_f = SchemaType.Storage) === null || _f === void 0 ? void 0 : _f.toLocaleLowerCase())) || ((_g = field === null || field === void 0 ? void 0 : field.name) === null || _g === void 0 ? void 0 : _g.includes('memory'))) {
    result = true;
  } else {
    result = false;
  }
  return result;
};
var formatPlanSchemaUnit = function formatPlanSchemaUnit(field, filedValue) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j;
  var result = '';
  var _k = (_a = NumberUnitReg === null || NumberUnitReg === void 0 ? void 0 : NumberUnitReg.exec(filedValue)) !== null && _a !== void 0 ? _a : ['', ''],
    value = _k[1],
    valueUnit = _k[2];
  if (((_b = field === null || field === void 0 ? void 0 : field.type) === null || _b === void 0 ? void 0 : _b.includes(SchemaType.CPU)) || ((_c = field === null || field === void 0 ? void 0 : field.name) === null || _c === void 0 ? void 0 : _c.includes((_d = SchemaType.CPU) === null || _d === void 0 ? void 0 : _d.toLocaleLowerCase()))) {
    result = valueUnit || cpuUnitKeyMap.C;
  } else if (((_e = field === null || field === void 0 ? void 0 : field.type) === null || _e === void 0 ? void 0 : _e.includes(SchemaType.Storage)) || ((_f = field === null || field === void 0 ? void 0 : field.name) === null || _f === void 0 ? void 0 : _f.includes((_g = SchemaType.Storage) === null || _g === void 0 ? void 0 : _g.toLocaleLowerCase())) || ((_h = field === null || field === void 0 ? void 0 : field.name) === null || _h === void 0 ? void 0 : _h.includes('memory'))) {
    result = ((_j = valueUnit === null || valueUnit === void 0 ? void 0 : valueUnit.replace(/Gi$/, 'G')) === null || _j === void 0 ? void 0 : _j.replace(/Mi$/, 'M')) || storageUnitKeyMap.G;
  } else {
    result = '';
  }
  return result;
};
var formatPlanSchemaSubmitData = function formatPlanSchemaSubmitData(field, filedValue, fieldUnit) {
  var _a;
  var result;
  if ((field === null || field === void 0 ? void 0 : field.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Integer)) {
    result = parseInt(filedValue);
  } else if ((field === null || field === void 0 ? void 0 : field.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Boolean)) {
    result = filedValue;
  } else if ((field === null || field === void 0 ? void 0 : field.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.List)) {
    result = filedValue === null || filedValue === void 0 ? void 0 : filedValue.split(',');
  } else if ((field === null || field === void 0 ? void 0 : field.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.Map)) {
    try {
      var valueMap_1 = JSON.parse(filedValue);
      var newValueMap = (_a = Object.keys(valueMap_1)) === null || _a === void 0 ? void 0 : _a.reduce(function (pre, cur) {
        var _a;
        var value;
        // 当前字符值的内容是数字，将内容格式为数字类型
        if (!isNaN(parseInt(valueMap_1 === null || valueMap_1 === void 0 ? void 0 : valueMap_1[cur]))) {
          value = parseInt(valueMap_1 === null || valueMap_1 === void 0 ? void 0 : valueMap_1[cur]);
        } else if ((valueMap_1 === null || valueMap_1 === void 0 ? void 0 : valueMap_1[cur]) === 'true') {
          // 当前字符值的内容是'true'，将内容格式为布尔类型
          value = true;
        } else if ((valueMap_1 === null || valueMap_1 === void 0 ? void 0 : valueMap_1[cur]) === 'false') {
          // 当前字符值的内容是 'false'，将内容格式为布尔类型
          value = false;
        } else {
          // 以上条件均不满足，则不进行格式话处理
          value = valueMap_1 === null || valueMap_1 === void 0 ? void 0 : valueMap_1[cur];
        }
        return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur] = value, _a));
      }, {});
      result = newValueMap;
    } catch (error) {
      result = '';
    }
  } else {
    result = filedValue + (fieldUnit !== null && fieldUnit !== void 0 ? fieldUnit : '');
  }
  return result;
};
var getUnitOptions = function getUnitOptions(field) {
  var _a, _b, _c, _d, _e, _f, _g;
  var result = [];
  if (((_a = field === null || field === void 0 ? void 0 : field.type) === null || _a === void 0 ? void 0 : _a.includes(SchemaType.CPU)) || ((_b = field === null || field === void 0 ? void 0 : field.name) === null || _b === void 0 ? void 0 : _b.includes((_c = SchemaType.CPU) === null || _c === void 0 ? void 0 : _c.toLocaleLowerCase()))) {
    result = unitOptions.CPU;
  } else if (((_d = field === null || field === void 0 ? void 0 : field.type) === null || _d === void 0 ? void 0 : _d.includes(SchemaType.Storage)) || ((_e = field === null || field === void 0 ? void 0 : field.name) === null || _e === void 0 ? void 0 : _e.includes((_f = SchemaType.Storage) === null || _f === void 0 ? void 0 : _f.toLocaleLowerCase())) || ((_g = field === null || field === void 0 ? void 0 : field.name) === null || _g === void 0 ? void 0 : _g.includes('memory'))) {
    result = unitOptions.Storage;
  } else {
    result = [];
  }
  return result;
};
var formatPlanSchemaValue = function formatPlanSchemaValue(value) {
  var newValue = '';
  if (!value) {
    newValue = '';
  } else if (!isNaN(parseInt(value))) {
    newValue = parseInt(value);
  } else if (value === 'true') {
    newValue = true;
  } else if (value === 'false') {
    newValue = false;
  } else if (_typeof(value) === 'object') {
    newValue = JSON.stringify(value);
  } else {
    newValue = value;
  }
  return newValue;
};
var NotSupportBindingSchemaVendors = ['etcd', 'zookeeper'];

var _a$3;
var PlanMustPropMap;
(function (PlanMustPropMap) {
  PlanMustPropMap["InstanceName"] = "instanceName";
  PlanMustPropMap["ClusterId"] = "clusterId";
})(PlanMustPropMap || (PlanMustPropMap = {}));
// 校验的相关配置
var PlanBaseValidateSchema = (_a$3 = {}, _a$3[PlanMustPropMap === null || PlanMustPropMap === void 0 ? void 0 : PlanMustPropMap.InstanceName] = {
  rules: [{
    type: ffValidator.RuleTypeEnum.custom,
    customFunc: function customFunc(value, store) {
      var instanceName = (store === null || store === void 0 ? void 0 : store.formData).instanceName;
      if (!instanceName || !NumberAndSymbolReg.test(instanceName)) {
        return {
          status: ffValidator.ValidatorStatusEnum.Failed,
          message: i18n.t('名称为2-63个字符，可包含数字、小写英文字以及短划线（-），且不能以短划线（-）开头')
        };
      }
      return {
        status: 1,
        message: ''
      };
    }
  }]
}, _a$3[PlanMustPropMap === null || PlanMustPropMap === void 0 ? void 0 : PlanMustPropMap.ClusterId] = {
  rules: [{
    type: ffValidator.RuleTypeEnum.custom,
    customFunc: function customFunc(value, store) {
      var clusterId = (store === null || store === void 0 ? void 0 : store.formData).clusterId;
      if (!clusterId) {
        return {
          status: ffValidator.ValidatorStatusEnum.Failed,
          message: i18n.t('目标集群不能为空')
        };
      }
      return {
        status: 1,
        message: ''
      };
    }
  }]
}, _a$3);
var isPlanRequireBaseProps = function isPlanRequireBaseProps(key) {
  return ['instanceName', 'clusterId'].includes(key);
};
var validateAllPlanProps = function validateAllPlanProps(edit, instanceSchemas) {
  var _a;
  var vm;
  var formData = edit.formData;
  (_a = Object.keys(formData)) === null || _a === void 0 ? void 0 : _a.forEach(function (key) {
    var _a;
    var _b, _c, _d, _e;
    var validator;
    if (isPlanRequireBaseProps(key)) {
      validator = (_d = (_c = (_b = PlanBaseValidateSchema === null || PlanBaseValidateSchema === void 0 ? void 0 : PlanBaseValidateSchema[key]) === null || _b === void 0 ? void 0 : _b.rules) === null || _c === void 0 ? void 0 : _c[0]) === null || _d === void 0 ? void 0 : _d.customFunc(formData[key], edit);
    } else if (isInSchemaProps(key, instanceSchemas)) {
      var schema = instanceSchemas === null || instanceSchemas === void 0 ? void 0 : instanceSchemas.find(function (item) {
        return (item === null || item === void 0 ? void 0 : item.name) === key;
      });
      validator = validatePlanSchema(schema, formData[key] + ((_e = formData === null || formData === void 0 ? void 0 : formData['unitMap']) === null || _e === void 0 ? void 0 : _e[key]), OperatorNumberUnitReg, formData);
    } else {
      validator = {
        status: ffValidator.ValidatorStatusEnum.Success,
        message: ''
      };
    }
    vm = tslib.__assign(tslib.__assign({}, vm), (_a = {}, _a[key] = validator, _a));
  });
  return vm;
};
var _validatePlan = function _validatePlan(data, schemas) {
  var validator = validateAllPlanProps(data, schemas);
  return Object.keys(validator).every(function (key) {
    var _a;
    return ((_a = validator === null || validator === void 0 ? void 0 : validator[key]) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success;
  });
};
var initServicePlanEdit = function initServicePlanEdit(data) {
  return {
    formData: data === null || data === void 0 ? void 0 : data.reduce(function (pre, cur) {
      var _a, _b;
      return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = getFormItemDefaultValue(cur), _a.unitMap = tslib.__assign(tslib.__assign({}, pre === null || pre === void 0 ? void 0 : pre.unitMap), (_b = {}, _b[cur === null || cur === void 0 ? void 0 : cur.name] = cur === null || cur === void 0 ? void 0 : cur.unit, _b)), _a));
    }, {
      instanceName: '',
      unitMap: {}
    }),
    validator: data === null || data === void 0 ? void 0 : data.reduce(function (pre, cur) {
      var _a;
      return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = {
        status: 1,
        message: ''
      }, _a));
    }, {
      instanceName: {
        status: 1,
        message: ''
      }
    })
  };
};

var HubCluster;
(function (HubCluster) {
  HubCluster.DefaultRegion = 1;
  var StatusEnum;
  (function (StatusEnum) {
    StatusEnum["Running"] = "Running";
    StatusEnum["Initializing"] = "Initializing";
  })(StatusEnum = HubCluster.StatusEnum || (HubCluster.StatusEnum = {}));
})(HubCluster || (HubCluster = {}));

var PaasMedium;
(function (PaasMedium) {
  var _a;
  var StorageTypeEnum;
  (function (StorageTypeEnum) {
    StorageTypeEnum["S3"] = "s3";
    StorageTypeEnum["NFS"] = "nfs";
  })(StorageTypeEnum = PaasMedium.StorageTypeEnum || (PaasMedium.StorageTypeEnum = {}));
  PaasMedium.StorageTypeMap = (_a = {}, _a[StorageTypeEnum.S3] = i18n.t("对象存储(S3)"), _a[StorageTypeEnum.NFS] = i18n.t("文件存储NFS"), _a);
  PaasMedium.storageTypeOptions = [{
    text: i18n.t("对象存储(S3)"),
    value: StorageTypeEnum.S3
  }, {
    text: i18n.t("文件存储NFS"),
    value: StorageTypeEnum.NFS
  }];
  PaasMedium.StorageTypeLabel = "tdcc.cloud.tencent.com/paas-storage-medium";
  PaasMedium.MeduimCreatorLabel = "tdcc.cloud.tencent.com/creato";
})(PaasMedium || (PaasMedium = {}));

var _a$4, _b$3, _c$2;
var DetailTabType;
(function (DetailTabType) {
  DetailTabType["Detail"] = "ServiceDetail";
  DetailTabType["BackUp"] = "ServiceOpsPlan";
  DetailTabType["Monitor"] = "ServiceMonitor";
  DetailTabType["ServiceBinding"] = "ServiceBinding";
})(DetailTabType || (DetailTabType = {}));
var detailTabs = [{
  id: DetailTabType.Detail,
  label: i18n.t('基本信息')
}, {
  id: DetailTabType.BackUp,
  label: i18n.t('备份')
},
// {
//   id: DetailTabType.Monitor,
//   label: t('监控')
// },
{
  id: DetailTabType.ServiceBinding,
  label: i18n.t('服务绑定')
}];
var BackupTypeNum;
(function (BackupTypeNum) {
  BackupTypeNum["Manual"] = "Manual";
  BackupTypeNum["Schedule"] = "Schedule";
  BackupTypeNum["Unknown"] = "unknown";
})(BackupTypeNum || (BackupTypeNum = {}));
var BackupTypeMap = (_a$4 = {}, _a$4[BackupTypeNum.Manual] = {
  text: i18n.t('手动备份'),
  className: 'text-warning'
}, _a$4[BackupTypeNum.Schedule] = {
  text: i18n.t('定时备份'),
  className: 'text-success'
}, _a$4[BackupTypeNum.Unknown] = {
  text: i18n.t('-'),
  className: 'text-weak'
}, _a$4);
var BackupStatusNum;
(function (BackupStatusNum) {
  BackupStatusNum["Waiting"] = "Running";
  BackupStatusNum["Success"] = "Succeed";
  BackupStatusNum["Failed"] = "Failed";
})(BackupStatusNum || (BackupStatusNum = {}));
var BackupStatusMap = (_b$3 = {}, _b$3[BackupStatusNum.Waiting] = {
  text: i18n.t('备份中'),
  className: 'text-warning'
}, _b$3[BackupStatusNum.Success] = {
  text: i18n.t('备份成功'),
  className: 'text-success'
}, _b$3[BackupStatusNum.Failed] = {
  text: i18n.t('备份失败'),
  className: 'text-danger'
}, _b$3);
var ServiceBindingStatusNum;
(function (ServiceBindingStatusNum) {
  ServiceBindingStatusNum["Waiting"] = "Waiting";
  ServiceBindingStatusNum["Ready"] = "Ready";
  ServiceBindingStatusNum["Failed"] = "ServiceInstanceNotReady";
})(ServiceBindingStatusNum || (ServiceBindingStatusNum = {}));
var ServiceBindingStatusMap = (_c$2 = {}, _c$2[ServiceBindingStatusNum.Waiting] = {
  text: i18n.t('创建中'),
  className: 'text-warning'
}, _c$2[ServiceBindingStatusNum.Ready] = {
  text: i18n.t('绑定成功'),
  className: 'text-success'
}, _c$2[ServiceBindingStatusNum.Failed] = {
  text: i18n.t('绑定失败'),
  className: 'text-danger'
}, _c$2);
var BackupStrategyEditInitValue = {
  formData: {
    enable: false,
    backupDate: [],
    backupTime: [],
    backupReserveDay: 30
  },
  validator: {
    enable: {
      message: '',
      status: ffValidator.ValidatorStatusEnum.Init
    },
    backupDate: {
      message: '',
      status: ffValidator.ValidatorStatusEnum.Init
    },
    backupTime: {
      message: '',
      status: ffValidator.ValidatorStatusEnum.Init
    },
    backupReserveDay: {
      message: '',
      status: ffValidator.ValidatorStatusEnum.Init
    }
  }
};

var Backup;
(function (Backup) {
  Backup.minReserveDay = 1;
  Backup.maxReserveDay = 30;
  Backup.weekConfig = [{
    value: '0',
    text: i18n.t('每周日')
  }, {
    value: '1',
    text: i18n.t('每周一')
  }, {
    value: '2',
    text: i18n.t('每周二')
  }, {
    value: '3',
    text: i18n.t('每周三')
  }, {
    value: '4',
    text: i18n.t('每周四')
  }, {
    value: '5',
    text: i18n.t('每周五')
  }, {
    value: '6',
    text: i18n.t('每周六')
  }];
  Backup.hourConfig = [{
    value: '0',
    text: '00:00'
  }, {
    value: '1',
    text: '01:00'
  }, {
    value: '2',
    text: '02:00'
  }, {
    value: '3',
    text: '03:00'
  }, {
    value: '4',
    text: '04:00'
  }, {
    value: '5',
    text: '05:00'
  }, {
    value: '6',
    text: '06:00'
  }, {
    value: '7',
    text: '07:00'
  }, {
    value: '8',
    text: '08:00'
  }, {
    value: '9',
    text: '09:00'
  }, {
    value: '10',
    text: '10:00'
  }, {
    value: '11',
    text: '11:00'
  }, {
    value: '12',
    text: '12:00'
  }, {
    value: '13',
    text: '13:00'
  }, {
    value: '14',
    text: '14:00'
  }, {
    value: '15',
    text: '15:00'
  }, {
    value: '16',
    text: '16:00'
  }, {
    value: '17',
    text: '17:00'
  }, {
    value: '18',
    text: '18:00'
  }, {
    value: '19',
    text: '19:00'
  }, {
    value: '20',
    text: '20:00'
  }, {
    value: '21',
    text: '21:00'
  }, {
    value: '22',
    text: '22:00'
  }, {
    value: '23',
    text: '23:00'
  }];
  Backup.invalidFormItemMsg = function (data) {
    var _a, _b, _c;
    return (_c = (_b = data === null || data === void 0 ? void 0 : data.validator[(_a = Object.keys(data === null || data === void 0 ? void 0 : data.validator)) === null || _a === void 0 ? void 0 : _a.find(function (item) {
      var _a;
      return ((_a = data === null || data === void 0 ? void 0 : data.validator[item]) === null || _a === void 0 ? void 0 : _a.status) === 2;
    })]) === null || _b === void 0 ? void 0 : _b.message) !== null && _c !== void 0 ? _c : '';
  };
  Backup._validateAll = function (data) {
    var _a;
    var allValidator = Backup.getValidatorModel(data);
    var key = (_a = Object.keys(allValidator)) === null || _a === void 0 ? void 0 : _a.find(function (key) {
      var _a;
      return ((_a = allValidator === null || allValidator === void 0 ? void 0 : allValidator[key]) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Failed;
    });
    var inValidItem = allValidator === null || allValidator === void 0 ? void 0 : allValidator[key];
    return {
      valid: !inValidItem,
      message: !inValidItem ? i18n.t('') : inValidItem === null || inValidItem === void 0 ? void 0 : inValidItem.message
    };
  };
  Backup.getValidatorModel = function (backupStrategyEdit) {
    var _a;
    var formData = backupStrategyEdit.formData;
    return (_a = Object === null || Object === void 0 ? void 0 : Object.keys(formData)) === null || _a === void 0 ? void 0 : _a.reduce(function (pre, key) {
      var _a;
      var _b, _c, _d, _e;
      var validator;
      if (!(formData === null || formData === void 0 ? void 0 : formData.enable)) {
        validator = {
          message: '',
          status: ffValidator.ValidatorStatusEnum.Success
        };
      } else {
        if ('backupDate' === key || 'backupTime' === key) {
          validator = {
            message: !!((_b = formData === null || formData === void 0 ? void 0 : formData['backupDate']) === null || _b === void 0 ? void 0 : _b.length) || !!((_c = formData === null || formData === void 0 ? void 0 : formData['backupTime']) === null || _c === void 0 ? void 0 : _c.length) ? i18n.t('') : i18n.t('备份日期和备份时间点至少选择一个'),
            status: !!((_d = formData === null || formData === void 0 ? void 0 : formData['backupDate']) === null || _d === void 0 ? void 0 : _d.length) || !!((_e = formData === null || formData === void 0 ? void 0 : formData['backupTime']) === null || _e === void 0 ? void 0 : _e.length) ? ffValidator.ValidatorStatusEnum.Success : ffValidator.ValidatorStatusEnum.Failed
          };
        } else if ('backupReserveTime' === key) {
          validator = {
            message: !!(formData === null || formData === void 0 ? void 0 : formData[key]) ? i18n.t('') : i18n.t('备份保留时间不能为空'),
            status: !!(formData === null || formData === void 0 ? void 0 : formData[key]) ? ffValidator.ValidatorStatusEnum.Success : ffValidator.ValidatorStatusEnum.Failed
          };
        } else {
          validator = {
            message: '',
            status: ffValidator.ValidatorStatusEnum.Success
          };
        }
      }
      return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[key] = validator, _a));
    }, {});
  };
  Backup.reduceBackupStrategyJson = function (data) {
    var _a, _b, _c, _d;
    var backupDate = data.backupDate,
      backupTime = data.backupTime,
      backupReserveDay = data.backupReserveDay,
      instanceName = data.instanceName,
      enable = data.enable,
      instanceId = data.instanceId,
      serviceName = data.serviceName,
      _e = data.mode,
      mode = _e === void 0 ? 'create' : _e,
      medium = data.medium;
    var weekStr = backupDate.sort().join(',') || '*';
    var hourStr = backupTime.sort().join(',') || '0';
    var minuteStr = '0';
    //cron 格式:`* * * * * *`,即:分 时 天 月 周 年(可省略)
    // const cron = `${minuteStr} ${hourStr} ? * ${weekStr}`;
    // 更改备份策略的资源空间
    var cron = "".concat(minuteStr, " ").concat(hourStr, " ? * ").concat(weekStr);
    var jsonData = {
      apiVersion: 'infra.tce.io/v1',
      kind: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup,
      metadata: {
        name: mode === 'create' ? 'backup-' + instanceName + '-' + new Date().getTime() : data === null || data === void 0 ? void 0 : data.name,
        namespace: SystemNamespace
      },
      spec: {
        enabled: enable,
        operation: {
          backup: {}
        },
        target: {
          instanceID: instanceId,
          serviceClass: serviceName
        }
      }
    };
    var operation = {
      backup: {}
    };
    if (medium) {
      if (((_b = (_a = medium === null || medium === void 0 ? void 0 : medium.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b[PaasMedium.StorageTypeLabel]) === PaasMedium.StorageTypeEnum.S3) {
        operation = {
          backup: {
            destination: {
              s3Objects: {
                secretRef: {
                  namespace: 'ssm',
                  name: (_c = medium === null || medium === void 0 ? void 0 : medium.metadata) === null || _c === void 0 ? void 0 : _c.name
                }
              }
            }
          }
        };
      } else {
        operation = {
          backup: {
            destination: {
              nfsObjects: {
                secretRef: {
                  namespace: 'ssm',
                  name: (_d = medium === null || medium === void 0 ? void 0 : medium.metadata) === null || _d === void 0 ? void 0 : _d.name
                }
              }
            }
          }
        };
      }
    }
    jsonData.spec['operation'] = operation;
    if (enable) {
      jsonData.spec['retain'] = {
        days: backupReserveDay
      };
      jsonData.spec['trigger'] = {
        type: BackupTypeNum === null || BackupTypeNum === void 0 ? void 0 : BackupTypeNum.Schedule,
        params: {
          cron: cron
        }
      };
    }
    return JSON.stringify(jsonData);
  };
})(Backup || (Backup = {}));

var ServiceBinding;
(function (ServiceBinding) {
  ServiceBinding.ServiceBindingEditInitValue = {
    formData: {
      name: '',
      namespace: ''
    },
    validator: {
      name: {
        message: '',
        status: ffValidator.ValidatorStatusEnum.Init
      },
      namespace: {
        message: '',
        status: ffValidator.ValidatorStatusEnum.Init
      }
    }
  };
  ServiceBinding.ServiceBindingEditSchema = {
    name: {
      rules: [{
        required: true,
        message: i18n.t('名称不能为空')
      }, {
        reg: /^[a-z][0-9a-z-]{1,62}$/,
        message: i18n.t('名称为2-63个字符，可包含数字、小写英文字以及短划线（-），且不能以短划线（-）开头')
      }, {
        max: 63,
        message: i18n.t('最长 63 个字符')
      }]
    },
    namespace: {
      rules: [{
        required: true,
        message: i18n.t('命名空间不能为空')
      }]
    }
  };
  ServiceBinding._validateServiceBindingItem = function (key, value) {
    var _a, _b;
    var rules = (_b = (_a = ServiceBinding.ServiceBindingEditSchema === null || ServiceBinding.ServiceBindingEditSchema === void 0 ? void 0 : ServiceBinding.ServiceBindingEditSchema[key]) === null || _a === void 0 ? void 0 : _a.rules) !== null && _b !== void 0 ? _b : [];
    var notMatchRule = rules === null || rules === void 0 ? void 0 : rules.find(function (item) {
      var _a;
      var isNotMatchRule = (item === null || item === void 0 ? void 0 : item.required) && !value || (item === null || item === void 0 ? void 0 : item.reg) && !((_a = item === null || item === void 0 ? void 0 : item.reg) === null || _a === void 0 ? void 0 : _a.test(value)) || !!(item === null || item === void 0 ? void 0 : item.max) && !!value && (value === null || value === void 0 ? void 0 : value.length) > (item === null || item === void 0 ? void 0 : item.max);
      return isNotMatchRule;
    });
    return {
      status: !notMatchRule ? ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success : ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed,
      message: !notMatchRule ? '' : notMatchRule === null || notMatchRule === void 0 ? void 0 : notMatchRule.message
    };
  };
  // export const _validatePlanSchema = (value:string | number | boolean,planSchema:PlanSchema):Validation =>{
  //   return  {
  //     status: planSchema?.optional ? ValidatorStatusEnum?.Success : ValidatorStatusEnum?.Failed,
  //     message: ''
  //   };
  // } 
  ServiceBinding._validateFormItem = function (editData, planSchemas, serviceName, instance) {
    var _a;
    var formData = editData.formData;
    return (_a = Object === null || Object === void 0 ? void 0 : Object.keys(formData)) === null || _a === void 0 ? void 0 : _a.reduce(function (pre, key) {
      var _a;
      var validator;
      var schema = planSchemas === null || planSchemas === void 0 ? void 0 : planSchemas.find(function (item) {
        return (item === null || item === void 0 ? void 0 : item.name) === key;
      });
      //当前属性是否是来自后端的schema属性
      if (schema) {
        validator = validatePlanSchema(schema, formData[key], OperatorNumberUnitReg, formData, serviceName, instance);
      } else {
        validator = ServiceBinding._validateServiceBindingItem(key, formData === null || formData === void 0 ? void 0 : formData[key]);
      }
      return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[key] = validator, _a));
    }, {});
  };
  ServiceBinding._validateAll = function (data, planSchemas, serviceName, instance) {
    var _a;
    var allValidator = ServiceBinding._validateFormItem(data, planSchemas, serviceName, instance);
    var key = (_a = Object.keys(allValidator)) === null || _a === void 0 ? void 0 : _a.find(function (key) {
      var _a;
      return ((_a = allValidator === null || allValidator === void 0 ? void 0 : allValidator[key]) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Failed;
    });
    var inValidItem = allValidator === null || allValidator === void 0 ? void 0 : allValidator[key];
    return {
      valid: !inValidItem,
      message: !inValidItem ? i18n.t('') : inValidItem === null || inValidItem === void 0 ? void 0 : inValidItem.message
    };
  };
  ServiceBinding.initServiceResourceEdit = function (data) {
    return {
      formData: data === null || data === void 0 ? void 0 : data.reduce(function (pre, cur) {
        var _a;
        return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = getFormItemDefaultValue(cur), _a));
      }, tslib.__assign({}, ServiceBinding.ServiceBindingEditInitValue === null || ServiceBinding.ServiceBindingEditInitValue === void 0 ? void 0 : ServiceBinding.ServiceBindingEditInitValue.formData)),
      validator: data === null || data === void 0 ? void 0 : data.reduce(function (pre, cur) {
        var _a;
        return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = {
          status: 1,
          message: ''
        }, _a));
      }, tslib.__assign({}, ServiceBinding.ServiceBindingEditInitValue === null || ServiceBinding.ServiceBindingEditInitValue === void 0 ? void 0 : ServiceBinding.ServiceBindingEditInitValue.validator))
    };
  };
})(ServiceBinding || (ServiceBinding = {}));

// import { yunApiFeedback } from './yunApiFeedback';
// import { Aegis } from './aegis';
function appendFunction(origin, append) {
  return function appended() {
    var args = [];
    for (var _i = 0; _i < arguments.length; _i++) {
      args[_i] = arguments[_i];
    }
    var result = origin.apply(this, args);
    return append.apply(this, [result].concat(args));
  };
}
/**
 * 判断是否和路由跳转之前的一样
 * @param prevpath: string  之前的路由
 * @param currentpath: string 当前的路由
 */
var isInSameModule = function isInSameModule(prevpath, currentpath) {
  var _a = prevpath.split('/').filter(function (item) {
      return item !== '';
    }),
    prevBusiness = _a[0],
    prevModule = _a[1],
    prevRest = _a.slice(2),
    _b = currentpath.split('/').filter(function (item) {
      return item !== '';
    }),
    currentBusiness = _b[0],
    currentModule = _b[1],
    currentRest = _b.slice(2);
  if (prevModule !== currentModule) {
    return false;
  }
  return true;
};
var nmcRouter = seajs.require('router');
var pageManager = seajs.require('nmc/main/pagemanager');
function getFragment() {
  var debug = nmcRouter.debug || '';
  var debugReg = new RegExp('^' + debug.replace('/', '\\/'));
  var fragment = nmcRouter.getFragment().split('?');
  return fragment[0].replace(debugReg, '');
}
function getQueryString() {
  var str = nmcRouter.getFragment().split('?');
  return str[1] ? '?' + str[1] : '';
}
function getInitialState(rule) {
  var fragment = getFragment();
  var params = nmcRouter.matchRoute(rule, fragment) || [];
  var queryString = getQueryString();
  var queries = queryString ? parseQueryString(queryString) : {};
  return {
    fragment: fragment,
    params: params,
    queryString: queryString,
    queries: queries
  };
}
var RouterNavigateAction = 'RouterNavigate';
function generateRouterReducer(rule) {
  function routerReducer(state, action) {
    if (state === void 0) {
      state = getInitialState(rule);
    }
    if (action.type === RouterNavigateAction) {
      return action.payload;
    }
    return state;
  }
  return routerReducer;
}
var curFragment = '',
  curQueryString = '',
  curQueries = {};
function navigateAction(fragment, params, queryString, queries) {
  return {
    type: RouterNavigateAction,
    payload: {
      fragment: fragment,
      params: params,
      queryString: queryString,
      queries: queries
    }
  };
}
function getCurrentRouterStatus(fragment, params, qString, queries) {
  var queryString = qString;
  var flag = false; /**记录是否变化， 默认无变化 */
  if (fragment === curFragment && curQueryString && !queryString) {
    /**path无变化，参数变化 */
    queryString = curQueryString;
  } else {
    /**path变化，重置当前缓存 */
    curFragment = fragment;
    curQueryString = queryString;
    curQueries = queries;
    flag = true;
  }
  return {
    flag: flag,
    queries: curQueries
  };
}
function startAction(rule) {
  return function (dispatch) {
    var _a = getInitialState(rule),
      fragment = _a.fragment,
      params = _a.params,
      queryString = _a.queryString,
      queries = _a.queries;
    dispatch(navigateAction(fragment, params, queryString, queries));
    nmcRouter.use(rule, function () {
      var args = [];
      for (var _i = 0; _i < arguments.length; _i++) {
        args[_i] = arguments[_i];
      }
      var params = args.slice();
      var fragment = getFragment();
      /* eslint-disable */
      _typeof(params[params.length - 1]) === 'object' ? params.pop() : {}; // nmcRouter parse error
      /* eslint-enable */
      var queryString = getQueryString();
      var queries = parseQueryString(queryString);
      dispatch(navigateAction(fragment, params, queryString, queries));
      var curStatus = getCurrentRouterStatus(fragment, params, queryString, queries);
      if (curStatus.flag) {
        // 更新导航状态
        var parts = fragment.split('/');
        parts.shift();
        var navArgs = [parts.shift(), parts.shift(), parts];
        pageManager.fragment = fragment;
        pageManager.changeNavStatus.apply(pageManager, navArgs);
      } else {
        nmcRouter.navigate(fragment + buildQueryString(curStatus.queries));
      }
    });
  };
}
function generateRouterDecorator(rule) {
  function decorator(target) {
    function onMount() {
      var _this = this;
      _this.props.dispatch(startAction(rule));
      // Aegis.init();
    }

    function onUnmount() {
      nmcRouter.unuse(rule);
      //当离开模块的时候，也清楚掉error
      // yunApiFeedback.clearApiError();
    }

    var proto = target.prototype;
    proto.componentWillMount = proto.componentWillMount ? appendFunction(proto.componentWillMount, onMount) : onMount;
    proto.componentWillUnmount = proto.componentWillUnmount ? appendFunction(proto.componentWillUnmount, onUnmount) : onUnmount;
  }
  return decorator;
}
var NAMED_PARAM_REGEX = /\(\/:(\w+)\)/g;
var RemoveExec4TdccCluster = /\(\/\(cluster\|startup\)\)/g;
var Router = /** @class */function () {
  /* eslint-disable */
  function Router(rule, defaults) {
    var _this_1 = this;
    this.rule = rule;
    this.defaults = defaults;
    this.paramNames = [];
    this.rule.replace(NAMED_PARAM_REGEX, function (match, name) {
      _this_1.paramNames.push(name);
      return match;
    });
  }
  /* eslint-enable */
  Router.prototype.getReducer = function () {
    return generateRouterReducer(this.rule);
  };
  Router.prototype.serve = function () {
    return generateRouterDecorator(this.rule);
  };
  Router.prototype.buildFragment = function (params) {
    var _this_1 = this;
    if (params === void 0) {
      params = {};
    }
    // TODO 优化默认路由生成
    var temp = this.rule.replace(NAMED_PARAM_REGEX, function (matched, name) {
      if (typeof params[name] === 'undefined' || params[name] === undefined || params[name] === _this_1.defaults[name]) {
        return '/' + _this_1.defaults[name];
      }
      return '/' + params[name];
    });
    temp = temp.replace(/([^:])[\/\\\\]{2,}/, '$1/');
    temp = temp.endsWith('/') ? temp.substring(0, temp.length - 1) : temp;
    temp = temp.replace(RemoveExec4TdccCluster, '');
    return temp;
  };
  Router.prototype.resolve = function (state) {
    var _this_1 = this;
    var resolved = {};
    var paramNames = this.paramNames.slice();
    var paramValues = state.params.slice();
    paramNames.forEach(function (x) {
      resolved[x] = paramValues.shift() || _this_1.defaults[x];
    });
    return resolved;
  };
  Router.prototype.buildUrl = function (params, queries) {
    if (params === void 0) {
      params = {};
    }
    return this.buildFragment(params) + buildQueryString(queries);
  };
  Router.prototype.navigate = function (params, queries, url) {
    if (params === void 0) {
      params = {};
    }
    if (url) {
      nmcRouter.navigate(url);
    } else {
      var nextLocationPath = this.buildFragment(params);
      var prevLocationPath = location.pathname;
      if (isInSameModule(prevLocationPath, nextLocationPath)) {
        nmcRouter.navigate(nextLocationPath + buildQueryString(queries));
      }
    }
  };
  return Router;
}();

var SubEnum;
(function (SubEnum) {
  SubEnum["List"] = "list";
  SubEnum["Detail"] = "detail";
  SubEnum["Create"] = "create";
  //编辑实例
  SubEnum["Edit"] = "edit";
})(SubEnum || (SubEnum = {}));
var TabEnum;
(function (TabEnum) {
  TabEnum["Info"] = "info";
  TabEnum["Instance"] = "instance";
  TabEnum["Yaml"] = "yaml";
})(TabEnum || (TabEnum = {}));
/**
 * @param sub 当前的模式，create | update | detail
 */
var router = new Router("".concat(Util === null || Util === void 0 ? void 0 : Util.getRouterPath(location === null || location === void 0 ? void 0 : location.pathname), "(/:sub)(/:tab)"), {
  sub: SubEnum.List,
  tab: ''
});

var routerSea = seajs.require('router');
function PaasHeader(props) {
  var _a, _b, _c;
  var route = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.base) === null || _a === void 0 ? void 0 : _a.route;
  });
  var tab = (router === null || router === void 0 ? void 0 : router.resolve(route)).tab;
  var platform = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.base) === null || _a === void 0 ? void 0 : _a.platform;
  });
  var services = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.services;
  });
  var urlParams = router === null || router === void 0 ? void 0 : router.resolve(route);
  var _d = route === null || route === void 0 ? void 0 : route.queries,
    _e = _d.servicename,
    servicename = _e === void 0 ? '-' : _e,
    _f = _d.instancename,
    instancename = _f === void 0 ? '-' : _f,
    mode = _d.mode;
  var headerTitle;
  var goBack = function goBack() {
    var _a;
    router.navigate({
      sub: 'list',
      tab: undefined
    }, {
      servicename: (_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.servicename,
      resourceType: ResourceTypeEnum.ServiceResource,
      mode: "list"
    });
  };
  if (platform === (PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TDCC)) {
    if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === 'create') {
      headerTitle = "\u65B0\u5EFA".concat(servicename !== null && servicename !== void 0 ? servicename : '-').concat((_a = ResourceTypeMap === null || ResourceTypeMap === void 0 ? void 0 : ResourceTypeMap[tab]) === null || _a === void 0 ? void 0 : _a.title);
    } else if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === 'detail') {
      headerTitle = "".concat(servicename, "\u5B9E\u4F8B / ").concat(instancename, "\u5B9E\u4F8B\u8BE6\u60C5");
    } else {
      headerTitle = '中间件列表';
    }
  } else {
    if (mode === 'create') {
      headerTitle = "\u65B0\u5EFA".concat(servicename !== null && servicename !== void 0 ? servicename : '-').concat((_b = ResourceTypeMap === null || ResourceTypeMap === void 0 ? void 0 : ResourceTypeMap[tab]) === null || _b === void 0 ? void 0 : _b.title);
    } else if (mode === 'detail') {
      headerTitle = "".concat(servicename, "\u5B9E\u4F8B / ").concat(instancename, "\u5B9E\u4F8B\u8BE6\u60C5");
    } else {
      headerTitle = '中间件列表';
    }
  }
  var showBackIcon = !((_c = ['', 'list']) === null || _c === void 0 ? void 0 : _c.includes(urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub));
  return React__default.createElement(teaComponent.Card, {
    style: {
      height: 50,
      display: "flex",
      alignItems: 'center',
      paddingLeft: 30,
      paddingRight: 30
    }
  }, React__default.createElement(teaComponent.Justify, {
    left: React__default.createElement(teaComponent.Row, {
      style: {
        display: "flex",
        alignItems: 'center'
      }
    }, showBackIcon && React__default.createElement(teaComponent.Icon, {
      type: "btnback",
      onClick: goBack,
      className: 'tea-mr-2n',
      style: {
        marginLeft: -2
      }
    }), React__default.createElement(teaComponent.H3, {
      className: 'tea-d-inline-block',
      style: {
        display: 'inline-block',
        padding: 0
      }
    }, i18n.t('{{headerTitle}}', {
      headerTitle: headerTitle
    })))
  }));
}

function LoadingPanel(props) {
  var _a = props.text,
    text = _a === void 0 ? '' : _a;
  return React__default.createElement("div", {
    style: {
      width: '100%'
    }
  }, React__default.createElement(teaComponent.Icon, {
    type: 'loading'
  }), React__default.createElement(teaComponent.Text, null, text, i18n.t('加载中')));
}

var TipInfo = /** @class */function (_super) {
  tslib.__extends(TipInfo, _super);
  function TipInfo() {
    return _super !== null && _super.apply(this, arguments) || this;
  }
  TipInfo.prototype.render = function () {
    var _a = this.props,
      _b = _a.style,
      style = _b === void 0 ? {} : _b,
      _c = _a.isShow,
      isShow = _c === void 0 ? true : _c,
      _d = _a.isForm,
      isForm = _d === void 0 ? false : _d,
      restProps = tslib.__rest(_a, ["style", "isShow", "isForm"]),
      renderStyle = style;
    // 用于在创建表单当中 展示错误信息
    if (isForm) {
      renderStyle = Object.assign({}, renderStyle, {
        //内部布局是flex
        display: 'inline-flex',
        marginLeft: '20px',
        marginBottom: '0px',
        maxWidth: '750px',
        maxHeight: '120px',
        overflow: 'auto'
      });
    }
    return isShow ? React__default.createElement(teaComponent.Alert, tslib.__assign({
      style: renderStyle
    }, restProps), this.props.children) : React__default.createElement("noscript", null);
  };
  return TipInfo;
}(React__default.Component);

var BComponent;
(function (BComponent) {
  BComponent.BaseActionType = {
    Clear: 'Clear',
    Validator: 'Validator'
  };
  BComponent.getActionType = function (moduleName, actionType) {
    return actionType ? moduleName + '_' + actionType : moduleName;
  };
  BComponent.createActionType = function (componentName, actionType) {
    for (var key in actionType) {
      var ns = [];
      componentName && ns.push(componentName);
      actionType[key] && ns.push(actionType[key]);
      actionType[key] = ns.join('_');
    }
  };
  BComponent.isNeedFetch = function (model, filter) {
    if (filter === void 0) {
      filter = null;
    }
    var isAllSame = true;
    if (filter) {
      var targetFilter_1 = [];
      var originFilter_1 = [];
      Object.keys(filter).forEach(function (key) {
        targetFilter_1.push(filter === null || filter === void 0 ? void 0 : filter[key]);
        originFilter_1.push(model.query.filter[key]);
      });
      isAllSame = JSON.stringify(targetFilter_1) === JSON.stringify(originFilter_1);
    }
    if (isAllSame) {
      //如果参数都没变，
      //1.正在加载的就不发请求
      //2.之前的请求没出错就不再发起请求
      if (typeof model['list'] !== 'undefined') {
        if (model['list'].loading || model['list'].fetched
        //&& model['list'].fetchState !== FetchState.Failed
        ) {
          return false;
        } else {
          return true;
        }
      }
      if (typeof model['object'] !== 'undefined') {
        if (model['object'].loading || model['object'].fetched
        // && model['object'].fetchState !== FetchState.Failed
        ) {
          return false;
        } else {
          return true;
        }
      }
    } else {
      return true;
    }
  };
  var ViewStyleEnum;
  (function (ViewStyleEnum) {
    ViewStyleEnum["Panel"] = "Panel";
    ViewStyleEnum["FormItem"] = "FormItem";
    ViewStyleEnum["Alert"] = "Alert";
    ViewStyleEnum["Text"] = "Text";
    ViewStyleEnum["Dialog"] = "Dialog";
  })(ViewStyleEnum = BComponent.ViewStyleEnum || (BComponent.ViewStyleEnum = {}));
  var EditStyleEnum;
  (function (EditStyleEnum) {
    EditStyleEnum["PopConfirm"] = "PopConfirm";
    EditStyleEnum["Panel"] = "Panel";
    EditStyleEnum["FormItem"] = "FormItem";
  })(EditStyleEnum = BComponent.EditStyleEnum || (BComponent.EditStyleEnum = {}));
  var OperationTypeEnum;
  (function (OperationTypeEnum) {
    OperationTypeEnum["Create"] = "Create";
    /**
     * @deprecated
     */
    OperationTypeEnum["CreateNode"] = "CreateNode";
    OperationTypeEnum["Delete"] = "Delete";
    OperationTypeEnum["Refund"] = "Refund";
    OperationTypeEnum["Modify"] = "Modify";
    OperationTypeEnum["View"] = "View";
    /**
     * @deprecated
     */
    OperationTypeEnum["EditSecret"] = "EditSecret";
    /**
     * @deprecated
     */
    OperationTypeEnum["IngressCreateSecret"] = "IngressCreateSecret";
  })(OperationTypeEnum = BComponent.OperationTypeEnum || (BComponent.OperationTypeEnum = {}));
  var ResourceTypeEnum;
  (function (ResourceTypeEnum) {
    ResourceTypeEnum["Asg"] = "Asg";
    ResourceTypeEnum["Node"] = "Node";
    ResourceTypeEnum["NodePool"] = "NodePool";
    ResourceTypeEnum["AutonodePool"] = "AutoNodePool";
    ResourceTypeEnum["VirtualNodePool"] = "VirtualNodePool";
    ResourceTypeEnum["VirtualNode"] = "VirtualNode";
    ResourceTypeEnum["EdgeCvm"] = "EdgeCvm";
    ResourceTypeEnum["EdgeDeploymentGridInstance"] = "EdgeDeploymentGridInstance";
    ResourceTypeEnum["LogListener"] = "LogListener";
    ResourceTypeEnum["Audit"] = "Audit";
    ResourceTypeEnum["Event"] = "Event";
    ResourceTypeEnum["Ingress"] = "Ingress";
    ResourceTypeEnum["Service"] = "Service";
    ResourceTypeEnum["HPA"] = "HPA";
    ResourceTypeEnum["Workload"] = "Workload";
    ResourceTypeEnum["Secret"] = "Secret";
    ResourceTypeEnum["EKSContainer"] = "EKSContainer";
    ResourceTypeEnum["External"] = "External";
    ResourceTypeEnum["ImageCache"] = "ImageCache";
    ResourceTypeEnum["Subscription"] = "Subscription";
    ResourceTypeEnum["Cluster"] = "Cluster";
    ResourceTypeEnum["ECluster"] = "ECluster";
    ResourceTypeEnum["HubCluster"] = "HubCluster";
    ResourceTypeEnum["EdgeCluster"] = "EdgeCluster";
    ResourceTypeEnum["PersistentVolume"] = "PersistentVolume";
    ResourceTypeEnum["PersistentVolumeClaim"] = "PersistentVolumeClaim";
    ResourceTypeEnum["StorageClass"] = "StorageClass";
    ResourceTypeEnum["Label"] = "Label";
    ResourceTypeEnum["Annotation"] = "Annotation";
  })(ResourceTypeEnum = BComponent.ResourceTypeEnum || (BComponent.ResourceTypeEnum = {}));
  BComponent.createHooks = function (_a) {
    var Context = _a.Context,
      itemsSelector = _a.itemsSelector,
      vKeyPrefix = _a.vKeyPrefix;
    var useModel = function useModel(modelSelector, eqFn) {
      if (eqFn === void 0) {
        eqFn = reactRedux.shallowEqual;
      }
      var context = React__default.useContext(Context);
      return reactRedux.useSelector(function (state) {
        return modelSelector(context.selector(function () {
          return state;
        }).model);
      }, eqFn);
    };
    var useFilter = function useFilter(filterSelector) {
      var context = React__default.useContext(Context);
      return reactRedux.useSelector(function (state) {
        return filterSelector(context.selector(function () {
          return state;
        }).filter);
      }, reactRedux.shallowEqual);
    };
    var useItem = function useItem(itemId, itemSelector) {
      return useModel(function (model) {
        var items = itemsSelector(model);
        var item = items.find(function (item) {
          return item.id === itemId;
        });
        return itemSelector(item);
      });
    };
    var useVkey = function useVkey(_a) {
      var key = _a.key,
        itemId = _a.itemId;
      var itemIndex = useModel(function (model) {
        var items = itemsSelector(model);
        return items.findIndex(function (item) {
          return item.id === itemId;
        });
      });
      return "".concat(vKeyPrefix, "[").concat(itemIndex, "].").concat(String(key));
    };
    return {
      useVkey: useVkey,
      useModel: useModel,
      useFilter: useFilter,
      useItem: useItem
    };
  };
  BComponent.createSimpleHooks = function (selector) {
    var useModel = function useModel(modelSelector, eqFn) {
      if (eqFn === void 0) {
        eqFn = reactRedux.shallowEqual;
      }
      return reactRedux.useSelector(function (state) {
        return modelSelector(selector(function () {
          return state;
        }).model);
      }, eqFn);
    };
    var useFilter = function useFilter(filterSelector) {
      return reactRedux.useSelector(function (state) {
        return filterSelector(selector(function () {
          return state;
        }).filter);
      }, reactRedux.shallowEqual);
    };
    return {
      useModel: useModel,
      useFilter: useFilter
    };
  };
  BComponent.createSlice = function (_a) {
    var pageName = _a.pageName;
  };
})(BComponent || (BComponent = {}));

var index = 10000;
var timeLead = 1e12;
function uuid() {
  return "app-tke-fe-" + (++index * timeLead + Math.random() * timeLead).toString(36);
}

var _a$5, _b$4;
var ClusterType;
(function (ClusterType) {
  ClusterType["TKE"] = "tke";
  ClusterType["EKS"] = "eks";
  ClusterType["EDGE"] = "tkeedge";
  ClusterType["External"] = "external";
})(ClusterType || (ClusterType = {}));
var ClusterResourceModuleName = (_a$5 = {}, _a$5[ClusterType.TKE] = 'tke', _a$5[ClusterType.EKS] = 'tke', _a$5[ClusterType.EDGE] = 'tke', _a$5[ClusterType.External] = 'tdcc', _a$5);
var ClusterResourceVersionName = (_b$4 = {}, _b$4[ClusterType.TKE] = '2018-05-25', _b$4[ClusterType.EKS] = '2018-05-25', _b$4[ClusterType.EDGE] = '2018-05-25', _b$4[ClusterType.External] = '2022-01-25', _b$4);
var ExternalCluster;
(function (ExternalCluster) {
  var StatusEnum;
  (function (StatusEnum) {
    StatusEnum["Running"] = "Running";
    StatusEnum["Failed"] = "Failed";
    StatusEnum["Waiting"] = "Waiting";
    // Initializing = 'Initializing',
    // Confined = 'Confined',
    // Idling = 'Idling',
    // Upgrading = 'Upgrading',
    StatusEnum["Terminating"] = "Terminating";
    // Upscaling = 'Upscaling',
    // Downscaling = 'Downscaling'
  })(StatusEnum = ExternalCluster.StatusEnum || (ExternalCluster.StatusEnum = {}));
  ExternalCluster.TKEStackDefaultCluster = 'global';
  ExternalCluster.getClusterDetailUrl = function (clusterInfo, regionId) {
    var _a = clusterInfo || {},
      clusterCategory = _a.clusterCategory,
      originalRegion = _a.originalRegion,
      originalClusterId = _a.originalClusterId;
    var url = "/tke2/external/sub/list/basic/info?rid=".concat(regionId, "&clusterId=").concat(clusterInfo === null || clusterInfo === void 0 ? void 0 : clusterInfo.clusterId);
    if (clusterCategory && originalRegion && originalClusterId) {
      switch (clusterCategory) {
        case 'tke':
          url = "/tke2/cluster/sub/list/basic/info?rid=".concat(originalRegion, "&clusterId=").concat(originalClusterId);
          break;
        case 'eks':
          url = "/tke2/ecluster/sub/list/basic/info?rid=".concat(originalRegion, "&clusterId=").concat(originalClusterId);
          break;
        case 'tkeedge':
          url = "/tke2/edge/sub/list/basic/info?rid=".concat(originalRegion, "&clusterId=").concat(originalClusterId);
          break;
      }
    }
    return url;
  };
})(ExternalCluster || (ExternalCluster = {}));

var getWorkflowError = function getWorkflowError(workflow) {
  var _a, _b, _c, _d, _e, _f, _g, _h;
  var message = ((_c = (_b = (_a = workflow === null || workflow === void 0 ? void 0 : workflow.results) === null || _a === void 0 ? void 0 : _a[0]) === null || _b === void 0 ? void 0 : _b.error) === null || _c === void 0 ? void 0 : _c.message) || ((_f = (_e = (_d = workflow === null || workflow === void 0 ? void 0 : workflow.results) === null || _d === void 0 ? void 0 : _d[0]) === null || _e === void 0 ? void 0 : _e.error) === null || _f === void 0 ? void 0 : _f.Message);
  var error = (_h = (_g = workflow === null || workflow === void 0 ? void 0 : workflow.results) === null || _g === void 0 ? void 0 : _g[0]) === null || _h === void 0 ? void 0 : _h.error;
  return message;
};

var BComponent$1;
(function (BComponent) {
  BComponent.BaseActionType = {
    PageName: '',
    Clear: 'Clear',
    Validator: 'Validator'
  };
  BComponent.getActionType = function (moduleName, actionType) {
    return actionType ? moduleName + '_' + actionType : moduleName;
  };
  BComponent.createActionType = function (componentName, actionType) {
    for (var key in actionType) {
      var ns = [];
      componentName && ns.push(componentName);
      actionType[key] && ns.push(actionType[key]);
      actionType[key] = ns.join('_');
    }
  };
  BComponent.isNeedFetch = function (model, filter) {
    if (filter === void 0) {
      filter = null;
    }
    var isAllSame = true;
    if (filter) {
      var targetFilter_1 = [];
      var originFilter_1 = [];
      Object.keys(filter).forEach(function (key) {
        targetFilter_1.push(filter === null || filter === void 0 ? void 0 : filter[key]);
        originFilter_1.push(model.query.filter[key]);
      });
      isAllSame = JSON.stringify(targetFilter_1) === JSON.stringify(originFilter_1);
    }
    if (isAllSame) {
      //如果参数都没变，
      //1.正在加载的就不发请求
      //2.之前的请求没出错就不再发起请求
      if (typeof model['list'] !== 'undefined') {
        if (model['list'].loading || model['list'].fetched
        //&& model['list'].fetchState !== FetchState.Failed
        ) {
          return false;
        } else {
          return true;
        }
      }
      if (typeof model['object'] !== 'undefined') {
        if (model['object'].loading || model['object'].fetched
        // && model['object'].fetchState !== FetchState.Failed
        ) {
          return false;
        } else {
          return true;
        }
      }
    } else {
      return true;
    }
  };
  var ViewStyleEnum;
  (function (ViewStyleEnum) {
    ViewStyleEnum["Panel"] = "Panel";
    ViewStyleEnum["FormItem"] = "FormItem";
    ViewStyleEnum["Alert"] = "Alert";
    ViewStyleEnum["Text"] = "Text";
    ViewStyleEnum["Dialog"] = "Dialog";
  })(ViewStyleEnum = BComponent.ViewStyleEnum || (BComponent.ViewStyleEnum = {}));
  var EditStyleEnum;
  (function (EditStyleEnum) {
    EditStyleEnum["PopConfirm"] = "PopConfirm";
    EditStyleEnum["Panel"] = "Panel";
    EditStyleEnum["FormItem"] = "FormItem";
    EditStyleEnum["Drawer"] = "Drawer";
  })(EditStyleEnum = BComponent.EditStyleEnum || (BComponent.EditStyleEnum = {}));
  var OperationTypeEnum;
  (function (OperationTypeEnum) {
    OperationTypeEnum["Create"] = "Create";
    /**
     * @deprecated
     */
    OperationTypeEnum["CreateNode"] = "CreateNode";
    OperationTypeEnum["Delete"] = "Delete";
    OperationTypeEnum["Refund"] = "Refund";
    OperationTypeEnum["Modify"] = "Modify";
    OperationTypeEnum["View"] = "View";
    /**
     * @deprecated
     */
    OperationTypeEnum["EditSecret"] = "EditSecret";
    /**
     * @deprecated
     */
    OperationTypeEnum["IngressCreateSecret"] = "IngressCreateSecret";
  })(OperationTypeEnum = BComponent.OperationTypeEnum || (BComponent.OperationTypeEnum = {}));
  var ResourceTypeEnum;
  (function (ResourceTypeEnum) {
    ResourceTypeEnum["Asg"] = "Asg";
    ResourceTypeEnum["Node"] = "Node";
    ResourceTypeEnum["MasterEtcdNode"] = "MasterEtcdNode";
    ResourceTypeEnum["NodePool"] = "NodePool";
    ResourceTypeEnum["NativeNodePool"] = "NativeNodePool";
    ResourceTypeEnum["VirtualNodePool"] = "VirtualNodePool";
    ResourceTypeEnum["VirtualNode"] = "VirtualNode";
    ResourceTypeEnum["EdgeCvm"] = "EdgeCvm";
    ResourceTypeEnum["EdgeDeploymentGridInstance"] = "EdgeDeploymentGridInstance";
    ResourceTypeEnum["LogListener"] = "LogListener";
    ResourceTypeEnum["Audit"] = "Audit";
    ResourceTypeEnum["Event"] = "Event";
    ResourceTypeEnum["Ingress"] = "Ingress";
    ResourceTypeEnum["Service"] = "Service";
    ResourceTypeEnum["HPA"] = "HPA";
    ResourceTypeEnum["Workload"] = "Workload";
    ResourceTypeEnum["Secret"] = "Secret";
    ResourceTypeEnum["EKSContainer"] = "EKSContainer";
    ResourceTypeEnum["External"] = "External";
    ResourceTypeEnum["ImageCache"] = "ImageCache";
    ResourceTypeEnum["Subscription"] = "Subscription";
    ResourceTypeEnum["Cluster"] = "Cluster";
    ResourceTypeEnum["ECluster"] = "ECluster";
    ResourceTypeEnum["HubCluster"] = "HubCluster";
    ResourceTypeEnum["EdgeCluster"] = "EdgeCluster";
    ResourceTypeEnum["PersistentVolume"] = "PersistentVolume";
    ResourceTypeEnum["PersistentVolumeClaim"] = "PersistentVolumeClaim";
    ResourceTypeEnum["StorageClass"] = "StorageClass";
    ResourceTypeEnum["Label"] = "Label";
    ResourceTypeEnum["Annotation"] = "Annotation";
  })(ResourceTypeEnum = BComponent.ResourceTypeEnum || (BComponent.ResourceTypeEnum = {}));
  BComponent.createHooks = function (_a) {
    var Context = _a.Context,
      itemsSelector = _a.itemsSelector,
      vKeyPrefix = _a.vKeyPrefix;
    var useModel = function useModel(modelSelector, eqFn) {
      if (eqFn === void 0) {
        eqFn = reactRedux.shallowEqual;
      }
      var context = React__default.useContext(Context);
      return reactRedux.useSelector(function (state) {
        return modelSelector(context.selector(function () {
          return state;
        }).model);
      }, eqFn);
    };
    var useFilter = function useFilter(filterSelector) {
      var context = React__default.useContext(Context);
      return reactRedux.useSelector(function (state) {
        return filterSelector(context.selector(function () {
          return state;
        }).filter);
      }, reactRedux.shallowEqual);
    };
    var useItem = function useItem(itemId, itemSelector) {
      return useModel(function (model) {
        var items = itemsSelector(model);
        var item = items.find(function (item) {
          return item.id === itemId;
        });
        return itemSelector(item);
      });
    };
    var useVkey = function useVkey(_a) {
      var key = _a.key,
        itemId = _a.itemId;
      var itemIndex = useModel(function (model) {
        var items = itemsSelector(model);
        return items.findIndex(function (item) {
          return item.id === itemId;
        });
      });
      return "".concat(vKeyPrefix, "[").concat(itemIndex, "].").concat(String(key));
    };
    return {
      useVkey: useVkey,
      useModel: useModel,
      useFilter: useFilter,
      useItem: useItem
    };
  };
  BComponent.createSimpleHooks = function (selector) {
    var useModel = function useModel(modelSelector, eqFn) {
      if (eqFn === void 0) {
        eqFn = reactRedux.shallowEqual;
      }
      return reactRedux.useSelector(function (state) {
        return modelSelector(selector(function () {
          return state;
        }).model);
      }, eqFn);
    };
    var useFilter = function useFilter(filterSelector) {
      return reactRedux.useSelector(function (state) {
        return filterSelector(selector(function () {
          return state;
        }).filter);
      }, reactRedux.shallowEqual);
    };
    return {
      useModel: useModel,
      useFilter: useFilter
    };
  };
  BComponent.createSlice = function (_a) {
    var pageName = _a.pageName;
  };
})(BComponent$1 || (BComponent$1 = {}));

var Method = {
  get: 'GET',
  "delete": 'DELETE',
  update: 'UPDATE',
  patch: 'PATCH',
  post: 'POST',
  put: 'PUT'
};

var CSRF_TOKEN = null;
function createCSRFHeader() {
  var _a;
  if (CSRF_TOKEN === null) {
    var tkeCookie = (_a = Cookies.get('tke')) !== null && _a !== void 0 ? _a : '';
    CSRF_TOKEN = SparkMD5.hash(tkeCookie);
  }
  return {
    'X-CSRF-TOKEN': CSRF_TOKEN
  };
}

var DefaultRegion = 1;
var util = seajs.require('util');
var getRegionId = function getRegionId() {
  return DefaultRegion;
};
var getProjectName = function getProjectName() {
  var rId = util.cookie.get('projectId');
  return rId;
};

// site 1 为 国内站，2为国际站
var isI18nSite = teaApp.i18n.site === 2;
// zh 为中文，en为英文
var isEngLng = teaApp.i18n.lng === 'en';

/**
 * 统一的请求处理
 * @param userParams: RequestParams
 */
var reduceYunApiNetworkRequest = function reduceYunApiNetworkRequest(userParams, clusterId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var method, url, _a, userDefinedHeader, _b, apiParams, rsp, module, regionId, _c, interfaceName, _d, restParams, opts, params, requestRegionId, error_1;
    var _e, _f;
    return tslib.__generator(this, function (_g) {
      switch (_g.label) {
        case 0:
          method = userParams.method, url = userParams.url, _a = userParams.userDefinedHeader, userDefinedHeader = _a === void 0 ? {} : _a, _b = userParams.data, apiParams = userParams.apiParams;
          // 请求tke-apiserver的 cluster的header
          if (clusterId) {
            userDefinedHeader = Object.assign({}, userDefinedHeader, {
              'X-TKE-ClusterName': clusterId
            });
          }
          module = apiParams.module, regionId = apiParams.regionId, _c = apiParams.interfaceName, interfaceName = _c === void 0 ? 'ForwardRequest' : _c, _d = apiParams.restParams, restParams = _d === void 0 ? {} : _d, opts = apiParams.opts;
          params = restParams;
          // 自定义header，只有forwardRequest 这个请求有用
          if (userDefinedHeader['Accept']) {
            params['Accept'] = userDefinedHeader['Accept'];
          }
          if (userDefinedHeader['Content-Type']) {
            params['ContentType'] = userDefinedHeader['Content-Type'];
          }
          if (userDefinedHeader['X-TKE-ClusterName']) {
            params['ClusterName'] = userDefinedHeader['X-TKE-ClusterName'];
          }
          params['Language'] = !isEngLng ? 'zh-CN' : 'en-US';
          requestRegionId = regionId;
          if (isNaN(requestRegionId)) {
            requestRegionId = getRegionId();
          } else {
            requestRegionId = requestRegionId ? requestRegionId : getRegionId();
          }
          _g.label = 1;
        case 1:
          _g.trys.push([1, 3,, 4]);
          return [4 /*yield*/, teaApp.app.capi.requestV3({
            regionId: requestRegionId,
            serviceType: module,
            cmd: interfaceName,
            data: params
          }, {
            tipErr: (_e = opts === null || opts === void 0 ? void 0 : opts.tipErr) !== null && _e !== void 0 ? _e : true,
            tipLoading: (_f = opts === null || opts === void 0 ? void 0 : opts.global) !== null && _f !== void 0 ? _f : false
          })];
        case 2:
          rsp = _g.sent();
          return [3 /*break*/, 4];
        case 3:
          error_1 = _g.sent();
          //在原来上包了一层data,通过tea.app调用的需要返回error.data
          throw error_1.data;
        case 4:
          return [2 /*return*/, reduceYunApiNetworkResponse(rsp, true)];
      }
    });
  });
};
/**
 * 处理返回的数据
 * @param type  判断当前控制台的类型
 * @param response  请求返回的数据
 */
var reduceYunApiNetworkResponse = function reduceYunApiNetworkResponse(response, isApp, interfaceName) {
  if (isApp === void 0) {
    isApp = false;
  }
  var result;
  // 平台的sdk进行
  if (isApp) {
    var finalData = response.data.Response;
    result = {
      code: 0,
      data: finalData,
      message: 'Success'
    };
  } else {
    // 如果是公有云的话，直接返回所有的请求，这里需要区分v2还是v3接口
    if (response.code !== undefined) {
      result = response;
    } else {
      // 此为V3接口的判断
      var rsp = response.Response;
      // 如果请求发生了错误，其kind为status
      result = {
        code: 0,
        data: rsp,
        message: 'Success'
      };
    }
  }
  return result;
};

/** 获取当前的uuid */
var uuid$1 = function uuid() {
  var s = [];
  var hexDigits = '0123456789abcdef';
  for (var i = 0; i < 36; i++) {
    s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
  }
  s[14] = '4'; // bits 12-15 of the time_hi_and_version field to 0010
  s[19] = hexDigits.substr(s[19] & 0x3 | 0x8, 1); // bits 6-7 of the clock_seq_hi_and_reserved to 01
  s[8] = s[13] = s[18] = s[23] = '-';
  var uuid = s.join('');
  return uuid;
};
/** 获取当前控制台modules 的域名匹配项 */
var GET_CONSOLE_MODULE_BASE_URL = location.origin || '';
/**
 * 统一的请求处理
 * @param userParams: RequestParams
 */
var reduceTkeStackNetworkRequest = function reduceTkeStackNetworkRequest(userParams, clusterId, projectId, keyword) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var method, url, _a, userDefinedHeader, _b, data, apiParams, _c,
      // baseURL = getConsoleAPIAddress(ConsoleModuleAddressEnum.PLATFORM)
      baseURL, rsp, pid, searchParams, params, error_1, uuid_1, response;
    return tslib.__generator(this, function (_d) {
      switch (_d.label) {
        case 0:
          method = userParams.method, url = userParams.url, _a = userParams.userDefinedHeader, userDefinedHeader = _a === void 0 ? {} : _a, _b = userParams.data, data = _b === void 0 ? {} : _b, apiParams = userParams.apiParams, _c = userParams.baseURL, baseURL = _c === void 0 ? GET_CONSOLE_MODULE_BASE_URL : _c;
          // 请求tke-apiserver的 cluster的header
          if (clusterId) {
            userDefinedHeader = Object.assign({}, userDefinedHeader, {
              'X-TKE-ClusterName': clusterId
            });
          }
          pid = projectId;
          try {
            searchParams = parseQueryString(location.search);
          } catch (error) {}
          // 这里指定为undefined而不是''，因为业务视图下helm仓库的逻辑有时候不需要传业务id，但会因为这里的逻辑从cookie中读取业务id并传到后端，
          // 导致过滤逻辑出现问题。调用方会显式指定 projectId = '' 来避免这种情况
          if (pid === undefined) {
            if (searchParams && (searchParams.projectName || searchParams.projectId)) {
              pid = searchParams.projectName || searchParams.projectId;
            } else {
              pid = getProjectName();
            }
          }
          /// #endif
          if (pid) {
            userDefinedHeader = Object.assign({}, {
              'X-TKE-ProjectName': pid
            }, userDefinedHeader);
          }
          if (keyword) {
            userDefinedHeader = Object.assign({}, userDefinedHeader, {
              'X-TKE-FuzzyResourceName': keyword
            });
          }
          console.log(userDefinedHeader, 'userDefinedHeader...');
          params = {
            method: method,
            baseURL: baseURL,
            url: url,
            withCredentials: true,
            headers: Object.assign({}, {
              'X-Remote-Extra-RequestID': uuid$1(),
              'Content-Type': 'application/json'
            }, userDefinedHeader)
          };
          if (data) {
            params = Object.assign({}, params, {
              data: data
            });
          }
          _d.label = 1;
        case 1:
          _d.trys.push([1, 3,, 4]);
          return [4 /*yield*/, axios(params)];
        case 2:
          rsp = _d.sent();
          return [3 /*break*/, 4];
        case 3:
          error_1 = _d.sent();
          // 如果返回是 401的话，自动登出，此时是鉴权不过，cookies失效了
          if (error_1.response && error_1.response.status === 401) {
            location.reload();
          } else if (error_1.response && error_1.response.status === 403) {
            // changeForbiddentConfig({
            //   isShow: true,
            //   message: error.response.data.message
            // });
            throw error_1;
          } else if (error_1.response === undefined) {
            uuid_1 = error_1.config && error_1.config.headers && error_1.config.headers['X-Remote-Extra-RequestID'] ? error_1.config.headers['X-Remote-Extra-RequestID'] : '';
            error_1.response = {
              data: {
                message: "\u7CFB\u7EDF\u5185\u90E8\u670D\u52A1\u9519\u8BEF\uFF08".concat(uuid_1, "\uFF09")
              }
            };
            throw error_1;
          } else {
            throw error_1;
          }
          return [3 /*break*/, 4];
        case 4:
          response = reduceTkeStackNetworkResponse(rsp);
          return [2 /*return*/, response];
      }
    });
  });
};
/**
 * 处理返回的数据
 * @param type  判断当前控制台的类型
 * @param response  请求返回的数据
 */
var reduceTkeStackNetworkResponse = function reduceTkeStackNetworkResponse(response) {
  if (response === void 0) {
    response = {};
  }
  var result;
  result = {
    code: response.status >= 200 && response.status < 300 ? 0 : response.status,
    data: response.data,
    message: response.statusText
  };
  return result;
};

var reduceNetworkRequest = function reduceNetworkRequest(params) {
  var newParams = tslib.__assign({}, params);
  var platform = params === null || params === void 0 ? void 0 : params.platform;
  if (platform === PlatformType.TDCC) {
    return reduceYunApiNetworkRequest(newParams);
  } else if (platform === PlatformType.TKESTACK) {
    var _a = newParams === null || newParams === void 0 ? void 0 : newParams.restParams,
      ClusterId = _a.ClusterId,
      ProjectId = _a.ProjectId;
    return reduceTkeStackNetworkRequest(tslib.__assign(tslib.__assign({}, newParams), {
      userDefinedHeader: tslib.__assign(tslib.__assign({}, newParams.userDefinedHeader), createCSRFHeader())
    }), ClusterId, ProjectId);
  } else {
    return reduceYunApiNetworkRequest(newParams);
  }
};
/**
 * 处理workflow的返回结果
 * @param target T[]
 * @param error any
 */
var operationResult = function operationResult(target, error, response) {
  if (target instanceof Array) {
    return target.map(function (x) {
      return {
        success: !error,
        target: x,
        error: error,
        response: response
      };
    });
  }
  return [{
    success: !error,
    target: target,
    error: error,
    response: response
  }];
};

var tips$1 = seajs.require('tips');
var RequestResult = /** @class */function () {
  function RequestResult() {}
  return RequestResult;
}();
var RequestApi;
(function (RequestApi) {
  var _this = this;
  RequestApi.SEND = function (args) {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var params, resp, response, error_1, errorMsg;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            params = tslib.__assign({}, args);
            resp = new RequestResult();
            _a.label = 1;
          case 1:
            _a.trys.push([1, 3,, 4]);
            return [4 /*yield*/, reduceNetworkRequest(params)];
          case 2:
            response = _a.sent();
            if (response.code !== 0 &&
            // (
            //   response.code !== 'ResourceNotFound' || 
            //   !response.message?.includes('404')
            // ) && 
            args.method !== Method.get) {
              tips$1.error(response === null || response === void 0 ? void 0 : response.message, 2000);
              resp.data = args.data;
              resp.error = response === null || response === void 0 ? void 0 : response.message;
            } else {
              resp.data = response.data;
              resp.error = null;
              resp.code = response.code;
            }
            return [2 /*return*/, resp];
          case 3:
            error_1 = _a.sent();
            errorMsg = void 0;
            if (args.method !== Method.get) {
              errorMsg = "\u64CD\u4F5C\u5931\u8D25:" + (error_1 === null || error_1 === void 0 ? void 0 : error_1.message);
            } else {
              errorMsg = error_1 === null || error_1 === void 0 ? void 0 : error_1.message;
            }
            if (
            // (
            //   error.code !== 'ResourceNotFound' || 
            //   !error.message?.includes('404')
            // ) && 
            args.method !== Method.get) {
              tips$1.error(errorMsg, 2000);
            }
            resp.data = args.data;
            resp.error = error_1 === null || error_1 === void 0 ? void 0 : error_1.message;
            resp.message = error_1 === null || error_1 === void 0 ? void 0 : error_1.message;
            resp.code = error_1 === null || error_1 === void 0 ? void 0 : error_1.code;
            return [2 /*return*/, resp];
          case 4:
            return [2 /*return*/];
        }
      });
    });
  };

  RequestApi.GET = function (args) {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var response;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            args.method = Method.get;
            args.data = null;
            return [4 /*yield*/, RequestApi.SEND(args)];
          case 1:
            response = _a.sent();
            return [2 /*return*/, response];
        }
      });
    });
  };
  RequestApi.DELETE = function (args) {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var response;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            args.method = Method["delete"];
            args.data = null;
            return [4 /*yield*/, RequestApi.SEND(args)];
          case 1:
            response = _a.sent();
            return [2 /*return*/, response];
        }
      });
    });
  };
  RequestApi.POST = function (args) {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var response;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            args.method = Method.post;
            args.data = JSON.stringify(args.data);
            return [4 /*yield*/, RequestApi.SEND(args)];
          case 1:
            response = _a.sent();
            return [2 /*return*/, response];
        }
      });
    });
  };
  RequestApi.PUT = function (args) {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var response;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            args.method = Method.put;
            args.data = JSON.stringify(args.data);
            return [4 /*yield*/, RequestApi.SEND(args)];
          case 1:
            response = _a.sent();
            return [2 /*return*/, response];
        }
      });
    });
  };
  RequestApi.PATCH = function (args) {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var response;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            args.method = Method.patch;
            args.data = JSON.stringify(args.data);
            return [4 /*yield*/, RequestApi.SEND(args)];
          case 1:
            response = _a.sent();
            return [2 /*return*/, response];
        }
      });
    });
  };
})(RequestApi || (RequestApi = {}));

var tips$2 = seajs.require('tips');
// 获取应用实例列表
var fetchHubCluster = function fetchHubCluster(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, regionId, platform, params, response, result;
    var _a, _b, _c, _d;
    return tslib.__generator(this, function (_e) {
      switch (_e.label) {
        case 0:
          userDefinedHeader = {};
          regionId = queryParams.regionId, platform = queryParams.platform;
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: YunApiName.DescribeHubClusters,
                regionId: regionId,
                restParams: {
                  Version: ResourceVersionName[PlatformType.TDCC],
                  Filters: [],
                  Offset: 0,
                  Limit: 1000
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              }
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: "/api/v1/namespaces/ssm/configmaps/supported-vendors",
              restParams: {
                ProjectId: ''
              }
            };
          }
          return [4 /*yield*/, RequestApi.GET(params)];
        case 1:
          response = _e.sent();
          result = {
            records: [],
            recordCount: 0
          };
          if (platform === PlatformType.TDCC) {
            result.records = (_c = (_b = (_a = response === null || response === void 0 ? void 0 : response.data) === null || _a === void 0 ? void 0 : _a.Clusters) === null || _b === void 0 ? void 0 : _b.map(function (cluster) {
              return {
                id: cluster === null || cluster === void 0 ? void 0 : cluster.ClusterId,
                clusterId: cluster === null || cluster === void 0 ? void 0 : cluster.ClusterId,
                name: cluster === null || cluster === void 0 ? void 0 : cluster.ClusterName,
                description: cluster === null || cluster === void 0 ? void 0 : cluster.ClusterDesc,
                k8sVersion: cluster === null || cluster === void 0 ? void 0 : cluster.K8SVersion,
                status: cluster === null || cluster === void 0 ? void 0 : cluster.Status,
                unVpcId: cluster === null || cluster === void 0 ? void 0 : cluster.VpcId,
                subnetId: cluster === null || cluster === void 0 ? void 0 : cluster.SubnetId,
                regionId: cluster === null || cluster === void 0 ? void 0 : cluster.RegionId
              };
            })) !== null && _c !== void 0 ? _c : [];
          } else {
            result.records = [{
              id: 'global',
              clusterId: 'global',
              name: 'global'
            }];
          }
          result.recordCount = (_d = result === null || result === void 0 ? void 0 : result.records) === null || _d === void 0 ? void 0 : _d.length;
          return [2 /*return*/, result];
      }
    });
  });
};
// 获取应用实例列表
var fetchOpenedServices = function fetchOpenedServices(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, regionId, platform, params, namespace, response, result, error_1;
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
    return tslib.__generator(this, function (_m) {
      switch (_m.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform;
          namespace = SystemNamespace;
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: YunApiName.DescribeServiceVendors,
                regionId: regionId,
                restParams: {
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: []
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: "/api/v1/namespaces/".concat(namespace, "/configmaps/supported-vendors"),
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          result = {
            records: [],
            recordCount: 0
          };
          _m.label = 1;
        case 1:
          _m.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _m.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              result.records = (_c = (_b = (_a = response === null || response === void 0 ? void 0 : response.data) === null || _a === void 0 ? void 0 : _a.Vendors) === null || _b === void 0 ? void 0 : _b.map(function (item) {
                return {
                  name: item === null || item === void 0 ? void 0 : item.ClassName,
                  enabled: item === null || item === void 0 ? void 0 : item.Enabled,
                  instanceCount: item === null || item === void 0 ? void 0 : item.InstanceCount,
                  version: item === null || item === void 0 ? void 0 : item.Version,
                  clusters: item === null || item === void 0 ? void 0 : item.Clusters
                };
              })) !== null && _c !== void 0 ? _c : [];
            } else {
              result.records = (_g = (_f = (_e = (_d = response === null || response === void 0 ? void 0 : response.data) === null || _d === void 0 ? void 0 : _d.data) === null || _e === void 0 ? void 0 : _e.vendors) === null || _f === void 0 ? void 0 : _f.split(',').map(function (item) {
                return {
                  name: item,
                  enabled: true,
                  instanceCount: 0,
                  version: '',
                  clusters: [ExternalCluster === null || ExternalCluster === void 0 ? void 0 : ExternalCluster.TKEStackDefaultCluster]
                };
              })) !== null && _g !== void 0 ? _g : [];
            }
            // 过滤掉未开启的中间件服务
            result.records = (_h = result === null || result === void 0 ? void 0 : result.records) === null || _h === void 0 ? void 0 : _h.filter(function (item) {
              return item === null || item === void 0 ? void 0 : item.enabled;
            });
            result.recordCount = (_j = result === null || result === void 0 ? void 0 : result.records) === null || _j === void 0 ? void 0 : _j.length;
          } else if ((response === null || response === void 0 ? void 0 : response.code) === 404 || ((_k = response === null || response === void 0 ? void 0 : response.message) === null || _k === void 0 ? void 0 : _k.includes('404'))) {
            result['hasError'] = platform === PlatformType.TKESTACK ? false : true;
          } else {
            result['hasError'] = true;
          }
          return [3 /*break*/, 4];
        case 3:
          error_1 = _m.sent();
          if ((error_1 === null || error_1 === void 0 ? void 0 : error_1.code) === 404 || ((_l = error_1 === null || error_1 === void 0 ? void 0 : error_1.message) === null || _l === void 0 ? void 0 : _l.includes('404'))) {
            result['hasError'] = platform === PlatformType.TKESTACK ? false : true;
          } else {
            result['hasError'] = true;
            throw error_1;
          }
          return [3 /*break*/, 4];
        case 4:
          return [2 /*return*/, result];
      }
    });
  });
};
var fetchServiceResources = function fetchServiceResources(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, regionId, platform, serviceName, _a, resourceType, paging, search, k8sQueryObj, params, extendParams, result, url, queryString, response, instances, error_2, keys, propertyKeys_1, pageIndex, pageSize, maxPageIndex;
    var _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v;
    return tslib.__generator(this, function (_w) {
      switch (_w.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, serviceName = queryParams.serviceName, _a = queryParams.resourceType, resourceType = _a === void 0 ? ResourceTypeEnum.ServiceResource : _a, paging = queryParams.paging, search = queryParams.search, k8sQueryObj = queryParams.k8sQueryObj;
          extendParams = {};
          result = {
            records: [],
            recordCount: 0
          };
          if ([ResourceTypeEnum.Backup, ResourceTypeEnum.ServiceBinding].includes(resourceType)) {
            extendParams = {
              Path: "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(SystemNamespace, "/").concat((_b = ResourceTypeMap[resourceType]) === null || _b === void 0 ? void 0 : _b.path),
              ServiceClass: undefined,
              Method: Method === null || Method === void 0 ? void 0 : Method.get
            };
          } else {
            extendParams = {
              Path: undefined,
              ServiceClass: serviceName,
              Method: undefined
            };
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ServiceMngYunApiName[resourceType],
                regionId: regionId,
                restParams: tslib.__assign({
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId
                }, extendParams),
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/").concat((_c = ResourceTypeMap[resourceType]) === null || _c === void 0 ? void 0 : _c.path);
            if (k8sQueryObj) {
              queryString = reduceK8sQueryString({
                k8sQueryObj: k8sQueryObj,
                restfulPath: url
              });
              url = url + queryString;
            }
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          _w.label = 1;
        case 1:
          _w.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _w.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              instances = JSON.parse(resourceType === ResourceTypeEnum.ServiceResource ? (_d = response === null || response === void 0 ? void 0 : response.data) === null || _d === void 0 ? void 0 : _d.Instances : (_e = response === null || response === void 0 ? void 0 : response.data) === null || _e === void 0 ? void 0 : _e.Plans);
              result.records = (_f = instances === null || instances === void 0 ? void 0 : instances.items) !== null && _f !== void 0 ? _f : [];
            } else {
              result.records = (_j = (_h = (_g = response === null || response === void 0 ? void 0 : response.data) === null || _g === void 0 ? void 0 : _g.items) === null || _h === void 0 ? void 0 : _h.filter(function (item) {
                var _a;
                return ((_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass) === serviceName;
              })) !== null && _j !== void 0 ? _j : [];
            }
          } else {
            result['hasError'] = true;
          }
          return [3 /*break*/, 4];
        case 3:
          error_2 = _w.sent();
          result['hasError'] = true;
          return [3 /*break*/, 4];
        case 4:
          // 前端过滤
          if (search) {
            result.records = (_k = result.records) === null || _k === void 0 ? void 0 : _k.filter(function (item) {
              var _a, _b;
              return (_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === null || _b === void 0 ? void 0 : _b.includes(search);
            });
          }
          // 非metadata?.name属性，需要前端自己进行过滤显示
          if (k8sQueryObj && (k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector) && !((_l = Object === null || Object === void 0 ? void 0 : Object.keys(k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector)) === null || _l === void 0 ? void 0 : _l.includes('metadata.name'))) {
            keys = Object === null || Object === void 0 ? void 0 : Object.keys(k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector);
            propertyKeys_1 = (_m = keys[0]) === null || _m === void 0 ? void 0 : _m.split('.');
            result.records = (_o = result.records) === null || _o === void 0 ? void 0 : _o.filter(function (item) {
              var _a;
              return (_a = Object.keys(k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector)) === null || _a === void 0 ? void 0 : _a.every(function (key) {
                var _a, _b;
                return ((_a = item === null || item === void 0 ? void 0 : item[propertyKeys_1 === null || propertyKeys_1 === void 0 ? void 0 : propertyKeys_1[0]]) === null || _a === void 0 ? void 0 : _a[propertyKeys_1[1]]) === ((_b = k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector) === null || _b === void 0 ? void 0 : _b[key]);
              });
            });
          }
          // 过滤
          // if(instanceId){
          //   result.records = result?.records?.filter((item) => {
          //     return resourceType === ResourceTypeEnum?.Backup ? 
          //       (item?.spec?.target?.instanceID === instanceId || 
          //         item?.metadata?.labels?.['ssm.infra.tce.io/instanceId'] === instanceId || 
          //         item?.metadata?.labels?.['ssm.infra.tce.io/instance-id'] === instanceId
          //       ) 
          //       :
          //       (
          //         item?.metadata?.labels?.['ssm.infra.tce.io/instanceId'] === instanceId ||
          //         item?.metadata?.labels?.['ssm.infra.tce.io/instance-id'] === instanceId
          //       );
          //   });
          // }
          //servicePlan首先按照规格类型自定义、预设顺序排序
          if (resourceType === ResourceTypeEnum.ServicePlan) {
            result.records = (_p = result.records) === null || _p === void 0 ? void 0 : _p.sort(function (pre, cur) {
              var _a, _b, _c, _d, _e;
              return (_c = (_b = (_a = pre === null || pre === void 0 ? void 0 : pre.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['ssm.infra.tce.io/owner']) === null || _c === void 0 ? void 0 : _c.localeCompare((_e = (_d = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _d === void 0 ? void 0 : _d.labels) === null || _e === void 0 ? void 0 : _e['ssm.infra.tce.io/owner'], 'zh-CN');
            });
          } else {
            // 按照创建时间降序排列
            result.records = (_q = result.records) === null || _q === void 0 ? void 0 : _q.sort(function (pre, cur) {
              var _a, _b, _c, _d;
              return ((_b = new Date((_a = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp)) === null || _b === void 0 ? void 0 : _b.getTime()) - ((_d = new Date((_c = pre === null || pre === void 0 ? void 0 : pre.metadata) === null || _c === void 0 ? void 0 : _c.creationTimestamp)) === null || _d === void 0 ? void 0 : _d.getTime());
            });
          }
          // 统计总数
          result.recordCount = (_s = (_r = result.records) === null || _r === void 0 ? void 0 : _r.length) !== null && _s !== void 0 ? _s : 0;
          // 前端分页
          if (paging) {
            pageIndex = (_t = paging === null || paging === void 0 ? void 0 : paging.pageIndex) !== null && _t !== void 0 ? _t : 1;
            pageSize = (_u = paging === null || paging === void 0 ? void 0 : paging.pageSize) !== null && _u !== void 0 ? _u : 20;
            maxPageIndex = Math.ceil(((_v = result.records) === null || _v === void 0 ? void 0 : _v.length) / pageSize);
            if (pageIndex > maxPageIndex) {
              pageIndex = maxPageIndex;
            }
            if (pageIndex === 1) {
              result.records = result.records.slice(pageIndex - 1, pageIndex * pageSize);
            } else {
              result.records = result.records.slice((pageIndex - 1) * pageSize, pageIndex * pageSize);
            }
          }
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 * 查询服务资源Schema
 * @param queryParams
 * @returns
 */
var fetchResourceSchemas = function fetchResourceSchemas(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, regionId, platform, serviceName, clusterId, params, url, response, result, responseBody, error_3;
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u;
    return tslib.__generator(this, function (_v) {
      switch (_v.label) {
        case 0:
          userDefinedHeader = {};
          regionId = queryParams.regionId, platform = queryParams.platform, serviceName = queryParams.serviceName, clusterId = queryParams.clusterId;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/serviceclasses/").concat(serviceName);
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.get,
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: []
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          _v.label = 1;
        case 1:
          _v.trys.push([1, 3,, 4]);
          responseBody = void 0;
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _v.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              responseBody = JSON.parse((_a = response === null || response === void 0 ? void 0 : response.data) === null || _a === void 0 ? void 0 : _a.ResponseBody);
            } else {
              responseBody = response === null || response === void 0 ? void 0 : response.data;
            }
            result = {
              // 后端暂不支持kafka的依赖 zk 的来源为new的场景，故临时过滤掉new选项
              instanceCreateParameterSchema: (_c = (_b = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _b === void 0 ? void 0 : _b.instanceCreateParameterSchema) === null || _c === void 0 ? void 0 : _c.map(function (item) {
                var _a;
                var candidates = (_a = item === null || item === void 0 ? void 0 : item.candidates) === null || _a === void 0 ? void 0 : _a.filter(function (candidate) {
                  if ((item === null || item === void 0 ? void 0 : item.name) === 'dep_zk_source') {
                    return candidate !== 'new';
                  } else {
                    return true;
                  }
                });
                var hideSchema = (item === null || item === void 0 ? void 0 : item.name) === 'dep_zk_source';
                return tslib.__assign(tslib.__assign({}, item), {
                  value: formatPlanSchemaData(item === null || item === void 0 ? void 0 : item.type, hideSchema ? candidates === null || candidates === void 0 ? void 0 : candidates[0] : !!(candidates === null || candidates === void 0 ? void 0 : candidates.length) ? candidates === null || candidates === void 0 ? void 0 : candidates[0] : item === null || item === void 0 ? void 0 : item["default"]),
                  unit: formatPlanSchemaUnit(item, hideSchema ? candidates === null || candidates === void 0 ? void 0 : candidates[0] : !!(candidates === null || candidates === void 0 ? void 0 : candidates.length) ? candidates === null || candidates === void 0 ? void 0 : candidates[0] : item === null || item === void 0 ? void 0 : item["default"]),
                  candidates: candidates,
                  "default": hideSchema ? candidates === null || candidates === void 0 ? void 0 : candidates[0] : !!(candidates === null || candidates === void 0 ? void 0 : candidates.length) ? candidates === null || candidates === void 0 ? void 0 : candidates[0] : item === null || item === void 0 ? void 0 : item["default"]
                });
              }),
              instanceUpdateParameterSchema: (_e = (_d = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _d === void 0 ? void 0 : _d.instanceUpdateParameterSchema) === null || _e === void 0 ? void 0 : _e.map(function (item) {
                return tslib.__assign(tslib.__assign({}, item), {
                  value: formatPlanSchemaData(item === null || item === void 0 ? void 0 : item.type, item === null || item === void 0 ? void 0 : item["default"]),
                  unit: formatPlanSchemaUnit(item, item === null || item === void 0 ? void 0 : item["default"])
                });
              }),
              planSchema: (_g = (_f = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _f === void 0 ? void 0 : _f.planSchema) === null || _g === void 0 ? void 0 : _g.map(function (item) {
                return tslib.__assign(tslib.__assign({}, item), {
                  value: formatPlanSchemaData(item === null || item === void 0 ? void 0 : item.type, item === null || item === void 0 ? void 0 : item["default"]),
                  unit: formatPlanSchemaUnit(item, item === null || item === void 0 ? void 0 : item["default"])
                });
              }),
              released: (_h = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _h === void 0 ? void 0 : _h.released,
              monitoring: (_j = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _j === void 0 ? void 0 : _j.monitoring,
              vendor: (_k = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _k === void 0 ? void 0 : _k.vendor,
              metadata: (_l = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _l === void 0 ? void 0 : _l.metadata,
              bindingResponseSchema: (_o = (_m = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _m === void 0 ? void 0 : _m.bindingResponseSchema) === null || _o === void 0 ? void 0 : _o.map(function (item) {
                return tslib.__assign(tslib.__assign({}, item), {
                  value: formatPlanSchemaData(item === null || item === void 0 ? void 0 : item.type, item === null || item === void 0 ? void 0 : item["default"]),
                  unit: formatPlanSchemaUnit(item, item === null || item === void 0 ? void 0 : item["default"])
                });
              }),
              bindingCreateParameterSchema: (_q = (_p = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _p === void 0 ? void 0 : _p.bindingCreateParameterSchema) === null || _q === void 0 ? void 0 : _q.map(function (item) {
                return tslib.__assign(tslib.__assign({}, item), {
                  value: formatPlanSchemaData(item === null || item === void 0 ? void 0 : item.type, item === null || item === void 0 ? void 0 : item["default"]),
                  unit: formatPlanSchemaUnit(item, item === null || item === void 0 ? void 0 : item["default"])
                });
              }),
              description: (_r = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _r === void 0 ? void 0 : _r.description,
              plans: (_s = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _s === void 0 ? void 0 : _s.plans,
              name: (_t = responseBody === null || responseBody === void 0 ? void 0 : responseBody.metadata) === null || _t === void 0 ? void 0 : _t.name,
              supportedOperations: (_u = responseBody === null || responseBody === void 0 ? void 0 : responseBody.spec) === null || _u === void 0 ? void 0 : _u.supportedOperations
            };
          } else {
            throw response;
          }
          return [3 /*break*/, 4];
        case 3:
          error_3 = _v.sent();
          throw error_3;
        case 4:
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 * 查询注册集群列表
 * @param queryParams
 * @returns
 */
var fetchExternalClusters = function fetchExternalClusters(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, regionId, platform, clusterIds, _a, regional, params, url, response, result, error_5;
    var _b, _c, _d, _e, _f, _g, _h;
    return tslib.__generator(this, function (_j) {
      switch (_j.label) {
        case 0:
          userDefinedHeader = {};
          regionId = queryParams.regionId, platform = queryParams.platform, clusterIds = queryParams.clusterIds, _a = queryParams.regional, regional = _a === void 0 ? false : _a;
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: YunApiName.DescribeExternalClusters,
                regionId: regionId,
                restParams: {
                  Version: ResourceVersionName[PlatformType.TDCC],
                  Filters: [],
                  Offset: 0,
                  Limit: 1000,
                  Regional: false,
                  ClusterIds: clusterIds
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            url = "/apis/platform.tkestack.io/v1/clusters";
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: undefined,
                ProjectId: ''
              },
              platform: platform
            };
          }
          result = {
            records: [],
            recordCount: 0
          };
          _j.label = 1;
        case 1:
          _j.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _j.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              result.records = (_d = (_c = (_b = response === null || response === void 0 ? void 0 : response.data) === null || _b === void 0 ? void 0 : _b.Clusters) === null || _c === void 0 ? void 0 : _c.map(function (item) {
                return {
                  clusterId: item === null || item === void 0 ? void 0 : item.ClusterId,
                  status: item === null || item === void 0 ? void 0 : item.Status,
                  clusterName: item === null || item === void 0 ? void 0 : item.ClusterName,
                  extraInfos: item === null || item === void 0 ? void 0 : item.ExtraInfos
                };
              })) !== null && _d !== void 0 ? _d : [];
            } else {
              result.records = (_g = (_f = (_e = response === null || response === void 0 ? void 0 : response.data) === null || _e === void 0 ? void 0 : _e.items) === null || _f === void 0 ? void 0 : _f.map(function (item) {
                var _a, _b, _c;
                return {
                  clusterId: (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name,
                  status: (_b = item === null || item === void 0 ? void 0 : item.status) === null || _b === void 0 ? void 0 : _b.phase,
                  clusterName: (_c = item === null || item === void 0 ? void 0 : item.spec) === null || _c === void 0 ? void 0 : _c.displayName
                };
              })) === null || _g === void 0 ? void 0 : _g.filter(function (item) {
                return (item === null || item === void 0 ? void 0 : item.clusterId) === (ExternalCluster === null || ExternalCluster === void 0 ? void 0 : ExternalCluster.TKEStackDefaultCluster);
              });
            }
          } else {
            result['hasError'] = true;
          }
          return [3 /*break*/, 4];
        case 3:
          error_5 = _j.sent();
          result['hasError'] = true;
          return [3 /*break*/, 4];
        case 4:
          result.recordCount = (_h = result === null || result === void 0 ? void 0 : result.records) === null || _h === void 0 ? void 0 : _h.length;
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 * 新建服务资源
 * @param queryParams
 * @returns
 */
var createServiceResource = function createServiceResource(resource, regionId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, _a, clusterId, jsonData, platform, _b, resourceType, instanceName, specificOperate, _c, namespace, params, url, response;
    var _d, _e;
    return tslib.__generator(this, function (_f) {
      switch (_f.label) {
        case 0:
          userDefinedHeader = {};
          _a = resource === null || resource === void 0 ? void 0 : resource[0], clusterId = _a.clusterId, jsonData = _a.jsonData, platform = _a.platform, _b = _a.resourceType, resourceType = _b === void 0 ? ResourceTypeEnum.ServiceResource : _b, instanceName = _a.instanceName, specificOperate = _a.specificOperate, _c = _a.namespace, namespace = _c === void 0 ? SystemNamespace : _c;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(namespace, "/").concat((_d = ResourceTypeMap[resourceType]) === null || _d === void 0 ? void 0 : _d.path);
          if (resourceType === ResourceTypeEnum.ServicePlan) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/").concat((_e = ResourceTypeMap[resourceType]) === null || _e === void 0 ? void 0 : _e.path);
          }
          if (instanceName) {
            url += "/".concat(instanceName);
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.post,
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: [],
                  RequestBody: jsonData
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform,
              method: Method.post,
              data: JSON.parse(jsonData)
            };
          }
          return [4 /*yield*/, RequestApi.POST(params)];
        case 1:
          response = _f.sent();
          try {
            if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
              specificOperate && specificOperate !== CreateSpecificOperatorEnum.CreateResource && tips$2.success(CreateSpecificOperatorMap === null || CreateSpecificOperatorMap === void 0 ? void 0 : CreateSpecificOperatorMap[specificOperate].msg);
              return [2 /*return*/, operationResult(resource)];
            }
            return [2 /*return*/, operationResult(resource, response)];
          } catch (error) {
            return [2 /*return*/, operationResult(resource, error)];
          }
          return [2 /*return*/];
      }
    });
  });
};
/**
 * 编辑服务资源
 * @param queryParams
 * @returns
 */
var updateServiceResource = function updateServiceResource(resource, regionId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, _a, clusterId, jsonData, platform, _b, resourceType, instanceName, _c, namespace, params, url, response;
    var _d, _e;
    return tslib.__generator(this, function (_f) {
      switch (_f.label) {
        case 0:
          userDefinedHeader = {
            'Content-Type': 'application/merge-patch+json'
          };
          _a = resource === null || resource === void 0 ? void 0 : resource[0], clusterId = _a.clusterId, jsonData = _a.jsonData, platform = _a.platform, _b = _a.resourceType, resourceType = _b === void 0 ? ResourceTypeEnum.ServiceResource : _b, instanceName = _a.instanceName, _c = _a.namespace, namespace = _c === void 0 ? SystemNamespace : _c;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(namespace, "/").concat((_d = ResourceTypeMap[resourceType]) === null || _d === void 0 ? void 0 : _d.path);
          if (resourceType === ResourceTypeEnum.ServicePlan) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/").concat((_e = ResourceTypeMap[resourceType]) === null || _e === void 0 ? void 0 : _e.path);
          }
          if (instanceName) {
            url += "/".concat(instanceName);
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.patch,
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: [],
                  RequestBody: jsonData
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              method: Method.patch,
              data: JSON.parse(jsonData),
              platform: platform
            };
          }
          return [4 /*yield*/, RequestApi.PATCH(params)];
        case 1:
          response = _f.sent();
          try {
            if (platform === PlatformType.TDCC) {} else {}
            if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
              tips$2.success(i18n.t('编辑成功'));
              return [2 /*return*/, operationResult(resource)];
            }
            return [2 /*return*/, operationResult(resource, response)];
          } catch (error) {
            return [2 /*return*/, operationResult(resource, error)];
          }
          return [2 /*return*/];
      }
    });
  });
};

var deleteMulServiceResource = function deleteMulServiceResource(resource, regionId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var _a, clusterId, platform, resourceType, resourceInfos, namespace, allRequests, result, errorResponse;
    return tslib.__generator(this, function (_b) {
      switch (_b.label) {
        case 0:
          _a = resource === null || resource === void 0 ? void 0 : resource[0], clusterId = _a.clusterId, platform = _a.platform, resourceType = _a.resourceType, resourceInfos = _a.resourceInfos, namespace = _a.namespace;
          allRequests = resourceInfos === null || resourceInfos === void 0 ? void 0 : resourceInfos.map(function (item) {
            var _a;
            return (_a = deleteServiceResource({
              platform: platform,
              clusterId: clusterId,
              regionId: regionId,
              resourceType: resourceType,
              resourceInfo: item,
              namespace: namespace
            }, regionId)) === null || _a === void 0 ? void 0 : _a.then(function (response) {
              return response;
            }, function (error) {
              return error;
            });
          });
          return [4 /*yield*/, Promise.all(allRequests)];
        case 1:
          result = _b.sent();
          errorResponse = result === null || result === void 0 ? void 0 : result.find(function (item) {
            var _a, _b;
            return ((_a = item === null || item === void 0 ? void 0 : item[0]) === null || _a === void 0 ? void 0 : _a.error) || !((_b = item === null || item === void 0 ? void 0 : item[0]) === null || _b === void 0 ? void 0 : _b.success);
          });
          if (!errorResponse) {
            tips$2.success(i18n.t('删除成功'));
            return [2 /*return*/, operationResult(resource)];
          }
          return [2 /*return*/, operationResult(resource, errorResponse === null || errorResponse === void 0 ? void 0 : errorResponse.error)];
      }
    });
  });
};
/**
 * 删除服务资源
 * @param queryParams
 * @returns
 */
var deleteServiceResource = function deleteServiceResource(resource, regionId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, platform, resourceInfo, resourceType, instanceName, namespace, params, url, response;
    var _a, _b, _c, _d, _e, _f, _g;
    return tslib.__generator(this, function (_h) {
      switch (_h.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = resource.clusterId, platform = resource.platform, resourceInfo = resource.resourceInfo;
          resourceType = (_a = resource === null || resource === void 0 ? void 0 : resource.resourceType) !== null && _a !== void 0 ? _a : resourceInfo === null || resourceInfo === void 0 ? void 0 : resourceInfo.kind;
          instanceName = (_b = resource === null || resource === void 0 ? void 0 : resource.resourceIns) !== null && _b !== void 0 ? _b : (_c = resourceInfo === null || resourceInfo === void 0 ? void 0 : resourceInfo.metadata) === null || _c === void 0 ? void 0 : _c.name;
          namespace = (_d = resource === null || resource === void 0 ? void 0 : resource.namespace) !== null && _d !== void 0 ? _d : (_e = resourceInfo === null || resourceInfo === void 0 ? void 0 : resourceInfo.metadata) === null || _e === void 0 ? void 0 : _e.namespace;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(namespace, "/").concat((_f = ResourceTypeMap[resourceType]) === null || _f === void 0 ? void 0 : _f.path);
          if (resourceType === ResourceTypeEnum.ServicePlan) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/").concat((_g = ResourceTypeMap[resourceType]) === null || _g === void 0 ? void 0 : _g.path);
          }
          if (instanceName) {
            url += "/".concat(instanceName);
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method["delete"],
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: []
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform,
              method: Method["delete"]
            };
          }
          return [4 /*yield*/, RequestApi.DELETE(params)];
        case 1:
          response = _h.sent();
          try {
            if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
              return [2 /*return*/, operationResult(resource)];
            }
            return [2 /*return*/, operationResult(resource, response)];
          } catch (error) {
            return [2 /*return*/, operationResult(resource, error)];
          }
          return [2 /*return*/];
      }
    });
  });
};
/**
 * 查询服务资源详情
 * @param queryParams
 * @returns
 */
var fetchServiceResourceDetail = function fetchServiceResourceDetail(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, regionId, platform, serviceName, resourceIns, _a, resourceType, _b, namespace, params, url, response, result, responseBody;
    var _c, _d, _e, _f, _g, _h, _j, _k;
    return tslib.__generator(this, function (_l) {
      switch (_l.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, serviceName = queryParams.serviceName, resourceIns = queryParams.resourceIns, _a = queryParams.resourceType, resourceType = _a === void 0 ? ResourceTypeEnum.ServiceResource : _a, _b = queryParams.namespace, namespace = _b === void 0 ? SystemNamespace : _b;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(namespace, "/").concat((_c = ResourceTypeMap[resourceType]) === null || _c === void 0 ? void 0 : _c.path);
          if (resourceType === ResourceTypeEnum.ServicePlan) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/").concat((_d = ResourceTypeMap[resourceType]) === null || _d === void 0 ? void 0 : _d.path);
          }
          if (resourceType === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Secret)) {
            url = "/api/v1/namespaces/".concat(namespace, "/secrets");
          }
          if (resourceIns) {
            url += "?fieldSelector=metadata.name=".concat(resourceIns);
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.get,
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: []
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          return [4 /*yield*/, RequestApi.GET(params)];
        case 1:
          response = _l.sent();
          try {
            if (platform === PlatformType.TDCC) {
              responseBody = JSON.parse((_e = response === null || response === void 0 ? void 0 : response.data) === null || _e === void 0 ? void 0 : _e.ResponseBody);
              result = (_g = (_f = responseBody === null || responseBody === void 0 ? void 0 : responseBody.items) === null || _f === void 0 ? void 0 : _f.filter(function (item) {
                var _a;
                return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === resourceIns;
              })) === null || _g === void 0 ? void 0 : _g[0];
            } else {
              result = (_k = (_j = (_h = response === null || response === void 0 ? void 0 : response.data) === null || _h === void 0 ? void 0 : _h.items) === null || _j === void 0 ? void 0 : _j.filter(function (item) {
                var _a;
                return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === resourceIns;
              })) === null || _k === void 0 ? void 0 : _k[0];
            }
          } catch (error) {}
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 * 开启控制台
 * @param queryParams
 * @returns
 */
var openInstanceConsole = function openInstanceConsole(resource, regionId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, _a, clusterId, platform, instanceName, jsonData, serviceName, isOpen, params, url, response;
    return tslib.__generator(this, function (_b) {
      switch (_b.label) {
        case 0:
          userDefinedHeader = {
            'Content-Type': 'application/merge-patch+json'
          };
          _a = resource === null || resource === void 0 ? void 0 : resource[0], clusterId = _a.clusterId, platform = _a.platform, instanceName = _a.instanceName, jsonData = _a.jsonData, serviceName = _a.serviceName, isOpen = _a.isOpen;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(SystemNamespace, "/serviceinstances/").concat(instanceName);
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.patch,
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: [],
                  RequestBody: jsonData
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          return [4 /*yield*/, RequestApi.POST(params)];
        case 1:
          response = _b.sent();
          try {
            if (platform === PlatformType.TDCC) {} else {}
            if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
              tips$2.success(isOpen ? i18n.t('安装成功') : i18n.t('卸载成功'));
              return [2 /*return*/, operationResult(resource)];
            }
            tips$2.error(response === null || response === void 0 ? void 0 : response.message);
            return [2 /*return*/, operationResult(resource, response)];
          } catch (error) {
            tips$2.error(error);
            return [2 /*return*/, operationResult(resource, error)];
          }
          return [2 /*return*/];
      }
    });
  });
};
/**
 *
 * @param queryParams
 * @returns
 */
var fetchInstanceResources = function fetchInstanceResources(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, regionId, platform, serviceName, _a, resourceType, resourceIns, instanceId, k8sQueryObj, paging, _b, namespace, params, url, queryString, response, result, responseBody, error_6, keys, propertyKeys_2, pageIndex, pageSize, maxPageIndex;
    var _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s;
    return tslib.__generator(this, function (_t) {
      switch (_t.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, serviceName = queryParams.serviceName, _a = queryParams.resourceType, resourceType = _a === void 0 ? ResourceTypeEnum.ServiceResource : _a, resourceIns = queryParams.resourceIns, instanceId = queryParams.instanceId, k8sQueryObj = queryParams.k8sQueryObj, paging = queryParams.paging, _b = queryParams.namespace, namespace = _b === void 0 ? SystemNamespace : _b;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(namespace, "/").concat((_c = ResourceTypeMap[resourceType]) === null || _c === void 0 ? void 0 : _c.path);
          if (resourceType === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding)) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/").concat((_d = ResourceTypeMap[resourceType]) === null || _d === void 0 ? void 0 : _d.path);
          } else if (resourceType === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup)) {
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/sso/serviceopsbackups");
          }
          //目前query参数查询只支持labelSelector
          if (k8sQueryObj && (k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.labelSelector)) {
            queryString = reduceK8sQueryString({
              k8sQueryObj: k8sQueryObj,
              restfulPath: url
            });
            url = url + queryString;
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.get,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: [],
                  Path: url
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              method: Method.get,
              platform: platform
            };
          }
          result = {
            records: [],
            recordCount: 0
          };
          _t.label = 1;
        case 1:
          _t.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _t.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              responseBody = JSON.parse((_e = response === null || response === void 0 ? void 0 : response.data) === null || _e === void 0 ? void 0 : _e.ResponseBody);
              result.records = (_f = responseBody === null || responseBody === void 0 ? void 0 : responseBody.items) !== null && _f !== void 0 ? _f : [];
            } else {
              result.records = (_h = (_g = response === null || response === void 0 ? void 0 : response.data) === null || _g === void 0 ? void 0 : _g.items) !== null && _h !== void 0 ? _h : [];
            }
          } else {
            result['hasError'] = true;
          }
          return [3 /*break*/, 4];
        case 3:
          error_6 = _t.sent();
          result['hasError'] = true;
          return [3 /*break*/, 4];
        case 4:
          // 非metadata?.name属性，需要前端自己进行过滤显示
          if (k8sQueryObj && (k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector) && !((_j = Object === null || Object === void 0 ? void 0 : Object.keys(k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector)) === null || _j === void 0 ? void 0 : _j.includes('metadata.name'))) {
            keys = Object === null || Object === void 0 ? void 0 : Object.keys(k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector);
            propertyKeys_2 = (_k = keys[0]) === null || _k === void 0 ? void 0 : _k.split('.');
            result.records = (_l = result.records) === null || _l === void 0 ? void 0 : _l.filter(function (item) {
              var _a;
              return (_a = Object.keys(k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector)) === null || _a === void 0 ? void 0 : _a.every(function (key) {
                var _a, _b;
                return ((_a = item === null || item === void 0 ? void 0 : item[propertyKeys_2 === null || propertyKeys_2 === void 0 ? void 0 : propertyKeys_2[0]]) === null || _a === void 0 ? void 0 : _a[propertyKeys_2[1]]) === ((_b = k8sQueryObj === null || k8sQueryObj === void 0 ? void 0 : k8sQueryObj.fieldSelector) === null || _b === void 0 ? void 0 : _b[key]);
              });
            });
          }
          // 过滤
          if (instanceId) {
            result.records = (_m = result === null || result === void 0 ? void 0 : result.records) === null || _m === void 0 ? void 0 : _m.filter(function (item) {
              var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
              return resourceType === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup) ? ((_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.target) === null || _b === void 0 ? void 0 : _b.instanceID) === instanceId || ((_d = (_c = item === null || item === void 0 ? void 0 : item.metadata) === null || _c === void 0 ? void 0 : _c.labels) === null || _d === void 0 ? void 0 : _d['ssm.infra.tce.io/instanceId']) === instanceId || ((_f = (_e = item === null || item === void 0 ? void 0 : item.metadata) === null || _e === void 0 ? void 0 : _e.labels) === null || _f === void 0 ? void 0 : _f['ssm.infra.tce.io/instance-id']) === instanceId : ((_h = (_g = item === null || item === void 0 ? void 0 : item.metadata) === null || _g === void 0 ? void 0 : _g.labels) === null || _h === void 0 ? void 0 : _h['ssm.infra.tce.io/instanceId']) === instanceId || ((_k = (_j = item === null || item === void 0 ? void 0 : item.metadata) === null || _j === void 0 ? void 0 : _j.labels) === null || _k === void 0 ? void 0 : _k['ssm.infra.tce.io/instance-id']) === instanceId;
            });
          }
          // 按照时间降序排序
          result.records = (_o = result.records) === null || _o === void 0 ? void 0 : _o.sort(function (pre, cur) {
            var _a, _b, _c, _d;
            return ((_b = new Date((_a = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp)) === null || _b === void 0 ? void 0 : _b.getTime()) - ((_d = new Date((_c = pre === null || pre === void 0 ? void 0 : pre.metadata) === null || _c === void 0 ? void 0 : _c.creationTimestamp)) === null || _d === void 0 ? void 0 : _d.getTime());
          });
          result.recordCount = (_p = result === null || result === void 0 ? void 0 : result.records) === null || _p === void 0 ? void 0 : _p.length;
          // 前端分页
          if (paging) {
            pageIndex = (_q = paging === null || paging === void 0 ? void 0 : paging.pageIndex) !== null && _q !== void 0 ? _q : 1;
            pageSize = (_r = paging === null || paging === void 0 ? void 0 : paging.pageSize) !== null && _r !== void 0 ? _r : 20;
            maxPageIndex = Math.ceil(((_s = result.records) === null || _s === void 0 ? void 0 : _s.length) / pageSize);
            if (pageIndex > maxPageIndex) {
              pageIndex = maxPageIndex;
            }
            if (pageIndex === 1) {
              result.records = result.records.slice(pageIndex - 1, pageIndex * pageSize);
            } else {
              result.records = result.records.slice((pageIndex - 1) * pageSize, pageIndex * pageSize);
            }
          }
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 *
 * @param queryParams
 * @returns
 */
var checkCosResource = function checkCosResource(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, platform, regionId, k8sQueryObj, params, url, queryString, response, result;
    var _a;
    return tslib.__generator(this, function (_b) {
      switch (_b.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, platform = queryParams.platform, regionId = queryParams.regionId, k8sQueryObj = queryParams.k8sQueryObj;
          url = "/api/v1/namespaces/ssm/secrets/backup-cos";
          if (k8sQueryObj) {
            queryString = reduceK8sQueryString({
              k8sQueryObj: k8sQueryObj,
              restfulPath: url
            });
            url = url + queryString;
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.get,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: [],
                  Path: url
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          return [4 /*yield*/, RequestApi.GET(params)];
        case 1:
          response = _b.sent();
          try {
            if (platform === PlatformType.TDCC) {
              result = JSON.parse((_a = response === null || response === void 0 ? void 0 : response.data) === null || _a === void 0 ? void 0 : _a.ResponseBody);
            } else {
              result = response === null || response === void 0 ? void 0 : response.data;
            }
          } catch (error) {}
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 * 查询命名空间
 * @param queryParams
 * @returns
 */
var fetchNamespaces = function fetchNamespaces(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, regionId, platform, serviceName, k8sQueryObj, params, url, queryString, response, result, responseBody, error_7;
    var _a, _b, _c, _d, _e, _f;
    return tslib.__generator(this, function (_g) {
      switch (_g.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, serviceName = queryParams.serviceName, k8sQueryObj = queryParams.k8sQueryObj;
          url = '/api/v1/namespaces';
          if (k8sQueryObj) {
            queryString = reduceK8sQueryString({
              k8sQueryObj: k8sQueryObj,
              restfulPath: url
            });
            url = url + queryString;
          }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.get,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: [],
                  Path: url
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          result = {
            records: [],
            recordCount: 0
          };
          _g.label = 1;
        case 1:
          _g.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _g.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              responseBody = JSON.parse((_a = response === null || response === void 0 ? void 0 : response.data) === null || _a === void 0 ? void 0 : _a.ResponseBody);
              result.records = (_b = responseBody === null || responseBody === void 0 ? void 0 : responseBody.items) !== null && _b !== void 0 ? _b : [];
            } else {
              result.records = (_d = (_c = response === null || response === void 0 ? void 0 : response.data) === null || _c === void 0 ? void 0 : _c.items) !== null && _d !== void 0 ? _d : [];
            }
          } else {
            result['hasError'] = true;
          }
          return [3 /*break*/, 4];
        case 3:
          error_7 = _g.sent();
          result['hasError'] = true;
          return [3 /*break*/, 4];
        case 4:
          //过滤掉指定的命名空间
          result.records = (_e = result === null || result === void 0 ? void 0 : result.records) === null || _e === void 0 ? void 0 : _e.filter(function (item) {
            var _a;
            return !(ExcludeNamespaces === null || ExcludeNamespaces === void 0 ? void 0 : ExcludeNamespaces.includes((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name));
          });
          result.recordCount = (_f = result === null || result === void 0 ? void 0 : result.records) === null || _f === void 0 ? void 0 : _f.length;
          return [2 /*return*/, result];
      }
    });
  });
};
var fetchUserInfo = function fetchUserInfo(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var url, userDefinedHeader, platform, params, result, response, error_8;
    var _a;
    return tslib.__generator(this, function (_b) {
      switch (_b.label) {
        case 0:
          url = '';
          userDefinedHeader = {};
          platform = queryParams.platform;
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                restParams: {
                  Method: Method.get,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  Filters: [],
                  Path: url
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            url = "/apis/gateway.tkestack.io/v1/tokens/info";
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ProjectId: ''
              },
              platform: platform
            };
          }
          _b.label = 1;
        case 1:
          _b.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _b.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) ; else {
              result = (_a = response === null || response === void 0 ? void 0 : response.data) !== null && _a !== void 0 ? _a : {};
            }
          }
          return [3 /*break*/, 4];
        case 3:
          error_8 = _b.sent();
          throw error_8;
        case 4:
          return [2 /*return*/, result];
      }
    });
  });
};
/**
 * 查询备份策略详情
 * @param queryParams
 * @returns
 */
var fetchBackStrategy = function fetchBackStrategy(queryParams) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var userDefinedHeader, clusterId, regionId, platform, serviceName, instanceId, _a, resourceType, k8sQueryObj, params, namespace, url, response, result, responseBody, error_9;
    var _b, _c, _d, _e, _f, _g, _h;
    return tslib.__generator(this, function (_j) {
      switch (_j.label) {
        case 0:
          userDefinedHeader = {};
          clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, serviceName = queryParams.serviceName, instanceId = queryParams.instanceId, _a = queryParams.resourceType, resourceType = _a === void 0 ? ResourceTypeEnum.ServiceResource : _a, k8sQueryObj = queryParams.k8sQueryObj;
          namespace = SystemNamespace;
          url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/namespaces/").concat(namespace, "/").concat((_b = ResourceTypeMap[resourceType]) === null || _b === void 0 ? void 0 : _b.path);
          //目前query参数查询只支持labelSelector
          // if (k8sQueryObj && k8sQueryObj?.labelSelector) {
          //   // 这里是去拼接，是否需要在k8s url后面拼接一些queryString
          //   const queryString = reduceK8sQueryString({ k8sQueryObj, restfulPath: url });
          //   url = url + queryString;
          // }
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ResourceModuleName[PlatformType.TDCC],
                interfaceName: ResourceApiName[PlatformType.TDCC],
                regionId: regionId,
                restParams: {
                  Method: Method.get,
                  Path: url,
                  Version: ResourceVersionName[PlatformType.TDCC],
                  ClusterId: clusterId,
                  Filters: []
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              url: url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          _j.label = 1;
        case 1:
          _j.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.GET(params)];
        case 2:
          response = _j.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            if (platform === PlatformType.TDCC) {
              responseBody = JSON.parse((_c = response === null || response === void 0 ? void 0 : response.data) === null || _c === void 0 ? void 0 : _c.ResponseBody);
              result = (_e = (_d = responseBody === null || responseBody === void 0 ? void 0 : responseBody.items) === null || _d === void 0 ? void 0 : _d.filter(function (item) {
                var _a, _b, _c, _d;
                return ((_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.trigger) === null || _b === void 0 ? void 0 : _b.type) === (BackupTypeNum === null || BackupTypeNum === void 0 ? void 0 : BackupTypeNum.Schedule) && ((_d = (_c = item === null || item === void 0 ? void 0 : item.spec) === null || _c === void 0 ? void 0 : _c.target) === null || _d === void 0 ? void 0 : _d.instanceID) === instanceId;
              })) === null || _e === void 0 ? void 0 : _e[0];
            } else {
              result = (_h = (_g = (_f = response === null || response === void 0 ? void 0 : response.data) === null || _f === void 0 ? void 0 : _f.items) === null || _g === void 0 ? void 0 : _g.filter(function (item) {
                var _a, _b, _c, _d;
                return ((_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.trigger) === null || _b === void 0 ? void 0 : _b.type) === (BackupTypeNum === null || BackupTypeNum === void 0 ? void 0 : BackupTypeNum.Schedule) && ((_d = (_c = item === null || item === void 0 ? void 0 : item.spec) === null || _c === void 0 ? void 0 : _c.target) === null || _d === void 0 ? void 0 : _d.instanceID) === instanceId;
              })) === null || _h === void 0 ? void 0 : _h[0];
            }
          } else {
            throw response;
          }
          return [3 /*break*/, 4];
        case 3:
          error_9 = _j.sent();
          return [3 /*break*/, 4];
        case 4:
          return [2 /*return*/, result];
      }
    });
  });
};
/**
  * 获取当前集群的admin角色
  * @param resource: CreateResource
  * @param regionid: number  当前的地域Id
  */
var getClusterAdminRole = function getClusterAdminRole(resource, regionId) {
  return tslib.__awaiter(void 0, void 0, void 0, function () {
    var _a, clusterId, _b, clusterType, platform, userDefinedHeader, InterfaceNameMap, params, response, error_10;
    var _c;
    return tslib.__generator(this, function (_d) {
      switch (_d.label) {
        case 0:
          _a = resource[0], clusterId = _a.clusterId, _b = _a.clusterType, clusterType = _b === void 0 ? ClusterType.TKE : _b, platform = _a.platform;
          userDefinedHeader = {};
          InterfaceNameMap = (_c = {}, _c[ClusterType.TKE] = 'AcquireClusterAdminRole', _c[ClusterType.EKS] = 'AcquireEKSClusterAdminRole', _c[ClusterType.External] = 'AcquireClusterAdminRole', _c);
          if (platform === PlatformType.TDCC) {
            params = {
              userDefinedHeader: userDefinedHeader,
              apiParams: {
                module: ClusterResourceModuleName === null || ClusterResourceModuleName === void 0 ? void 0 : ClusterResourceModuleName[clusterType],
                interfaceName: InterfaceNameMap === null || InterfaceNameMap === void 0 ? void 0 : InterfaceNameMap[clusterType],
                regionId: regionId,
                restParams: {
                  Version: ClusterResourceVersionName === null || ClusterResourceVersionName === void 0 ? void 0 : ClusterResourceVersionName[clusterType],
                  ClusterId: clusterId,
                  Filters: []
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
          } else if (platform === PlatformType.TKESTACK) {
            params = {
              userDefinedHeader: userDefinedHeader,
              // url,
              restParams: {
                ClusterId: clusterId,
                ProjectId: ''
              },
              platform: platform
            };
          }
          _d.label = 1;
        case 1:
          _d.trys.push([1, 3,, 4]);
          return [4 /*yield*/, RequestApi.SEND(params)];
        case 2:
          response = _d.sent();
          if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
            return [2 /*return*/, operationResult(resource)];
          } else {
            return [2 /*return*/, operationResult(resource, response)];
          }
        case 3:
          error_10 = _d.sent();
          return [2 /*return*/, operationResult(resource, error_10)];
        case 4:
          return [2 /*return*/];
      }
    });
  });
};

var reduceK8sQueryStringForTdcc = function reduceK8sQueryStringForTdcc(_a) {
  var _b;
  var _c = _a.k8sQueryObj,
    k8sQueryObj = _c === void 0 ? {} : _c,
    _d = _a.restfulPath,
    restfulPath = _d === void 0 ? '' : _d,
    _e = _a.isShadow,
    isShadow = _e === void 0 ? false : _e;
  var operator = '?';
  var queryString = '';
  if (!isEmpty(k8sQueryObj)) {
    var queryKeys = Object.keys(k8sQueryObj);
    queryKeys.forEach(function (queryKey, index) {
      if (index !== 0) {
        queryString += '&';
      }
      // 这里去判断每种资源的query，eg：fieldSelector、limit等
      var specificQuery = k8sQueryObj[queryKey];
      if (_typeof(specificQuery) === 'object') {
        // 这里是对于 query的字段里面，还有多种过滤条件，比如fieldSelector支持 involvedObject.name=*,involvedObject.kind=*
        var specificKeys = Object.keys(specificQuery),
          specificString_1 = '';
        specificKeys.forEach(function (speKey, index) {
          if (index !== 0) {
            specificString_1 += ',';
          }
          if (Array.isArray(specificQuery[speKey]) && specificQuery[speKey].length !== 0) {
            // tdcc 应用管理中的接口需要编码
            if (isShadow) {
              specificString_1 += encodeURIComponent("".concat(speKey, "+in+(").concat(specificQuery[speKey].join(','), ")"));
            } else {
              specificString_1 += "".concat(speKey, "+in+(").concat(specificQuery[speKey].join(','), ")");
            }
          } else {
            specificString_1 += speKey + (specificQuery[speKey] ? "=".concat(specificQuery[speKey]) : '');
          }
        });
        if (specificString_1) {
          queryString += "".concat(queryKey, "=").concat(specificString_1);
        }
      } else {
        queryString += "".concat(queryKey, "=").concat(k8sQueryObj[queryKey]);
      }
    });
  }
  /** 如果原本的url里面已经有 ? 了，则我们这里的query的内容，必须是拼接在后面，而不能直接加多一个 ? */
  /**
   * 分布式云中前面已经有个？了，所以需要继续往后面拼Method=GET&amp;Path=/apis/platform.tke/v1/clusters/cls-afg8xkzg/proxy?path=/apis/shadow/v1alpha1/namespaces/foo/deployments?labelSelector=clusternet-app%3Dmulti-cluster-nginx
   */
  if (!isShadow && restfulPath.includes('?')) {
    operator = '&';
  } else if (((_b = restfulPath === null || restfulPath === void 0 ? void 0 : restfulPath.split('?')) === null || _b === void 0 ? void 0 : _b.length) === 3) {
    operator = '&';
  }
  return queryString ? "".concat(operator).concat(queryString) : '';
};
var medium = {
  fetchResourceList: function fetchResourceList(queryParams) {
    return tslib.__awaiter(void 0, void 0, void 0, function () {
      var userDefinedHeader, clusterId, regionId, platform, searchFilter, k8sQueryObj, params, result, url, queryString, response, responseBody, error_1;
      var _a, _b, _c, _d, _e;
      return tslib.__generator(this, function (_f) {
        switch (_f.label) {
          case 0:
            userDefinedHeader = {};
            clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, searchFilter = queryParams.searchFilter, k8sQueryObj = queryParams.k8sQueryObj;
            result = {
              records: [],
              recordCount: 0
            };
            url = '/api/v1/namespaces/ssm/secrets';
            if (searchFilter === null || searchFilter === void 0 ? void 0 : searchFilter.resourceName) {
              k8sQueryObj.fieldSelector = tslib.__assign(tslib.__assign({}, k8sQueryObj.fieldSelector), {
                'metadata.name': (_a = searchFilter.resourceName) === null || _a === void 0 ? void 0 : _a[0]
              });
            }
            if (k8sQueryObj) {
              queryString = reduceK8sQueryStringForTdcc({
                k8sQueryObj: k8sQueryObj,
                restfulPath: url
              });
              url = url + queryString;
            }
            if (platform === PlatformType.TDCC) {
              params = {
                userDefinedHeader: userDefinedHeader,
                apiParams: {
                  module: ResourceModuleName[PlatformType.TDCC],
                  interfaceName: ResourceApiName[PlatformType.TDCC],
                  regionId: regionId,
                  restParams: {
                    Method: Method.get,
                    Version: ResourceVersionName[PlatformType.TDCC],
                    ClusterId: clusterId,
                    Filters: [],
                    Path: url
                  },
                  opts: {
                    tipErr: false,
                    global: false
                  }
                },
                platform: platform
              };
            } else if (platform === PlatformType.TKESTACK) {
              params = {
                userDefinedHeader: userDefinedHeader,
                url: url,
                restParams: {
                  ClusterId: clusterId,
                  ProjectId: ''
                },
                platform: platform
              };
            }
            _f.label = 1;
          case 1:
            _f.trys.push([1, 3,, 4]);
            return [4 /*yield*/, RequestApi.GET(params)];
          case 2:
            response = _f.sent();
            if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
              if (platform === PlatformType.TDCC) {
                responseBody = JSON.parse((_b = response === null || response === void 0 ? void 0 : response.data) === null || _b === void 0 ? void 0 : _b.ResponseBody);
                result.records = (_c = responseBody === null || responseBody === void 0 ? void 0 : responseBody.items) !== null && _c !== void 0 ? _c : [];
              } else {
                result.records = (_e = (_d = response === null || response === void 0 ? void 0 : response.data) === null || _d === void 0 ? void 0 : _d.items) !== null && _e !== void 0 ? _e : [];
              }
              result['code'] = response;
            } else {
              throw response;
            }
            return [3 /*break*/, 4];
          case 3:
            error_1 = _f.sent();
            result['hasError'] = true;
            result['code'] = response;
            throw error_1;
          case 4:
            return [2 /*return*/, result];
        }
      });
    });
  },
  createCosConfigForTdcc: function createCosConfigForTdcc(resource, options) {
    return tslib.__awaiter(void 0, void 0, void 0, function () {
      var _a, mode, yamlData, resourceInfo, namespace, jsonData, _b, base64encode, regionId, clusterId, platform, url, RequestBody, EncodedBody, params, response;
      return tslib.__generator(this, function (_c) {
        switch (_c.label) {
          case 0:
            _a = resource[0], mode = _a.mode, yamlData = _a.yamlData, resourceInfo = _a.resourceInfo, namespace = _a.namespace, jsonData = _a.jsonData, _b = _a.base64encode, base64encode = _b === void 0 ? false : _b;
            regionId = options.regionId, clusterId = options.clusterId, platform = options.platform;
            url = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/apply");
            EncodedBody = false;
            if (yamlData) {
              RequestBody = EncodedBody ? jsBase64.Base64.encode(yamlData) : yamlData;
            } else {
              RequestBody = EncodedBody ? jsBase64.Base64.encode(jsonData) : jsonData;
            }
            params = {
              method: Method.post,
              url: url,
              data: yamlData ? yamlData : jsonData,
              apiParams: {
                module: 'tdcc',
                interfaceName: 'ForwardRequestTDCC',
                regionId: regionId,
                restParams: {
                  Method: Method.post,
                  ClusterId: clusterId,
                  Path: url,
                  Version: '2022-01-25',
                  RequestBody: jsonData,
                  EncodedBody: EncodedBody
                },
                opts: {
                  tipErr: false,
                  global: false
                }
              },
              platform: platform
            };
            return [4 /*yield*/, RequestApi.POST(params)];
          case 1:
            response = _c.sent();
            try {
              if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
                return [2 /*return*/, operationResult(resource)];
              }
              bridge.tips.error(response === null || response === void 0 ? void 0 : response.message);
              return [2 /*return*/, operationResult(resource, response)];
            } catch (error) {
              bridge.tips.error(error);
              return [2 /*return*/, operationResult(resource, error)];
            }
            return [2 /*return*/];
        }
      });
    });
  },

  deleteResource: function deleteResource(resource, regionId) {
    return tslib.__awaiter(void 0, void 0, void 0, function () {
      var userDefinedHeader, clusterId, platform, resourceInfo, resourceIns, params, url, response;
      return tslib.__generator(this, function (_a) {
        switch (_a.label) {
          case 0:
            userDefinedHeader = {};
            clusterId = resource.clusterId, platform = resource.platform, resourceInfo = resource.resourceInfo, resourceIns = resource.resourceIns;
            url = "/api/v1/namespaces/ssm/secrets/".concat(resourceIns);
            if (platform === PlatformType.TDCC) {
              params = {
                userDefinedHeader: userDefinedHeader,
                apiParams: {
                  module: ResourceModuleName[PlatformType.TDCC],
                  interfaceName: ResourceApiName[PlatformType.TDCC],
                  regionId: regionId,
                  restParams: {
                    Method: Method["delete"],
                    Path: url,
                    Version: ResourceVersionName[PlatformType.TDCC],
                    ClusterId: clusterId,
                    Filters: []
                  },
                  opts: {
                    tipErr: false,
                    global: false
                  }
                },
                platform: platform
              };
            } else if (platform === PlatformType.TKESTACK) {
              params = {
                userDefinedHeader: userDefinedHeader,
                url: url,
                restParams: {
                  ClusterId: clusterId,
                  ProjectId: ''
                },
                platform: platform,
                method: Method["delete"]
              };
            }
            return [4 /*yield*/, RequestApi.DELETE(params)];
          case 1:
            response = _a.sent();
            try {
              if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
                return [2 /*return*/, operationResult(resource)];
              }
              return [2 /*return*/, operationResult(resource, response)];
            } catch (error) {
              return [2 /*return*/, operationResult(resource, error)];
            }
            return [2 /*return*/];
        }
      });
    });
  },

  fetchResourceDetail: function fetchResourceDetail(queryParams) {
    return tslib.__awaiter(void 0, void 0, void 0, function () {
      var userDefinedHeader, clusterId, regionId, platform, instanceName, params, result, url, url_1, response, responseBody;
      var _a, _b, _c, _d;
      return tslib.__generator(this, function (_e) {
        switch (_e.label) {
          case 0:
            userDefinedHeader = {};
            clusterId = queryParams.clusterId, regionId = queryParams.regionId, platform = queryParams.platform, instanceName = queryParams.instanceName;
            result = {
              records: [],
              recordCount: 0
            };
            url = "/api/v1/namespaces/ssm/secrets/".concat(instanceName);
            if (platform === PlatformType.TDCC) {
              params = {
                userDefinedHeader: userDefinedHeader,
                apiParams: {
                  module: ResourceModuleName[PlatformType.TDCC],
                  interfaceName: ResourceApiName[PlatformType.TDCC],
                  regionId: regionId,
                  restParams: {
                    Method: Method.get,
                    Version: ResourceVersionName[PlatformType.TDCC],
                    ClusterId: clusterId,
                    Filters: [],
                    Path: url
                  },
                  opts: {
                    tipErr: false,
                    global: false
                  }
                },
                platform: platform
              };
            } else if (platform === PlatformType.TKESTACK) {
              url_1 = "/apis/platform.tkestack.io/v1/clusters/".concat(clusterId, "/proxy?path=/apis/infra.tce.io/v1/test");
              params = {
                userDefinedHeader: userDefinedHeader,
                url: url_1,
                restParams: {
                  ClusterId: clusterId,
                  ProjectId: ''
                },
                platform: platform
              };
            }
            return [4 /*yield*/, RequestApi.GET(params)];
          case 1:
            response = _e.sent();
            if ((response === null || response === void 0 ? void 0 : response.code) === 0) {
              if (platform === PlatformType.TDCC) {
                responseBody = JSON.parse((_a = response === null || response === void 0 ? void 0 : response.data) === null || _a === void 0 ? void 0 : _a.ResponseBody);
                result.records = (_b = responseBody === null || responseBody === void 0 ? void 0 : responseBody.items) !== null && _b !== void 0 ? _b : [];
              } else {
                result.records = (_d = (_c = response === null || response === void 0 ? void 0 : response.data) === null || _c === void 0 ? void 0 : _c.items) !== null && _d !== void 0 ? _d : [];
              }
            } else {
              result['hasError'] = true;
            }
            return [2 /*return*/, result];
        }
      });
    });
  }
};

var MediumSelectPanel;
(function (MediumSelectPanel) {
  var _this = this;
  var ComponentName = 'MediumSelectPanel';
  var ActionTypes;
  (function (ActionTypes) {
    ActionTypes["VALIDATOR"] = "VALIDATOR";
    ActionTypes["Clear"] = "Clear";
    ActionTypes["Mediums"] = "Mediums";
  })(ActionTypes = MediumSelectPanel.ActionTypes || (MediumSelectPanel.ActionTypes = {}));
  BComponent$1.createActionType(ComponentName, ActionTypes);
  MediumSelectPanel.createValidateSchema = function (_a) {
    var pageName = _a.pageName;
    var schema = {
      formKey: BComponent$1.getActionType(pageName, ActionTypes.VALIDATOR),
      fields: [{
        label: i18n.t('备份介质'),
        vKey: 'mediums',
        modelType: ffValidator.ModelTypeEnum.FFRedux,
        rules: [ffValidator.RuleTypeEnum.isRequire]
      }]
    };
    return schema;
  };
  MediumSelectPanel.createActions = function (_a) {
    var pageName = _a.pageName,
      _getRecord = _a.getRecord;
    var actions = {
      validator: ffValidator.createValidatorActions({
        userDefinedSchema: MediumSelectPanel.createValidateSchema({
          pageName: pageName
        }),
        validateStateLocator: function validateStateLocator(store) {
          return _getRecord(function () {
            return store;
          });
        },
        validatorStateLocation: function validatorStateLocation(store) {
          return _getRecord(function () {
            return store;
          }).validator;
        }
      }),
      mediums: ffRedux.createFFListActions({
        actionName: BComponent$1.getActionType(pageName, ActionTypes.Mediums),
        fetcher: function fetcher(query) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var k8sQueryObj, result;
            return tslib.__generator(this, function (_a) {
              switch (_a.label) {
                case 0:
                  k8sQueryObj = {
                    labelSelector: {
                      'tdcc.cloud.tencent.com/paas-storage-medium': ['s3', 'nfs']
                    }
                  };
                  if (!(query === null || query === void 0 ? void 0 : query.filter)) return [3 /*break*/, 2];
                  return [4 /*yield*/, medium.fetchResourceList(tslib.__assign(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter), {
                    k8sQueryObj: k8sQueryObj
                  }))];
                case 1:
                  result = _a.sent();
                  _a.label = 2;
                case 2:
                  return [2 /*return*/, result];
              }
            });
          });
        },
        getRecord: function getRecord(getState) {
          return _getRecord(getState).mediums;
        },
        selectFirst: true,
        onFinish: function onFinish(record, dispatch) {}
      }),
      clearState: function clearState() {
        return {
          type: BComponent$1.getActionType(pageName, ActionTypes.Clear)
        };
      }
    };
    return actions;
  };
  MediumSelectPanel.createReducer = function (_a) {
    var pageName = _a.pageName;
    var reducer = redux.combineReducers({
      validator: ffValidator.createValidatorReducer(MediumSelectPanel.createValidateSchema({
        pageName: pageName
      })),
      mediums: ffRedux.createFFListReducer({
        actionName: BComponent$1.getActionType(pageName, ActionTypes.Mediums),
        displayField: function displayField(r) {
          return r.metadata.name;
        },
        valueField: function valueField(r) {
          return r.metadata.name;
        }
      })
    });
    var finalReducer = function finalReducer(state, action) {
      var newState = state;
      // 销毁组件
      if (action.type === BComponent$1.getActionType(pageName, ActionTypes.Clear)) {
        newState = undefined;
      }
      return reducer(newState, action);
    };
    return finalReducer;
  };
  MediumSelectPanel.Component = function (props) {
    var _a, _b, _c, _d, _e, _f;
    var _g = props.clearData,
      clearData = _g === void 0 ? false : _g,
      model = props.model,
      action = props.action,
      clusterId = props.clusterId,
      regionId = props.regionId,
      platform = props.platform;
    var mediums = model.mediums,
      validator = model.validator;
    React.useEffect(function () {
      return function () {
        var _a;
        if (clearData) {
          (_a = action === null || action === void 0 ? void 0 : action.mediums) === null || _a === void 0 ? void 0 : _a.select(null);
        }
      };
    }, [clearData, action]);
    React.useEffect(function () {
      var filter = {
        platform: platform,
        regionId: regionId,
        clusterId: clusterId,
        namespace: 'ssm'
      };
      if (clusterId && regionId && BComponent$1.isNeedFetch(mediums, filter)) {
        action.mediums.changeFilter(filter);
        action.mediums.fetch({
          fetchAll: true
        });
      }
    }, [clusterId, regionId, mediums, action.mediums, platform]);
    var RBACForbidden = 'FailedOperation.RBACForbidden';
    var RBACForbidden403 = 403;
    if (mediums.list.fetchState === ffRedux.FetchState.Failed) {
      if (ffComponent.isCamRefused(mediums.list.error)) {
        return React__default.createElement(ffComponent.CamBox, {
          message: mediums.list.error.message
        });
      }
    }
    return React__default.createElement(ffComponent.FormPanel.Item, {
      formvalidator: validator,
      vactions: action.validator,
      vkey: 'mediums',
      label: i18n.t('备份介质'),
      select: {
        model: mediums,
        action: action.mediums,
        showRefreshBtn: true
      },
      message: i18n.t('自定义设置外部存储介质，例如用户自己的对象存储或者文件存储，适合用户创建的中间件实例备份场景'),
      errorTips: ((_b = (_a = mediums === null || mediums === void 0 ? void 0 : mediums.list) === null || _a === void 0 ? void 0 : _a.error) === null || _b === void 0 ? void 0 : _b.code) === RBACForbidden || ((_d = (_c = mediums === null || mediums === void 0 ? void 0 : mediums.list) === null || _c === void 0 ? void 0 : _c.error) === null || _d === void 0 ? void 0 : _d.code) === RBACForbidden403 ? i18n.t('权限不足，请联系集群管理员添加权限') : (_f = (_e = mediums === null || mediums === void 0 ? void 0 : mediums.list) === null || _e === void 0 ? void 0 : _e.error) === null || _f === void 0 ? void 0 : _f.code,
      after: React__default.createElement(teaComponent$1.ExternalLink, {
        href: platform === PlatformType.TKESTACK ? "https://console.cloud.tencent.com/tdcc/medium?clusterId=".concat(clusterId) : "/tdcc/medium?clusterId=".concat(clusterId)
      }, i18n.t('前往添加备份介质'))
    });
  };
})(MediumSelectPanel || (MediumSelectPanel = {}));

var GetRbacAdminDialog;
(function (GetRbacAdminDialog) {
  GetRbacAdminDialog.ComponentName = "GetRbacAdminDialog";
  GetRbacAdminDialog.ActionType = Object.assign({}, BComponent$1.BaseActionType, {
    GetClusterAdminRoleFlow: "GetClusterAdminRoleFlow"
  });
  BComponent$1.createActionType(GetRbacAdminDialog.ComponentName, GetRbacAdminDialog.ActionType);
  GetRbacAdminDialog.createActions = function (_a) {
    var pageName = _a.pageName,
      getRecord = _a.getRecord;
    var actions = {
      getClusterAdminRole: ffRedux.generateWorkflowActionCreator({
        actionType: BComponent$1.getActionType(pageName, GetRbacAdminDialog.ActionType.GetClusterAdminRoleFlow),
        workflowStateLocator: function workflowStateLocator(state) {
          return getRecord(function () {
            return state;
          }).getClusterAdminRoleFlow;
        },
        operationExecutor: getClusterAdminRole,
        after: {}
      })
    };
    return actions;
  };
  GetRbacAdminDialog.createReducer = function (_a) {
    var pageName = _a.pageName;
    var TempReducer = redux.combineReducers({
      getClusterAdminRoleFlow: ffRedux.generateWorkflowReducer({
        actionType: BComponent$1.getActionType(pageName, GetRbacAdminDialog.ActionType.GetClusterAdminRoleFlow)
      })
    });
    var Reducer = function Reducer(state, action) {
      var newState = state;
      // 销毁页面
      if (action.type === BComponent$1.getActionType(pageName, GetRbacAdminDialog.ActionType.Clear)) {
        newState = undefined;
      }
      return TempReducer(newState, action);
    };
    return Reducer;
  };
  GetRbacAdminDialog.Component = function (props) {
    var action = props.action,
      getClusterAdminRoleFlow = props.model.getClusterAdminRoleFlow,
      _a = props.filter,
      clusterId = _a.clusterId,
      regionId = _a.regionId,
      platform = _a.platform,
      onSuccess = props.onSuccess,
      resourceType = props.resourceType;
    var workflow = getClusterAdminRoleFlow;
    var workFlowAction = action.getClusterAdminRole;
    React.useEffect(function () {
      if (workflow.operationState === ffRedux.OperationState.Done && ffRedux.isSuccessWorkflow(workflow)) {
        workFlowAction.reset();
        if (typeof onSuccess === "function") {
          onSuccess();
        }
      }
    }, [workflow.operationState, onSuccess, workflow, workFlowAction]);
    var cancel = function cancel() {
      if (workflow.operationState === ffRedux.OperationState.Done) {
        workFlowAction.reset();
      } else {
        workFlowAction.cancel();
      }
    };
    var perform = function perform() {
      // 提交操作
      var params = {
        id: uuid(),
        platform: platform,
        clusterId: clusterId,
        clusterType: ClusterType.External
      };
      workFlowAction.start([params], regionId);
      workFlowAction.perform();
    };
    var failed = workflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(workflow);
    if (getClusterAdminRoleFlow.operationState === ffRedux.OperationState.Pending) {
      return React__default.createElement("noscript", null);
    }
    return React__default.createElement(ModalMain.Modal, {
      caption: i18n.t("获取集群Admin角色"),
      disableEscape: true,
      visible: true,
      onClose: cancel
    }, React__default.createElement(ModalMain.Modal.Body, null, React__default.createElement(i18n.Trans, null, React__default.createElement(teaComponent.Text, {
      theme: "text",
      parent: "p"
    }, "\u5B50\u8D26\u6237\u5728CAM\u62E5\u6709AcquireClusterAdminRole Action\u6743\u9650\u4E4B\u540E\u53EF\u4EE5\u901A\u8FC7\u6B64\u65B9\u5F0F\u83B7\u53D6\u96C6\u7FA4Kubernetes RBAC\u7684tke:admin\u89D2\u8272"), React__default.createElement(teaComponent.Text, {
      theme: "text",
      parent: "p"
    }, "\u70B9\u51FB\u786E\u5B9A\u4E4B\u540E\uFF0C\u540E\u53F0\u4F1A\u4E3A\u60A8\u521B\u5EFAtke:admin\u89D2\u8272\u7684ClusterRolebinding\uFF0C\u4E4B\u540E\u60A8\u5C06\u62E5\u6709\u6B64\u96C6\u7FA4\u5185\u8D44\u6E90\u7684\u7BA1\u7406\u5458\u6743\u9650\u3002")), failed && (ffComponent.isCamRefused(workflow.results[0].error) ? React__default.createElement(ffComponent.CamBox, {
      message: workflow.results[0].error.message
    }) : React__default.createElement(TipInfo, {
      isForm: true,
      type: "error"
    }, getWorkflowError(workflow)))), React__default.createElement(ModalMain.Modal.Footer, null, React__default.createElement(teaComponent.Button, {
      disabled: workflow.operationState === ffRedux.OperationState.Performing,
      type: "primary",
      onClick: perform
    }, failed ? i18n.t("重试") : i18n.t("确定")), React__default.createElement(teaComponent.Button, {
      onClick: cancel
    }, i18n.t("取消"))));
  };
})(GetRbacAdminDialog || (GetRbacAdminDialog = {}));

function RetryPanel(props) {
  var _a = props.style,
    style = _a === void 0 ? {} : _a,
    action = props.action,
    _b = props.loadingText,
    loadingText = _b === void 0 ? i18n.t('加载失败') : _b,
    _c = props.retryText,
    retryText = _c === void 0 ? i18n.t('刷新重试') : _c,
    _d = props.loadingTextTheme,
    loadingTextTheme = _d === void 0 ? 'danger' : _d;
  return React__default.createElement("div", {
    style: tslib.__assign({
      width: '100%',
      display: 'flex',
      alignItems: 'center'
    }, style)
  }, React__default.createElement(teaComponent.Text, {
    theme: loadingTextTheme,
    className: 'tea-mr-2n'
  }, loadingText), React__default.createElement(teaComponent.Button, {
    type: "link",
    onClick: function onClick() {
      action && action();
    }
  }, retryText));
}

var initRecords = [{
  id: uuid(),
  key: '',
  value: ''
}];
var MapField = function MapField(_a) {
  var plan = _a.plan,
    onChange = _a.onChange;
  var _b = React.useState(initRecords),
    records = _b[0],
    setRecords = _b[1];
  React.useEffect(function () {
    onChange && onChange({
      field: plan === null || plan === void 0 ? void 0 : plan.name,
      value: records
    });
  }, [records]);
  var _delete = function _delete(id) {
    setRecords(records === null || records === void 0 ? void 0 : records.filter(function (item) {
      return (item === null || item === void 0 ? void 0 : item.id) !== id;
    }));
  };
  var _add = function _add() {
    var newItem = {
      id: uuid(),
      key: '',
      value: ''
    };
    setRecords(records === null || records === void 0 ? void 0 : records.concat([newItem]));
  };
  var _update = function _update(fieldName, selectItem, value) {
    if (value === void 0) {
      value = '';
    }
    var newRecords = records === null || records === void 0 ? void 0 : records.map(function (item) {
      var _a;
      return tslib.__assign(tslib.__assign({}, item), (_a = {}, _a[fieldName] = (item === null || item === void 0 ? void 0 : item.id) === (selectItem === null || selectItem === void 0 ? void 0 : selectItem.id) ? value : item === null || item === void 0 ? void 0 : item[fieldName], _a));
    });
    setRecords(newRecords);
  };
  return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Table, {
    bordered: true,
    recordKey: "id",
    columns: [{
      key: "key",
      header: 'Key',
      render: function render(item) {
        return React__default.createElement(teaComponent.Input, {
          value: item === null || item === void 0 ? void 0 : item.key,
          onChange: function onChange(value) {
            _update('key', item, value);
          }
        });
      }
    }, {
      key: "value",
      header: 'Value',
      render: function render(item) {
        return React__default.createElement(teaComponent.Input, {
          value: item === null || item === void 0 ? void 0 : item.value,
          onChange: function onChange(value) {
            _update('value', item, value);
          }
        });
      }
    }, {
      key: "operate",
      header: '操作',
      render: function render(item) {
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
          type: "link",
          onClick: function onClick() {
            _delete(item === null || item === void 0 ? void 0 : item.id);
          }
        }, i18n.t('删除')));
      }
    }],
    records: records
  }), React__default.createElement(teaComponent.Button, {
    onClick: _add,
    type: 'link',
    className: 'tea-mt-2n'
  }, i18n.t('添加配置')));
};

function UpdateResource(props) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o;
  var _p = props.base,
    platform = _p.platform,
    isI18n = _p.isI18n,
    route = _p.route,
    _q = props.list,
    services = _q.services,
    servicePlanEdit = _q.servicePlanEdit,
    updateResourceWorkflow = _q.updateResourceWorkflow,
    resourceDetail = props.detail.resourceDetail,
    actions = props.actions;
  var servicename = (route === null || route === void 0 ? void 0 : route.queries).servicename;
  var instanceParamsFields = (_b = (_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.instanceSchema;
  var reduceCreateServiceResourceDataJson = function reduceCreateServiceResourceDataJson(data) {
    var _a, _b, _c, _d, _e, _f;
    var _g = data === null || data === void 0 ? void 0 : data.formData,
      instanceName = _g.instanceName,
      description = _g.description,
      clusterId = _g.clusterId;
    //拼接中间件实例参数部分属性值
    var parameters = (_c = (_b = (_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.instanceSchema) === null || _c === void 0 ? void 0 : _c.reduce(function (pre, cur) {
      var _a;
      var _b, _c, _d, _e;
      // return Object.assign(pre,!!data?.formData?.[cur?.name] ? {[cur?.name]:data?.formData?.[cur?.name] + (data?.formData?.['unitMap']?.[cur?.name] ?? '')} : {});
      return Object.assign(pre, ((_b = data === null || data === void 0 ? void 0 : data.formData) === null || _b === void 0 ? void 0 : _b[cur === null || cur === void 0 ? void 0 : cur.name]) !== '' ? (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = formatPlanSchemaSubmitData(cur, (_c = data === null || data === void 0 ? void 0 : data.formData) === null || _c === void 0 ? void 0 : _c[cur === null || cur === void 0 ? void 0 : cur.name], (_e = (_d = data === null || data === void 0 ? void 0 : data.formData) === null || _d === void 0 ? void 0 : _d['unitMap']) === null || _e === void 0 ? void 0 : _e[cur === null || cur === void 0 ? void 0 : cur.name]), _a) : {});
    }, {});
    var json = {
      id: uuid(),
      apiVersion: 'infra.tce.io/v1',
      kind: (_d = ResourceTypeMap === null || ResourceTypeMap === void 0 ? void 0 : ResourceTypeMap[ResourceTypeEnum.ServicePlan]) === null || _d === void 0 ? void 0 : _d.resourceKind,
      metadata: {
        labels: {
          'ssm.infra.tce.io/owner': ServicePlanTypeEnum.Custom
        },
        name: instanceName,
        namespace: DefaultNamespace
      },
      spec: {
        serviceClass: ((_e = services === null || services === void 0 ? void 0 : services.selection) === null || _e === void 0 ? void 0 : _e.name) || ((_f = route === null || route === void 0 ? void 0 : route.queries) === null || _f === void 0 ? void 0 : _f.servicename),
        metadata: tslib.__assign({}, parameters)
      }
    };
    if (description) {
      json.spec['description'] = description;
    }
    return JSON.stringify(json);
  };
  var _submit = function _submit() {
    var _a, _b, _c, _d;
    var formData = servicePlanEdit.formData;
    actions.create.validatePlan((_b = (_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.instanceSchema);
    if (_validatePlan(servicePlanEdit, (_d = (_c = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.instanceSchema)) {
      var params = {
        platform: platform,
        regionId: HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion,
        clusterId: formData === null || formData === void 0 ? void 0 : formData.clusterId,
        jsonData: reduceCreateServiceResourceDataJson(servicePlanEdit),
        resourceType: ResourceTypeEnum.ServicePlan,
        instanceName: formData === null || formData === void 0 ? void 0 : formData.instanceName
      };
      actions.create.updateResource.start([params], HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion);
      actions.create.updateResource.perform();
    }
  };
  var _cancel = function _cancel() {
    actions.create.updateResource.reset();
    actions.list.showCreateResourceDialog(false);
  };
  var loading = !((_c = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _c === void 0 ? void 0 : _c.fetched) || (resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object.fetchState) === ffRedux.FetchState.Fetching;
  var loadSchemaFailed = ((_d = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _d === void 0 ? void 0 : _d.error) || ((_e = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _e === void 0 ? void 0 : _e.fetchState) === ffRedux.FetchState.Failed;
  var isSubmitting = (updateResourceWorkflow === null || updateResourceWorkflow === void 0 ? void 0 : updateResourceWorkflow.operationState) === ffRedux.OperationState.Performing;
  var failed = updateResourceWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(updateResourceWorkflow);
  return React__default.createElement(ffComponent.FormPanel, {
    isNeedCard: false
  }, loading && React__default.createElement(LoadingPanel, null), !loading && loadSchemaFailed && React__default.createElement(RetryPanel, {
    style: {
      minWidth: 150
    },
    loadingText: i18n.t('加载失败'),
    action: (_g = (_f = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _f === void 0 ? void 0 : _f.instanceDetail) === null || _g === void 0 ? void 0 : _g.fetch
  }), !loading && !loadSchemaFailed && React__default.createElement(React__default.Fragment, null, React__default.createElement(ffComponent.FormPanel.Item, {
    label: React__default.createElement(teaComponent.Text, {
      style: {
        display: 'flex',
        alignItems: 'center'
      }
    }, React__default.createElement(teaComponent.Text, null, i18n.t('规格名称')), React__default.createElement(teaComponent.Text, {
      className: 'text-danger tea-pt-1n'
    }, "*")),
    validator: (_h = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.validator) === null || _h === void 0 ? void 0 : _h['instanceName'],
    input: {
      value: (_j = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _j === void 0 ? void 0 : _j['instanceName'],
      onChange: function onChange(e) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updatePlan('instanceName', e);
      },
      onBlur: function onBlur() {},
      disabled: true
    }
  }), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('中间件类型'),
    text: true
  }, React__default.createElement(ffComponent.FormPanel.Text, null, i18n.t('{{name}}', {
    name: (_k = services === null || services === void 0 ? void 0 : services.selection) === null || _k === void 0 ? void 0 : _k.name
  }))), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('目标集群:')
  }, React__default.createElement(ffComponent.FormPanel.Text, null, i18n.t('{{clusterId}}', {
    clusterId: (_l = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _l === void 0 ? void 0 : _l['clusterId']
  }))), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('描述'),
    validator: (_m = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.validator) === null || _m === void 0 ? void 0 : _m['description'],
    input: {
      value: (_o = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _o === void 0 ? void 0 : _o['description'],
      onChange: function onChange(e) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updatePlan('description', e);
      },
      onBlur: function onBlur() {}
    }
  }), instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.map(function (item, index) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
    var _l = (_b = (_a = item.description) === null || _a === void 0 ? void 0 : _a.split('---')) !== null && _b !== void 0 ? _b : [],
      english = _l[0],
      chinese = _l[1];
    return React__default.createElement(ffComponent.FormPanel.Item, {
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: 'flex',
          alignItems: 'center'
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t('{{name}}', {
        name: prefixForSchema(item, servicename) + (item === null || item === void 0 ? void 0 : item.label) + suffixUnitForSchema(item)
      })), React__default.createElement(teaComponent.Icon, {
        type: "info",
        tooltip: isI18n ? english : chinese
      }), !(item === null || item === void 0 ? void 0 : item.optional) && React__default.createElement(teaComponent.Text, {
        className: 'text-danger tea-pt-1n'
      }, "*")),
      key: "".concat(item === null || item === void 0 ? void 0 : item.name).concat(index),
      validator: (_c = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.validator) === null || _c === void 0 ? void 0 : _c[item === null || item === void 0 ? void 0 : item.name]
    }, getFormItemType(item) === FormItemType.Select && React__default.createElement(ffComponent.FormPanel.Select, {
      placeholder: i18n.t('{{title}}', {
        title: '请选择' + (item === null || item === void 0 ? void 0 : item.label)
      }),
      options: (_d = item === null || item === void 0 ? void 0 : item.candidates) === null || _d === void 0 ? void 0 : _d.map(function (candidate) {
        return {
          value: candidate,
          text: candidate
        };
      }),
      value: (_e = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _e === void 0 ? void 0 : _e[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.Switch && React__default.createElement(teaComponent.Switch, {
      value: (_f = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _f === void 0 ? void 0 : _f[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.Input && React__default.createElement(teaComponent.Input, {
      value: (_g = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _g === void 0 ? void 0 : _g[item === null || item === void 0 ? void 0 : item.name],
      placeholder: i18n.t('{{title}}', {
        title: '请输入' + (item === null || item === void 0 ? void 0 : item.label)
      }),
      onChange: function onChange(e) {
        actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.Paasword && React__default.createElement(InputPassword.InputPassword, {
      value: (_h = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _h === void 0 ? void 0 : _h[item === null || item === void 0 ? void 0 : item.name],
      placeholder: i18n.t('{{title}}', {
        title: '请输入' + (item === null || item === void 0 ? void 0 : item.label)
      }),
      onChange: function onChange(e) {
        actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.MapField && React__default.createElement(MapField, {
      plan: item,
      onChange: function onChange(_a) {
        var field = _a.field,
          value = _a.value;
        actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, JSON.stringify(value === null || value === void 0 ? void 0 : value.reduce(function (pre, cur) {
          var _a;
          return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.key] = cur === null || cur === void 0 ? void 0 : cur.value, _a));
        }, {})));
      }
    }), showUnitOptions(item) && React__default.createElement(ffComponent.FormPanel.Select, {
      size: 's',
      className: 'tea-ml-2n',
      placeholder: i18n.t('{{title}}', {
        title: '请选择unit'
      }),
      options: getUnitOptions(item),
      value: (_k = (_j = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _j === void 0 ? void 0 : _j['unitMap']) === null || _k === void 0 ? void 0 : _k[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        var _a;
        var _b;
        actions === null || actions === void 0 ? void 0 : actions.create.updatePlan('unitMap', tslib.__assign(tslib.__assign({}, (_b = servicePlanEdit === null || servicePlanEdit === void 0 ? void 0 : servicePlanEdit.formData) === null || _b === void 0 ? void 0 : _b['unitMap']), (_a = {}, _a[item === null || item === void 0 ? void 0 : item.name] = e, _a)));
      }
    }));
  }), React__default.createElement(ffComponent.FormPanel.Item, {
    className: 'tea-mr-2n'
  }, React__default.createElement(teaComponent.Col, {
    span: 24
  }, React__default.createElement(teaComponent.Button, {
    type: 'primary',
    className: 'tea-mr-2n',
    onClick: _submit,
    loading: isSubmitting,
    disabled: isSubmitting
  }, failed ? i18n.t('重试') : i18n.t('确定')), React__default.createElement(teaComponent.Button, {
    onClick: _cancel
  }, i18n.t('取消')), React__default.createElement(TipInfo, {
    isShow: failed,
    type: "error",
    isForm: true
  }, getWorkflowError(updateResourceWorkflow))))));
}

function ServiceCreateDialog(props) {
  var _a, _b;
  var showCreateResourceDialog = props.list.showCreateResourceDialog,
    resourceDetail = props.detail.resourceDetail,
    actions = props.actions,
    _c = props.mode,
    mode = _c === void 0 ? 'edit' : _c;
  var loading = ((_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _a === void 0 ? void 0 : _a.loading) || (resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object.fetchState) === ffRedux.FetchState.Fetching;
  var failed = (resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object.fetchState) === ffRedux.FetchState.Failed;
  var resource = (_b = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _b === void 0 ? void 0 : _b.data;
  return React__default.createElement(teaComponent.Modal, {
    visible: showCreateResourceDialog,
    caption: mode === 'create' ? i18n.t('新建规格') : i18n.t('编辑规格'),
    onClose: function onClose() {
      actions.list.showCreateResourceDialog(false);
    },
    size: 'm'
  }, React__default.createElement(teaComponent.Modal.Body, null, React__default.createElement(UpdateResource, tslib.__assign({}, props))));
}

var routerSea$1 = seajs.require('router');
function ServiceDetailHeader(props) {
  var _a, _b, _c, _d, _e, _f;
  var services = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.services;
  });
  var selectedTab = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.base) === null || _a === void 0 ? void 0 : _a.selectedTab;
  });
  var resource = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.serviceResources;
  });
  var platform = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.base) === null || _a === void 0 ? void 0 : _a.platform;
  });
  var actions = props.actions;
  var _g = React.useState([]),
    searchBoxValues = _g[0],
    setSearchBoxValues = _g[1];
  var _h = React.useState(0),
    searchBoxLength = _h[0],
    setSearchBoxLength = _h[1];
  var onCreate = function onCreate() {
    var _a, _b, _c, _d, _e;
    (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.serviceResources) === null || _b === void 0 ? void 0 : _b.clearPolling();
    if (selectedTab === ResourceTypeEnum.ServiceResource) {
      router.navigate({
        sub: 'create',
        tab: ResourceTypeEnum.ServiceResource
      }, {
        servicename: (_c = services === null || services === void 0 ? void 0 : services.selection) === null || _c === void 0 ? void 0 : _c.name,
        resourceType: ResourceTypeEnum.ServiceResource,
        mode: "create"
      });
    } else if (selectedTab === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan)) {
      router.navigate({
        sub: 'create',
        tab: ResourceTypeEnum.ServicePlan
      }, {
        servicename: (_d = services === null || services === void 0 ? void 0 : services.selection) === null || _d === void 0 ? void 0 : _d.name,
        resourceType: ResourceTypeEnum.ServicePlan,
        mode: "create"
      });
    } else {
      router.navigate({
        sub: 'create',
        tab: ResourceTypeEnum.ServiceResource
      }, {
        servicename: (_e = services === null || services === void 0 ? void 0 : services.selection) === null || _e === void 0 ? void 0 : _e.name,
        resourceType: ResourceTypeEnum.ServiceResource,
        mode: "create"
      });
    }
  };
  var btnText = selectedTab === ResourceTypeEnum.ServiceResource ? i18n.t('创建中间件实例') : i18n.t('创建规格');
  /** 生成手动刷新按钮 */
  var _renderManualRenew = function _renderManualRenew() {
    var loading = resource.list.loading || resource.list.fetchState === ffRedux.FetchState.Fetching;
    return  React__default.createElement(teaComponent.Button, {
      icon: "refresh",
      disabled: loading,
      onClick: function onClick(e) {
        // actions.list.serviceResources.reset();
        actions.list.serviceResources.fetch({
          noCache: true
        });
      },
      title: i18n.t('刷新')
    }) ;
  };
  /** 生成搜索框 */
  var _renderTagSearchBox = function _renderTagSearchBox() {
    // tagSearch的过滤选项
    var attributes = [{
      type: 'input',
      key: 'resourceName',
      name: i18n.t('名称')
    }];
    var values = resource.query.search ? searchBoxValues : [];
    return  React__default.createElement("div", {
      style: {
        width: 350,
        display: 'inline-block'
      }
    }, React__default.createElement(teaComponent.TagSearchBox, {
      className: "myTagSearchBox",
      attributes: attributes,
      value: values,
      onChange: function onChange(tags) {
        _handleClickForTagSearch(tags);
      }
    })) ;
  };
  /** 搜索框的操作，不同的搜索进行相对应的操作 */
  var _handleClickForTagSearch = function _handleClickForTagSearch(tags) {
    // 这里是控制tagSearch的展示
    setSearchBoxValues(tags);
    setSearchBoxLength(tags.length);
    // 如果检测到 tags的长度变化，并且key为 resourceName 去掉了，则清除搜索条件
    if (tags.length === 0 || tags.length === 1 && resource.query.search && tags[0].attr && tags[0].attr.key !== 'resourceName') {
      actions.list.serviceResources.changeKeyword('');
      actions.list.serviceResources.performSearch('');
    }
    tags.forEach(function (tagItem) {
      var attrKey = tagItem.attr ? tagItem.attr.key : null;
      if (attrKey === 'resourceName' || attrKey === null) {
        var search = tagItem.values[0].name;
        actions.list.serviceResources.changeKeyword(search);
        actions.list.serviceResources.performSearch(search);
      }
    });
  };
  var isDisabledCreate = !((_b = (_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.clusters) === null || _b === void 0 ? void 0 : _b.length);
  var title = !((_c = services === null || services === void 0 ? void 0 : services.selection) === null || _c === void 0 ? void 0 : _c.name) ? i18n.t('请您选择左侧的服务') : !((_e = (_d = services === null || services === void 0 ? void 0 : services.selection) === null || _d === void 0 ? void 0 : _d.clusters) === null || _e === void 0 ? void 0 : _e.length) ? i18n.t('{{title}}', {
    title: "".concat((_f = services === null || services === void 0 ? void 0 : services.selection) === null || _f === void 0 ? void 0 : _f.name, "\u5C1A\u672A\u5728\u76EE\u6807\u96C6\u7FA4\u5F00\u542F\uFF0C\u8BF7\u60A8\u524D\u5F80\u5206\u5E03\u5F0F\u4E91\u4E2D\u5FC3\u6982\u89C8\u9875\u5F00\u542F")
  }) : null;
  return React__default.createElement("div", {
    className: 'tea-pb-5n'
  }, React__default.createElement(teaComponent.Justify, {
    left: React__default.createElement(teaComponent.Button, {
      type: "primary",
      onClick: onCreate,
      disabled: isDisabledCreate,
      title: title
    }, btnText),
    right: React__default.createElement(React__default.Fragment, null, _renderTagSearchBox(), _renderManualRenew())
  }));
}

/*
 * @File: 这是文件的描述
 * @Description: 这是文件的描述
 * @Version: 1.0
 * @Autor: brycewwang
 * @Date: 2022-06-25 20:24:07
 * @LastEditors: brycewwang
 * @LastEditTime: 2022-06-28 20:24:57
 */
function dateFormatter(date, format) {
  if (date.toString() === 'Invalid Date') {
    return '-';
  }
  var o = {
    /**
     * 完整年份
     * @example 2015 2016 2017 2018
     */
    YYYY: function YYYY() {
      return date.getFullYear().toString();
    },
    /**
     * 年份后两位
     * @example 15 16 17 18
     */
    YY: function YY() {
      return this.YYYY().slice(-2);
    },
    /**
     * 月份，保持两位数
     * @example 01 02 03 .... 11 12
     */
    MM: function MM() {
      return leftPad(this.M(), 2);
    },
    /**
     * 月份
     * @example 1 2 3 .... 11 12
     */
    M: function M() {
      return (date.getMonth() + 1).toString();
    },
    /**
     * 每月中的日期，保持两位数
     * @example 01 02 03 .... 30 31
     */
    DD: function DD() {
      return leftPad(this.D(), 2);
    },
    /**
     * 每月中的日期
     * @example 1 2 3 .... 30 31
     */
    D: function D() {
      return date.getDate().toString();
    },
    /**
     * 小时，24 小时制，保持两位数
     * @example 00 01 02 .... 22 23
     */
    HH: function HH() {
      return leftPad(this.H(), 2);
    },
    /**
     * 小时，24 小时制
     * @example 0 1 2 .... 22 23
     */
    H: function H() {
      return date.getHours().toString();
    },
    /**
     * 小时，12 小时制，保持两位数
     * @example 00 01 02 .... 22 23
     */
    hh: function hh() {
      return leftPad(this.h(), 2);
    },
    /**
     * 小时，12 小时制
     * @example 0 1 2 .... 22 23
     */
    h: function h() {
      var h = (date.getHours() % 12).toString();
      return h === '0' ? '12' : h;
    },
    /**
     * 分钟，保持两位数
     * @example 00 01 02 .... 59 60
     */
    mm: function mm() {
      return leftPad(this.m(), 2);
    },
    /**
     * 分钟
     * @example 0 1 2 .... 59 60
     */
    m: function m() {
      return date.getMinutes().toString();
    },
    /**
     * 秒，保持两位数
     * @example 00 01 02 .... 59 60
     */
    ss: function ss() {
      return leftPad(this.s(), 2);
    },
    /**
     * 秒
     * @example 0 1 2 .... 59 60
     */
    s: function s() {
      return date.getSeconds().toString();
    }
  };
  return Object.keys(o).reduce(function (pre, cur) {
    return pre.replace(new RegExp(cur), function (match) {
      /* eslint-disable */
      return o[match].call(o);
      /* eslint-enable */
    });
  }, format);
}
function leftPad(num, width, c) {
  if (c === void 0) {
    c = '0';
  }
  var numStr = num.toString();
  var padWidth = width - numStr.length;
  return padWidth > 0 ? new Array(padWidth + 1).join(c.toString()) + numStr : numStr;
}

function ServiceDetailTable(props) {
  var _this = this;
  var _a, _b, _c;
  var actions = props.actions,
    _d = props.base,
    platform = _d.platform,
    route = _d.route,
    regionId = _d.regionId,
    _e = props.list,
    serviceResourceList = _e.serviceResourceList,
    externalClusters = _e.externalClusters;
  var serviceResources = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.serviceResources;
  });
  var selectedTab = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.base) === null || _a === void 0 ? void 0 : _a.selectedTab;
  });
  var _navigateDetail = function _navigateDetail(item) {
    var _a, _b, _c, _d;
    var serviceName = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass;
    var instanceName = (_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name;
    var clusterId = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route);
    var resourceType = item === null || item === void 0 ? void 0 : item.kind;
    actions.list.serviceResources.select(item);
    (_d = (_c = actions === null || actions === void 0 ? void 0 : actions.list) === null || _c === void 0 ? void 0 : _c.serviceResources) === null || _d === void 0 ? void 0 : _d.clearPolling();
    router.navigate({
      sub: 'detail'
    }, {
      servicename: serviceName,
      instancename: instanceName,
      clusterid: clusterId,
      resourceType: resourceType,
      mode: 'detail'
    });
  };
  var showInstance = function showInstance(item) {
    var _a, _b, _c, _d;
    var clusterId = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route);
    var serviceName = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass;
    var servicePlan = (_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name;
    // 过滤条件
    var k8sQueryObj = {
      fieldSelector: {
        'spec.servicePlan': servicePlan
      },
      labelSelector: {
        'ssm.infra.tce.io/cluster-id': clusterId
      }
    };
    actions.list.showInstanceDialog(true);
    // 查看实例列表
    (_c = actions.detail.instanceResource) === null || _c === void 0 ? void 0 : _c.clear();
    actions.detail.instanceResource.applyFilter({
      platform: platform,
      clusterId: clusterId,
      serviceName: serviceName,
      regionId: regionId,
      resourceType: ResourceTypeEnum.ServiceResource,
      k8sQueryObj: k8sQueryObj,
      instanceId: (_d = item === null || item === void 0 ? void 0 : item.spec) === null || _d === void 0 ? void 0 : _d.externalID,
      namespace: DefaultNamespace
    });
  };
  var resourceMap = React.useMemo(function () {
    var _a, _b, _c, _d, _e, _f, _g, _h;
    var result = {};
    var resources = (_b = (_a = serviceResourceList === null || serviceResourceList === void 0 ? void 0 : serviceResourceList.list) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.records;
    if (selectedTab === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan)) {
      result = (_e = (_d = (_c = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.list) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.records) === null || _e === void 0 ? void 0 : _e.reduce(function (pre, cur) {
        var _a;
        var _b, _c;
        var instances = resources === null || resources === void 0 ? void 0 : resources.filter(function (item) {
          var _a, _b;
          return ((_a = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = item === null || item === void 0 ? void 0 : item.spec) === null || _b === void 0 ? void 0 : _b.servicePlan) && (Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, cur)) === (Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item));
        });
        return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[(_b = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _b === void 0 ? void 0 : _b.uid] = (_c = instances === null || instances === void 0 ? void 0 : instances.length) !== null && _c !== void 0 ? _c : 0, _a));
      }, {});
    } else if (selectedTab === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource)) {
      result = (_h = (_g = (_f = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.list) === null || _f === void 0 ? void 0 : _f.data) === null || _g === void 0 ? void 0 : _g.records) === null || _h === void 0 ? void 0 : _h.reduce(function (pre, cur) {
        var _a;
        var _b, _c;
        var plan = resources === null || resources === void 0 ? void 0 : resources.find(function (item) {
          var _a, _b;
          return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = cur === null || cur === void 0 ? void 0 : cur.spec) === null || _b === void 0 ? void 0 : _b.servicePlan) && (Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, cur)) === (Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item));
        });
        var map = ((_b = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _b === void 0 ? void 0 : _b.uid) && plan ? (_a = {}, _a[(_c = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _c === void 0 ? void 0 : _c.uid] = plan, _a) : {};
        return tslib.__assign(tslib.__assign({}, pre), map);
      }, {});
    }
    return result;
  }, [(_a = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.list) === null || _a === void 0 ? void 0 : _a.data, (_b = serviceResourceList === null || serviceResourceList === void 0 ? void 0 : serviceResourceList.list) === null || _b === void 0 ? void 0 : _b.data, selectedTab]);
  var resourceMapLoading = ((_c = serviceResourceList === null || serviceResourceList === void 0 ? void 0 : serviceResourceList.list) === null || _c === void 0 ? void 0 : _c.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Fetching);
  var columns = React.useMemo(function () {
    var cols = [];
    if (selectedTab === ResourceTypeEnum.ServicePlan) {
      cols = [{
        key: 'name',
        header: i18n.t('名称'),
        render: function render(item) {
          var _a, _b, _c;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
            type: 'link',
            onClick: function onClick() {
              showInstance(item);
            }
          }, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name), ((_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name) && React__default.createElement(teaComponent.Copy, {
            text: (_c = item === null || item === void 0 ? void 0 : item.metadata) === null || _c === void 0 ? void 0 : _c.name
          }));
        }
      }, {
        key: 'description',
        header: i18n.t('描述'),
        render: function render(item) {
          var _a, _b;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, (_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.description) !== null && _b !== void 0 ? _b : '-'));
        }
      }, {
        key: 'clusterId',
        // header: t('集群ID(集群名称)'),
        header: i18n.t('集群ID'),
        render: function render(item) {
          return React__default.createElement(teaComponent.Text, null, i18n.t('{{clusterId}}', {
            clusterId: Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route)
          }));
        }
      }, {
        key: 'planType',
        header: i18n.t('规格类型'),
        render: function render(item) {
          var _a, _b, _c, _d, _e, _f;
          return React__default.createElement(teaComponent.Text, {
            className: (_c = ServicePlanTypeMap === null || ServicePlanTypeMap === void 0 ? void 0 : ServicePlanTypeMap[((_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['ssm.infra.tce.io/owner']) || ServicePlanTypeEnum.Custom]) === null || _c === void 0 ? void 0 : _c.className
          }, (_f = ServicePlanTypeMap === null || ServicePlanTypeMap === void 0 ? void 0 : ServicePlanTypeMap[((_e = (_d = item === null || item === void 0 ? void 0 : item.metadata) === null || _d === void 0 ? void 0 : _d.labels) === null || _e === void 0 ? void 0 : _e['ssm.infra.tce.io/owner']) || ServicePlanTypeEnum.Custom]) === null || _f === void 0 ? void 0 : _f.text);
        }
      }, {
        key: 'instanceNum',
        header: i18n.t('实例个数'),
        render: function render(item) {
          var _a, _b;
          var data = resourceMap && ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.uid) ? resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.uid] : 0;
          return React__default.createElement(React__default.Fragment, null, resourceMapLoading ? React__default.createElement(teaComponent.Icon, {
            type: "loading"
          }) : React__default.createElement(teaComponent.Text, null, data));
        }
      }, {
        key: 'cpu',
        header: i18n.t('CPU'),
        render: function render(item) {
          var _a;
          return React__default.createElement(teaComponent.Text, null, (_a = item === null || item === void 0 ? void 0 : item.spec.metadata) === null || _a === void 0 ? void 0 : _a.cpu);
        }
      }, {
        key: 'memory',
        header: i18n.t('内存'),
        render: function render(item) {
          var _a;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, (_a = item === null || item === void 0 ? void 0 : item.spec.metadata) === null || _a === void 0 ? void 0 : _a.memory));
        }
      }, {
        key: 'storage',
        header: i18n.t('存储'),
        render: function render(item) {
          var _a;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, (_a = item === null || item === void 0 ? void 0 : item.spec.metadata) === null || _a === void 0 ? void 0 : _a.storage));
        }
      }, {
        key: 'user',
        header: i18n.t('创建人'),
        render: function render(item) {
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, Util.getCreator(platform, item)));
        }
      }, {
        key: 'creationTimestamp',
        header: i18n.t('时间戳'),
        render: function render(item) {
          var _a;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp), 'YYYY-MM-DD HH:mm:ss') || '-'));
        }
      }, {
        key: 'operate',
        header: i18n.t('操作'),
        render: function render(item) {
          return renderOperateButtons(item);
        }
      }];
    } else {
      cols = [{
        key: 'name',
        header: i18n.t('名称'),
        render: function render(item) {
          var _a, _b;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
            type: 'link',
            onClick: function onClick() {
              _navigateDetail(item);
            }
          }, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name), React__default.createElement(teaComponent.Copy, {
            text: (_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name
          }));
        }
      }, {
        key: 'version',
        header: i18n.t('版本'),
        render: function render(item) {
          var _a, _b, _c;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, (_c = (_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.parameters) === null || _b === void 0 ? void 0 : _b.version) !== null && _c !== void 0 ? _c : '-'));
        }
      }, {
        key: 'clusterId',
        // header: t('集群ID(集群名称)'),
        header: i18n.t('集群ID'),
        render: function render(item) {
          return React__default.createElement(teaComponent.Text, null, i18n.t('{{clusterId}}', {
            clusterId: Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route)
          }));
        }
      }, {
        key: 'status',
        header: i18n.t('状态'),
        render: function render(item) {
          var _a, _b, _c, _d, _e, _f, _g, _h, _j;
          var state = ((_a = item === null || item === void 0 ? void 0 : item.status) === null || _a === void 0 ? void 0 : _a.state) || '-';
          var className = (_c = ServiceInstanceMap === null || ServiceInstanceMap === void 0 ? void 0 : ServiceInstanceMap[((_b = item === null || item === void 0 ? void 0 : item.status) === null || _b === void 0 ? void 0 : _b.state) || (ServiceInstanceStatusEnum === null || ServiceInstanceStatusEnum === void 0 ? void 0 : ServiceInstanceStatusEnum.Unknown)]) === null || _c === void 0 ? void 0 : _c.className;
          var isDeleting = showResourceDeleteLoading(item, (_e = (_d = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.list) === null || _d === void 0 ? void 0 : _d.data) === null || _e === void 0 ? void 0 : _e.records, platform);
          return React__default.createElement(React__default.Fragment, null, isDeleting && React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Icon, {
            type: 'loading'
          }), i18n.t('删除中')), !isDeleting && React__default.createElement(teaComponent.Text, {
            className: className + " tea-mr-1n"
          }, state), !isDeleting && !!((_g = (_f = item === null || item === void 0 ? void 0 : item.status) === null || _f === void 0 ? void 0 : _f.conditions) === null || _g === void 0 ? void 0 : _g.length) && React__default.createElement(teaComponent.Bubble, {
            content: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.List, {
              type: "number",
              style: {
                width: '100%'
              }
            }, (_j = (_h = item === null || item === void 0 ? void 0 : item.status) === null || _h === void 0 ? void 0 : _h.conditions) === null || _j === void 0 ? void 0 : _j.map(function (item) {
              return React__default.createElement(teaComponent.List.Item, null, React__default.createElement(teaComponent.Text, {
                className: 'tea-mr-2n'
              }, " ", "".concat(item === null || item === void 0 ? void 0 : item.type, " : ").concat((item === null || item === void 0 ? void 0 : item.reason) || (item === null || item === void 0 ? void 0 : item.message))), React__default.createElement(teaComponent.Icon, {
                type: (item === null || item === void 0 ? void 0 : item.status) === 'True' ? 'success' : 'error'
              }));
            })))
          }, React__default.createElement(teaComponent.Icon, {
            type: "info",
            className: "tea-mr-2n"
          })));
        }
      }, {
        key: 'namespace',
        header: i18n.t('命名空间'),
        render: function render(item) {
          var _a;
          return React__default.createElement(teaComponent.Text, null, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.namespace);
        }
      }, {
        key: 'specific',
        header: i18n.t('规格'),
        render: function render(item) {
          var _a, _b, _c, _d, _e, _f, _g, _h;
          var resource = resourceMap && ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.uid) ? resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.uid] : {};
          if (resourceMapLoading) {
            return React__default.createElement(teaComponent.Icon, {
              type: 'loading'
            });
          }
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, i18n.t('CPU: '), (_d = (_c = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _c === void 0 ? void 0 : _c.metadata) === null || _d === void 0 ? void 0 : _d.cpu), React__default.createElement("p", null, i18n.t('内存: '), (_f = (_e = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _e === void 0 ? void 0 : _e.metadata) === null || _f === void 0 ? void 0 : _f.memory), React__default.createElement("p", null, i18n.t('存储: '), (_h = (_g = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _g === void 0 ? void 0 : _g.metadata) === null || _h === void 0 ? void 0 : _h.storage));
        }
      }, {
        key: 'user',
        header: i18n.t('创建人'),
        render: function render(item) {
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, Util.getCreator(platform, item)));
        }
      }, {
        key: 'creationTimestamp',
        header: i18n.t('时间戳'),
        render: function render(item) {
          var _a;
          return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp), 'YYYY-MM-DD HH:mm:ss') || '-'));
        }
      }, {
        key: 'operate',
        header: i18n.t('操作'),
        render: function render(item) {
          return renderOperateButtons(item);
        }
      }];
    }
    return cols;
  }, [selectedTab, serviceResourceList === null || serviceResourceList === void 0 ? void 0 : serviceResourceList.list]);
  var renderOperateButtons = function renderOperateButtons(item) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t;
    var buttons;
    if (selectedTab === ResourceTypeEnum.ServicePlan) {
      buttons = [{
        text: i18n.t('查看实例'),
        handleClick: function handleClick() {
          showInstance(item);
        }
      }, {
        text: i18n.t('编辑'),
        handleClick: function handleClick() {
          var _a, _b;
          var clusterId = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route);
          var serviceName = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass;
          var resourceIns = (_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name;
          actions.detail.instanceDetail.applyFilter({
            platform: platform,
            serviceName: serviceName,
            clusterId: clusterId,
            resourceIns: resourceIns,
            regionId: regionId,
            resourceType: selectedTab
          });
          actions.list.showCreateResourceDialog(true);
        },
        hidden: ((_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['ssm.infra.tce.io/owner']) === ServicePlanTypeEnum.System,
        disabled: resourceMap && ((_c = item === null || item === void 0 ? void 0 : item.metadata) === null || _c === void 0 ? void 0 : _c.uid) && (resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_d = item === null || item === void 0 ? void 0 : item.metadata) === null || _d === void 0 ? void 0 : _d.uid]) > 0,
        tooltips: resourceMap && ((_e = item === null || item === void 0 ? void 0 : item.metadata) === null || _e === void 0 ? void 0 : _e.uid) && (resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_f = item === null || item === void 0 ? void 0 : item.metadata) === null || _f === void 0 ? void 0 : _f.uid]) > 0 ? i18n.t('当前规格正在被实例引用，不可进行编辑') : i18n.t('')
      }, {
        text: i18n.t('删除'),
        handleClick: function handleClick() {
          actions.list.selectDeleteResources([item]);
        },
        hidden: ((_h = (_g = item === null || item === void 0 ? void 0 : item.metadata) === null || _g === void 0 ? void 0 : _g.labels) === null || _h === void 0 ? void 0 : _h['ssm.infra.tce.io/owner']) === ServicePlanTypeEnum.System,
        disabled: ((_k = (_j = item === null || item === void 0 ? void 0 : item.metadata) === null || _j === void 0 ? void 0 : _j.labels) === null || _k === void 0 ? void 0 : _k['ssm.infra.tce.io/owner']) === ServicePlanTypeEnum.System || resourceMap && ((_l = item === null || item === void 0 ? void 0 : item.metadata) === null || _l === void 0 ? void 0 : _l.uid) && (resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_m = item === null || item === void 0 ? void 0 : item.metadata) === null || _m === void 0 ? void 0 : _m.uid]) > 0,
        tooltips: ((_p = (_o = item === null || item === void 0 ? void 0 : item.metadata) === null || _o === void 0 ? void 0 : _o.labels) === null || _p === void 0 ? void 0 : _p['ssm.infra.tce.io/owner']) === ServicePlanTypeEnum.System ? i18n.t('预设规格') : resourceMap && ((_q = item === null || item === void 0 ? void 0 : item.metadata) === null || _q === void 0 ? void 0 : _q.uid) && (resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_r = item === null || item === void 0 ? void 0 : item.metadata) === null || _r === void 0 ? void 0 : _r.uid]) > 0 ? i18n.t('当前规格正在被实例引用，不可进行删除') : i18n.t('')
      }];
    } else {
      buttons = [{
        text: i18n.t('管理'),
        disabled: !!((_s = item === null || item === void 0 ? void 0 : item.metadata) === null || _s === void 0 ? void 0 : _s.deletionTimestamp),
        handleClick: function handleClick() {
          _navigateDetail(item);
        }
      }, {
        text: i18n.t('删除'),
        handleClick: function handleClick() {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var clusterId, serviceName, resourceIns, instanceId, result;
            var _a, _b, _c, _d, _e;
            return tslib.__generator(this, function (_f) {
              switch (_f.label) {
                case 0:
                  clusterId = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route);
                  serviceName = ((_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass) || ((_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.servicename);
                  resourceIns = (_c = item === null || item === void 0 ? void 0 : item.metadata) === null || _c === void 0 ? void 0 : _c.name;
                  instanceId = ((_d = item === null || item === void 0 ? void 0 : item.spec) === null || _d === void 0 ? void 0 : _d.externalID) || (Util === null || Util === void 0 ? void 0 : Util.getInstanceId(platform, (_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.instancename, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection));
                  return [4 /*yield*/, fetchInstanceResources({
                    platform: platform,
                    clusterId: clusterId,
                    resourceIns: resourceIns,
                    serviceName: serviceName,
                    regionId: regionId,
                    resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding,
                    k8sQueryObj: {
                      labelSelector: {
                        'ssm.infra.tce.io/instance-id': instanceId
                      }
                    }
                  })];
                case 1:
                  result = _f.sent();
                  if (!(result === null || result === void 0 ? void 0 : result.hasError) && (result === null || result === void 0 ? void 0 : result.recordCount) > 0) {
                    teaComponent.message.warning({
                      content: "\u5B9E\u4F8B\u3010".concat(resourceIns, "\u3011\u5173\u8054\u4E86\u670D\u52A1\u7ED1\u5B9A,\u8BF7\u60A8\u89E3\u7ED1\u540E\u518D\u5220\u9664!")
                    });
                    return [2 /*return*/];
                  }

                  !(result === null || result === void 0 ? void 0 : result.hasError) && actions.list.selectDeleteResources([item]);
                  return [2 /*return*/];
              }
            });
          });
        },

        disabled: !!((_t = item === null || item === void 0 ? void 0 : item.metadata) === null || _t === void 0 ? void 0 : _t.deletionTimestamp),
        tooltips: ''
      }];
    }
    return buttons === null || buttons === void 0 ? void 0 : buttons.map(function (item, index) {
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Bubble, {
        key: index,
        content: (item === null || item === void 0 ? void 0 : item.tooltips) ? React__default.createElement(teaComponent.Text, null, item === null || item === void 0 ? void 0 : item.tooltips) : null
      }, !(item === null || item === void 0 ? void 0 : item.hidden) && React__default.createElement(teaComponent.Button, {
        type: 'link',
        key: index,
        onClick: item === null || item === void 0 ? void 0 : item.handleClick,
        className: "tea-mr-2n",
        disabled: item === null || item === void 0 ? void 0 : item.disabled,
        hidden: item === null || item === void 0 ? void 0 : item.hidden
      }, i18n.t('{{text}}', {
        text: item === null || item === void 0 ? void 0 : item.text
      }))));
    });
  };
  return React__default.createElement(ffComponent.TablePanel, {
    style: {
      maxHeight: 600,
      overflow: 'auto'
    },
    cardProps: {
      style: {
        boxShadow: '0 0px 0px transparent',
        border: '1px solid #eee'
      }
    },
    key: "".concat(selectedTab, "table"),
    recordKey: function recordKey(record) {
      var _a;
      return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.uid;
    },
    columns: columns,
    model: serviceResources,
    action: actions.list.serviceResources,
    isNeedPagination: true,
    rowDisabled: function rowDisabled(record) {
      var _a;
      return !!((_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp);
    }
  });
}

function ServiceInstanceTable(props) {
  var _a, _b, _c;
  var actions = props.actions,
    _d = props.base,
    selectedTab = _d.selectedTab,
    platform = _d.platform,
    route = _d.route,
    instanceResource = props.detail.instanceResource,
    servicePlans = props.list.servicePlans;
  var _navigateDetail = function _navigateDetail(item) {
    var _a, _b;
    var serviceName = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass;
    var instanceName = (_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.name;
    var clusterId = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item, route);
    actions.list.serviceResources.select(item);
    router.navigate({
      sub: 'detail'
    }, {
      servicename: serviceName,
      instancename: instanceName,
      clusterid: clusterId
    });
  };
  var resourceMap = React.useMemo(function () {
    var _a, _b, _c, _d;
    var result = {};
    var resources = (_b = (_a = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.list) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.records;
    var serviceResources = (_d = (_c = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.list) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.records;
    result = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.reduce(function (pre, cur) {
      var _a;
      var _b, _c;
      var plan = resources === null || resources === void 0 ? void 0 : resources.find(function (item) {
        var _a, _b;
        return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = cur === null || cur === void 0 ? void 0 : cur.spec) === null || _b === void 0 ? void 0 : _b.servicePlan) && (Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, cur)) === (Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, item));
      });
      var map = ((_b = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _b === void 0 ? void 0 : _b.uid) && plan ? (_a = {}, _a[(_c = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _c === void 0 ? void 0 : _c.uid] = plan, _a) : {};
      return tslib.__assign(tslib.__assign({}, pre), map);
    }, {});
    return result;
  }, [(_a = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.list) === null || _a === void 0 ? void 0 : _a.data, (_b = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.list) === null || _b === void 0 ? void 0 : _b.data, selectedTab]);
  var resourceMapLoading = ((_c = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.list) === null || _c === void 0 ? void 0 : _c.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Fetching);
  var columns = [{
    key: 'name',
    header: i18n.t('名称'),
    render: function render(item) {
      var _a;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
        type: 'link',
        onClick: function onClick() {
          _navigateDetail(item);
        }
      }, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name));
    }
  }, {
    key: 'version',
    header: i18n.t('版本'),
    render: function render(item) {
      var _a, _b, _c;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, (_c = (_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.parameters) === null || _b === void 0 ? void 0 : _b.version) !== null && _c !== void 0 ? _c : '-'));
    }
  }, {
    key: 'clusterId',
    // header: t('集群ID(集群名称)'),
    header: i18n.t('集群ID'),
    render: function render(item) {
      return React__default.createElement(teaComponent.Text, null, i18n.t('{{clusterId}}', {
        clusterId: Util.getClusterId(platform, item, route)
      }));
    }
  }, {
    key: 'status',
    header: i18n.t('状态'),
    render: function render(item) {
      var _a, _b, _c, _d, _e, _f, _g;
      var state = (_b = (_a = item === null || item === void 0 ? void 0 : item.status) === null || _a === void 0 ? void 0 : _a.state) !== null && _b !== void 0 ? _b : '-';
      var className = ((_c = item === null || item === void 0 ? void 0 : item.status) === null || _c === void 0 ? void 0 : _c.state) === (ServiceInstanceStatusEnum === null || ServiceInstanceStatusEnum === void 0 ? void 0 : ServiceInstanceStatusEnum.Ready) ? 'text-success' : '';
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, {
        className: className + " tea-mr-1n"
      }, state), !!((_e = (_d = item === null || item === void 0 ? void 0 : item.status) === null || _d === void 0 ? void 0 : _d.conditions) === null || _e === void 0 ? void 0 : _e.length) && React__default.createElement(teaComponent.Bubble, {
        content: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.List, {
          type: "number"
        }, (_g = (_f = item === null || item === void 0 ? void 0 : item.status) === null || _f === void 0 ? void 0 : _f.conditions) === null || _g === void 0 ? void 0 : _g.map(function (item) {
          return React__default.createElement(teaComponent.List.Item, null, React__default.createElement(teaComponent.Text, {
            className: 'tea-mr-2n'
          }, " ", "".concat(item === null || item === void 0 ? void 0 : item.type, " : ").concat((item === null || item === void 0 ? void 0 : item.reason) || (item === null || item === void 0 ? void 0 : item.message))), React__default.createElement(teaComponent.Icon, {
            type: (item === null || item === void 0 ? void 0 : item.status) === 'True' ? 'success' : 'error'
          }));
        })))
      }, React__default.createElement(teaComponent.Icon, {
        type: "info",
        className: "tea-mr-2n"
      })));
    }
  }, {
    key: 'namespace',
    header: i18n.t('命名空间'),
    render: function render(item) {
      var _a;
      return React__default.createElement(teaComponent.Text, null, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.namespace);
    }
  }, {
    key: 'specific',
    header: i18n.t('规格'),
    render: function render(item) {
      var _a, _b, _c, _d, _e, _f, _g, _h;
      var resource = resourceMap && ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.uid) ? resourceMap === null || resourceMap === void 0 ? void 0 : resourceMap[(_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.uid] : {};
      if (resourceMapLoading) {
        return React__default.createElement(teaComponent.Icon, {
          type: 'loading'
        });
      }
      return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, i18n.t('CPU: '), ((_d = (_c = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _c === void 0 ? void 0 : _c.metadata) === null || _d === void 0 ? void 0 : _d.cpu) || '-'), React__default.createElement("p", null, i18n.t('内存: '), ((_f = (_e = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _e === void 0 ? void 0 : _e.metadata) === null || _f === void 0 ? void 0 : _f.memory) || '-'), React__default.createElement("p", null, i18n.t('存储: '), ((_h = (_g = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _g === void 0 ? void 0 : _g.metadata) === null || _h === void 0 ? void 0 : _h.storage) || '-'));
    }
  }, {
    key: 'user',
    header: i18n.t('创建人'),
    render: function render(item) {
      return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, Util.getCreator(platform, item)));
    }
  }, {
    key: 'creationTimestamp',
    header: i18n.t('时间戳'),
    render: function render(item) {
      var _a;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement("p", null, dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp), 'YYYY-MM-DD HH:mm:ss') || '-'));
    }
  }];
  return React__default.createElement(ffComponent.TablePanel, {
    isNeedCard: false,
    recordKey: function recordKey(record) {
      var _a;
      return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name;
    },
    columns: columns,
    model: instanceResource,
    action: actions.detail.instanceResource,
    isNeedPagination: true
  });
}

function ServiceInstanceTableDialog(props) {
  var showInstanceTableDialog = props.list.showInstanceTableDialog,
    actions = props.actions,
    children = props.children;
  return React__default.createElement(teaComponent.Modal, {
    visible: showInstanceTableDialog && children !== null,
    caption: i18n.t('查看实例'),
    onClose: function onClose() {
      actions.list.showInstanceDialog(false);
    },
    size: 'l'
  }, React__default.createElement(teaComponent.Modal.Body, null, children));
}

function ServiceDetail(props) {
  return React__default.createElement(React__default.Fragment, null, React__default.createElement(ServiceDetailHeader, tslib.__assign({}, props)), React__default.createElement(ServiceDetailTable, tslib.__assign({}, props)), React__default.createElement(ServiceInstanceTableDialog, tslib.__assign({}, props, {
    children: React__default.createElement(React__default.Fragment, null, React__default.createElement(ServiceInstanceTable, tslib.__assign({}, props)))
  })), React__default.createElement(ServiceCreateDialog, tslib.__assign({}, props)));
}

var routerSea$2 = seajs.require('router');
function PaasContentPanel(props) {
  var base = props.base,
    actions = props.actions;
  var selectedTab = base === null || base === void 0 ? void 0 : base.selectedTab;
  var serviceName = reactRedux.useSelector(function (state) {
    var _a, _b, _c;
    return (_c = (_b = (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.services) === null || _b === void 0 ? void 0 : _b.selection) === null || _c === void 0 ? void 0 : _c.name;
  });
  return React__default.createElement(teaComponent.Layout.Content, {
    className: 'tea-bg-layout--border',
    style: {
      height: '100%',
      backgroundColor: '#fff',
      paddingLeft: 20,
      paddingRight: 20,
      overflow: 'hidden'
    }
  }, React__default.createElement(teaComponent.Row, {
    style: {
      margin: 0
    }
  }, React__default.createElement(teaComponent.Col, {
    span: 24,
    className: 'tea-pt-4n tea-pb-4n',
    style: {
      paddingLeft: 0,
      paddingRight: 0
    }
  }, React__default.createElement(teaComponent.H3, null, i18n.t('{{serviceName}}', {
    serviceName: serviceName
  })))), React__default.createElement(teaComponent.Card, {
    style: {
      boxShadow: '0 0px 0px transparent'
    }
  }, React__default.createElement(teaComponent.Tabs, {
    activeId: selectedTab,
    defaultActiveId: selectedTab,
    tabs: serviceMngTabs,
    placement: 'top',
    onActive: function onActive(tab) {
      actions.base.selectTab(tab === null || tab === void 0 ? void 0 : tab.id);
    }
  }, serviceMngTabs === null || serviceMngTabs === void 0 ? void 0 : serviceMngTabs.map(function (tab) {
    return React__default.createElement(teaComponent.TabPanel, {
      id: tab === null || tab === void 0 ? void 0 : tab.id
    });
  })), React__default.createElement(ServiceDetail, tslib.__assign({}, props))));
}

function PaasSiderPanel(props) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
  var _m = React.useState(''),
    selected = _m[0],
    setSelected = _m[1];
  var _o = React.useState({
      records: [],
      recordCount: 0,
      fetched: false
    }),
    menus = _o[0],
    setMenus = _o[1];
  var services = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.services;
  });
  var actions = props.actions;
  function LoadingPanel() {
    return React__default.createElement(teaComponent.Card, {
      bordered: false,
      style: {
        height: '100%',
        display: 'flex',
        alignItems: "center",
        justifyContent: 'center',
        width: 200,
        maxWidth: 200,
        boxShadow: '0 0px 0px transparent',
        border: '1px solid #eee'
      }
    }, React__default.createElement(teaComponent.Card.Body, null, React__default.createElement(teaComponent.Icon, {
      type: 'loading'
    }), React__default.createElement(teaComponent.Text, null, i18n.t('加载中'))));
  }
  function EmptyPanel() {
    return React__default.createElement(teaComponent.Card, {
      bordered: false,
      style: {
        height: '100%',
        display: 'flex',
        alignItems: "center",
        justifyContent: 'center',
        width: 200,
        maxWidth: 200,
        border: '1px solid #eee'
      }
    }, React__default.createElement(teaComponent.Card.Body, null, React__default.createElement(i18n.Trans, null, "\u60A8\u5C1A\u672A\u5F00\u901A\u5206\u5E03\u5F0F\u4E91\u4E2D\u5FC3\u6570\u636E\u670D\u52A1\uFF0C\u8BF7\u524D\u5F80\u5206\u5E03\u5F0F\u4E91\u4E2D\u5FC3", React__default.createElement("a", {
      href: 'https://console.cloud.tencent.com/tdcc/paasoverview',
      target: "_blank"
    }, "\u5F00\u901A\u670D\u52A1"))));
  }
  function RetryPanel(retryProps) {
    return React__default.createElement(teaComponent.Card, {
      bordered: false,
      style: {
        height: '100%',
        display: 'flex',
        alignItems: "center",
        justifyContent: 'center',
        width: 200,
        maxWidth: 200,
        boxShadow: '0 0px 0px transparent',
        border: '1px solid #eee'
      }
    }, React__default.createElement(teaComponent.Text, {
      className: 'text-danger tea-mr-2n'
    }, i18n.t('加载失败')), React__default.createElement(teaComponent.Button, {
      type: 'link',
      onClick: function onClick() {
        retryProps === null || retryProps === void 0 ? void 0 : retryProps.onRetry();
      }
    }, i18n.t('刷新重试')));
  }
  var _renderMenuList = function _renderMenuList(subRouters) {
    return React__default.createElement(teaComponent.Menu, null, subRouters === null || subRouters === void 0 ? void 0 : subRouters.map(function (item, index) {
      var _a, _b, _c, _d;
      return React__default.createElement("div", {
        key: item === null || item === void 0 ? void 0 : item.name
      }, !((_a = item === null || item === void 0 ? void 0 : item.sub) === null || _a === void 0 ? void 0 : _a.length) && React__default.createElement(teaComponent.Menu.Item, {
        style: index === 0 ? {
          paddingTop: 6
        } : {},
        title: item === null || item === void 0 ? void 0 : item.name,
        selected: ((_b = services === null || services === void 0 ? void 0 : services.selection) === null || _b === void 0 ? void 0 : _b.name) === (item === null || item === void 0 ? void 0 : item.name),
        onClick: function onClick() {
          var _a, _b;
          (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.services) === null || _b === void 0 ? void 0 : _b.select(item !== null && item !== void 0 ? item : {});
        }
      }), ((_c = item === null || item === void 0 ? void 0 : item.sub) === null || _c === void 0 ? void 0 : _c.length) > 0 && React__default.createElement(teaComponent.Menu.SubMenu, {
        title: item === null || item === void 0 ? void 0 : item.name,
        key: item === null || item === void 0 ? void 0 : item.name
      }, (_d = item === null || item === void 0 ? void 0 : item.sub) === null || _d === void 0 ? void 0 : _d.map(function (subItem) {
        var _a;
        return React__default.createElement(teaComponent.Menu.Item, {
          key: subItem === null || subItem === void 0 ? void 0 : subItem.name,
          title: subItem === null || subItem === void 0 ? void 0 : subItem.name,
          selected: ((_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.name) === (item === null || item === void 0 ? void 0 : item.name),
          onClick: function onClick() {
            var _a, _b;
            (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.services) === null || _b === void 0 ? void 0 : _b.select(item !== null && item !== void 0 ? item : {});
          }
        });
      })));
    }));
  };
  var isEmpty = ((_a = services === null || services === void 0 ? void 0 : services.list) === null || _a === void 0 ? void 0 : _a.fetched) && ((_d = (_c = (_b = services === null || services === void 0 ? void 0 : services.list) === null || _b === void 0 ? void 0 : _b.data) === null || _c === void 0 ? void 0 : _c.records) === null || _d === void 0 ? void 0 : _d.length) === 0;
  var isLoading = ((_e = services === null || services === void 0 ? void 0 : services.list) === null || _e === void 0 ? void 0 : _e.fetchState) === ffRedux.FetchState.Fetching || !((_f = services === null || services === void 0 ? void 0 : services.list) === null || _f === void 0 ? void 0 : _f.fetched);
  var failed = ((_g = services === null || services === void 0 ? void 0 : services.list) === null || _g === void 0 ? void 0 : _g.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Failed) || ((_j = (_h = services === null || services === void 0 ? void 0 : services.list) === null || _h === void 0 ? void 0 : _h.data) === null || _j === void 0 ? void 0 : _j['hasError']);
  return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Layout.Sider, {
    className: 'tea-bg-layout--border'
  }, isLoading && React__default.createElement(LoadingPanel, null), !isLoading && failed && React__default.createElement(RetryPanel, {
    onRetry: function onRetry() {
      var _a, _b;
      (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.services) === null || _b === void 0 ? void 0 : _b.fetch();
    }
  }), !isLoading && !failed && isEmpty && React__default.createElement(EmptyPanel, null), !isLoading && !failed && !isEmpty && _renderMenuList((_l = (_k = services === null || services === void 0 ? void 0 : services.list) === null || _k === void 0 ? void 0 : _k.data) === null || _l === void 0 ? void 0 : _l.records)));
}

function PaasContent(props) {
  return React__default.createElement(teaComponent.Card, {
    style: {
      padding: 0
    }
  }, React__default.createElement(teaComponent.Layout, {
    style: {
      height: 840,
      minHeight: 840
    }
  }, React__default.createElement(teaComponent.Layout.Body, null, React__default.createElement(PaasSiderPanel, tslib.__assign({}, props)), React__default.createElement(PaasContentPanel, tslib.__assign({}, props)))));
}

function prefixPageId(obj, pageId) {
  Object.keys(obj).forEach(function (key) {
    obj[key] = pageId + '_' + obj[key];
  });
}
var Base;
(function (Base) {
  Base["IsI18n"] = "IsI18n";
  Base["FETCH_PLATFORM"] = "FETCH_PLATFORM";
  Base["FETCH_REGION"] = "FETCH_REGION";
  Base["HubCluster"] = "HubCluster";
  Base["ClusterVersion"] = "ClusterVersion";
  Base["SELECT_TAB"] = "SELECT_TAB";
  Base["Clear"] = "Clear";
  Base["FETCH_UserInfo"] = "FETCH_UserInfo";
  Base["UPDATE_ROUTE"] = "UPDATE_ROUTE";
  Base["GetClusterAdminRoleFlow"] = "GetClusterAdminRoleFlow";
})(Base || (Base = {}));
var Create;
(function (Create) {
  Create["SERVICE_INSTANCE_EDIT"] = "SERVICE_INSTANCE_EDIT";
  Create["SERVICE_INSTANCE_EDIT_VALIDATOR"] = "SERVICE_INSTANCE_EDIT_VALIDATOR";
  Create["CREATE_SERVICE_RESOURCE"] = "CREATE_SERVICE_RESOURCE";
  Create["UPDATE_SERVICE_RESOURCE"] = "UPDATE_SERVICE_RESOURCE";
  Create["SERVICE_PLAN_EDIT"] = "SERVICE_PLAN_EDIT";
  Create["CREATE_SERVICE_INSTANCE"] = "CREATE_SERVICE_INSTANCE";
  Create["BackupNowWorkflow"] = "BackupNowWorkflow";
})(Create || (Create = {}));
var List;
(function (List) {
  List["FETCH_SERVICES"] = "FETCH_SERVICES";
  List["FETCH_CREATE_RESOURCE_SCHEMAS"] = "FETCH_CREATE_RESOURCE_SCHEMAS";
  List["FETCH_SERVICE_RESOURCE"] = "FETCH_SERVICE_RESOURCE";
  List["FETCH_SERVICE_PLANS"] = "FETCH_SERVICE_PLANS";
  List["FETCH_Service_Plan_Map"] = "FETCH_Service_Plan_Map";
  List["FETCH_SERVICE_RESOURCE_LIST"] = "FETCH_SERVICE_RESOURCE_LIST";
  List["SELECT_DELETE_RESOURCE"] = "SELECT_DELETE_RESOURCE";
  List["DELETE_SERVICE_RESOURCE"] = "DELETE_SERVICE_RESOURCE";
  List["SHOW_INSTANCE_TABLE_DIALOG"] = "SHOW_INSTANCE_TABLE_DIALOG";
  List["SHOW_CREATE_RESOURCE_DIALOG"] = "SHOW_CREATE_RESOURCE_DIALOG";
  List["FETCH_EXTERNAL_CLUSTERS"] = "FETCH_EXTERNAL_CLUSTERS";
  List["FETCH_Affinity_Localization"] = "FETCH_Affinity_Localization";
  List["Clear"] = "Clear";
})(List || (List = {}));
var Detail;
(function (Detail) {
  Detail["RESOURCE_DETAIL"] = "RESOURCE_DETAIL";
  Detail["Clear"] = "Clear";
  Detail["Select_Detail_Tab"] = "Select_Detail_Tab";
  Detail["FETCH_INSTANCE_RESOURCE"] = "FETCH_INSTANCE_RESOURCE";
  Detail["BACKUP_RESOURCE"] = "BACKUP_RESOURCE";
  Detail["CHECK_COS"] = "CHECK_COS";
  Detail["BACKUP_RESOURCE_LOADING"] = "BACKUP_RESOURCE_LOADING";
  Detail["SHOW_BACKUP_STRATEGY_DIALOG"] = "SHOW_BACKUP_STRATEGY_DIALOG";
  Detail["BACKUP_STRATEGY_EDIT"] = "BACKUP_STRATEGY_EDIT";
  Detail["OPEN_CONSOLE_WORKFLOW"] = "OPEN_CONSOLE_WORKFLOW";
  Detail["SHOW_CREATE_RESOURCE_DIALOG"] = "SHOW_CREATE_RESOURCE_DIALOG";
  Detail["FETCH_NAMESPACES"] = "FETCH_NAMESPACES";
  Detail["SERVICE_BINDING_EDIT"] = "SERVICE_BINDING_EDIT";
  Detail["FETCH_SERVICE_INSTANCE_SCHEMA"] = "FETCH_SERVICE_INSTANCE_SCHEMA";
  Detail["SELECT_DETAIL_RESOURCE"] = "SELECT_DETAIL_RESOURCE";
  Detail["FETCH_BACKUP_STRATEGY"] = "FETCH_BACKUP_STRATEGY";
  Detail["BackupStrategyMedium"] = "BackupStrategyMedium";
})(Detail || (Detail = {}));
prefixPageId(Base, 'Base');
prefixPageId(Create, 'Create');
prefixPageId(List, 'List');
prefixPageId(Detail, 'Detail');
var ActionType = {
  Base: Base,
  Create: Create,
  List: List,
  Detail: Detail
};

var _a$6;
var _b$5, _c$3;
var detailActions = {
  selectDetailTab: function selectDetailTab(tabId) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, _b, route, platform, selectedTab, regionId, serviceResources, _c, servicename, instancename, clusterid, instanceId;
        var _d, _e, _f, _g, _h;
        return tslib.__generator(this, function (_j) {
          dispatch({
            type: ActionType.Detail.Select_Detail_Tab,
            payload: tabId
          });
          _a = getState(), _b = _a.base, route = _b.route, platform = _b.platform, selectedTab = _b.selectedTab, regionId = _b.regionId, serviceResources = _a.list.serviceResources;
          _c = route === null || route === void 0 ? void 0 : route.queries, servicename = _c.servicename, instancename = _c.instancename, clusterid = _c.clusterid;
          instanceId = Util === null || Util === void 0 ? void 0 : Util.getInstanceId(platform, instancename);
          dispatch((_d = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceDetail) === null || _d === void 0 ? void 0 : _d.clear());
          dispatch((_e = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _e === void 0 ? void 0 : _e.clear());
          dispatch((_f = detailActions === null || detailActions === void 0 ? void 0 : detailActions.serviceInstanceSchema) === null || _f === void 0 ? void 0 : _f.clear());
          if (tabId === DetailTabType.Detail) {
            dispatch(detailActions.instanceDetail.applyFilter({
              platform: platform,
              clusterId: clusterid,
              resourceIns: instancename,
              serviceName: servicename,
              regionId: regionId,
              resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource,
              namespace: (_h = (_g = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection) === null || _g === void 0 ? void 0 : _g.metadata) === null || _h === void 0 ? void 0 : _h.namespace
            }));
          } else if ([DetailTabType.BackUp].includes(tabId)) {
            dispatch(detailActions.instanceResource.applyFilter({
              platform: platform,
              clusterId: clusterid,
              resourceIns: instancename,
              serviceName: servicename,
              regionId: regionId,
              resourceType: tabId,
              instanceId: instanceId,
              k8sQueryObj: {
                labelSelector: {
                  "ssm.infra.tce.io/instance-id": instanceId
                }
              }
            }));
          } else if ([DetailTabType.ServiceBinding].includes(tabId)) {
            dispatch(detailActions.instanceResource.applyFilter({
              platform: platform,
              clusterId: clusterid,
              resourceIns: instancename,
              serviceName: servicename,
              regionId: regionId,
              resourceType: tabId,
              instanceId: instanceId,
              k8sQueryObj: {
                labelSelector: {
                  "ssm.infra.tce.io/instance-id": instanceId
                }
              }
            }));
            dispatch(detailActions.serviceInstanceSchema.applyFilter({
              platform: platform,
              clusterId: clusterid,
              serviceName: servicename,
              regionId: regionId,
              resourceType: tabId
            }));
          }
          return [2 /*return*/];
        });
      });
    };
  },

  instanceDetail: ffRedux.createFFObjectActions({
    actionName: ActionType.Detail.RESOURCE_DETAIL,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, clusterId, regionId, resourceIns, serviceName, resourceType, _b, namespace, _c, _d, platform, route, selectDetailResource, response, resource, secretResource, resourceSchemas, instanceSchema;
        var _e, _f;
        return tslib.__generator(this, function (_g) {
          switch (_g.label) {
            case 0:
              _a = query.filter, clusterId = _a.clusterId, regionId = _a.regionId, resourceIns = _a.resourceIns, serviceName = _a.serviceName, resourceType = _a.resourceType, _b = _a.namespace, namespace = _b === void 0 ? DefaultNamespace : _b;
              _c = getState(), _d = _c.base, platform = _d.platform, route = _d.route, selectDetailResource = _c.detail.selectDetailResource;
              response = {
                resource: null,
                instanceSchema: []
              };
              return [4 /*yield*/, fetchServiceResourceDetail({
                clusterId: clusterId,
                regionId: regionId,
                serviceName: serviceName,
                platform: platform,
                resourceIns: resourceIns,
                resourceType: resourceType ? resourceType : ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource,
                namespace: namespace
              })];
            case 1:
              resource = _g.sent();
              if (!(resourceType === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding))) return [3 /*break*/, 3];
              return [4 /*yield*/, fetchServiceResourceDetail({
                clusterId: clusterId,
                regionId: regionId,
                serviceName: serviceName,
                platform: platform,
                resourceIns: (_e = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _e === void 0 ? void 0 : _e.secretName,
                resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Secret,
                namespace: namespace
              })];
            case 2:
              secretResource = _g.sent();
              _g.label = 3;
            case 3:
              return [4 /*yield*/, fetchResourceSchemas({
                clusterId: clusterId,
                regionId: regionId,
                serviceName: serviceName,
                platform: platform,
                resourceType: resourceType
              })];
            case 4:
              resourceSchemas = _g.sent();
              instanceSchema = [];
              if (resourceType === ResourceTypeEnum.ServicePlan) {
                instanceSchema = resourceSchemas === null || resourceSchemas === void 0 ? void 0 : resourceSchemas.planSchema;
              } else if (resourceType === ResourceTypeEnum.ServiceBinding) {
                instanceSchema = resourceSchemas === null || resourceSchemas === void 0 ? void 0 : resourceSchemas.bindingResponseSchema;
              } else {
                instanceSchema = resourceSchemas === null || resourceSchemas === void 0 ? void 0 : resourceSchemas.instanceCreateParameterSchema;
              }
              response = {
                resource: resource,
                instanceSchema: (_f = instanceSchema === null || instanceSchema === void 0 ? void 0 : instanceSchema.filter(function (item) {
                  var _a;
                  if (item === null || item === void 0 ? void 0 : item.enabledCondition) {
                    var _b = (_a = item === null || item === void 0 ? void 0 : item.enabledCondition) === null || _a === void 0 ? void 0 : _a.split("=="),
                      conditionKey_1 = _b[0],
                      conditionValue = _b[1];
                    var values = instanceSchema === null || instanceSchema === void 0 ? void 0 : instanceSchema.find(function (schema) {
                      return (schema === null || schema === void 0 ? void 0 : schema.name) === conditionKey_1;
                    });
                    var value = values === null || values === void 0 ? void 0 : values.value;
                    return String(value) === String(conditionValue);
                  } else {
                    return true;
                  }
                }).map(function (item) {
                  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m;
                  var schemaValue = "";
                  var schemaName = item === null || item === void 0 ? void 0 : item.name;
                  switch (resourceType) {
                    case ResourceTypeEnum.ServiceResource:
                      schemaValue = ((_a = resource === null || resource === void 0 ? void 0 : resource.status) === null || _a === void 0 ? void 0 : _a.metadata) && schemaName || !isEmpty((_b = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _b === void 0 ? void 0 : _b.parameters) ? formatPlanSchemaValue(((_d = (_c = resource === null || resource === void 0 ? void 0 : resource.status) === null || _c === void 0 ? void 0 : _c.metadata) === null || _d === void 0 ? void 0 : _d[schemaName]) || ((_e = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _e === void 0 ? void 0 : _e.parameters[schemaName])) || "-" : "";
                      break;
                    case ResourceTypeEnum.ServiceBinding:
                      schemaValue = (secretResource === null || secretResource === void 0 ? void 0 : secretResource.data) && schemaName && ((_f = secretResource === null || secretResource === void 0 ? void 0 : secretResource.data) === null || _f === void 0 ? void 0 : _f[schemaName]) ? ((_h = jsBase64.Base64.atob((_g = secretResource === null || secretResource === void 0 ? void 0 : secretResource.data) === null || _g === void 0 ? void 0 : _g[schemaName])) === null || _h === void 0 ? void 0 : _h.toString()) || "-" : "";
                      break;
                    default:
                      schemaValue = ((_j = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _j === void 0 ? void 0 : _j.metadata) && schemaName ? formatPlanSchemaValue((_m = (_l = (_k = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _k === void 0 ? void 0 : _k.metadata) === null || _l === void 0 ? void 0 : _l[schemaName]) === null || _m === void 0 ? void 0 : _m.toString()) || "-" : "";
                      break;
                  }
                  return tslib.__assign(tslib.__assign({}, item), {
                    value: formatPlanSchemaData(item === null || item === void 0 ? void 0 : item.type, schemaValue),
                    unit: formatPlanSchemaUnit(item, schemaValue)
                  });
                })) === null || _f === void 0 ? void 0 : _f.filter(function (item) {
                  return !hideSchema(item);
                })
              };
              return [2 /*return*/, response];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().detail.resourceDetail;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o;
      if (record === null || record === void 0 ? void 0 : record.data) {
        // 更新服务实例编辑数据
        var _p = getState(),
          resourceDetail = _p.detail.resourceDetail,
          selectedTab = _p.base.selectedTab;
        var basePro = {
          clusterId: (_b = (_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.query) === null || _a === void 0 ? void 0 : _a.filter) === null || _b === void 0 ? void 0 : _b.clusterId,
          instanceName: (_e = (_d = (_c = record === null || record === void 0 ? void 0 : record.data) === null || _c === void 0 ? void 0 : _c.resource) === null || _d === void 0 ? void 0 : _d.metadata) === null || _e === void 0 ? void 0 : _e.name,
          description: (_h = (_g = (_f = record === null || record === void 0 ? void 0 : record.data) === null || _f === void 0 ? void 0 : _f.resource) === null || _g === void 0 ? void 0 : _g.spec) === null || _h === void 0 ? void 0 : _h.description,
          serviceName: (_l = (_k = (_j = record === null || record === void 0 ? void 0 : record.data) === null || _j === void 0 ? void 0 : _j.resource) === null || _k === void 0 ? void 0 : _k.spec) === null || _l === void 0 ? void 0 : _l.serviceClass,
          unitMap: {}
        };
        var instanceFormData = (_o = (_m = record === null || record === void 0 ? void 0 : record.data) === null || _m === void 0 ? void 0 : _m.instanceSchema) === null || _o === void 0 ? void 0 : _o.reduce(function (pre, cur) {
          var _a, _b;
          return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = cur === null || cur === void 0 ? void 0 : cur.value, _a.unitMap = tslib.__assign(tslib.__assign({}, pre === null || pre === void 0 ? void 0 : pre.unitMap), (_b = {}, _b[cur === null || cur === void 0 ? void 0 : cur.name] = cur === null || cur === void 0 ? void 0 : cur.unit, _b)), _a));
        }, basePro);
        var instanceEdit = {
          formData: instanceFormData
        };
        if (selectedTab === ResourceTypeEnum.ServiceResource) {
          dispatch(createActions.serviceInstanceEdit(instanceEdit));
        } else {
          dispatch(createActions.servicePlanEdit(instanceEdit));
        }
      }
    }
  }),
  // 实例资源列表
  instanceResource: ffRedux.createFFListActions({
    actionName: ActionType.Detail.FETCH_INSTANCE_RESOURCE,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var result;
        var _a;
        return tslib.__generator(this, function (_b) {
          switch (_b.label) {
            case 0:
              if (!(query === null || query === void 0 ? void 0 : query.filter)) return [3 /*break*/, 5];
              if (!(((_a = query === null || query === void 0 ? void 0 : query.filter) === null || _a === void 0 ? void 0 : _a.resourceType) === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource))) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchServiceResources(tslib.__assign(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter), {
                paging: query === null || query === void 0 ? void 0 : query.paging
              }))];
            case 1:
              result = _b.sent();
              return [3 /*break*/, 4];
            case 2:
              return [4 /*yield*/, fetchInstanceResources(tslib.__assign(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter), {
                paging: query === null || query === void 0 ? void 0 : query.paging
              }))];
            case 3:
              result = _b.sent();
              _b.label = 4;
            case 4:
              return [3 /*break*/, 6];
            case 5:
              result.records = [];
              _b.label = 6;
            case 6:
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().detail.instanceResource;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b, _c, _d, _e, _f;
      var instanceResource = getState().detail.instanceResource;
      //判断是否开启轮训
      if ((_b = (_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.records) === null || _b === void 0 ? void 0 : _b.some(function (item) {
        var _a;
        return (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp;
      })) {
        dispatch((_c = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _c === void 0 ? void 0 : _c.polling({
          delayTime: 8000
        }));
      } else {
        dispatch((_d = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _d === void 0 ? void 0 : _d.clearPolling());
      }
      if ((_e = record === null || record === void 0 ? void 0 : record.data) === null || _e === void 0 ? void 0 : _e.recordCount) {
        var _g = (_f = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.query) === null || _f === void 0 ? void 0 : _f.filter,
          platform = _g.platform,
          clusterId = _g.clusterId,
          regionId = _g.regionId,
          serviceName = _g.serviceName;
        // 查看实例关联的规格资源
        dispatch(listActions === null || listActions === void 0 ? void 0 : listActions.servicePlans.applyFilter({
          platform: platform,
          clusterId: clusterId,
          regionId: regionId,
          serviceName: serviceName,
          resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan
        }));
      }
    },
    onSelect: function onSelect(record, dispatch, getState) {
      var route = getState().base.route;
    }
  }),
  checkCosResource: ffRedux.createFFObjectActions({
    actionName: ActionType.Detail.CHECK_COS,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var filter, _a, platform, route, response;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              filter = query.filter;
              _a = (_b = getState()) === null || _b === void 0 ? void 0 : _b.base, platform = _a.platform, route = _a.route;
              if (!filter) return [3 /*break*/, 2];
              return [4 /*yield*/, checkCosResource(filter)];
            case 1:
              response = _c.sent();
              _c.label = 2;
            case 2:
              return [2 /*return*/, response];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().detail.checkCosResource;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a;
      if (record === null || record === void 0 ? void 0 : record.data) {
        var checkCosResource_1 = getState().detail.checkCosResource;
        var filter = (_a = checkCosResource_1 === null || checkCosResource_1 === void 0 ? void 0 : checkCosResource_1.query) === null || _a === void 0 ? void 0 : _a.filter;
        if (filter) {
          dispatch(createActions.createResource.start([filter], filter === null || filter === void 0 ? void 0 : filter.regionId));
        }
      } else {
        bridge.tips.error(i18n.t("{{msg}}", {
          msg: ErrorMsgEnum === null || ErrorMsgEnum === void 0 ? void 0 : ErrorMsgEnum.COS_Resource_Not_Found
        }));
      }
    }
  }),
  showBackupDialog: function showBackupDialog(show) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        return tslib.__generator(this, function (_a) {
          dispatch({
            type: ActionType.Detail.SHOW_BACKUP_STRATEGY_DIALOG,
            payload: show
          });
          if (!show) {
            dispatch(detailActions.resetBackupStrategyEdit());
          }
          return [2 /*return*/];
        });
      });
    };
  },

  selectDetailResource: function selectDetailResource(resources) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, _b, platform, route, regionId, services, clusterId, resourceIns, namespace;
        var _c, _d, _e, _f, _g, _h;
        return tslib.__generator(this, function (_j) {
          _a = getState(), _b = _a.base, platform = _b.platform, route = _b.route, regionId = _b.regionId, services = _a.list.services;
          dispatch({
            type: ActionType.Detail.SELECT_DETAIL_RESOURCE,
            payload: resources
          });
          if (resources === null || resources === void 0 ? void 0 : resources[0]) {
            clusterId = Util.getClusterId(platform, resources === null || resources === void 0 ? void 0 : resources[0], route);
            resourceIns = (_d = (_c = resources === null || resources === void 0 ? void 0 : resources[0]) === null || _c === void 0 ? void 0 : _c.metadata) === null || _d === void 0 ? void 0 : _d.name;
            namespace = (_f = (_e = resources === null || resources === void 0 ? void 0 : resources[0]) === null || _e === void 0 ? void 0 : _e.metadata) === null || _f === void 0 ? void 0 : _f.namespace;
            dispatch(detailActions.instanceDetail.applyFilter({
              platform: platform,
              clusterId: clusterId,
              resourceIns: resourceIns,
              serviceName: ((_g = services === null || services === void 0 ? void 0 : services.selection) === null || _g === void 0 ? void 0 : _g.name) || ((_h = route === null || route === void 0 ? void 0 : route.queries) === null || _h === void 0 ? void 0 : _h.servicename),
              regionId: regionId,
              resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding,
              namespace: namespace
            }));
          }
          return [2 /*return*/];
        });
      });
    };
  },

  initBackUpStrategy: function initBackUpStrategy(formData) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var newEdit;
        var _a, _b;
        return tslib.__generator(this, function (_c) {
          newEdit = {
            validator: (_a = Object.keys(formData)) === null || _a === void 0 ? void 0 : _a.reduce(function (pre, key) {
              var _a;
              return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[key] = {
                message: "",
                status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Success
              }, _a));
            }, {}),
            formData: tslib.__assign({}, formData)
          };
          dispatch({
            type: (_b = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _b === void 0 ? void 0 : _b.BACKUP_STRATEGY_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  updateBackUpStrategy: function updateBackUpStrategy(key, value) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b;
        var _c, _d, _e;
        return tslib.__generator(this, function (_f) {
          if (key) {
            _a = (_d = (_c = getState()) === null || _c === void 0 ? void 0 : _c.detail) === null || _d === void 0 ? void 0 : _d.backupStrategyEdit, formData = _a.formData, validator = _a.validator;
            newEdit = {
              validator: validator,
              formData: tslib.__assign(tslib.__assign({}, formData), (_b = {}, _b[key] = value, _b))
            };
            dispatch({
              type: (_e = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _e === void 0 ? void 0 : _e.BACKUP_STRATEGY_EDIT,
              payload: newEdit
            });
          }
          return [2 /*return*/];
        });
      });
    };
  },

  mediums: MediumSelectPanel.createActions({
    pageName: (_b$5 = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _b$5 === void 0 ? void 0 : _b$5.BackupStrategyMedium,
    getRecord: function getRecord(getState) {
      var _a, _b;
      return (_b = (_a = getState()) === null || _a === void 0 ? void 0 : _a.detail) === null || _b === void 0 ? void 0 : _b.mediums;
    }
  }),
  validateAll: function validateAll() {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var backupStrategyEdit, _a, formData, validator, newValidatorModel, newEdit;
        var _b, _c, _d;
        return tslib.__generator(this, function (_e) {
          switch (_e.label) {
            case 0:
              backupStrategyEdit = getState().detail.backupStrategyEdit;
              if (!backupStrategyEdit) return [3 /*break*/, 2];
              _a = (_c = (_b = getState()) === null || _b === void 0 ? void 0 : _b.detail) === null || _c === void 0 ? void 0 : _c.backupStrategyEdit, formData = _a.formData, validator = _a.validator;
              return [4 /*yield*/, Backup.getValidatorModel(backupStrategyEdit)];
            case 1:
              newValidatorModel = _e.sent();
              newEdit = {
                formData: formData,
                validator: tslib.__assign(tslib.__assign({}, validator), newValidatorModel)
              };
              dispatch({
                type: (_d = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _d === void 0 ? void 0 : _d.BACKUP_STRATEGY_EDIT,
                payload: newEdit
              });
              _e.label = 2;
            case 2:
              return [2 /*return*/];
          }
        });
      });
    };
  },

  resetBackupStrategyEdit: function resetBackupStrategyEdit() {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var backupStrategyEdit, newEdit;
        var _a;
        return tslib.__generator(this, function (_b) {
          backupStrategyEdit = getState().detail.backupStrategyEdit;
          newEdit = tslib.__assign({}, BackupStrategyEditInitValue);
          dispatch({
            type: (_a = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _a === void 0 ? void 0 : _a.BACKUP_STRATEGY_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  openInstanceConsole: ffRedux.generateWorkflowActionCreator({
    actionType: (_c$3 = ActionType.Detail) === null || _c$3 === void 0 ? void 0 : _c$3.OPEN_CONSOLE_WORKFLOW,
    workflowStateLocator: function workflowStateLocator(state) {
      var _a;
      return (_a = state.detail) === null || _a === void 0 ? void 0 : _a.openConsoleWorkflow;
    },
    operationExecutor: openInstanceConsole,
    after: (_a$6 = {}, _a$6[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {
      var openConsoleWorkflow = getState().detail.openConsoleWorkflow;
      if (ffRedux.isSuccessWorkflow(openConsoleWorkflow)) {
        dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceDetail.fetch());
      }
    }, _a$6)
  }),
  showCreateResourceDialog: function showCreateResourceDialog(show) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, serviceInstanceSchema, route, _b, route_1, platform, regionId;
        var _c, _d, _e, _f, _g, _h;
        return tslib.__generator(this, function (_j) {
          _a = getState(), serviceInstanceSchema = _a.detail.serviceInstanceSchema, route = _a.base.route;
          dispatch({
            type: ActionType.Detail.SHOW_CREATE_RESOURCE_DIALOG,
            payload: show
          });
          if (!show) {
            dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.initServiceBindingEdit((_e = (_d = (_c = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.bindingCreateParameterSchema) !== null && _e !== void 0 ? _e : []));
          } else {
            _b = getState().base, route_1 = _b.route, platform = _b.platform, regionId = _b.regionId;
            dispatch(detailActions.serviceInstanceSchema.applyFilter({
              platform: platform,
              clusterId: (_f = route_1 === null || route_1 === void 0 ? void 0 : route_1.queries) === null || _f === void 0 ? void 0 : _f.clusterid,
              serviceName: (_g = route_1 === null || route_1 === void 0 ? void 0 : route_1.queries) === null || _g === void 0 ? void 0 : _g.servicename,
              regionId: Util === null || Util === void 0 ? void 0 : Util.getDefaultRegion(platform, route_1)
            }));
            dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.namespaces.applyFilter({
              clusterId: (_h = route_1 === null || route_1 === void 0 ? void 0 : route_1.queries) === null || _h === void 0 ? void 0 : _h.clusterid,
              regionId: regionId,
              platform: platform
            }));
          }
          return [2 /*return*/];
        });
      });
    };
  },

  // 命名空间资源
  namespaces: ffRedux.createFFListActions({
    actionName: ActionType.Detail.FETCH_NAMESPACES,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var result;
        return tslib.__generator(this, function (_a) {
          switch (_a.label) {
            case 0:
              if (!(query === null || query === void 0 ? void 0 : query.filter)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchNamespaces(query === null || query === void 0 ? void 0 : query.filter)];
            case 1:
              result = _a.sent();
              _a.label = 2;
            case 2:
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().detail.namespaces;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b;
      if ((_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a[0]) {
        dispatch(listActions.services.select((_b = record === null || record === void 0 ? void 0 : record.data) === null || _b === void 0 ? void 0 : _b[0]));
      }
    },
    onSelect: function onSelect(record, dispatch, getState) {
      var route = getState().base.route;
    }
  }),
  updateServiceBinding: function updateServiceBinding(key, value) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b, _c;
        var _d, _e, _f;
        return tslib.__generator(this, function (_g) {
          if (key) {
            _a = (_e = (_d = getState()) === null || _d === void 0 ? void 0 : _d.detail) === null || _e === void 0 ? void 0 : _e.serviceBindingEdit, formData = _a.formData, validator = _a.validator;
            newEdit = {
              validator: tslib.__assign(tslib.__assign({}, validator), (_b = {}, _b[key] = ServiceBinding._validateServiceBindingItem(key, value), _b)),
              formData: tslib.__assign(tslib.__assign({}, formData), (_c = {}, _c[key] = value, _c))
            };
            dispatch({
              type: (_f = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _f === void 0 ? void 0 : _f.SERVICE_BINDING_EDIT,
              payload: newEdit
            });
          }
          return [2 /*return*/];
        });
      });
    };
  },

  initServiceBindingEdit: function initServiceBindingEdit(data) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, editData, newEdit;
        var _b, _c;
        return tslib.__generator(this, function (_d) {
          _a = (_c = (_b = getState()) === null || _b === void 0 ? void 0 : _b.detail) === null || _c === void 0 ? void 0 : _c.serviceBindingEdit, formData = _a.formData, validator = _a.validator;
          editData = ServiceBinding === null || ServiceBinding === void 0 ? void 0 : ServiceBinding.initServiceResourceEdit(data);
          newEdit = {
            validator: tslib.__assign(tslib.__assign({}, validator), editData === null || editData === void 0 ? void 0 : editData.validator),
            formData: tslib.__assign(tslib.__assign({}, formData), editData === null || editData === void 0 ? void 0 : editData.formData)
          };
          dispatch({
            type: ActionType.Detail.SERVICE_BINDING_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  serviceInstanceSchema: ffRedux.createFFObjectActions({
    actionName: ActionType.Detail.FETCH_SERVICE_INSTANCE_SCHEMA,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var filter, _a, platform, route, response;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              filter = query.filter;
              _a = (_b = getState()) === null || _b === void 0 ? void 0 : _b.base, platform = _a.platform, route = _a.route;
              if (!filter) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchResourceSchemas(filter)];
            case 1:
              response = _c.sent();
              _c.label = 2;
            case 2:
              return [2 /*return*/, response];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().detail.serviceInstanceSchema;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a;
      if (record === null || record === void 0 ? void 0 : record.data) {
        var selectedDetailTab = getState().detail.selectedDetailTab;
        if (selectedDetailTab === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding)) {
          dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.initServiceBindingEdit((_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.bindingCreateParameterSchema));
        }
      }
    }
  }),
  backupStrategy: ffRedux.createFFObjectActions({
    actionName: ActionType.Detail.FETCH_BACKUP_STRATEGY,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var filter, _a, platform, route, result, response, cronArr;
        var _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m;
        return tslib.__generator(this, function (_o) {
          switch (_o.label) {
            case 0:
              filter = query.filter;
              _a = (_b = getState()) === null || _b === void 0 ? void 0 : _b.base, platform = _a.platform, route = _a.route;
              return [4 /*yield*/, fetchBackStrategy(filter)];
            case 1:
              response = _o.sent();
              if (response) {
                cronArr = (_f = (_e = (_d = (_c = response === null || response === void 0 ? void 0 : response.spec) === null || _c === void 0 ? void 0 : _c.trigger) === null || _d === void 0 ? void 0 : _d.params) === null || _e === void 0 ? void 0 : _e.cron) === null || _f === void 0 ? void 0 : _f.split(" ");
                result = {
                  enable: (_g = response === null || response === void 0 ? void 0 : response.spec) === null || _g === void 0 ? void 0 : _g.enabled,
                  backupReserveDay: (_j = (_h = response === null || response === void 0 ? void 0 : response.spec) === null || _h === void 0 ? void 0 : _h.retain) === null || _j === void 0 ? void 0 : _j.days,
                  backupDate: (_k = cronArr === null || cronArr === void 0 ? void 0 : cronArr[4]) === null || _k === void 0 ? void 0 : _k.split(","),
                  backupTime: (_l = cronArr === null || cronArr === void 0 ? void 0 : cronArr[1]) === null || _l === void 0 ? void 0 : _l.split(","),
                  name: (_m = response === null || response === void 0 ? void 0 : response.metadata) === null || _m === void 0 ? void 0 : _m.name,
                  spec: response === null || response === void 0 ? void 0 : response.spec,
                  metadata: response === null || response === void 0 ? void 0 : response.metadata
                };
              }
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().detail.backupStrategy;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      if (record === null || record === void 0 ? void 0 : record.data) {
        dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.initBackUpStrategy(record === null || record === void 0 ? void 0 : record.data));
      }
    }
  }),
  validateAllServiceBinding: function validateAllServiceBinding() {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, _b, editData, serviceInstanceSchema, _c, services, serviceResources, route, _d, formData, validator, instance, newValidatorModel, newEdit;
        var _e, _f, _g, _h, _j, _k, _l, _m, _o;
        return tslib.__generator(this, function (_p) {
          switch (_p.label) {
            case 0:
              _a = getState(), _b = _a.detail, editData = _b.serviceBindingEdit, serviceInstanceSchema = _b.serviceInstanceSchema, _c = _a.list, services = _c.services, serviceResources = _c.serviceResources, route = _a.base.route;
              if (!editData) return [3 /*break*/, 2];
              _d = (_f = (_e = getState()) === null || _e === void 0 ? void 0 : _e.detail) === null || _f === void 0 ? void 0 : _f.serviceBindingEdit, formData = _d.formData, validator = _d.validator;
              instance = (_j = (_h = (_g = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.list) === null || _g === void 0 ? void 0 : _g.data) === null || _h === void 0 ? void 0 : _h.records) === null || _j === void 0 ? void 0 : _j.find(function (item) {
                var _a, _b;
                return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.instancename);
              });
              return [4 /*yield*/, ServiceBinding._validateFormItem(editData, (_l = (_k = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _k === void 0 ? void 0 : _k.data) === null || _l === void 0 ? void 0 : _l.bindingCreateParameterSchema, (_m = services === null || services === void 0 ? void 0 : services.selection) === null || _m === void 0 ? void 0 : _m.name, instance)];
            case 1:
              newValidatorModel = _p.sent();
              newEdit = {
                formData: formData,
                validator: tslib.__assign(tslib.__assign({}, validator), newValidatorModel)
              };
              dispatch({
                type: (_o = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Detail) === null || _o === void 0 ? void 0 : _o.SERVICE_BINDING_EDIT,
                payload: newEdit
              });
              _p.label = 2;
            case 2:
              return [2 /*return*/];
          }
        });
      });
    };
  }
};

var _a$7, _b$6, _c$4, _d$2;
var createActions = {
  updateInstanceEditSchema: function updateInstanceEditSchema(data) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, serviceInstanceEdit, newEdit;
        var _b, _c;
        return tslib.__generator(this, function (_d) {
          _a = (_c = (_b = getState()) === null || _b === void 0 ? void 0 : _b.list) === null || _c === void 0 ? void 0 : _c.serviceInstanceEdit, formData = _a.formData, validator = _a.validator;
          serviceInstanceEdit = initServiceInstanceEdit(data);
          newEdit = {
            validator: tslib.__assign(tslib.__assign({}, validator), serviceInstanceEdit === null || serviceInstanceEdit === void 0 ? void 0 : serviceInstanceEdit.validator),
            formData: tslib.__assign(tslib.__assign({}, formData), serviceInstanceEdit === null || serviceInstanceEdit === void 0 ? void 0 : serviceInstanceEdit.formData)
          };
          dispatch({
            type: ActionType.Create.SERVICE_INSTANCE_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  serviceInstanceEdit: function serviceInstanceEdit(data) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b, _c;
        return tslib.__generator(this, function (_d) {
          _a = (_c = (_b = getState()) === null || _b === void 0 ? void 0 : _b.list) === null || _c === void 0 ? void 0 : _c.serviceInstanceEdit, formData = _a.formData, validator = _a.validator;
          newEdit = {
            validator: tslib.__assign({}, validator),
            formData: tslib.__assign(tslib.__assign({}, formData), data === null || data === void 0 ? void 0 : data.formData)
          };
          dispatch({
            type: ActionType.Create.SERVICE_INSTANCE_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  updateInstance: function updateInstance(key, value) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b;
        var _c, _d;
        return tslib.__generator(this, function (_e) {
          if (key) {
            _a = (_d = (_c = getState()) === null || _c === void 0 ? void 0 : _c.list) === null || _d === void 0 ? void 0 : _d.serviceInstanceEdit, formData = _a.formData, validator = _a.validator;
            newEdit = {
              validator: validator,
              formData: tslib.__assign(tslib.__assign({}, formData), (_b = {}, _b[key] = value, _b))
            };
            dispatch({
              type: ActionType.Create.SERVICE_INSTANCE_EDIT,
              payload: newEdit
            });
          }
          return [2 /*return*/];
        });
      });
    };
  },

  validateTimeBackup: function validateTimeBackup(key, message) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b;
        var _c, _d;
        return tslib.__generator(this, function (_e) {
          if (key) {
            _a = (_d = (_c = getState()) === null || _c === void 0 ? void 0 : _c.list) === null || _d === void 0 ? void 0 : _d.serviceInstanceEdit, formData = _a.formData, validator = _a.validator;
            newEdit = {
              validator: tslib.__assign(tslib.__assign({}, validator), (_b = {}, _b[key] = {
                status: ffValidator.ValidatorStatusEnum === null || ffValidator.ValidatorStatusEnum === void 0 ? void 0 : ffValidator.ValidatorStatusEnum.Failed,
                message: message
              }, _b)),
              formData: formData
            };
            dispatch({
              type: ActionType.Create.SERVICE_INSTANCE_EDIT,
              payload: newEdit
            });
          }
          return [2 /*return*/];
        });
      });
    };
  },

  validateInstance: function validateInstance() {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, serviceInstanceEdit, servicesInstance, formData, instanceSchemas, vm;
        var _b, _c;
        return tslib.__generator(this, function (_d) {
          _a = getState().list, serviceInstanceEdit = _a.serviceInstanceEdit, servicesInstance = _a.servicesInstance, formData = serviceInstanceEdit.formData;
          instanceSchemas = (_c = (_b = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _b === void 0 ? void 0 : _b.data) === null || _c === void 0 ? void 0 : _c.instanceCreateParameterSchema;
          vm = validateAllProps(serviceInstanceEdit, instanceSchemas);
          dispatch({
            type: ActionType.Create.SERVICE_INSTANCE_EDIT,
            payload: tslib.__assign(tslib.__assign({}, serviceInstanceEdit), {
              validator: vm
            })
          });
          return [2 /*return*/];
        });
      });
    };
  },

  updateServicePlanSchema: function updateServicePlanSchema(data) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, editData, newEdit;
        var _b, _c;
        return tslib.__generator(this, function (_d) {
          _a = (_c = (_b = getState()) === null || _b === void 0 ? void 0 : _b.list) === null || _c === void 0 ? void 0 : _c.servicePlanEdit, formData = _a.formData, validator = _a.validator;
          editData = initServicePlanEdit(data);
          newEdit = {
            validator: tslib.__assign(tslib.__assign({}, validator), editData === null || editData === void 0 ? void 0 : editData.validator),
            formData: tslib.__assign(tslib.__assign({}, formData), editData === null || editData === void 0 ? void 0 : editData.formData)
          };
          dispatch({
            type: ActionType.Create.SERVICE_PLAN_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  servicePlanEdit: function servicePlanEdit(data) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b, _c;
        return tslib.__generator(this, function (_d) {
          _a = (_c = (_b = getState()) === null || _b === void 0 ? void 0 : _b.list) === null || _c === void 0 ? void 0 : _c.servicePlanEdit, formData = _a.formData, validator = _a.validator;
          newEdit = {
            validator: tslib.__assign({}, validator),
            formData: tslib.__assign(tslib.__assign({}, formData), data === null || data === void 0 ? void 0 : data.formData)
          };
          dispatch({
            type: ActionType.Create.SERVICE_PLAN_EDIT,
            payload: newEdit
          });
          return [2 /*return*/];
        });
      });
    };
  },

  updatePlan: function updatePlan(key, value) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, formData, validator, newEdit;
        var _b;
        var _c, _d;
        return tslib.__generator(this, function (_e) {
          if (key) {
            _a = (_d = (_c = getState()) === null || _c === void 0 ? void 0 : _c.list) === null || _d === void 0 ? void 0 : _d.servicePlanEdit, formData = _a.formData, validator = _a.validator;
            newEdit = {
              validator: validator,
              formData: tslib.__assign(tslib.__assign({}, formData), (_b = {}, _b[key] = value, _b))
            };
            dispatch({
              type: ActionType.Create.SERVICE_PLAN_EDIT,
              payload: newEdit
            });
          }
          return [2 /*return*/];
        });
      });
    };
  },

  validatePlan: function validatePlan(instanceSchemas) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var servicePlanEdit, formData, vm;
        return tslib.__generator(this, function (_a) {
          servicePlanEdit = getState().list.servicePlanEdit, formData = servicePlanEdit.formData;
          vm = validateAllPlanProps(servicePlanEdit, instanceSchemas);
          dispatch({
            type: ActionType.Create.SERVICE_PLAN_EDIT,
            payload: tslib.__assign(tslib.__assign({}, servicePlanEdit), {
              validator: vm
            })
          });
          return [2 /*return*/];
        });
      });
    };
  },

  createResource: ffRedux.generateWorkflowActionCreator({
    actionType: ActionType.Create.CREATE_SERVICE_RESOURCE,
    workflowStateLocator: function workflowStateLocator(state) {
      var _a;
      return (_a = state.list) === null || _a === void 0 ? void 0 : _a.createResourceWorkflow;
    },
    operationExecutor: createServiceResource,
    after: (_a$7 = {}, _a$7[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {
      var _a, _b, _c, _d, _e, _f;
      var _g = getState(),
        route = _g.base.route,
        createResourceWorkflow = _g.list.createResourceWorkflow;
      var specificOperate = ((_a = createResourceWorkflow === null || createResourceWorkflow === void 0 ? void 0 : createResourceWorkflow.targets) === null || _a === void 0 ? void 0 : _a[0]).specificOperate;
      if (ffRedux.isSuccessWorkflow(createResourceWorkflow)) {
        if ((_b = [CreateSpecificOperatorEnum.BackupNow, CreateSpecificOperatorEnum.BackupStrategy, CreateSpecificOperatorEnum === null || CreateSpecificOperatorEnum === void 0 ? void 0 : CreateSpecificOperatorEnum.CreateServiceBinding]) === null || _b === void 0 ? void 0 : _b.includes(specificOperate)) {
          dispatch((_c = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _c === void 0 ? void 0 : _c.fetch());
          dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.showBackupDialog(false));
          dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.showCreateResourceDialog(false));
        } else {
          router.navigate({
            sub: 'list',
            tab: undefined
          }, {
            servicename: (_d = route === null || route === void 0 ? void 0 : route.queries) === null || _d === void 0 ? void 0 : _d.servicename,
            resourceType: (_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.resourceType,
            mode: 'list'
          });
        }
        dispatch(createActions.createResource.reset());
      } else {
        specificOperate === CreateSpecificOperatorEnum.BackupNow && dispatch(createActions.createResource.reset());
      }
      dispatch({
        type: (_f = ActionType.Detail) === null || _f === void 0 ? void 0 : _f.BACKUP_RESOURCE_LOADING,
        payload: false
      });
    }, _a$7)
  }),
  backupNowWorkflow: ffRedux.generateWorkflowActionCreator({
    actionType: ActionType.Create.BackupNowWorkflow,
    workflowStateLocator: function workflowStateLocator(state) {
      var _a;
      return (_a = state.list) === null || _a === void 0 ? void 0 : _a.backupNowWorkflow;
    },
    operationExecutor: createServiceResource,
    after: (_b$6 = {}, _b$6[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {}, _b$6)
  }),
  updateResource: ffRedux.generateWorkflowActionCreator({
    actionType: ActionType.Create.CREATE_SERVICE_RESOURCE,
    workflowStateLocator: function workflowStateLocator(state) {
      var _a;
      return (_a = state.list) === null || _a === void 0 ? void 0 : _a.updateResourceWorkflow;
    },
    operationExecutor: updateServiceResource,
    after: (_c$4 = {}, _c$4[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {
      var _a, _b, _c;
      var updateResourceWorkflow = getState().list.updateResourceWorkflow;
      if (ffRedux.isSuccessWorkflow(updateResourceWorkflow)) {
        dispatch(createActions.updateResource.reset());
        var specificOperate = ((_a = updateResourceWorkflow === null || updateResourceWorkflow === void 0 ? void 0 : updateResourceWorkflow.targets) === null || _a === void 0 ? void 0 : _a[0]).specificOperate;
        if ((_b = [CreateSpecificOperatorEnum.BackupNow, CreateSpecificOperatorEnum.BackupStrategy, CreateSpecificOperatorEnum === null || CreateSpecificOperatorEnum === void 0 ? void 0 : CreateSpecificOperatorEnum.CreateServiceBinding]) === null || _b === void 0 ? void 0 : _b.includes(specificOperate)) {
          dispatch(detailActions === null || detailActions === void 0 ? void 0 : detailActions.showBackupDialog(false));
          dispatch((_c = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _c === void 0 ? void 0 : _c.fetch());
        } else {
          dispatch(listActions.showCreateResourceDialog(false));
          dispatch(listActions.serviceResources.fetch());
        }
      }
    }, _c$4)
  }),
  backupResourceNow: function backupResourceNow(params) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var hubCluster, platform, clusterId, regionId, cosParam, result;
        var _a, _b, _c;
        return tslib.__generator(this, function (_d) {
          switch (_d.label) {
            case 0:
              hubCluster = getState().base.hubCluster;
              platform = params.platform, clusterId = params.clusterId, regionId = params.regionId;
              if (!params) return [3 /*break*/, 2];
              cosParam = {
                platform: platform,
                clusterId: Util === null || Util === void 0 ? void 0 : Util.getCOSClusterId(platform, (_a = hubCluster === null || hubCluster === void 0 ? void 0 : hubCluster.object) === null || _a === void 0 ? void 0 : _a.data),
                regionId: regionId
              };
              dispatch({
                type: (_b = ActionType.Detail) === null || _b === void 0 ? void 0 : _b.BACKUP_RESOURCE_LOADING,
                payload: true
              });
              return [4 /*yield*/, checkCosResource(cosParam)];
            case 1:
              result = _d.sent();
              if (result) {
                dispatch(createActions.createResource.start([params], regionId));
                dispatch(createActions.createResource.perform());
              } else {
                dispatch({
                  type: (_c = ActionType.Detail) === null || _c === void 0 ? void 0 : _c.BACKUP_RESOURCE_LOADING,
                  payload: false
                });
                bridge.tips.error(ErrorMsgEnum.COS_Resource_Not_Found);
              }
              _d.label = 2;
            case 2:
              return [2 /*return*/];
          }
        });
      });
    };
  },

  createServiceInstance: ffRedux.generateWorkflowActionCreator({
    actionType: ActionType.Create.CREATE_SERVICE_INSTANCE,
    workflowStateLocator: function workflowStateLocator(state) {
      var _a;
      return (_a = state.list) === null || _a === void 0 ? void 0 : _a.createServiceInstanceWorkflow;
    },
    operationExecutor: function operationExecutor(targets, params, dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, instance, backupStrategy, result, error, response;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              _a = targets === null || targets === void 0 ? void 0 : targets[0], instance = _a.instance, backupStrategy = _a.backupStrategy;
              return [4 /*yield*/, Promise.all([createServiceResource(instance, params).then(function (response) {
                return response;
              }, function (error) {
                return error;
              }), createServiceResource(backupStrategy, params).then(function (response) {
                return response;
              }, function (error) {
                return error;
              })])];
            case 1:
              result = _c.sent();
              error = result === null || result === void 0 ? void 0 : result.find(function (item) {
                var _a;
                return (_a = item === null || item === void 0 ? void 0 : item[0]) === null || _a === void 0 ? void 0 : _a.error;
              });
              response = [];
              result === null || result === void 0 ? void 0 : result.forEach(function (item) {
                response.push(item === null || item === void 0 ? void 0 : item[0]);
              });
              if (!(error === null || error === void 0 ? void 0 : error.length)) {
                return [2 /*return*/, operationResult(targets, [], response)];
              } else {
                return [2 /*return*/, operationResult(targets, (_b = error === null || error === void 0 ? void 0 : error[0]) === null || _b === void 0 ? void 0 : _b.error)];
              }
          }
        });
      });
    },

    after: (_d$2 = {}, _d$2[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {
      var _a, _b, _c, _d, _e, _f;
      var _g = getState(),
        route = _g.base.route,
        createServiceInstanceWorkflow = _g.list.createServiceInstanceWorkflow;
      if ((_c = (_b = (_a = createServiceInstanceWorkflow === null || createServiceInstanceWorkflow === void 0 ? void 0 : createServiceInstanceWorkflow.results) === null || _a === void 0 ? void 0 : _a[0]) === null || _b === void 0 ? void 0 : _b.response) === null || _c === void 0 ? void 0 : _c.every(function (item) {
        return item === null || item === void 0 ? void 0 : item.success;
      })) {
        dispatch((_d = createActions === null || createActions === void 0 ? void 0 : createActions.createServiceInstance) === null || _d === void 0 ? void 0 : _d.reset());
        router.navigate({
          sub: 'list',
          tab: undefined
        }, {
          servicename: (_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.servicename,
          resourceType: (_f = route === null || route === void 0 ? void 0 : route.queries) === null || _f === void 0 ? void 0 : _f.resourceType,
          mode: 'list'
        });
      }
    }, _d$2)
  })
};

var _a$8;
var listActions = {
  resourceSchemas: ffRedux.createFFObjectActions({
    actionName: ActionType.List.FETCH_CREATE_RESOURCE_SCHEMAS,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, clusterId, regionId, serviceName, platform, response;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              _a = query.filter, clusterId = _a.clusterId, regionId = _a.regionId, serviceName = _a.serviceName;
              platform = ((_b = getState()) === null || _b === void 0 ? void 0 : _b.base).platform;
              if (!(platform && clusterId && regionId && serviceName)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchResourceSchemas({
                clusterId: clusterId,
                regionId: regionId,
                serviceName: serviceName,
                platform: platform
              })];
            case 1:
              response = _c.sent();
              _c.label = 2;
            case 2:
              return [2 /*return*/, response];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().list.servicesInstance;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b;
      if (record === null || record === void 0 ? void 0 : record.data) {
        var selectedTab = getState().base.selectedTab;
        if (selectedTab === ResourceTypeEnum.ServiceResource) {
          dispatch(createActions.updateInstanceEditSchema((_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.instanceCreateParameterSchema));
        } else if (selectedTab === ResourceTypeEnum.ServicePlan) {
          dispatch(createActions.updateServicePlanSchema((_b = record === null || record === void 0 ? void 0 : record.data) === null || _b === void 0 ? void 0 : _b.planSchema));
        }
      }
    }
  }),
  externalClusters: ffRedux.createFFListActions({
    actionName: ActionType.List.FETCH_EXTERNAL_CLUSTERS,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, regionId, _b, _c, regional, _d, platform, services, result;
        var _e, _f;
        return tslib.__generator(this, function (_g) {
          switch (_g.label) {
            case 0:
              _a = query.filter, regionId = _a.regionId, _b = _a.clusterIds, _c = _a.regional, regional = _c === void 0 ? false : _c;
              _d = getState(), platform = _d.base.platform, services = _d.list.services;
              if (!(platform && regionId)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchExternalClusters({
                platform: platform,
                regionId: regionId
              })];
            case 1:
              result = _g.sent();
              _g.label = 2;
            case 2:
              // 过滤状态为运行状态的注册集群
              result.records = (_e = result === null || result === void 0 ? void 0 : result.records) === null || _e === void 0 ? void 0 : _e.filter(function (item) {
                var _a;
                return (item === null || item === void 0 ? void 0 : item.status) === ((_a = ExternalCluster === null || ExternalCluster === void 0 ? void 0 : ExternalCluster.StatusEnum) === null || _a === void 0 ? void 0 : _a.Running);
              });
              // // 集群排序
              // result.records = result.records?.sort((a,b) => {
              //   if(services.selection?.clusters?.includes(a?.clusterId)){
              //     return -1
              //   }else if(services.selection?.clusters?.includes(b?.clusterId)){
              //     return 1
              //   }else{
              //     return 0
              //   }
              // })
              result.recordCount = (_f = result === null || result === void 0 ? void 0 : result.records) === null || _f === void 0 ? void 0 : _f.length;
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().list.externalClusters;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b;
      var services = getState().list.services;
      var openVendorClusters = (_b = (_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.records) === null || _b === void 0 ? void 0 : _b.filter(function (item) {
        var _a, _b;
        return (_b = (_a = services.selection) === null || _a === void 0 ? void 0 : _a.clusters) === null || _b === void 0 ? void 0 : _b.includes(item === null || item === void 0 ? void 0 : item.clusterId);
      });
      if (!!(openVendorClusters === null || openVendorClusters === void 0 ? void 0 : openVendorClusters.length)) {
        dispatch(listActions.externalClusters.select(openVendorClusters === null || openVendorClusters === void 0 ? void 0 : openVendorClusters[0]));
      }
    },
    onSelect: function onSelect(record, dispatch, getState) {
      if (record) {
        var _a = getState().base,
          selectedTab = _a.selectedTab,
          route = _a.route,
          platform = _a.platform;
        var sub = (router === null || router === void 0 ? void 0 : router.resolve(route)).sub;
        var mode = (route === null || route === void 0 ? void 0 : route.queries).mode;
        if (sub === 'create' || mode === 'create') {
          if (selectedTab === ResourceTypeEnum.ServiceResource) {
            dispatch(createActions.updateInstance('clusterId', record === null || record === void 0 ? void 0 : record.clusterId));
          } else if (selectedTab === ResourceTypeEnum.ServicePlan) {
            dispatch(createActions.updatePlan('clusterId', record === null || record === void 0 ? void 0 : record.clusterId));
          }
        }
      }
    }
  }),
  serviceResources: ffRedux.createFFListActions({
    actionName: ActionType.List.FETCH_SERVICE_RESOURCE,
    fetcher: function fetcher(query, getState, fetchOptions, dispatch) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, clusterId, regionId, serviceName, _b, resourceType, search, paging, platform, result;
        var _c;
        return tslib.__generator(this, function (_d) {
          switch (_d.label) {
            case 0:
              _a = query.filter, clusterId = _a.clusterId, regionId = _a.regionId, serviceName = _a.serviceName, _b = _a.resourceType, resourceType = _b === void 0 ? ResourceTypeEnum.ServiceResource : _b, search = query.search, paging = query.paging;
              platform = ((_c = getState()) === null || _c === void 0 ? void 0 : _c.base).platform;
              result = {
                recordCount: 0,
                records: []
              };
              if (!(platform && clusterId && serviceName && resourceType)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchServiceResources(tslib.__assign(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter), {
                search: search,
                paging: paging
              }))];
            case 1:
              result = _d.sent();
              _d.label = 2;
            case 2:
              result.recordCount = result === null || result === void 0 ? void 0 : result.recordCount;
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().list.serviceResources;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t;
      var _u = getState(),
        route = _u.base.route,
        _v = _u.list,
        serviceResources = _v.serviceResources,
        services = _v.services;
      var sub = (router === null || router === void 0 ? void 0 : router.resolve(route)).sub;
      if ((_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.recordCount) {
        var instancename_1 = (route === null || route === void 0 ? void 0 : route.queries).instancename;
        var resource = void 0;
        if (instancename_1) {
          resource = (_c = (_b = record === null || record === void 0 ? void 0 : record.data) === null || _b === void 0 ? void 0 : _b.records) === null || _c === void 0 ? void 0 : _c.find(function (item) {
            var _a;
            return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === instancename_1;
          });
        } else {
          resource = (_e = (_d = record === null || record === void 0 ? void 0 : record.data) === null || _d === void 0 ? void 0 : _d.records) === null || _e === void 0 ? void 0 : _e[0];
        }
        if (resource) {
          dispatch((_f = listActions === null || listActions === void 0 ? void 0 : listActions.serviceResources) === null || _f === void 0 ? void 0 : _f.select(resource));
        }
        // 拉取规格列表关联的实例资源/拉取实例列表中关联的规格资源
        var filter = tslib.__assign(tslib.__assign({}, (_g = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.query) === null || _g === void 0 ? void 0 : _g.filter), {
          resourceType: ((_j = (_h = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.query) === null || _h === void 0 ? void 0 : _h.filter) === null || _j === void 0 ? void 0 : _j.resourceType) === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan) ? ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource : ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan
        });
        dispatch((_k = listActions === null || listActions === void 0 ? void 0 : listActions.serviceResourceList) === null || _k === void 0 ? void 0 : _k.clear());
        dispatch((_l = listActions === null || listActions === void 0 ? void 0 : listActions.serviceResourceList) === null || _l === void 0 ? void 0 : _l.applyFilter(filter));
      }
      //判断是否开启轮训
      if (((_m = ['', 'list']) === null || _m === void 0 ? void 0 : _m.includes(sub)) && ((_p = (_o = record === null || record === void 0 ? void 0 : record.data) === null || _o === void 0 ? void 0 : _o.records) === null || _p === void 0 ? void 0 : _p.some(function (item) {
        var _a;
        return (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp;
      })) && ((_r = (_q = record === null || record === void 0 ? void 0 : record.data) === null || _q === void 0 ? void 0 : _q.records) === null || _r === void 0 ? void 0 : _r.every(function (item) {
        var _a, _b;
        return ((_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.serviceClass) === ((_b = services === null || services === void 0 ? void 0 : services.selection) === null || _b === void 0 ? void 0 : _b.name);
      }))) {
        dispatch((_s = listActions === null || listActions === void 0 ? void 0 : listActions.serviceResources) === null || _s === void 0 ? void 0 : _s.polling({
          delayTime: 8000
        }));
      } else {
        dispatch((_t = listActions === null || listActions === void 0 ? void 0 : listActions.serviceResources) === null || _t === void 0 ? void 0 : _t.clearPolling());
      }
    }
  }),
  services: ffRedux.createFFListActions({
    actionName: ActionType.List.FETCH_SERVICES,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, clusterId, regionId, platform, result;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              _a = query.filter, clusterId = _a.clusterId, regionId = _a.regionId;
              platform = ((_b = getState()) === null || _b === void 0 ? void 0 : _b.base).platform;
              if (!(platform && clusterId)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchOpenedServices({
                platform: platform,
                clusterId: clusterId,
                regionId: regionId
              })];
            case 1:
              result = _c.sent();
              _c.label = 2;
            case 2:
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().list.services;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b, _c, _d, _e, _f, _g;
      if ((_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.recordCount) {
        var route = getState().base.route;
        var service = void 0;
        var servicename_1 = ((_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.servicename) || ((_c = parseQueryString(location === null || location === void 0 ? void 0 : location.search)) === null || _c === void 0 ? void 0 : _c['servicename']);
        if (servicename_1 && servicename_1 !== 'undefined') {
          service = (_e = (_d = record === null || record === void 0 ? void 0 : record.data) === null || _d === void 0 ? void 0 : _d.records) === null || _e === void 0 ? void 0 : _e.find(function (item) {
            return (item === null || item === void 0 ? void 0 : item.name) === servicename_1;
          });
        } else {
          service = (_g = (_f = record === null || record === void 0 ? void 0 : record.data) === null || _f === void 0 ? void 0 : _f.records) === null || _g === void 0 ? void 0 : _g[0];
        }
        dispatch(listActions.services.select(service !== null && service !== void 0 ? service : {}));
      }
    },
    onSelect: function onSelect(record, dispatch, getState) {
      var _a, _b, _c;
      var _d = getState(),
        route = _d.base.route,
        services = _d.list.services;
      var resourceType = (route === null || route === void 0 ? void 0 : route.queries).resourceType;
      var selectedTab;
      if (resourceType) {
        selectedTab = resourceType;
      } else {
        selectedTab = (_a = serviceMngTabs === null || serviceMngTabs === void 0 ? void 0 : serviceMngTabs[0]) === null || _a === void 0 ? void 0 : _a.id;
      }
      // 重置服务资源列表
      dispatch(listActions.serviceResources.clear());
      dispatch(listActions.serviceResources.clearPolling());
      dispatch(baseActions.selectTab(selectedTab));
      var _e = (_b = services === null || services === void 0 ? void 0 : services.query) === null || _b === void 0 ? void 0 : _b.filter,
        regionId = _e.regionId,
        platform = _e.platform;
      var clusterId = (_c = route === null || route === void 0 ? void 0 : route.queries) === null || _c === void 0 ? void 0 : _c.clusterid;
      // 查询resourceSchema配置
      if (platform && record && regionId && clusterId) {
        dispatch(listActions === null || listActions === void 0 ? void 0 : listActions.resourceSchemas.applyFilter({
          platform: platform,
          serviceName: record === null || record === void 0 ? void 0 : record.name,
          clusterId: clusterId,
          regionId: regionId
        }));
      }
    }
  }),
  servicePlans: ffRedux.createFFListActions({
    actionName: ActionType.List.FETCH_SERVICE_PLANS,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, clusterId, regionId, serviceName, result;
        var _b, _c, _d;
        return tslib.__generator(this, function (_e) {
          switch (_e.label) {
            case 0:
              _a = query.filter, clusterId = _a.clusterId, regionId = _a.regionId, serviceName = _a.serviceName;
              if (!(query === null || query === void 0 ? void 0 : query.filter)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchServiceResources(query === null || query === void 0 ? void 0 : query.filter)];
            case 1:
              result = _e.sent();
              _e.label = 2;
            case 2:
              // 过滤掉为空的规格数据
              result.records = (_b = result === null || result === void 0 ? void 0 : result.records) === null || _b === void 0 ? void 0 : _b.filter(function (item) {
                var _a;
                return !isEmpty((_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.metadata);
              });
              // 按照创建时间降序排列
              result.records = (_c = result.records) === null || _c === void 0 ? void 0 : _c.sort(function (pre, cur) {
                var _a, _b, _c, _d;
                return ((_b = new Date((_a = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp)) === null || _b === void 0 ? void 0 : _b.getTime()) - ((_d = new Date((_c = pre === null || pre === void 0 ? void 0 : pre.metadata) === null || _c === void 0 ? void 0 : _c.creationTimestamp)) === null || _d === void 0 ? void 0 : _d.getTime());
              });
              result.recordCount = (_d = result === null || result === void 0 ? void 0 : result.records) === null || _d === void 0 ? void 0 : _d.length;
              return [2 /*return*/, result];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().list.servicePlans;
    },
    onFinish: function onFinish(record, dispatch) {
      var _a, _b, _c;
      if ((_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.recordCount) {
        dispatch(listActions.servicePlans.select((_c = (_b = record === null || record === void 0 ? void 0 : record.data) === null || _b === void 0 ? void 0 : _b.records) === null || _c === void 0 ? void 0 : _c[0]));
      }
    },
    onSelect: function onSelect(record, dispatch, get) {
      var _a;
      if (record) {
        dispatch(createActions.updateInstance('plan', (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name));
      }
    }
  }),
  /**
   * 实例资源列表
   */
  serviceResourceList: ffRedux.createFFListActions({
    actionName: ActionType.List.FETCH_SERVICE_RESOURCE_LIST,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var search, paging, serviceResources, resources;
        var _a, _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              search = query.search, paging = query.paging;
              serviceResources = getState().list.serviceResources;
              if (!(query === null || query === void 0 ? void 0 : query.filter)) return [3 /*break*/, 2];
              return [4 /*yield*/, fetchServiceResources(tslib.__assign(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter), {
                search: search,
                paging: paging
              }))];
            case 1:
              resources = _c.sent();
              _c.label = 2;
            case 2:
              // 按照创建时间降序排列
              resources.records = (_a = resources.records) === null || _a === void 0 ? void 0 : _a.sort(function (pre, cur) {
                var _a, _b, _c, _d;
                return ((_b = new Date((_a = cur === null || cur === void 0 ? void 0 : cur.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp)) === null || _b === void 0 ? void 0 : _b.getTime()) - ((_d = new Date((_c = pre === null || pre === void 0 ? void 0 : pre.metadata) === null || _c === void 0 ? void 0 : _c.creationTimestamp)) === null || _d === void 0 ? void 0 : _d.getTime());
              });
              resources.recordCount = (_b = resources === null || resources === void 0 ? void 0 : resources.records) === null || _b === void 0 ? void 0 : _b.length;
              return [2 /*return*/, resources];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().list.serviceResourceList;
    }
  }),
  selectDeleteResources: function selectDeleteResources(resources) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        return tslib.__generator(this, function (_a) {
          dispatch({
            type: ActionType.List.SELECT_DELETE_RESOURCE,
            payload: resources
          });
          return [2 /*return*/];
        });
      });
    };
  },

  showInstanceDialog: function showInstanceDialog(isShowInstanceDialog) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a;
        return tslib.__generator(this, function (_b) {
          if (!isShowInstanceDialog) {
            dispatch((_a = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _a === void 0 ? void 0 : _a.clear());
          }
          dispatch({
            type: ActionType.List.SHOW_INSTANCE_TABLE_DIALOG,
            payload: isShowInstanceDialog
          });
          return [2 /*return*/];
        });
      });
    };
  },

  showCreateResourceDialog: function showCreateResourceDialog(isShowDialog) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        return tslib.__generator(this, function (_a) {
          dispatch({
            type: ActionType.List.SHOW_CREATE_RESOURCE_DIALOG,
            payload: isShowDialog
          });
          return [2 /*return*/];
        });
      });
    };
  },

  deleteResource: ffRedux.generateWorkflowActionCreator({
    actionType: ActionType.List.DELETE_SERVICE_RESOURCE,
    workflowStateLocator: function workflowStateLocator(state) {
      return state.list.deleteResourceWorkflow;
    },
    operationExecutor: deleteMulServiceResource,
    after: (_a$8 = {}, _a$8[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {
      var _a, _b, _c, _d, _e, _f;
      var deleteResourceWorkflow = getState().list.deleteResourceWorkflow;
      var resourceInfos = ((_a = deleteResourceWorkflow === null || deleteResourceWorkflow === void 0 ? void 0 : deleteResourceWorkflow.targets) === null || _a === void 0 ? void 0 : _a[0]).resourceInfos;
      if (ffRedux.isSuccessWorkflow(deleteResourceWorkflow)) {
        dispatch(listActions.deleteResource.reset());
        dispatch(listActions.selectDeleteResources([]));
        if ((_b = [ResourceTypeEnum.ServiceResource, ResourceTypeEnum.ServicePlan]) === null || _b === void 0 ? void 0 : _b.includes((_c = resourceInfos === null || resourceInfos === void 0 ? void 0 : resourceInfos[0]) === null || _c === void 0 ? void 0 : _c.kind)) {
          dispatch((_d = listActions === null || listActions === void 0 ? void 0 : listActions.serviceResources) === null || _d === void 0 ? void 0 : _d.fetch());
        } else {
          dispatch((_e = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _e === void 0 ? void 0 : _e.fetch());
          dispatch((_f = detailActions === null || detailActions === void 0 ? void 0 : detailActions.instanceResource) === null || _f === void 0 ? void 0 : _f.selects([]));
        }
      }
    }, _a$8)
  })
};

var routerSea$3 = seajs.require('router');
var baseActions = {
  /** 国际版 */
  toggleIsI18n: function toggleIsI18n(isI18n) {
    return {
      type: ActionType.Base.IsI18n,
      payload: isI18n
    };
  },
  fetchRegion: function fetchRegion(regionId) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, route, platform;
        var _b;
        return tslib.__generator(this, function (_c) {
          dispatch({
            type: ActionType.Base.FETCH_REGION,
            payload: regionId
          });
          _a = getState().base, route = _a.route, platform = _a.platform;
          // tdcc首先拉取hub集群,然后拉取已开通的服务列表;tkeStack同时拉取目标集群和拉取已开通的服务列表
          if (platform === (PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TDCC)) {
            dispatch(baseActions.hubCluster.applyFilter({
              regionId: regionId,
              platform: platform
            }));
          } else if (platform === (PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TKESTACK)) {
            dispatch(listActions.services.applyFilter({
              platform: platform,
              clusterId: ExternalCluster.TKEStackDefaultCluster,
              regionId: regionId
            }));
            dispatch(listActions === null || listActions === void 0 ? void 0 : listActions.externalClusters.applyFilter({
              platform: platform,
              clusterIds: [],
              regionId: regionId
            }));
          }
          //拉取用户信息
          dispatch((_b = baseActions.userInfo) === null || _b === void 0 ? void 0 : _b.applyFilter({
            platform: platform
          }));
          return [2 /*return*/];
        });
      });
    };
  },

  fetchPlatform: function fetchPlatform(platform, region) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var route, regionId;
        return tslib.__generator(this, function (_a) {
          dispatch({
            type: ActionType.Base.FETCH_PLATFORM,
            payload: platform
          });
          window['platform'] = platform;
          route = getState().base.route;
          regionId = region !== null && region !== void 0 ? region : HubCluster.DefaultRegion;
          dispatch(baseActions === null || baseActions === void 0 ? void 0 : baseActions.fetchRegion(regionId));
          return [2 /*return*/];
        });
      });
    };
  },

  hubCluster: ffRedux.createFFObjectActions({
    actionName: ActionType.Base.HubCluster,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, platform, route, response, data;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              _a = getState().base, platform = _a.platform, route = _a.route;
              return [4 /*yield*/, fetchHubCluster(query === null || query === void 0 ? void 0 : query.filter)];
            case 1:
              response = _c.sent();
              data = (_b = response === null || response === void 0 ? void 0 : response.records) === null || _b === void 0 ? void 0 : _b[0];
              if (data && !(data === null || data === void 0 ? void 0 : data.regionId)) {
                data.regionId = HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
              }
              return [2 /*return*/, data];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().base.hubCluster;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b, _c, _d, _e;
      var _f = getState().base,
        route = _f.route,
        platform = _f.platform,
        regionId = _f.regionId;
      dispatch({
        type: ActionType.Base.ClusterVersion,
        payload: '1.18.4'
      });
      if (!record.data) {
        if (platform === PlatformType.TDCC) {
          routerSea$3.navigate('/tdcc/paasoverview/startup');
        }
      } else {
        // 更新region
        dispatch({
          payload: (_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.regionId,
          type: (_b = ActionType === null || ActionType === void 0 ? void 0 : ActionType.Base) === null || _b === void 0 ? void 0 : _b.FETCH_REGION
        });
        dispatch(listActions.services.applyFilter({
          platform: platform,
          clusterId: (_c = record.data) === null || _c === void 0 ? void 0 : _c.clusterId,
          regionId: (_d = record === null || record === void 0 ? void 0 : record.data) === null || _d === void 0 ? void 0 : _d.regionId
        }));
        dispatch(listActions === null || listActions === void 0 ? void 0 : listActions.externalClusters.applyFilter({
          platform: platform,
          clusterIds: [],
          regionId: (_e = record === null || record === void 0 ? void 0 : record.data) === null || _e === void 0 ? void 0 : _e.regionId
        }));
      }
    }
  }),
  selectTab: function selectTab(tabId) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, _b, platform, route, hubCluster, regionId, _c, serviceResources, services, filter;
        var _d, _e;
        return tslib.__generator(this, function (_f) {
          dispatch({
            type: ActionType.Base.SELECT_TAB,
            payload: tabId !== null && tabId !== void 0 ? tabId : ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource
          });
          _a = getState(), _b = _a.base, platform = _b.platform, route = _b.route, hubCluster = _b.hubCluster, regionId = _b.regionId, _c = _a.list, serviceResources = _c.serviceResources, services = _c.services;
          filter = {
            platform: platform,
            clusterId: Util.getCOSClusterId(platform, (_d = hubCluster === null || hubCluster === void 0 ? void 0 : hubCluster.object) === null || _d === void 0 ? void 0 : _d.data),
            serviceName: (_e = services === null || services === void 0 ? void 0 : services.selection) === null || _e === void 0 ? void 0 : _e.name,
            resourceType: tabId,
            regionId: regionId
          };
          // 重新查询服务资源列表
          dispatch(listActions.serviceResources.clear());
          dispatch(listActions.serviceResources.applyFilter(filter));
          return [2 /*return*/];
        });
      });
    };
  },

  userInfo: ffRedux.createFFObjectActions({
    actionName: ActionType.Base.FETCH_UserInfo,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var response;
        var _a;
        return tslib.__generator(this, function (_b) {
          switch (_b.label) {
            case 0:
              if (!(((_a = query === null || query === void 0 ? void 0 : query.filter) === null || _a === void 0 ? void 0 : _a.platform) === PlatformType.TDCC)) return [3 /*break*/, 1];
              response = {
                name: Util === null || Util === void 0 ? void 0 : Util.getUserName()
              };
              return [3 /*break*/, 3];
            case 1:
              return [4 /*yield*/, fetchUserInfo(query === null || query === void 0 ? void 0 : query.filter)];
            case 2:
              response = _b.sent();
              _b.label = 3;
            case 3:
              return [2 /*return*/, response !== null && response !== void 0 ? response : {}];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().base.userInfo;
    }
  }),
  // 获取集群Admin权限
  getClusterAdminRole: GetRbacAdminDialog.createActions({
    pageName: ActionType.Base.GetClusterAdminRoleFlow,
    getRecord: function getRecord(getState) {
      return getState().base.getClusterAdminRole;
    }
  })
  // generateWorkflowActionCreator<RBACResource, number>({
  //   actionType: ActionType.Base.GetClusterAdminRoleFlow,
  //   workflowStateLocator: state => state.base.getClusterAdminRoleFlow,
  //   operationExecutor: getClusterAdminRole,
  //   after: {
  //     [OperationTrigger.Done]: (dispatch, getState: GetState) => {
  //       let {
  //         base: {
  //           getClusterAdminRoleFlow
  //         }
  //       } = getState();
  //       if(isSuccessWorkflow(getClusterAdminRoleFlow)){
  //         dispatch(baseActions.getClusterAdminRole.reset());
  //         dispatch(listActions.servicePlans.fetch());
  //         dispatch(listActions.resourceSchemas.fetch());
  //       }
  //     }
  //   }
  // })
};

var allActions = {
  base: baseActions,
  create: createActions,
  list: listActions,
  detail: detailActions
};

/**
 * 重置redux store，用于离开页面时清空状态
 */
var ResetStoreAction = 'ResetStore';
/**
 * 生成可重置的reducer，用于rootReducer简单包装
 * @return 可重置的reducer，当接收到 ResetStoreAction 时重置之
 */
var generateResetableReducer = function generateResetableReducer(rootReducer) {
  return function (state, action) {
    var newState = state;
    // 销毁页面
    if (action.type === ResetStoreAction) {
      newState = undefined;
    }
    return rootReducer(newState, action);
  };
};

var TempReducer = redux.combineReducers({
  route: router.getReducer(),
  isI18n: ffRedux.reduceToPayload(ActionType.Base.IsI18n, false),
  platform: ffRedux.reduceToPayload(ActionType.Base.FETCH_PLATFORM, PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TDCC),
  regionId: ffRedux.reduceToPayload(ActionType.Base.FETCH_REGION, HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion),
  hubCluster: ffRedux.createFFObjectReducer({
    actionName: ActionType.Base.HubCluster
  }),
  selectedTab: ffRedux.reduceToPayload(ActionType.Base.SELECT_TAB, ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource),
  clusterVersion: ffRedux.reduceToPayload(ActionType.Base.ClusterVersion, ''),
  userInfo: ffRedux.createFFObjectReducer({
    actionName: ActionType.Base.FETCH_UserInfo
  }),
  getClusterAdminRole: GetRbacAdminDialog.createReducer({
    pageName: ActionType.Base.GetClusterAdminRoleFlow
  })
});
var baseReducer = function baseReducer(state, action) {
  var newState = state;
  // 销毁页面
  if (action.type === ActionType.Base.Clear) {
    newState = undefined;
  }
  return TempReducer(newState, action);
};

var TempReducer$1 = redux.combineReducers({
  resourceDetail: ffRedux.createFFObjectReducer({
    actionName: ActionType.Detail.RESOURCE_DETAIL
  }),
  selectedDetailTab: ffRedux.reduceToPayload(ActionType.Detail.Select_Detail_Tab, DetailTabType === null || DetailTabType === void 0 ? void 0 : DetailTabType.Detail),
  instanceResource: ffRedux.createFFListReducer({
    actionName: ActionType.Detail.FETCH_INSTANCE_RESOURCE
  }),
  /**检测COS存储是否配置 */
  checkCosResource: ffRedux.createFFObjectReducer({
    actionName: ActionType.Detail.CHECK_COS
  }),
  backupResourceLoading: ffRedux.reduceToPayload(ActionType.Detail.BACKUP_RESOURCE_LOADING, false),
  showBackupStrategyDialog: ffRedux.reduceToPayload(ActionType.Detail.SHOW_BACKUP_STRATEGY_DIALOG, false),
  showCreateResourceDialog: ffRedux.reduceToPayload(ActionType.Detail.SHOW_CREATE_RESOURCE_DIALOG, false),
  selectDetailResource: ffRedux.reduceToPayload(ActionType.Detail.SELECT_DETAIL_RESOURCE, []),
  backupStrategyEdit: ffRedux.reduceToPayload(ActionType.Detail.BACKUP_STRATEGY_EDIT, tslib.__assign({}, BackupStrategyEditInitValue)),
  /** 开启/关闭控制台 */
  openConsoleWorkflow: ffRedux.generateWorkflowReducer({
    actionType: ActionType.Detail.OPEN_CONSOLE_WORKFLOW
  }),
  mediums: MediumSelectPanel.createReducer({
    pageName: ActionType.Detail.BackupStrategyMedium
  }),
  namespaces: ffRedux.createFFListReducer({
    actionName: ActionType.Detail.FETCH_NAMESPACES
  }),
  serviceBindingEdit: ffRedux.reduceToPayload(ActionType.Detail.SERVICE_BINDING_EDIT, tslib.__assign({}, ServiceBinding === null || ServiceBinding === void 0 ? void 0 : ServiceBinding.ServiceBindingEditInitValue)),
  serviceInstanceSchema: ffRedux.createFFObjectReducer({
    actionName: ActionType.Detail.FETCH_SERVICE_INSTANCE_SCHEMA
  }),
  backupStrategy: ffRedux.createFFObjectReducer({
    actionName: ActionType.Detail.FETCH_BACKUP_STRATEGY
  })
});
var detailReducer = function detailReducer(state, action) {
  var newState = state;
  // 销毁页面
  if (action.type === ActionType.List.Clear) {
    newState = undefined;
  }
  return TempReducer$1(newState, action);
};

var TempReducer$2 = redux.combineReducers({
  services: ffRedux.createFFListReducer({
    actionName: ActionType.List.FETCH_SERVICES
  }),
  servicesInstance: ffRedux.createFFObjectReducer({
    actionName: ActionType.List.FETCH_CREATE_RESOURCE_SCHEMAS
  }),
  serviceResources: ffRedux.createFFListReducer({
    actionName: ActionType.List.FETCH_SERVICE_RESOURCE
  }),
  servicePlans: ffRedux.createFFListReducer({
    actionName: ActionType.List.FETCH_SERVICE_PLANS
  }),
  externalClusters: ffRedux.createFFListReducer({
    actionName: ActionType.List.FETCH_EXTERNAL_CLUSTERS
  }),
  backupNowWorkflow: ffRedux.generateWorkflowReducer({
    actionType: ActionType.Create.BackupNowWorkflow
  }),
  serviceInstanceEdit: ffRedux.reduceToPayload(ActionType.Create.SERVICE_INSTANCE_EDIT, {
    formData: {
      clusterId: "",
      instanceName: "",
      plan: "",
      timeBackup: false,
      isSetParams: false,
      backupDate: [],
      backupTime: [],
      backupReserveDay: 30
    }
  }),
  servicePlanEdit: ffRedux.reduceToPayload(ActionType.Create.SERVICE_PLAN_EDIT, {
    formData: {
      clusterId: "",
      instanceName: "",
      description: ""
    }
  }),
  /** 新建资源工作流 */
  createResourceWorkflow: ffRedux.generateWorkflowReducer({
    actionType: ActionType.Create.CREATE_SERVICE_RESOURCE
  }),
  /** 编辑资源工作流 */
  updateResourceWorkflow: ffRedux.generateWorkflowReducer({
    actionType: ActionType.Create.CREATE_SERVICE_RESOURCE
  }),
  /** 创建备份策略工作流 */
  createServiceInstanceWorkflow: ffRedux.generateWorkflowReducer({
    actionType: ActionType.Create.CREATE_SERVICE_INSTANCE
  }),
  /** 已关联的资源列表(实例或者规格) */
  serviceResourceList: ffRedux.createFFListReducer({
    actionName: ActionType.List.FETCH_SERVICE_RESOURCE_LIST
  }),
  deleteResourceSelection: ffRedux.reduceToPayload(ActionType.List.SELECT_DELETE_RESOURCE, []),
  showInstanceTableDialog: ffRedux.reduceToPayload(ActionType.List.SHOW_INSTANCE_TABLE_DIALOG, false),
  showCreateResourceDialog: ffRedux.reduceToPayload(ActionType.List.SHOW_CREATE_RESOURCE_DIALOG, false),
  /** 删除资源工作流 */
  deleteResourceWorkflow: ffRedux.generateWorkflowReducer({
    actionType: ActionType.List.DELETE_SERVICE_RESOURCE
  })
});
var listReducer = function listReducer(state, action) {
  var newState = state;
  // 销毁页面
  if (action.type === ActionType.List.Clear) {
    newState = undefined;
  }
  return TempReducer$2(newState, action);
};

var RootReducer = redux.combineReducers({
  base: baseReducer,
  list: listReducer,
  detail: detailReducer
});

var createStore = process.env.NODE_ENV === 'development' ? redux.applyMiddleware(thunk, reduxLogger.createLogger({
  collapsed: true,
  diff: true
}))(redux.createStore) : redux.applyMiddleware(thunk)(redux.createStore);

function configStore() {
  var store = createStore(generateResetableReducer(RootReducer));
  // hot reloading
  // if (typeof module !== 'undefined' && (module as any).hot) {
  //   (module as any).hot.accept('../reducers/RootReducer', () => {
  //     store.replaceReducer(generateResetableReducer(require('../reducers/RootReducer').RootReducer));
  //   });
  // }
  return store;
}

var ErrorEnum;
(function (ErrorEnum) {
  var Code;
  (function (Code) {
    Code["RBACForbidden"] = "FailedOperation.RBACForbidden";
    Code[Code["RBACForbidden403"] = 403] = "RBACForbidden403";
  })(Code = ErrorEnum.Code || (ErrorEnum.Code = {}));
})(ErrorEnum || (ErrorEnum = {}));

var _a$9, _b$7;
/** 节点亲和性调度 亲和性调度操作符 */
var NodeAffinityOperatorEnum;
(function (NodeAffinityOperatorEnum) {
  NodeAffinityOperatorEnum["In"] = "In";
  NodeAffinityOperatorEnum["NotIn"] = "NotIn";
  NodeAffinityOperatorEnum["Exists"] = "Exists";
  NodeAffinityOperatorEnum["DoesNotExist"] = "DoesNotExist";
  NodeAffinityOperatorEnum["Gt"] = "Gt";
  NodeAffinityOperatorEnum["Lt"] = "Lt";
})(NodeAffinityOperatorEnum || (NodeAffinityOperatorEnum = {}));
/**
 * 服务调度的操作符
 */
var NodeAffinityRuleOperatorList = [{
  value: NodeAffinityOperatorEnum.In,
  tip: i18n.t('Label的value在列表中')
}, {
  value: NodeAffinityOperatorEnum.NotIn,
  tip: i18n.t('Label的value不在列表中')
}, {
  value: NodeAffinityOperatorEnum.Exists,
  tip: i18n.t('Label的key存在')
}, {
  value: NodeAffinityOperatorEnum.DoesNotExist,
  tip: i18n.t('Labe的key不存在')
}, {
  value: NodeAffinityOperatorEnum.Gt,
  tip: i18n.t('Label的值大于列表值（字符串匹配）')
}, {
  value: NodeAffinityOperatorEnum.Lt,
  tip: i18n.t('Label的值小于列表值（字符串匹配）')
}];
var initRecords$1 = [{
  id: uuid(),
  key: '',
  value: '',
  operator: NodeAffinityOperatorEnum.In,
  v_key: {
    status: ffValidator.ValidatorStatusEnum.Success,
    message: ''
  },
  v_value: {
    status: ffValidator.ValidatorStatusEnum.Success,
    message: ''
  }
}];
var NodeAffinityMustNeedKey = 'topology.loopdevice.csi.infra.tce.io/hostname';
var NoeAffinityInitRecords = [{
  id: uuid(),
  key: NodeAffinityMustNeedKey,
  value: '',
  operator: NodeAffinityOperatorEnum.Exists,
  v_key: (_a$9 = initRecords$1 === null || initRecords$1 === void 0 ? void 0 : initRecords$1[0]) === null || _a$9 === void 0 ? void 0 : _a$9.v_key,
  v_value: (_b$7 = initRecords$1 === null || initRecords$1 === void 0 ? void 0 : initRecords$1[0]) === null || _b$7 === void 0 ? void 0 : _b$7.v_value
}];
var stylize = teaComponent.Table.addons.stylize;
var AffinityMapField = function AffinityMapField(_a) {
  var _b;
  var plan = _a.plan,
    onChange = _a.onChange;
  var _c = React.useState({
      'nodeSelector': initRecords$1,
      'nodeAffinity': initRecords$1
    }),
    map = _c[0],
    setMap = _c[1];
  var _d = React.useState(true),
    enable = _d[0],
    setEnable = _d[1];
  var _e = React.useState(''),
    editType = _e[0],
    setEditType = _e[1];
  React.useEffect(function () {
    var _a, _b, _c, _d, _e;
    var data = {
      nodeSelector: {
        records: (_a = map === null || map === void 0 ? void 0 : map['nodeSelector']) !== null && _a !== void 0 ? _a : [],
        isValid: true
      },
      nodeAffinity: {
        records: ((_b = map === null || map === void 0 ? void 0 : map['nodeAffinity']) === null || _b === void 0 ? void 0 : _b.some(function (item) {
          return item.key === NodeAffinityMustNeedKey;
        })) ? map === null || map === void 0 ? void 0 : map['nodeAffinity'] : NoeAffinityInitRecords.concat(map === null || map === void 0 ? void 0 : map['nodeAffinity']),
        isValid: true
      }
    };
    data.nodeSelector.records = enable ? data.nodeSelector.records : [];
    data.nodeSelector.isValid = enable ? (_c = data.nodeSelector.records) === null || _c === void 0 ? void 0 : _c.every(function (item) {
      var _a, _b;
      return ((_a = item === null || item === void 0 ? void 0 : item.v_key) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success && ((_b = item === null || item === void 0 ? void 0 : item.v_value) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Success;
    }) : true;
    data.nodeAffinity.records = enable ? (_d = data.nodeAffinity.records) === null || _d === void 0 ? void 0 : _d.map(function (item) {
      return tslib.__assign(tslib.__assign({}, item), {
        value: item.operator === NodeAffinityOperatorEnum.Exists || item.operator === NodeAffinityOperatorEnum.DoesNotExist ? '' : item.value,
        v_value: item.operator === NodeAffinityOperatorEnum.Exists || item.operator === NodeAffinityOperatorEnum.DoesNotExist ? {
          status: ffValidator.ValidatorStatusEnum.Success,
          message: ''
        } : item === null || item === void 0 ? void 0 : item['v_value']
      });
    }) : [];
    data.nodeAffinity.isValid = enable ? (_e = data.nodeAffinity.records) === null || _e === void 0 ? void 0 : _e.every(function (item) {
      var _a, _b;
      return ((_a = item === null || item === void 0 ? void 0 : item.v_key) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success && ((_b = item === null || item === void 0 ? void 0 : item.v_value) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Success;
    }) : true;
    onChange(tslib.__assign({
      enable: enable
    }, data));
  }, [map, editType]);
  React.useEffect(function () {
    var _a, _b;
    var map;
    if (!enable) {
      map = {
        'nodeSelector': [],
        'nodeAffinity': []
      };
    } else {
      var affinityRecords = [{
        id: uuid(),
        key: NodeAffinityMustNeedKey,
        value: '',
        operator: NodeAffinityOperatorEnum.Exists,
        v_key: (_a = initRecords$1 === null || initRecords$1 === void 0 ? void 0 : initRecords$1[0]) === null || _a === void 0 ? void 0 : _a.v_key,
        v_value: (_b = initRecords$1 === null || initRecords$1 === void 0 ? void 0 : initRecords$1[0]) === null || _b === void 0 ? void 0 : _b.v_value
      }];
      map = {
        'nodeSelector': [],
        'nodeAffinity': affinityRecords
      };
    }
    setMap(tslib.__assign({}, map));
  }, [enable]);
  var _delete = function _delete(type, id) {
    var _a;
    setEditType(type);
    var records = map === null || map === void 0 ? void 0 : map[type];
    setMap(tslib.__assign(tslib.__assign({}, map), (_a = {}, _a[type] = records === null || records === void 0 ? void 0 : records.filter(function (item) {
      return (item === null || item === void 0 ? void 0 : item.id) !== id;
    }), _a)));
  };
  var _add = function _add(type) {
    var _a;
    setEditType(type);
    var newItem = tslib.__assign(tslib.__assign({}, initRecords$1 === null || initRecords$1 === void 0 ? void 0 : initRecords$1[0]), {
      id: uuid()
    });
    var records = map === null || map === void 0 ? void 0 : map[type];
    setMap(tslib.__assign(tslib.__assign({}, map), (_a = {}, _a[type] = records === null || records === void 0 ? void 0 : records.concat([newItem]), _a)));
  };
  var _update = function _update(type, fieldName, selectItem, value) {
    var _a;
    if (value === void 0) {
      value = '';
    }
    setEditType(type);
    var records = map === null || map === void 0 ? void 0 : map[type];
    var validation = {};
    if (fieldName === 'key') {
      validation = !value ? {
        status: ffValidator.ValidatorStatusEnum.Failed,
        message: i18n.t("".concat(fieldName, "\u4E0D\u80FD\u4E3A\u7A7A"))
      } : (records === null || records === void 0 ? void 0 : records.some(function (item) {
        return item.key === value;
      })) ? {
        status: ffValidator.ValidatorStatusEnum.Failed,
        message: i18n.t("".concat(fieldName, "\u4E0D\u80FD\u91CD\u590D"))
      } : {
        status: ffValidator.ValidatorStatusEnum.Success,
        message: ''
      };
    } else if (fieldName === 'value') {
      validation = !value ? {
        status: ffValidator.ValidatorStatusEnum.Failed,
        message: i18n.t("".concat(fieldName, "\u4E0D\u80FD\u4E3A\u7A7A"))
      } : {
        status: ffValidator.ValidatorStatusEnum.Success,
        message: ''
      };
    } else {
      validation = {
        status: ffValidator.ValidatorStatusEnum.Success,
        message: ''
      };
    }
    var newRecords = records === null || records === void 0 ? void 0 : records.map(function (item) {
      var _a;
      return tslib.__assign(tslib.__assign({}, item), (_a = {}, _a[fieldName] = (item === null || item === void 0 ? void 0 : item.id) === (selectItem === null || selectItem === void 0 ? void 0 : selectItem.id) ? value : item === null || item === void 0 ? void 0 : item[fieldName], _a["v_".concat(fieldName)] = item.id === (selectItem === null || selectItem === void 0 ? void 0 : selectItem.id) ? tslib.__assign({}, validation) : item === null || item === void 0 ? void 0 : item["v_".concat(fieldName)], _a));
    });
    setMap(tslib.__assign(tslib.__assign({}, map), (_a = {}, _a[type] = newRecords, _a)));
  };
  return React__default.createElement(React__default.Fragment, null, enable && React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('节点选择')
  }, !!((_b = map === null || map === void 0 ? void 0 : map.nodeSelector) === null || _b === void 0 ? void 0 : _b.length) && React__default.createElement(teaComponent.Table, {
    addons: [stylize({
      headClassName: "nodeSelector-head",
      headStyle: {
        backgroundColor: "rgb(231, 234, 239)"
      }
    })],
    bordered: true,
    recordKey: "id",
    columns: [{
      key: "key",
      header: 'Key',
      width: '45%',
      render: function render(item) {
        var _a, _b;
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Input, {
          value: item === null || item === void 0 ? void 0 : item.key,
          onChange: function onChange(value) {
            _update('nodeSelector', 'key', item, value);
          },
          placeholder: i18n.t('Label Key'),
          style: ((_a = item === null || item === void 0 ? void 0 : item['v_key']) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Failed ? {
            border: '1px solid red',
            width: '90%'
          } : {
            width: '100%'
          },
          className: 'tea-mr-2n',
          width: '100%'
        }), ((_b = item === null || item === void 0 ? void 0 : item['v_key']) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Failed && React__default.createElement(teaComponent.Bubble, {
          content: i18n.t('{{msg}}', {
            msg: item === null || item === void 0 ? void 0 : item['v_key'].message
          })
        }, React__default.createElement(teaComponent.Icon, {
          type: 'error'
        })));
      }
    }, {
      key: "value",
      header: 'Value',
      width: '45%',
      render: function render(item) {
        var _a, _b;
        return React__default.createElement("div", null, React__default.createElement(teaComponent.Input, {
          value: item === null || item === void 0 ? void 0 : item.value,
          onChange: function onChange(value) {
            _update('nodeSelector', 'value', item, value);
          },
          placeholder: i18n.t('多个Label Value请以 ; 分隔符隔开'),
          style: ((_a = item === null || item === void 0 ? void 0 : item['v_value']) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Failed ? {
            border: '1px solid red',
            width: '90%'
          } : {
            width: '100%'
          },
          className: 'tea-mr-2n',
          width: '100%'
        }), ((_b = item === null || item === void 0 ? void 0 : item['v_value']) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Failed && React__default.createElement(teaComponent.Bubble, {
          content: i18n.t('{{msg}}', {
            msg: item === null || item === void 0 ? void 0 : item['v_value'].message
          })
        }, React__default.createElement(teaComponent.Icon, {
          type: 'error'
        })));
      }
    }, {
      key: 'operate',
      header: null,
      width: '10%',
      render: function render(item) {
        return React__default.createElement(teaComponent.Button, {
          type: 'link',
          onClick: function onClick() {
            _delete('nodeSelector', item === null || item === void 0 ? void 0 : item.id);
          }
        }, i18n.t('删除'));
      }
    }],
    records: map.nodeSelector
  }), React__default.createElement(teaComponent.Button, {
    onClick: function onClick() {
      _add('nodeSelector');
    },
    type: 'link',
    className: 'tea-mt-2n'
  }, i18n.t('添加'))), enable && React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('节点亲和性')
  }, React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Table, {
    addons: [stylize({
      headClassName: "nodeAffinity-head",
      headStyle: {
        backgroundColor: "rgb(231, 234, 239)"
      }
    })],
    bordered: true,
    recordKey: "id",
    columns: [{
      key: "key",
      header: 'Key',
      width: '30%',
      render: function render(item) {
        var _a, _b, _c, _d;
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Tooltip, {
          title: (item === null || item === void 0 ? void 0 : item.key) === NodeAffinityMustNeedKey && ((_a = item === null || item === void 0 ? void 0 : item.v_key) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success ? i18n.t('该数据不可进行操作') : ''
        }, React__default.createElement(teaComponent.Input, {
          value: item === null || item === void 0 ? void 0 : item.key,
          onChange: function onChange(value) {
            _update('nodeAffinity', 'key', item, value);
          },
          placeholder: i18n.t('Label Key'),
          style: ((_b = item === null || item === void 0 ? void 0 : item['v_key']) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Failed ? {
            border: '1px solid red',
            width: '90%'
          } : {
            width: '100%'
          },
          className: 'tea-mr-2n',
          disabled: (item === null || item === void 0 ? void 0 : item.key) === NodeAffinityMustNeedKey && ((_c = item === null || item === void 0 ? void 0 : item.v_key) === null || _c === void 0 ? void 0 : _c.status) === ffValidator.ValidatorStatusEnum.Success
        }), ((_d = item === null || item === void 0 ? void 0 : item['v_key']) === null || _d === void 0 ? void 0 : _d.status) === ffValidator.ValidatorStatusEnum.Failed && React__default.createElement(teaComponent.Bubble, {
          content: i18n.t('{{msg}}', {
            msg: item === null || item === void 0 ? void 0 : item['v_key'].message
          })
        }, React__default.createElement(teaComponent.Icon, {
          type: 'error'
        }))));
      }
    }, {
      key: "operator",
      header: '操作',
      width: '30%',
      render: function render(item) {
        var _a;
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Select, {
          virtual: true,
          value: item === null || item === void 0 ? void 0 : item.operator,
          key: 'value',
          options: NodeAffinityRuleOperatorList.map(function (item) {
            return tslib.__assign(tslib.__assign({}, item), {
              tooltip: item === null || item === void 0 ? void 0 : item.tip,
              text: item === null || item === void 0 ? void 0 : item.value
            });
          }),
          onChange: function onChange(value) {
            _update('nodeAffinity', 'operator', item, value);
          },
          style: {
            width: '100%'
          },
          disabled: (item === null || item === void 0 ? void 0 : item.key) === NodeAffinityMustNeedKey && ((_a = item === null || item === void 0 ? void 0 : item.v_key) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success,
          matchButtonWidth: true,
          appearance: "button"
        }));
      }
    }, {
      key: "value",
      header: 'Value',
      width: '30%',
      render: function render(item) {
        var _a, _b;
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Tooltip, {
          title: item.operator === NodeAffinityOperatorEnum.Exists || item.operator === NodeAffinityOperatorEnum.DoesNotExist ? i18n.t('DoesNotExist,Exists操作符不需要填写value') : null
        }, React__default.createElement(teaComponent.Input, {
          value: item === null || item === void 0 ? void 0 : item.value,
          onChange: function onChange(value) {
            _update('nodeAffinity', 'value', item, value);
          },
          placeholder: item.operator === NodeAffinityOperatorEnum.Exists || item.operator === NodeAffinityOperatorEnum.DoesNotExist ? i18n.t('DoesNotExist,Exists操作符不需要填写value') : i18n.t('多个Label Value请以 ; 分隔符隔开'),
          style: ((_a = item === null || item === void 0 ? void 0 : item['v_value']) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Failed ? {
            border: '1px solid red',
            width: '90%'
          } : {
            width: '100%'
          },
          disabled: item.operator === NodeAffinityOperatorEnum.Exists || item.operator === NodeAffinityOperatorEnum.DoesNotExist,
          className: 'tea-mr-2n',
          size: 'full'
        })), ((_b = item === null || item === void 0 ? void 0 : item['v_value']) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Failed && React__default.createElement(teaComponent.Bubble, {
          content: i18n.t('{{msg}}', {
            msg: item === null || item === void 0 ? void 0 : item['v_value'].message
          })
        }, React__default.createElement(teaComponent.Icon, {
          type: 'error'
        })));
      }
    }, {
      key: 'operate',
      header: null,
      width: '10%',
      render: function render(item) {
        var _a, _b;
        return React__default.createElement(teaComponent.Button, {
          type: 'link',
          onClick: function onClick() {
            _delete('nodeAffinity', item === null || item === void 0 ? void 0 : item.id);
          },
          disabled: (item === null || item === void 0 ? void 0 : item.key) === NodeAffinityMustNeedKey && ((_a = item === null || item === void 0 ? void 0 : item.v_key) === null || _a === void 0 ? void 0 : _a.status) === ffValidator.ValidatorStatusEnum.Success,
          tooltip: (item === null || item === void 0 ? void 0 : item.key) === NodeAffinityMustNeedKey && ((_b = item === null || item === void 0 ? void 0 : item.v_key) === null || _b === void 0 ? void 0 : _b.status) === ffValidator.ValidatorStatusEnum.Success ? i18n.t('该数据不可进行删除操作') : ''
        }, i18n.t('删除'));
      }
    }],
    records: map === null || map === void 0 ? void 0 : map.nodeAffinity
  }), React__default.createElement(teaComponent.Button, {
    onClick: function onClick() {
      _add('nodeAffinity');
    },
    type: 'link',
    className: 'tea-mt-2n'
  }, i18n.t('添加')))));
};

function ServiceInstanceCreatePanel(props) {
  var _this = this;
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p;
  var actions = props.actions,
    _q = props.base,
    platform = _q.platform,
    route = _q.route,
    hubCluster = _q.hubCluster,
    isI18n = _q.isI18n,
    userInfo = _q.userInfo,
    regionId = _q.regionId,
    getClusterAdminRole = _q.getClusterAdminRole,
    _r = props.list,
    services = _r.services,
    servicePlans = _r.servicePlans,
    createResourceWorkflow = _r.createResourceWorkflow,
    createServiceInstanceWorkflow = _r.createServiceInstanceWorkflow,
    externalClusters = _r.externalClusters;
  var servicename = (route === null || route === void 0 ? void 0 : route.queries).servicename;
  var servicesInstance = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.servicesInstance;
  });
  var servicesInstanceEdit = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.serviceInstanceEdit;
  });
  var isLoadingSchema = ((_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _a === void 0 ? void 0 : _a.fetched) && (externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) ? !((_b = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _b === void 0 ? void 0 : _b.fetched) : !((_c = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _c === void 0 ? void 0 : _c.fetched);
  var loadSchemaFailed = ((_d = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _d === void 0 ? void 0 : _d.fetched) && ((_e = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _e === void 0 ? void 0 : _e.error) && (!((_f = ['ResourceNotFound', 404]) === null || _f === void 0 ? void 0 : _f.includes((_h = (_g = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _g === void 0 ? void 0 : _g.error) === null || _h === void 0 ? void 0 : _h.code)) || !((_l = (_k = (_j = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _j === void 0 ? void 0 : _j.error) === null || _k === void 0 ? void 0 : _k.message) === null || _l === void 0 ? void 0 : _l.includes('404')));
  var createWorkflow = ((_m = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _m === void 0 ? void 0 : _m.timeBackup) ? createServiceInstanceWorkflow : createResourceWorkflow;
  React.useEffect(function () {
    var _a;
    var clusterId = (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId;
    if (actions.list.resourceSchemas && platform && servicename && clusterId) {
      actions.list.resourceSchemas.reset();
      actions.list.resourceSchemas.applyFilter({
        platform: platform,
        serviceName: servicename,
        clusterId: clusterId,
        regionId: regionId
      });
    }
  }, [actions.list.resourceSchemas, regionId, platform, servicename, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  React.useEffect(function () {
    var _a;
    var clusterId = (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId;
    if (actions.list.servicePlans && platform && servicename && clusterId) {
      actions.list.servicePlans.reset();
      actions.list.servicePlans.applyFilter({
        platform: platform,
        serviceName: servicename,
        clusterId: clusterId,
        regionId: regionId,
        resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan
      });
    }
  }, [actions.list.servicePlans, platform, servicename, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  //切换集群时重置创建实例数据流状态
  React.useEffect(function () {
    var _a;
    if (createWorkflow.operationState === ffRedux.OperationState.Done) {
      if (!((_a = servicesInstanceEdit.formData) === null || _a === void 0 ? void 0 : _a.timeBackup)) {
        actions.create.createResource.reset();
      } else {
        actions.create.createServiceInstance.reset();
      }
    }
  }, [externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  React.useEffect(function () {
    var _a, _b;
    if (platform && ((_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.externalClusters) && (services === null || services === void 0 ? void 0 : services.selection)) {
      (_b = actions === null || actions === void 0 ? void 0 : actions.list) === null || _b === void 0 ? void 0 : _b.externalClusters.applyFilter({
        platform: platform,
        clusterIds: [],
        regionId: regionId
      });
    }
  }, [(_o = actions === null || actions === void 0 ? void 0 : actions.list) === null || _o === void 0 ? void 0 : _o.externalClusters, platform, services === null || services === void 0 ? void 0 : services.selection]);
  // 由于vendor在目标集群上开启为异步操作，因此若vendor在集群上执行了开启且schema加载resource not found,则认为vendor正在开启中
  var openingVendor = React__default === null || React__default === void 0 ? void 0 : React__default.useMemo(function () {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
    return ((_b = (_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.clusters) === null || _b === void 0 ? void 0 : _b.includes((_c = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _c === void 0 ? void 0 : _c.clusterId)) && ((_d = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _d === void 0 ? void 0 : _d.fetched) && (((_f = (_e = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _e === void 0 ? void 0 : _e.error) === null || _f === void 0 ? void 0 : _f.code) === 'ResourceNotFound' || ((_h = (_g = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _g === void 0 ? void 0 : _g.error) === null || _h === void 0 ? void 0 : _h.code) === 404 || ((_l = (_k = (_j = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _j === void 0 ? void 0 : _j.error) === null || _k === void 0 ? void 0 : _k.message) === null || _l === void 0 ? void 0 : _l.includes('404')));
  }, [servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object, services === null || services === void 0 ? void 0 : services.selection, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  var reduceNodeAffinityJson = function reduceNodeAffinityJson(data) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
    var _m = (_a = data.formData) === null || _a === void 0 ? void 0 : _a['nodeSchedule'],
      enable = _m.enable,
      nodeSelector = _m.nodeSelector,
      nodeAffinity = _m.nodeAffinity;
    var nodeSelectorJson = ((_c = (_b = nodeSelector === null || nodeSelector === void 0 ? void 0 : nodeSelector.records) === null || _b === void 0 ? void 0 : _b.filter(function (item) {
      return item === null || item === void 0 ? void 0 : item.key;
    })) === null || _c === void 0 ? void 0 : _c.length) ? {
      nodeSelector: (_e = (_d = nodeSelector === null || nodeSelector === void 0 ? void 0 : nodeSelector.records) === null || _d === void 0 ? void 0 : _d.filter(function (item) {
        return item === null || item === void 0 ? void 0 : item.key;
      })) === null || _e === void 0 ? void 0 : _e.reduce(function (pre, cur) {
        var _a;
        return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.key] = cur === null || cur === void 0 ? void 0 : cur.value, _a));
      }, {})
    } : {};
    var nodeAffinityJson = ((_g = (_f = nodeAffinity === null || nodeAffinity === void 0 ? void 0 : nodeAffinity.records) === null || _f === void 0 ? void 0 : _f.filter(function (item) {
      return item === null || item === void 0 ? void 0 : item.key;
    })) === null || _g === void 0 ? void 0 : _g.length) ? {
      affinity: {
        nodeAffinity: {
          requiredDuringSchedulingIgnoredDuringExecution: {
            nodeSelectorTerms: [{
              matchExpressions: (_j = (_h = nodeAffinity === null || nodeAffinity === void 0 ? void 0 : nodeAffinity.records) === null || _h === void 0 ? void 0 : _h.filter(function (item) {
                return item === null || item === void 0 ? void 0 : item.key;
              })) === null || _j === void 0 ? void 0 : _j.map(function (item) {
                var _a, _b;
                if ([NodeAffinityOperatorEnum.DoesNotExist, NodeAffinityOperatorEnum.Exists].includes(item === null || item === void 0 ? void 0 : item.operator)) {
                  return {
                    key: item === null || item === void 0 ? void 0 : item.key,
                    operator: item === null || item === void 0 ? void 0 : item.operator
                  };
                } else {
                  return {
                    key: item === null || item === void 0 ? void 0 : item.key,
                    operator: item === null || item === void 0 ? void 0 : item.operator,
                    value: (_b = (_a = item === null || item === void 0 ? void 0 : item.value) === null || _a === void 0 ? void 0 : _a.split(';')) !== null && _b !== void 0 ? _b : []
                  };
                }
              })
            }]
          }
        }
      }
    } : {};
    var dataJson = enable && (((_k = nodeSelector === null || nodeSelector === void 0 ? void 0 : nodeSelector.records) === null || _k === void 0 ? void 0 : _k.length) || ((_l = nodeAffinity === null || nodeAffinity === void 0 ? void 0 : nodeAffinity.records) === null || _l === void 0 ? void 0 : _l.length)) ? {
      scheduling: tslib.__assign(tslib.__assign({}, nodeSelectorJson), nodeAffinityJson)
    } : {};
    return dataJson;
  };
  var reduceCreateServiceResourceDataJson = function reduceCreateServiceResourceDataJson(data) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
    var _l = data === null || data === void 0 ? void 0 : data.formData,
      instanceName = _l.instanceName,
      plan = _l.plan,
      timeBackup = _l.timeBackup,
      clusterId = _l.clusterId;
    //拼接中间件实例参数部分属性值
    var parameters = (_c = (_b = (_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.instanceCreateParameterSchema) === null || _c === void 0 ? void 0 : _c.reduce(function (pre, cur) {
      var _a;
      var _b, _c, _d, _e;
      return Object.assign(pre, ((_b = data === null || data === void 0 ? void 0 : data.formData) === null || _b === void 0 ? void 0 : _b[cur === null || cur === void 0 ? void 0 : cur.name]) !== '' ? (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = formatPlanSchemaSubmitData(cur, (_c = data === null || data === void 0 ? void 0 : data.formData) === null || _c === void 0 ? void 0 : _c[cur === null || cur === void 0 ? void 0 : cur.name], (_e = (_d = data === null || data === void 0 ? void 0 : data.formData) === null || _d === void 0 ? void 0 : _d['unitMap']) === null || _e === void 0 ? void 0 : _e[cur === null || cur === void 0 ? void 0 : cur.name]), _a) : {});
    }, {});
    //拼接节点调度选项参数
    parameters = tslib.__assign(tslib.__assign({}, parameters), reduceNodeAffinityJson(data));
    var json = {
      id: uuid$1(),
      apiVersion: 'infra.tce.io/v1',
      kind: (_d = ResourceTypeMap === null || ResourceTypeMap === void 0 ? void 0 : ResourceTypeMap[ResourceTypeEnum.ServiceResource]) === null || _d === void 0 ? void 0 : _d.resourceKind,
      metadata: {
        labels: {
          'ssm.infra.tce.io/cluster-id': clusterId
        },
        annotations: {
          'ssm.infra.tce.io/creator': decodeURIComponent ? decodeURIComponent((_f = (_e = userInfo === null || userInfo === void 0 ? void 0 : userInfo.object) === null || _e === void 0 ? void 0 : _e.data) === null || _f === void 0 ? void 0 : _f.name) : (_h = (_g = userInfo === null || userInfo === void 0 ? void 0 : userInfo.object) === null || _g === void 0 ? void 0 : _g.data) === null || _h === void 0 ? void 0 : _h.name
        },
        name: instanceName,
        namespace: DefaultNamespace
      },
      spec: {
        serviceClass: ((_j = services === null || services === void 0 ? void 0 : services.selection) === null || _j === void 0 ? void 0 : _j.name) || ((_k = route === null || route === void 0 ? void 0 : route.queries) === null || _k === void 0 ? void 0 : _k.servicename),
        servicePlan: plan,
        enableBackup: timeBackup,
        parameters: parameters,
        externalID: DefaultNamespace + '-' + instanceName
        // ...backupParams,
      }
    };

    return JSON.stringify(json);
  };
  var _cancel = function _cancel() {
    var _a, _b, _c;
    router === null || router === void 0 ? void 0 : router.navigate({
      sub: 'list',
      tab: undefined
    }, {
      servicename: (_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.servicename,
      resourceType: (_c = (_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.resourceType) !== null && _c !== void 0 ? _c : ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource,
      mode: 'list'
    });
  };
  var _submit = function _submit() {
    var _a, _b, _c, _d;
    var formData = servicesInstanceEdit.formData;
    actions.create.validateInstance();
    if (_validateInstance(servicesInstanceEdit, (_b = (_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.instanceCreateParameterSchema)) {
      if (!(formData === null || formData === void 0 ? void 0 : formData.timeBackup)) {
        var params = {
          platform: platform,
          regionId: HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion,
          clusterId: formData === null || formData === void 0 ? void 0 : formData.clusterId,
          jsonData: reduceCreateServiceResourceDataJson(servicesInstanceEdit),
          namespace: DefaultNamespace
        };
        actions.create.createResource.start([params], HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion);
        actions.create.createResource.perform();
      } else {
        var instanceParams = {
          platform: platform,
          regionId: HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion,
          clusterId: formData === null || formData === void 0 ? void 0 : formData.clusterId,
          jsonData: reduceCreateServiceResourceDataJson(servicesInstanceEdit),
          specificOperate: CreateSpecificOperatorEnum === null || CreateSpecificOperatorEnum === void 0 ? void 0 : CreateSpecificOperatorEnum.CreateResource,
          namespace: DefaultNamespace
        };
        //更改备份策略命名空间
        var backupStrategyParams = {
          platform: platform,
          regionId: HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion,
          clusterId: formData === null || formData === void 0 ? void 0 : formData.clusterId,
          jsonData: Backup.reduceBackupStrategyJson({
            enable: formData === null || formData === void 0 ? void 0 : formData.timeBackup,
            backupDate: formData === null || formData === void 0 ? void 0 : formData.backupDate,
            backupTime: formData === null || formData === void 0 ? void 0 : formData.backupTime,
            backupReserveDay: formData === null || formData === void 0 ? void 0 : formData.backupReserveDay,
            instanceId: DefaultNamespace + '-' + (formData === null || formData === void 0 ? void 0 : formData.instanceName),
            instanceName: formData === null || formData === void 0 ? void 0 : formData.instanceName,
            serviceName: ((_c = services === null || services === void 0 ? void 0 : services.selection) === null || _c === void 0 ? void 0 : _c.name) || ((_d = route === null || route === void 0 ? void 0 : route.queries) === null || _d === void 0 ? void 0 : _d.servicename)
          }),
          resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup,
          namespace: SystemNamespace
        };
        var params = {
          instance: [instanceParams],
          backupStrategy: [backupStrategyParams]
        };
        actions.create.createServiceInstance.start([params], HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion);
        actions.create.createServiceInstance.perform();
      }
    }
  };
  var renderContent = function renderContent(data) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v, _w, _x, _y, _z, _0, _1, _2, _3, _4, _5, _6, _7, _8, _9, _10, _11, _12, _13, _14, _15, _16, _17, _18, _19, _20, _21, _22, _23, _24, _25, _26, _27, _28, _29;
    var versionFields = (_a = data === null || data === void 0 ? void 0 : data.instanceCreateParameterSchema) === null || _a === void 0 ? void 0 : _a.filter(function (item) {
      return (item === null || item === void 0 ? void 0 : item.name) === 'version';
    });
    var instanceParamsFields = (_b = data === null || data === void 0 ? void 0 : data.instanceCreateParameterSchema) === null || _b === void 0 ? void 0 : _b.filter(function (item) {
      return (item === null || item === void 0 ? void 0 : item.name) !== 'version';
    });
    var isSubmitting = (createWorkflow === null || createWorkflow === void 0 ? void 0 : createWorkflow.operationState) === ffRedux.OperationState.Performing;
    var failed = createWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(createWorkflow);
    var showBackUpOperation = (_e = (_d = (_c = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.supportedOperations) === null || _e === void 0 ? void 0 : _e.some(function (operation) {
      return (operation === null || operation === void 0 ? void 0 : operation.operation) === (SupportedOperationsEnum === null || SupportedOperationsEnum === void 0 ? void 0 : SupportedOperationsEnum.Backup);
    });
    return React__default.createElement(ffComponent.FormPanel, null, React__default.createElement(ffComponent.FormPanel.Item, {
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: 'flex',
          alignItems: 'center'
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t('中间件类型')), React__default.createElement(teaComponent.Text, {
        className: 'text-danger tea-pt-1n'
      }, "*")),
      text: true
    }, React__default.createElement(ffComponent.FormPanel.Text, null, i18n.t('{{name}}', {
      name: (_f = route === null || route === void 0 ? void 0 : route.queries) === null || _f === void 0 ? void 0 : _f.servicename
    }))), React__default.createElement(ffComponent.FormPanel.Item, {
      after: React__default.createElement(teaComponent.Button, {
        icon: "refresh",
        onClick: function onClick() {
          actions.list.externalClusters.fetch();
        }
      }),
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: 'flex',
          alignItems: 'center'
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t('目标集群')), React__default.createElement(teaComponent.Text, {
        className: 'text-danger tea-pt-1n'
      }, "*")),
      validator: (_g = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _g === void 0 ? void 0 : _g['clusterId'],
      select: {
        model: externalClusters,
        valueField: 'clusterId',
        displayField: function displayField(record) {
          return (record === null || record === void 0 ? void 0 : record.clusterId) + ' (' + (record === null || record === void 0 ? void 0 : record.clusterName) + ') ';
        },
        action: (_h = actions === null || actions === void 0 ? void 0 : actions.list) === null || _h === void 0 ? void 0 : _h.externalClusters,
        disabledField: function disabledField(record) {
          var _a, _b, _c;
          var disabled = !((_b = (_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.clusters) === null || _b === void 0 ? void 0 : _b.includes(record === null || record === void 0 ? void 0 : record.clusterId));
          return {
            disabled: disabled,
            tooltip: disabled ? i18n.t('{{tooltip}}', {
              tooltip: "".concat((_c = services === null || services === void 0 ? void 0 : services.selection) === null || _c === void 0 ? void 0 : _c.name, "\u5C1A\u672A\u5728\u8BE5\u76EE\u6807\u96C6\u7FA4\u5F00\u542F,\u8BF7\u60A8\u524D\u5F80\u5206\u5E03\u5F0F\u4E91\u4E2D\u5FC3\u6982\u89C8\u9875\u5F00\u542F")
            }) : null
          };
        }
      },
      message: React__default.createElement(React__default.Fragment, null, ((_j = services === null || services === void 0 ? void 0 : services.list) === null || _j === void 0 ? void 0 : _j.fetched) && !((_l = (_k = services.selection) === null || _k === void 0 ? void 0 : _k.clusters) === null || _l === void 0 ? void 0 : _l.length) ? React__default.createElement(RetryPanel, {
        loadingText: i18n.t('集群列表尚无已开启vendor的集群'),
        retryText: i18n.t('刷新重试'),
        action: actions.list.services.fetch
      }) : ((_m = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _m === void 0 ? void 0 : _m.fetched) && !((_p = (_o = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _o === void 0 ? void 0 : _o.data) === null || _p === void 0 ? void 0 : _p.recordCount) ? React__default.createElement(RetryPanel, {
        loadingText: i18n.t('集群列表尚无运行状态的目标集群'),
        retryText: i18n.t('前往确认'),
        action: function action() {
          router.navigate({}, {}, '/tdcc/cluster');
        }
      }) : null)
    }), React__default.createElement(ffComponent.FormPanel.Item, {
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: 'flex',
          alignItems: 'center'
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t('中间件实例名称')), React__default.createElement(teaComponent.Text, {
        className: 'text-danger tea-pt-1n'
      }, "*")),
      style: {
        paddingTop: 0
      },
      validator: (_q = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _q === void 0 ? void 0 : _q['instanceName'],
      input: {
        value: (_r = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _r === void 0 ? void 0 : _r['instanceName'],
        onChange: function onChange(e) {
          var _a;
          (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updateInstance('instanceName', e);
        },
        onBlur: function onBlur() {}
      }
    }), React__default.createElement(ffComponent.FormPanel.Item, {
      after: React__default.createElement(teaComponent.Button, {
        icon: "refresh",
        onClick: function onClick() {
          actions.list.servicePlans.fetch();
        }
      }),
      key: 'plan',
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: 'flex',
          alignItems: 'center'
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t('规格')), React__default.createElement(teaComponent.Text, {
        className: 'text-danger tea-pt-1n'
      }, "*")),
      validator: (_s = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _s === void 0 ? void 0 : _s['plan'],
      message: React__default.createElement(React__default.Fragment, null, ((_t = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.list) === null || _t === void 0 ? void 0 : _t.error) ? React__default.createElement(RetryPanel, {
        loadingText: "\u52A0\u8F7D\u89C4\u683C\u5931\u8D25",
        action: actions.list.servicePlans
      }) : ((_u = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.list) === null || _u === void 0 ? void 0 : _u.fetched) && !((_w = (_v = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.list) === null || _v === void 0 ? void 0 : _v.data) === null || _w === void 0 ? void 0 : _w.recordCount) ? React__default.createElement(teaComponent.Text, null, i18n.t('无可用规格,'), React__default.createElement(teaComponent.Button, {
        type: "link",
        onClick: function onClick() {
          var _a, _b;
          router.navigate({
            sub: 'create',
            tab: ResourceTypeEnum.ServicePlan
          }, {
            servicename: ((_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.name) || ((_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.servicename),
            mode: 'create',
            resourceType: ResourceTypeEnum.ServicePlan
          });
        }
      }, i18n.t('立即新建'))) : null)
    }, React__default.createElement(ffComponent.FormPanel.Select, {
      placeholder: i18n.t('{{title}}', {
        title: '请选择规格'
      }),
      model: servicePlans,
      action: (_x = actions === null || actions === void 0 ? void 0 : actions.list) === null || _x === void 0 ? void 0 : _x.servicePlans,
      valueField: function valueField(record) {
        var _a;
        return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name;
      },
      displayField: function displayField(record) {
        var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
        return "".concat((_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name, "(cpu:").concat((_d = (_c = (_b = record === null || record === void 0 ? void 0 : record.spec) === null || _b === void 0 ? void 0 : _b.metadata) === null || _c === void 0 ? void 0 : _c.cpu) !== null && _d !== void 0 ? _d : '-', ",memory:").concat((_g = (_f = (_e = record === null || record === void 0 ? void 0 : record.spec) === null || _e === void 0 ? void 0 : _e.metadata) === null || _f === void 0 ? void 0 : _f.memory) !== null && _g !== void 0 ? _g : '-', ",storage:").concat((_k = (_j = (_h = record === null || record === void 0 ? void 0 : record.spec) === null || _h === void 0 ? void 0 : _h.metadata) === null || _j === void 0 ? void 0 : _j.storage) !== null && _k !== void 0 ? _k : '-', ")");
      },
      value: (_z = (_y = servicePlans === null || servicePlans === void 0 ? void 0 : servicePlans.selection) === null || _y === void 0 ? void 0 : _y.metadata) === null || _z === void 0 ? void 0 : _z.name
    })), showBackUpOperation && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('定时备份'),
      required: false,
      validator: (_0 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _0 === void 0 ? void 0 : _0['timeBackup'],
      message: ((_2 = (_1 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _1 === void 0 ? void 0 : _1['backupDate']) === null || _2 === void 0 ? void 0 : _2.status) === 2 || ((_4 = (_3 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _3 === void 0 ? void 0 : _3['backupTime']) === null || _4 === void 0 ? void 0 : _4.status) === 2 ? React__default.createElement(teaComponent.Text, {
        className: "text-danger"
      }, i18n.t('{{message}}', {
        message: (_6 = (_5 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _5 === void 0 ? void 0 : _5['backupDate']) === null || _6 === void 0 ? void 0 : _6.message
      })) : null
    }, React__default.createElement(teaComponent.Switch, {
      value: (_7 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _7 === void 0 ? void 0 : _7['timeBackup'],
      onChange: function onChange(value) {
        return tslib.__awaiter(_this, void 0, void 0, function () {
          var cosResource;
          var _a;
          return tslib.__generator(this, function (_b) {
            switch (_b.label) {
              case 0:
                if (!value) return [3 /*break*/, 2];
                return [4 /*yield*/, checkCosResource({
                  platform: platform,
                  clusterId: Util === null || Util === void 0 ? void 0 : Util.getCOSClusterId(platform, (_a = hubCluster === null || hubCluster === void 0 ? void 0 : hubCluster.object) === null || _a === void 0 ? void 0 : _a.data),
                  regionId: regionId
                })];
              case 1:
                cosResource = _b.sent();
                if (cosResource) {
                  actions === null || actions === void 0 ? void 0 : actions.create.updateInstance('timeBackup', value);
                } else {
                  actions === null || actions === void 0 ? void 0 : actions.create.validateTimeBackup('backupDate', i18n.t('{{msg}}', {
                    msg: ErrorMsgEnum.COS_Resource_Not_Found
                  }));
                }
                return [3 /*break*/, 3];
              case 2:
                actions === null || actions === void 0 ? void 0 : actions.create.updateInstance('timeBackup', value);
                _b.label = 3;
              case 3:
                return [2 /*return*/];
            }
          });
        });
      }
    })), showBackUpOperation && ((_8 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _8 === void 0 ? void 0 : _8['timeBackup']) ? React__default.createElement(React__default.Fragment, null, React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('备份日期')
    }, React__default.createElement(teaComponent.Checkbox.Group, {
      value: (_9 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _9 === void 0 ? void 0 : _9.backupDate,
      onChange: function onChange(value) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updateInstance('backupDate', value);
      }
    }, (_10 = Backup === null || Backup === void 0 ? void 0 : Backup.weekConfig) === null || _10 === void 0 ? void 0 : _10.map(function (item) {
      return React__default.createElement(teaComponent.Checkbox, {
        key: item === null || item === void 0 ? void 0 : item.value,
        name: item === null || item === void 0 ? void 0 : item.value,
        className: "tea-mb-2n",
        style: {
          borderRadius: 5
        }
      }, i18n.t(item === null || item === void 0 ? void 0 : item.text));
    }))), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('备份时间点')
    }, React__default.createElement(teaComponent.Checkbox.Group, {
      value: (_11 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _11 === void 0 ? void 0 : _11.backupTime,
      onChange: function onChange(value) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updateInstance('backupTime', value);
      }
    }, (_12 = Backup === null || Backup === void 0 ? void 0 : Backup.hourConfig) === null || _12 === void 0 ? void 0 : _12.map(function (item) {
      return React__default.createElement(teaComponent.Checkbox, {
        key: item === null || item === void 0 ? void 0 : item.value,
        name: item === null || item === void 0 ? void 0 : item.value,
        className: "tea-mb-2n",
        style: {
          borderRadius: 5
        }
      }, i18n.t(item === null || item === void 0 ? void 0 : item.text));
    }))), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('备份保留时间(天)')
    }, React__default.createElement(ffComponent.FormPanel.InputNumber, {
      min: Backup.minReserveDay,
      max: Backup.maxReserveDay,
      value: (_13 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _13 === void 0 ? void 0 : _13.backupReserveDay,
      onChange: function onChange(value) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updateInstance('backupReserveDay', value);
      }
    }), React__default.createElement(teaComponent.Text, null, i18n.t('天后自动删除')))) : null, React__default.createElement(AffinityMapField, {
      plan: {
        name: 'nodeAffinity',
        type: SchemaType.Custom,
        label: i18n.t('节点调度策略')
      },
      onChange: function onChange(data) {
        var _a;
        actions === null || actions === void 0 ? void 0 : actions.create.updateInstance('nodeSchedule', tslib.__assign(tslib.__assign({}, (_a = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _a === void 0 ? void 0 : _a['nodeSchedule']), data));
      }
    }), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('设置实例参数'),
      required: false,
      validator: (_14 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _14 === void 0 ? void 0 : _14['isSetParams']
    }, React__default.createElement(teaComponent.Switch, {
      value: (_15 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _15 === void 0 ? void 0 : _15['isSetParams'],
      onChange: function onChange(e) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updateInstance('isSetParams', e);
      }
    })), isLoadingSchema && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t(''),
      text: true
    }, React__default.createElement(LoadingPanel, null)), !isLoadingSchema && (externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) && loadSchemaFailed && ((_16 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _16 === void 0 ? void 0 : _16['isSetParams']) && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t(''),
      text: true
    }, React__default.createElement(i18n.Trans, null, React__default.createElement("div", {
      style: {
        display: 'flex',
        alignItems: 'center'
      }
    }, React__default.createElement(teaComponent.Text, {
      theme: "danger",
      className: "tea-mr-2n",
      style: {
        minWidth: 100
      }
    }, "\u52A0\u8F7DSchema\u5931\u8D25:"), React__default.createElement(i18n.Slot, {
      content: ((_17 = [ErrorEnum.Code.RBACForbidden, ErrorEnum.Code.RBACForbidden403]) === null || _17 === void 0 ? void 0 : _17.includes((_19 = (_18 = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _18 === void 0 ? void 0 : _18.error) === null || _19 === void 0 ? void 0 : _19.code)) ? React__default.createElement(i18n.Trans, null, React__default.createElement(teaComponent.Text, {
        verticalAlign: "middle",
        theme: "danger"
      }, "\u6743\u9650\u4E0D\u8DB3\uFF0C\u8BF7\u8054\u7CFB\u96C6\u7FA4\u7BA1\u7406\u5458\u6DFB\u52A0\u6743\u9650\uFF1B\u82E5\u60A8\u672C\u8EAB\u662F\u96C6\u7FA4\u7BA1\u7406\u5458\uFF0C\u53EF\u76F4\u63A5"), React__default.createElement(teaComponent.Button, {
        type: "link",
        onClick: function onClick() {
          actions.base.getClusterAdminRole.getClusterAdminRole.start([]);
        },
        className: 'tea-mr-2n'
      }, "\u83B7\u53D6\u96C6\u7FA4admin\u89D2\u8272"), React__default.createElement(teaComponent.Button, {
        type: "link",
        onClick: function onClick() {
          var _a, _b;
          actions === null || actions === void 0 ? void 0 : actions.list.servicePlans.fetch();
          (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.resourceSchemas) === null || _b === void 0 ? void 0 : _b.fetch();
        }
      }, React__default.createElement(i18n.Slot, {
        content: i18n.t('重试')
      }))) : React__default.createElement(RetryPanel, {
        loadingText: (_21 = (_20 = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _20 === void 0 ? void 0 : _20.error) === null || _21 === void 0 ? void 0 : _21.message,
        action: function action() {
          var _a, _b;
          actions === null || actions === void 0 ? void 0 : actions.list.servicePlans.fetch();
          (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.resourceSchemas) === null || _b === void 0 ? void 0 : _b.fetch();
        }
      })
    })))), !isLoadingSchema && (externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) && !loadSchemaFailed && React__default.createElement(React__default.Fragment, null, openingVendor && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t(''),
      text: true
    }, React__default.createElement(RetryPanel, {
      loadingText: i18n.t('{{text}}', {
        text: "".concat((_22 = services === null || services === void 0 ? void 0 : services.selection) === null || _22 === void 0 ? void 0 : _22.name, "\u5F00\u542F\u4E2D...")
      }),
      loadingTextTheme: 'text',
      action: function action() {
        var _a, _b;
        actions === null || actions === void 0 ? void 0 : actions.list.servicePlans.fetch();
        (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.resourceSchemas) === null || _b === void 0 ? void 0 : _b.fetch();
      }
    })), !openingVendor && React__default.createElement(React__default.Fragment, null, versionFields === null || versionFields === void 0 ? void 0 : versionFields.map(function (item) {
      var _a, _b, _c, _d, _e;
      var _f = (_b = (_a = item.description) === null || _a === void 0 ? void 0 : _a.split('---')) !== null && _b !== void 0 ? _b : [],
        english = _f[0],
        chinese = _f[1];
      return React__default.createElement(ffComponent.FormPanel.Item, {
        key: item === null || item === void 0 ? void 0 : item.name,
        label: React__default.createElement(teaComponent.Text, {
          style: {
            display: 'flex',
            alignItems: 'center'
          }
        }, React__default.createElement(teaComponent.Text, null, i18n.t('版本')), React__default.createElement(teaComponent.Icon, {
          type: "info",
          tooltip: isI18n ? english : chinese
        }), !(item === null || item === void 0 ? void 0 : item.optional) && React__default.createElement(teaComponent.Text, {
          className: 'text-danger tea-pt-1n'
        }, "*")),
        validator: (_c = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _c === void 0 ? void 0 : _c[item === null || item === void 0 ? void 0 : item.name]
      }, React__default.createElement(ffComponent.FormPanel.Select, {
        placeholder: i18n.t('{{title}}', {
          title: '请选择' + (item === null || item === void 0 ? void 0 : item.label)
        }),
        options: (_d = item === null || item === void 0 ? void 0 : item.candidates) === null || _d === void 0 ? void 0 : _d.map(function (candidate) {
          return {
            value: candidate,
            text: candidate
          };
        }),
        value: (_e = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _e === void 0 ? void 0 : _e[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }));
    }), ((_23 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _23 === void 0 ? void 0 : _23['isSetParams']) && React__default.createElement(React__default.Fragment, null, instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.map(function (item, index) {
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
      var values = tslib.__assign({}, servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData);
      if (item.enabledCondition && values) {
        var _m = item === null || item === void 0 ? void 0 : item.enabledCondition.split('=='),
          conditionKey = _m[0],
          conditionValue = _m[1];
        var value = values[conditionKey];
        if (String(value) !== String(conditionValue)) {
          return null;
        }
      }
      var _o = (_b = (_a = item.description) === null || _a === void 0 ? void 0 : _a.split('---')) !== null && _b !== void 0 ? _b : [],
        english = _o[0],
        chinese = _o[1];
      return !hideSchema(item) ? React__default.createElement(ffComponent.FormPanel.Item, {
        label: React__default.createElement(teaComponent.Text, {
          style: {
            display: 'flex',
            alignItems: 'center'
          }
        }, React__default.createElement(teaComponent.Text, null, i18n.t('{{name}}', {
          name: prefixForSchema(item, servicename) + (item === null || item === void 0 ? void 0 : item.label) + suffixUnitForSchema(item)
        })), React__default.createElement(teaComponent.Icon, {
          type: "info",
          tooltip: isI18n ? english : chinese
        }), !(item === null || item === void 0 ? void 0 : item.optional) && React__default.createElement(teaComponent.Text, {
          className: 'text-danger tea-pt-1n'
        }, "*")),
        key: "".concat(item === null || item === void 0 ? void 0 : item.name).concat(index),
        validator: (_c = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _c === void 0 ? void 0 : _c[item === null || item === void 0 ? void 0 : item.name]
      }, getFormItemType(item) === FormItemType.Select && React__default.createElement(ffComponent.FormPanel.Select, {
        placeholder: i18n.t('{{title}}', {
          title: '请选择' + (item === null || item === void 0 ? void 0 : item.label)
        }),
        options: (_d = item === null || item === void 0 ? void 0 : item.candidates) === null || _d === void 0 ? void 0 : _d.map(function (candidate) {
          return {
            value: candidate,
            text: candidate
          };
        }),
        value: (_e = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _e === void 0 ? void 0 : _e[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.Switch && React__default.createElement(teaComponent.Switch, {
        value: (_f = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _f === void 0 ? void 0 : _f[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.Input && React__default.createElement(teaComponent.Input, {
        value: (_g = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _g === void 0 ? void 0 : _g[item === null || item === void 0 ? void 0 : item.name],
        placeholder: i18n.t('{{title}}', {
          title: '请输入' + (item === null || item === void 0 ? void 0 : item.label)
        }),
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.Paasword && React__default.createElement(InputPassword.InputPassword, {
        value: (_h = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _h === void 0 ? void 0 : _h[item === null || item === void 0 ? void 0 : item.name],
        placeholder: i18n.t('{{title}}', {
          title: '请输入' + (item === null || item === void 0 ? void 0 : item.label)
        }),
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, e);
        },
        rules: false
      }), getFormItemType(item) === FormItemType.InputNumber && React__default.createElement(teaComponent.InputNumber, {
        value: (_j = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _j === void 0 ? void 0 : _j[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, e);
        },
        min: SchemaInputNumOption === null || SchemaInputNumOption === void 0 ? void 0 : SchemaInputNumOption.min
      }), getFormItemType(item) === FormItemType.MapField && React__default.createElement(MapField, {
        plan: item,
        onChange: function onChange(_a) {
          var field = _a.field,
            value = _a.value;
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, JSON.stringify(value === null || value === void 0 ? void 0 : value.reduce(function (pre, cur) {
            var _a;
            return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.key] = cur === null || cur === void 0 ? void 0 : cur.value, _a));
          }, {})));
        }
      }), showUnitOptions(item) && React__default.createElement(ffComponent.FormPanel.Select, {
        size: "s",
        className: "tea-ml-2n",
        placeholder: i18n.t('{{title}}', {
          title: '请选择unit'
        }),
        options: getUnitOptions(item),
        value: (_l = (_k = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _k === void 0 ? void 0 : _k['unitMap']) === null || _l === void 0 ? void 0 : _l[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          var _a;
          var _b;
          actions === null || actions === void 0 ? void 0 : actions.create.updateInstance('unitMap', tslib.__assign(tslib.__assign({}, (_b = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _b === void 0 ? void 0 : _b['unitMap']), (_a = {}, _a[item === null || item === void 0 ? void 0 : item.name] = e, _a)));
        }
      })) : null;
    })))), React__default.createElement(GetRbacAdminDialog.Component, {
      model: getClusterAdminRole,
      action: actions.base.getClusterAdminRole,
      filter: {
        platform: platform,
        regionId: regionId,
        clusterId: (_24 = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _24 === void 0 ? void 0 : _24.clusterId
      },
      onSuccess: function onSuccess() {
        var _a, _b;
        (_a = actions.list) === null || _a === void 0 ? void 0 : _a.servicePlans.fetch();
        (_b = actions.list) === null || _b === void 0 ? void 0 : _b.resourceSchemas.fetch();
      }
    }), React__default.createElement(ffComponent.FormPanel.Footer, null, React__default.createElement(teaComponent.Button, {
      type: "primary",
      style: {
        marginRight: 10
      },
      onClick: _submit,
      loading: isSubmitting,
      disabled: isSubmitting || !(externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) || ((_25 = servicesInstanceEdit.formData) === null || _25 === void 0 ? void 0 : _25['isSetParams']) && loadSchemaFailed || openingVendor
    }, failed ? i18n.t('重试') : i18n.t('确定')), React__default.createElement(teaComponent.Button, {
      onClick: _cancel
    }, i18n.t('取消')), React__default.createElement(TipInfo, {
      isShow: failed,
      type: "error",
      isForm: true
    }, !((_26 = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _26 === void 0 ? void 0 : _26.timeBackup) ? getWorkflowError(createWorkflow) : (_29 = (_28 = (_27 = createWorkflow === null || createWorkflow === void 0 ? void 0 : createWorkflow.results) === null || _27 === void 0 ? void 0 : _27[0]) === null || _28 === void 0 ? void 0 : _28.error) === null || _29 === void 0 ? void 0 : _29.message)));
  };
  return React__default.createElement(React__default.Fragment, null, renderContent((_p = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _p === void 0 ? void 0 : _p.data));
}

function ServicePlanCreatePanel(props) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o;
  var actions = props.actions,
    _p = props.base,
    platform = _p.platform,
    route = _p.route,
    isI18n = _p.isI18n,
    userInfo = _p.userInfo,
    regionId = _p.regionId,
    getClusterAdminRole = _p.getClusterAdminRole,
    _q = props.list,
    services = _q.services,
    servicePlans = _q.servicePlans,
    createResourceWorkflow = _q.createResourceWorkflow,
    externalClusters = _q.externalClusters;
  var servicename = (route === null || route === void 0 ? void 0 : route.queries).servicename;
  var _r = React.useState({}),
    planSchemaUnitMap = _r[0],
    setPlanSchemaUnitMap = _r[1];
  React.useEffect(function () {
    var _a;
    var clusterId = (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId;
    if (actions.list.resourceSchemas && platform && servicename && clusterId) {
      actions.list.resourceSchemas.applyFilter({
        platform: platform,
        serviceName: servicename,
        clusterId: clusterId,
        regionId: regionId
      });
    }
  }, [actions.list.resourceSchemas, regionId, platform, servicename, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  React.useEffect(function () {
    var _a, _b;
    if (platform && ((_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.externalClusters) && (services === null || services === void 0 ? void 0 : services.selection)) {
      (_b = actions === null || actions === void 0 ? void 0 : actions.list) === null || _b === void 0 ? void 0 : _b.externalClusters.applyFilter({
        platform: platform,
        clusterIds: [],
        regionId: regionId
      });
    }
  }, [(_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.externalClusters, platform, services === null || services === void 0 ? void 0 : services.selection]);
  React.useEffect(function () {
    var _a;
    var clusterId = (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId;
    if (actions.list.servicePlans && platform && servicename && clusterId) {
      actions.list.servicePlans.applyFilter({
        platform: platform,
        serviceName: servicename,
        clusterId: clusterId,
        regionId: regionId,
        resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan
      });
    }
  }, [actions.list.servicePlans, regionId, platform, servicename, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  var servicesInstance = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.servicesInstance;
  });
  var servicesInstanceEdit = reactRedux.useSelector(function (state) {
    var _a;
    return (_a = state === null || state === void 0 ? void 0 : state.list) === null || _a === void 0 ? void 0 : _a.servicePlanEdit;
  });
  var isLoadingSchema = ((_b = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _b === void 0 ? void 0 : _b.fetched) && (externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) ? !((_c = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _c === void 0 ? void 0 : _c.fetched) : !((_d = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _d === void 0 ? void 0 : _d.fetched);
  var loadSchemaFailed = ((_e = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _e === void 0 ? void 0 : _e.fetched) && ((_f = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _f === void 0 ? void 0 : _f.error) && (!((_g = ["ResourceNotFound", 404]) === null || _g === void 0 ? void 0 : _g.includes((_j = (_h = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _h === void 0 ? void 0 : _h.error) === null || _j === void 0 ? void 0 : _j.code)) || !((_m = (_l = (_k = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _k === void 0 ? void 0 : _k.error) === null || _l === void 0 ? void 0 : _l.message) === null || _m === void 0 ? void 0 : _m.includes("404")));
  var openingVendor = React__default === null || React__default === void 0 ? void 0 : React__default.useMemo(function () {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
    return ((_b = (_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.clusters) === null || _b === void 0 ? void 0 : _b.includes((_c = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _c === void 0 ? void 0 : _c.clusterId)) && ((_d = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _d === void 0 ? void 0 : _d.fetched) && (((_f = (_e = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _e === void 0 ? void 0 : _e.error) === null || _f === void 0 ? void 0 : _f.code) === "ResourceNotFound" || ((_h = (_g = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _g === void 0 ? void 0 : _g.error) === null || _h === void 0 ? void 0 : _h.code) === 404 || ((_l = (_k = (_j = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _j === void 0 ? void 0 : _j.error) === null || _k === void 0 ? void 0 : _k.message) === null || _l === void 0 ? void 0 : _l.includes("404")));
  }, [servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object, services === null || services === void 0 ? void 0 : services.selection, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection]);
  var reduceCreateServiceResourceDataJson = function reduceCreateServiceResourceDataJson(data) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
    var _l = data === null || data === void 0 ? void 0 : data.formData,
      instanceName = _l.instanceName,
      description = _l.description,
      clusterId = _l.clusterId;
    //拼接中间件实例参数部分属性值
    var parameters = (_c = (_b = (_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.planSchema) === null || _c === void 0 ? void 0 : _c.reduce(function (pre, cur) {
      var _a;
      var _b, _c, _d, _e;
      // return Object.assign(pre,!!data?.formData?.[cur?.name] ? {[cur?.name]:data?.formData?.[cur?.name] + (data?.formData?.['unitMap']?.[cur?.name] ?? '')} : {});
      return Object.assign(pre, ((_b = data === null || data === void 0 ? void 0 : data.formData) === null || _b === void 0 ? void 0 : _b[cur === null || cur === void 0 ? void 0 : cur.name]) !== "" ? (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = formatPlanSchemaSubmitData(cur, (_c = data === null || data === void 0 ? void 0 : data.formData) === null || _c === void 0 ? void 0 : _c[cur === null || cur === void 0 ? void 0 : cur.name], (_e = (_d = data === null || data === void 0 ? void 0 : data.formData) === null || _d === void 0 ? void 0 : _d["unitMap"]) === null || _e === void 0 ? void 0 : _e[cur === null || cur === void 0 ? void 0 : cur.name]), _a) : {});
    }, {});
    var json = {
      id: uuid$1(),
      apiVersion: "infra.tce.io/v1",
      kind: (_d = ResourceTypeMap === null || ResourceTypeMap === void 0 ? void 0 : ResourceTypeMap[ResourceTypeEnum.ServicePlan]) === null || _d === void 0 ? void 0 : _d.resourceKind,
      metadata: {
        labels: {
          "ssm.infra.tce.io/cluster-id": clusterId,
          "ssm.infra.tce.io/owner": ServicePlanTypeEnum.Custom
        },
        annotations: {
          "ssm.infra.tce.io/creator": decodeURIComponent ? decodeURIComponent((_f = (_e = userInfo === null || userInfo === void 0 ? void 0 : userInfo.object) === null || _e === void 0 ? void 0 : _e.data) === null || _f === void 0 ? void 0 : _f.name) : (_h = (_g = userInfo === null || userInfo === void 0 ? void 0 : userInfo.object) === null || _g === void 0 ? void 0 : _g.data) === null || _h === void 0 ? void 0 : _h.name
        },
        name: instanceName,
        namespace: DefaultNamespace
      },
      spec: {
        serviceClass: ((_j = services === null || services === void 0 ? void 0 : services.selection) === null || _j === void 0 ? void 0 : _j.name) || ((_k = route === null || route === void 0 ? void 0 : route.queries) === null || _k === void 0 ? void 0 : _k.servicename),
        metadata: tslib.__assign({}, parameters)
      }
    };
    if (description) {
      json.spec["description"] = description;
    }
    return JSON.stringify(json);
  };
  var _cancel = function _cancel() {
    var _a, _b;
    router.navigate({
      sub: "list",
      tab: undefined
    }, {
      servicename: (_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.servicename,
      resourceType: (_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.resourceType,
      mode: "list"
    });
  };
  var _submit = function _submit() {
    var _a, _b, _c, _d;
    var formData = servicesInstanceEdit.formData;
    actions.create.validatePlan((_b = (_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.planSchema);
    if (_validatePlan(servicesInstanceEdit, (_d = (_c = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.planSchema)) {
      var params = {
        platform: platform,
        regionId: HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion,
        clusterId: formData === null || formData === void 0 ? void 0 : formData.clusterId,
        jsonData: reduceCreateServiceResourceDataJson(servicesInstanceEdit),
        resourceType: ResourceTypeEnum.ServicePlan
      };
      actions.create.createResource.start([params], HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion);
      actions.create.createResource.perform();
    }
  };
  var renderContent = function renderContent(data) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v;
    var instanceParamsFields = data === null || data === void 0 ? void 0 : data.planSchema;
    var isSubmitting = (createResourceWorkflow === null || createResourceWorkflow === void 0 ? void 0 : createResourceWorkflow.operationState) === ffRedux.OperationState.Performing;
    var failed = createResourceWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(createResourceWorkflow);
    return React__default.createElement(ffComponent.FormPanel, null, React__default.createElement(ffComponent.FormPanel.Item, {
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: "flex",
          alignItems: "center"
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t("中间件类型")), React__default.createElement(teaComponent.Text, {
        className: "text-danger tea-pt-1n"
      }, "*")),
      text: true
    }, React__default.createElement(ffComponent.FormPanel.Text, null, i18n.t("{{name}}", {
      name: (_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.servicename
    }))), React__default.createElement(ffComponent.FormPanel.Item, {
      after: React__default.createElement(teaComponent.Button, {
        icon: "refresh",
        onClick: function onClick() {
          actions.list.externalClusters.fetch();
        }
      }),
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: "flex",
          alignItems: "center"
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t("目标集群")), React__default.createElement(teaComponent.Text, {
        className: "text-danger tea-pt-1n"
      }, "*")),
      validator: (_b = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _b === void 0 ? void 0 : _b["clusterId"],
      select: {
        model: externalClusters,
        valueField: "clusterId",
        displayField: function displayField(record) {
          return (record === null || record === void 0 ? void 0 : record.clusterId) + " (" + (record === null || record === void 0 ? void 0 : record.clusterName) + ") ";
        },
        action: (_c = actions === null || actions === void 0 ? void 0 : actions.list) === null || _c === void 0 ? void 0 : _c.externalClusters,
        disabledField: function disabledField(record) {
          var _a, _b, _c;
          var disabled = !((_b = (_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.clusters) === null || _b === void 0 ? void 0 : _b.includes(record === null || record === void 0 ? void 0 : record.clusterId));
          return {
            disabled: disabled,
            tooltip: disabled ? i18n.t("{{tooltip}}", {
              tooltip: "".concat((_c = services === null || services === void 0 ? void 0 : services.selection) === null || _c === void 0 ? void 0 : _c.name, "\u5C1A\u672A\u5728\u8BE5\u76EE\u6807\u96C6\u7FA4\u5F00\u542F,\u8BF7\u60A8\u524D\u5F80\u5206\u5E03\u5F0F\u4E91\u4E2D\u5FC3\u6982\u89C8\u9875\u5F00\u542F")
            }) : null
          };
        }
      },
      message: React__default.createElement(React__default.Fragment, null, ((_d = services === null || services === void 0 ? void 0 : services.list) === null || _d === void 0 ? void 0 : _d.fetched) && !((_f = (_e = services.selection) === null || _e === void 0 ? void 0 : _e.clusters) === null || _f === void 0 ? void 0 : _f.length) ? React__default.createElement(RetryPanel, {
        loadingText: i18n.t("集群列表尚无已开启vendor的集群"),
        retryText: i18n.t("刷新重试"),
        action: actions.list.services.fetch
      }) : ((_g = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _g === void 0 ? void 0 : _g.fetched) && !((_j = (_h = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _h === void 0 ? void 0 : _h.data) === null || _j === void 0 ? void 0 : _j.recordCount) ? React__default.createElement(RetryPanel, {
        loadingText: i18n.t("集群列表尚无运行状态的目标集群"),
        retryText: i18n.t("前往确认"),
        action: function action() {
          router.navigate({}, {}, "/tdcc/cluster");
        }
      }) : null)
    }), React__default.createElement(ffComponent.FormPanel.Item, {
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: "flex",
          alignItems: "center"
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t("规格名称")), React__default.createElement(teaComponent.Text, {
        className: "text-danger tea-pt-1n"
      }, "*")),
      validator: (_k = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _k === void 0 ? void 0 : _k["instanceName"],
      input: {
        value: (_l = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _l === void 0 ? void 0 : _l["instanceName"],
        onChange: function onChange(e) {
          var _a;
          (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updatePlan("instanceName", e);
        },
        onBlur: function onBlur() {}
      }
    }), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t("描述"),
      validator: (_m = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _m === void 0 ? void 0 : _m["description"],
      input: {
        value: (_o = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _o === void 0 ? void 0 : _o["description"],
        onChange: function onChange(e) {
          var _a;
          (_a = actions === null || actions === void 0 ? void 0 : actions.create) === null || _a === void 0 ? void 0 : _a.updatePlan("description", e);
        },
        onBlur: function onBlur() {}
      }
    }), isLoadingSchema && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t(""),
      text: true
    }, React__default.createElement(LoadingPanel, null)), !isLoadingSchema && (externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) && loadSchemaFailed && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t(""),
      text: true
    }, React__default.createElement(i18n.Trans, null, React__default.createElement("div", {
      style: {
        display: "flex",
        alignItems: "center"
      }
    }, React__default.createElement(teaComponent.Text, {
      theme: "danger",
      className: "tea-mr-2n",
      style: {
        minWidth: 100
      }
    }, "\u52A0\u8F7DSchema\u5931\u8D25:"), React__default.createElement(i18n.Slot, {
      content: ((_p = [ErrorEnum.Code.RBACForbidden, ErrorEnum.Code.RBACForbidden403]) === null || _p === void 0 ? void 0 : _p.includes((_r = (_q = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _q === void 0 ? void 0 : _q.error) === null || _r === void 0 ? void 0 : _r.code)) ? React__default.createElement(i18n.Trans, null, React__default.createElement(teaComponent.Text, {
        verticalAlign: "middle",
        theme: "danger"
      }, "\u6743\u9650\u4E0D\u8DB3\uFF0C\u8BF7\u8054\u7CFB\u96C6\u7FA4\u7BA1\u7406\u5458\u6DFB\u52A0\u6743\u9650\uFF1B\u82E5\u60A8\u672C\u8EAB\u662F\u96C6\u7FA4\u7BA1\u7406\u5458\uFF0C\u53EF\u76F4\u63A5"), React__default.createElement(teaComponent.Button, {
        type: "link",
        onClick: function onClick() {
          actions.base.getClusterAdminRole.getClusterAdminRole.start([]);
        },
        className: "tea-mr-2n"
      }, "\u83B7\u53D6\u96C6\u7FA4admin\u89D2\u8272"), React__default.createElement(teaComponent.Button, {
        type: "link",
        onClick: function onClick() {
          var _a, _b;
          (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.resourceSchemas) === null || _b === void 0 ? void 0 : _b.fetch();
        }
      }, React__default.createElement(i18n.Slot, {
        content: i18n.t("重试")
      }))) : React__default.createElement(RetryPanel, {
        loadingText: (_t = (_s = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _s === void 0 ? void 0 : _s.error) === null || _t === void 0 ? void 0 : _t.message,
        action: function action() {
          var _a, _b;
          (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.resourceSchemas) === null || _b === void 0 ? void 0 : _b.fetch();
        }
      })
    })))), !isLoadingSchema && (externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) && !loadSchemaFailed && React__default.createElement(React__default.Fragment, null, openingVendor && React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t(""),
      text: true
    }, React__default.createElement(RetryPanel, {
      style: {
        minWidth: 170,
        width: 170
      },
      loadingText: i18n.t("{{text}}", {
        text: "".concat((_u = services === null || services === void 0 ? void 0 : services.selection) === null || _u === void 0 ? void 0 : _u.name, "\u5F00\u542F\u4E2D...")
      }),
      loadingTextTheme: "text",
      action: function action() {
        var _a, _b;
        (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.resourceSchemas) === null || _b === void 0 ? void 0 : _b.fetch();
      }
    })), !openingVendor && React__default.createElement(React__default.Fragment, null, instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.map(function (item, index) {
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l;
      var values = tslib.__assign({}, servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData);
      if (item.enabledCondition && values) {
        var _m = item === null || item === void 0 ? void 0 : item.enabledCondition.split("=="),
          conditionKey = _m[0],
          conditionValue = _m[1];
        var value = values[conditionKey];
        if (String(value) !== String(conditionValue)) {
          return null;
        }
      }
      var _o = (_b = (_a = item.description) === null || _a === void 0 ? void 0 : _a.split("---")) !== null && _b !== void 0 ? _b : [],
        english = _o[0],
        chinese = _o[1];
      return !hideSchema(item) ? React__default.createElement(ffComponent.FormPanel.Item, {
        label: React__default.createElement(teaComponent.Text, {
          style: {
            display: "flex",
            alignItems: "center"
          }
        }, React__default.createElement(teaComponent.Text, null, i18n.t("{{name}}", {
          name: prefixForSchema(item, servicename) + (item === null || item === void 0 ? void 0 : item.label) + suffixUnitForSchema(item)
        })), React__default.createElement(teaComponent.Icon, {
          type: "info",
          tooltip: isI18n ? english : chinese
        }), !(item === null || item === void 0 ? void 0 : item.optional) && React__default.createElement(teaComponent.Text, {
          className: "text-danger tea-pt-1n"
        }, "*")),
        key: "".concat(item === null || item === void 0 ? void 0 : item.name).concat(index),
        validator: (_c = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.validator) === null || _c === void 0 ? void 0 : _c[item === null || item === void 0 ? void 0 : item.name]
      }, getFormItemType(item) === FormItemType.Select && React__default.createElement(ffComponent.FormPanel.Select, {
        placeholder: i18n.t("{{title}}", {
          title: "请选择" + (item === null || item === void 0 ? void 0 : item.label)
        }),
        options: (_d = item === null || item === void 0 ? void 0 : item.candidates) === null || _d === void 0 ? void 0 : _d.map(function (candidate) {
          return {
            value: candidate,
            text: candidate
          };
        }),
        value: (_e = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _e === void 0 ? void 0 : _e[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.Switch && React__default.createElement(teaComponent.Switch, {
        value: (_f = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _f === void 0 ? void 0 : _f[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.Input && React__default.createElement(teaComponent.Input, {
        value: (_g = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _g === void 0 ? void 0 : _g[item === null || item === void 0 ? void 0 : item.name],
        placeholder: i18n.t("{{title}}", {
          title: "请输入" + (item === null || item === void 0 ? void 0 : item.label)
        }),
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.Paasword && React__default.createElement(InputPassword.InputPassword, {
        value: (_h = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _h === void 0 ? void 0 : _h[item === null || item === void 0 ? void 0 : item.name],
        placeholder: i18n.t("{{title}}", {
          title: "请输入" + (item === null || item === void 0 ? void 0 : item.label)
        }),
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
        }
      }), getFormItemType(item) === FormItemType.InputNumber && React__default.createElement(teaComponent.InputNumber, {
        value: (_j = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _j === void 0 ? void 0 : _j[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, e);
        },
        min: SchemaInputNumOption === null || SchemaInputNumOption === void 0 ? void 0 : SchemaInputNumOption.min
      }), getFormItemType(item) === FormItemType.MapField && React__default.createElement(MapField, {
        plan: item,
        onChange: function onChange(_a) {
          var field = _a.field,
            value = _a.value;
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan(item === null || item === void 0 ? void 0 : item.name, JSON.stringify(value === null || value === void 0 ? void 0 : value.reduce(function (pre, cur) {
            var _a;
            return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.key] = cur === null || cur === void 0 ? void 0 : cur.value, _a));
          }, {})));
        }
      }), showUnitOptions(item) && React__default.createElement(ffComponent.FormPanel.Select, {
        size: "s",
        className: "tea-ml-2n",
        placeholder: i18n.t("{{title}}", {
          title: "请选择unit"
        }),
        options: getUnitOptions(item),
        value: (_l = (_k = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _k === void 0 ? void 0 : _k["unitMap"]) === null || _l === void 0 ? void 0 : _l[item === null || item === void 0 ? void 0 : item.name],
        onChange: function onChange(e) {
          var _a;
          var _b;
          actions === null || actions === void 0 ? void 0 : actions.create.updatePlan("unitMap", tslib.__assign(tslib.__assign({}, (_b = servicesInstanceEdit === null || servicesInstanceEdit === void 0 ? void 0 : servicesInstanceEdit.formData) === null || _b === void 0 ? void 0 : _b["unitMap"]), (_a = {}, _a[item === null || item === void 0 ? void 0 : item.name] = e, _a)));
        }
      })) : null;
    }))), React__default.createElement(GetRbacAdminDialog.Component, {
      model: getClusterAdminRole,
      action: actions.base.getClusterAdminRole,
      filter: {
        platform: platform,
        regionId: regionId,
        clusterId: (_v = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _v === void 0 ? void 0 : _v.clusterId
      },
      onSuccess: function onSuccess() {
        var _a, _b;
        (_a = actions.list) === null || _a === void 0 ? void 0 : _a.servicePlans.fetch();
        (_b = actions.list) === null || _b === void 0 ? void 0 : _b.resourceSchemas.fetch();
      }
    }), React__default.createElement(ffComponent.FormPanel.Footer, null, React__default.createElement(teaComponent.Button, {
      type: "primary",
      className: "tea-mr-2n",
      onClick: _submit,
      loading: isSubmitting,
      disabled: isSubmitting || !(externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) || loadSchemaFailed || openingVendor
    }, failed ? i18n.t("重试") : i18n.t("确定")), React__default.createElement(teaComponent.Button, {
      onClick: _cancel
    }, i18n.t("取消")), React__default.createElement(TipInfo, {
      isShow: failed,
      type: "error",
      isForm: true
    }, getWorkflowError(createResourceWorkflow))));
  };
  return React__default.createElement(React__default.Fragment, null, renderContent((_o = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _o === void 0 ? void 0 : _o.data));
}

function ServiceCreate(props) {
  var _a = props.base,
    selectedTab = _a.selectedTab,
    route = _a.route;
  var tab = (router === null || router === void 0 ? void 0 : router.resolve(route)).tab;
  if (tab === ResourceTypeEnum.ServiceResource) {
    return React__default.createElement(ServiceInstanceCreatePanel, tslib.__assign({}, props));
  } else if (tab === ResourceTypeEnum.ServicePlan) {
    return React__default.createElement(ServicePlanCreatePanel, tslib.__assign({}, props));
  } else {
    return React__default.createElement(React__default.Fragment, null, React__default.createElement(LoadingPanel, null));
  }
}

function ResourceDeleteDialog(props) {
  var _this = this;
  var _a;
  var _b = props.base,
    platform = _b.platform,
    route = _b.route,
    regionId = _b.regionId,
    _c = props.list,
    deleteResourceSelection = _c.deleteResourceSelection,
    deleteResourceWorkflow = _c.deleteResourceWorkflow,
    serviceResources = _c.serviceResources,
    actions = props.actions;
  var finalResourceInfo = deleteResourceSelection === null || deleteResourceSelection === void 0 ? void 0 : deleteResourceSelection[0];
  var loading = deleteResourceWorkflow.operationState === ffRedux.OperationState.Performing;
  var failed = deleteResourceWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(deleteResourceWorkflow);
  if (!(deleteResourceSelection === null || deleteResourceSelection === void 0 ? void 0 : deleteResourceSelection.length)) {
    return React__default.createElement("noscript", null);
  }
  var _submit = function _submit() {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var basePrams, externalParams, resource;
      var _a, _b, _c, _d;
      return tslib.__generator(this, function (_e) {
        basePrams = {
          id: uuid(),
          resourceInfos: deleteResourceSelection,
          regionId: regionId,
          platform: platform,
          clusterId: Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, finalResourceInfo, route)
        };
        externalParams = {};
        if ((finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.kind) === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceOpsBackup)) {
          externalParams = ((_a = finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.spec) === null || _a === void 0 ? void 0 : _a.trigger) === (BackupTypeNum === null || BackupTypeNum === void 0 ? void 0 : BackupTypeNum.Schedule) ? {
            resourceInfos: deleteResourceSelection
          } : {
            resourceInfos: (_b = []) === null || _b === void 0 ? void 0 : _b.concat(deleteResourceSelection, [{
              metadata: {
                name: (_d = (_c = finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.metadata) === null || _c === void 0 ? void 0 : _c.labels) === null || _d === void 0 ? void 0 : _d['ssm.infra.tce.io/ops-plan'],
                namespace: SystemNamespace
              },
              kind: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup
            }])
          };
        } else {
          externalParams = {
            resourceInfos: deleteResourceSelection,
            namespace: (finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.kind) !== (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding) ? DefaultNamespace : undefined
          };
        }
        resource = tslib.__assign(tslib.__assign({}, basePrams), externalParams);
        actions.list.deleteResource.start([resource], HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion);
        actions.list.deleteResource.perform();
        return [2 /*return*/];
      });
    });
  };

  var _cancel = function _cancel() {
    actions.list.selectDeleteResources([]);
    if (deleteResourceWorkflow.operationState === ffRedux.OperationState.Done) {
      actions.list.deleteResource.reset();
    }
    if (deleteResourceWorkflow.operationState === ffRedux.OperationState.Started) {
      actions.list.deleteResource.cancel();
    }
  };
  var caption = (finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.kind) === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding) ? i18n.t('解绑确认') : i18n.t('删除资源');
  var resourceIns = (_a = deleteResourceSelection === null || deleteResourceSelection === void 0 ? void 0 : deleteResourceSelection.map(function (item) {
    var _a;
    return (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name;
  })) === null || _a === void 0 ? void 0 : _a.join(',');
  var content = (finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.kind) === (ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceBinding) ? React__default.createElement(React__default.Fragment, null, i18n.t('您确定要解绑{{headTitle}}:', {
    headTitle: finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.kind
  }), React__default.createElement("strong", null, i18n.t('{{resourceIns}}', {
    resourceIns: resourceIns
  })), i18n.t('吗,解除绑定关系后，应用侧可能无法继续使用本服务实例。')) : React__default.createElement(React__default.Fragment, null, i18n.t('您确定要删除{{headTitle}}:', {
    headTitle: finalResourceInfo === null || finalResourceInfo === void 0 ? void 0 : finalResourceInfo.kind
  }), React__default.createElement("strong", null, i18n.t('{{resourceIns}}', {
    resourceIns: resourceIns
  })), i18n.t('吗?'));
  return React__default.createElement(teaComponent.Modal, {
    visible: true,
    caption: caption,
    onClose: _cancel
  }, React__default.createElement(teaComponent.Modal.Body, null, React__default.createElement("div", {
    style: {
      fontSize: '14px',
      lineHeight: '20px'
    }
  }, React__default.createElement("p", {
    style: {
      wordWrap: 'break-word'
    }
  }, content)), React__default.createElement(TipInfo, {
    isShow: failed,
    type: "error",
    isForm: true
  }, getWorkflowError(deleteResourceWorkflow))), React__default.createElement(teaComponent.Modal.Footer, null, React__default.createElement(teaComponent.Button, {
    onClick: _submit,
    type: failed ? 'pay' : 'primary',
    loading: loading
  }, failed ? i18n.t('重试') : i18n.t('确定')), React__default.createElement(teaComponent.Button, {
    onClick: _cancel
  }, i18n.t('取消'))));
}

function InstanceBaseDetail(props) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v, _w, _x, _y, _z, _0, _1, _2, _3, _4, _5, _6, _7, _8, _9;
  var _10 = props.base,
    platform = _10.platform,
    route = _10.route,
    regionId = _10.regionId,
    servicesInstance = props.list.servicesInstance,
    _11 = props.detail,
    resourceDetail = _11.resourceDetail,
    openConsoleWorkflow = _11.openConsoleWorkflow,
    actions = props.actions;
  var isBaseLoading = !(resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object.fetched) || (resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object.loading);
  var resource = (_b = (_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.resource;
  var openConsoleLoading = (openConsoleWorkflow === null || openConsoleWorkflow === void 0 ? void 0 : openConsoleWorkflow.operationState) === ffRedux.OperationState.Performing;
  var failed = openConsoleWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(openConsoleWorkflow);
  var instanceParamsFields = (_d = (_c = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.instanceSchema;
  var domain = ((_f = (_e = resource === null || resource === void 0 ? void 0 : resource.status) === null || _e === void 0 ? void 0 : _e.metadata) === null || _f === void 0 ? void 0 : _f.host) || i18n.t('未设置');
  var vip = ((_h = (_g = resource === null || resource === void 0 ? void 0 : resource.status) === null || _g === void 0 ? void 0 : _g.metadata) === null || _h === void 0 ? void 0 : _h.ip) || ((_k = (_j = resource === null || resource === void 0 ? void 0 : resource.status) === null || _j === void 0 ? void 0 : _j.metadata) === null || _k === void 0 ? void 0 : _k.port) ? "".concat((_m = (_l = resource === null || resource === void 0 ? void 0 : resource.status) === null || _l === void 0 ? void 0 : _l.metadata) === null || _m === void 0 ? void 0 : _m.ip, ":").concat((_p = (_o = resource === null || resource === void 0 ? void 0 : resource.status) === null || _o === void 0 ? void 0 : _o.metadata) === null || _p === void 0 ? void 0 : _p.port) : i18n.t('未设置');
  var rsIps = getRsIps((_q = resource === null || resource === void 0 ? void 0 : resource.status) === null || _q === void 0 ? void 0 : _q.deploy);

  var state = ((_r = resource === null || resource === void 0 ? void 0 : resource.status) === null || _r === void 0 ? void 0 : _r.state) || '-';
  var className = ServiceInstanceMap === null || ServiceInstanceMap === void 0 ? void 0 : ServiceInstanceMap[((_s = resource === null || resource === void 0 ? void 0 : resource.status) === null || _s === void 0 ? void 0 : _s.state) || (ServiceInstanceStatusEnum === null || ServiceInstanceStatusEnum === void 0 ? void 0 : ServiceInstanceStatusEnum.Unknown)].className;
  var isDeleting = showResourceDeleteLoading(resource, [resource]);
  return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Row, {
    style: {
      width: '100%',
      margin: 0
    }
  }, React__default.createElement(teaComponent.Col, {
    span: 16,
    style: {
      padding: 0
    },
    className: 'tea-pr-5n'
  }, React__default.createElement("div", null, React__default.createElement(ffComponent.FormPanel, {
    title: '基本信息',
    style: {
      width: '100%'
    }
  }, isBaseLoading && React__default.createElement(LoadingPanel, null), !isBaseLoading && React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Row, null, React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n',
    overflow: true
  }, i18n.t('名称:')), React__default.createElement(teaComponent.Text, null, (_u = (_t = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _t === void 0 ? void 0 : _t.name) !== null && _u !== void 0 ? _u : '-')), React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n'
  }, i18n.t('ID:')), React__default.createElement(teaComponent.Text, null, (_w = (_v = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _v === void 0 ? void 0 : _v.externalID) !== null && _w !== void 0 ? _w : '-')), React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n'
  }, i18n.t('状态:')), React__default.createElement(React__default.Fragment, null, isDeleting && React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Icon, {
    type: 'loading'
  }), i18n.t('删除中')), !isDeleting && React__default.createElement(teaComponent.Text, {
    className: className + " tea-mr-1n"
  }, state), !isDeleting && !!((_y = (_x = resource === null || resource === void 0 ? void 0 : resource.status) === null || _x === void 0 ? void 0 : _x.conditions) === null || _y === void 0 ? void 0 : _y.length) && React__default.createElement(teaComponent.Bubble, {
    content: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.List, {
      type: "number",
      style: {
        width: '100%'
      }
    }, (_0 = (_z = resource === null || resource === void 0 ? void 0 : resource.status) === null || _z === void 0 ? void 0 : _z.conditions) === null || _0 === void 0 ? void 0 : _0.map(function (item) {
      return React__default.createElement(teaComponent.List.Item, {
        key: item === null || item === void 0 ? void 0 : item.type
      }, React__default.createElement(teaComponent.Text, {
        className: 'tea-mr-2n'
      }, " ", "".concat(item === null || item === void 0 ? void 0 : item.type, " : ").concat((item === null || item === void 0 ? void 0 : item.reason) || (item === null || item === void 0 ? void 0 : item.message))), React__default.createElement(teaComponent.Icon, {
        type: (item === null || item === void 0 ? void 0 : item.status) === 'True' ? 'success' : 'error'
      }));
    })))
  }, React__default.createElement(teaComponent.Icon, {
    type: "info",
    className: "tea-mr-2n"
  })))), React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n'
  }, i18n.t('版本:')), React__default.createElement(teaComponent.Text, null, (_3 = (_2 = (_1 = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _1 === void 0 ? void 0 : _1.parameters) === null || _2 === void 0 ? void 0 : _2.version) !== null && _3 !== void 0 ? _3 : '-'))), React__default.createElement(teaComponent.Row, null, React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n'
  }, i18n.t('规格:')), React__default.createElement(teaComponent.Button, {
    type: 'link',
    onClick: function onClick() {
      var _a;
      router.navigate({
        sub: 'list',
        tab: undefined
      }, {
        servicename: (_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.servicename,
        resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServicePlan,
        mode: 'list'
      });
    }
  }, (_5 = (_4 = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _4 === void 0 ? void 0 : _4.servicePlan) !== null && _5 !== void 0 ? _5 : '-')), React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n',
    overflow: true
  }, i18n.t('创建人:')), React__default.createElement(teaComponent.Text, {
    overflow: true,
    tooltip: true
  }, Util.getCreator(platform, resource))), React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n'
  }, i18n.t('命名空间:')), React__default.createElement(teaComponent.Text, {
    overflow: true,
    tooltip: true
  }, (_7 = (_6 = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _6 === void 0 ? void 0 : _6.namespace) !== null && _7 !== void 0 ? _7 : '-')), React__default.createElement(teaComponent.Col, {
    span: 6
  }, React__default.createElement(teaComponent.Text, {
    className: 'tea-mr-2n'
  }, i18n.t('创建时间:')), React__default.createElement(teaComponent.Text, null, ((_8 = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _8 === void 0 ? void 0 : _8.creationTimestamp) ? dateFormatter(new Date((_9 = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _9 === void 0 ? void 0 : _9.creationTimestamp), 'YYYY-MM-DD HH:mm:ss') : '-'))))), React__default.createElement(ffComponent.FormPanel, {
    title: '服务入口',
    style: {
      width: '100%'
    }
  }, isBaseLoading && React__default.createElement(LoadingPanel, null), !isBaseLoading && React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Row, null, React__default.createElement(teaComponent.Col, {
    span: 4
  }, React__default.createElement(teaComponent.Text, null, i18n.t('域名:'))), React__default.createElement(teaComponent.Col, null, React__default.createElement(teaComponent.Text, null, i18n.t('{{domain}}', {
    domain: domain
  })), domain && React__default.createElement(teaComponent.Copy, {
    text: domain
  }))), React__default.createElement(teaComponent.Row, null, React__default.createElement(teaComponent.Col, {
    span: 4
  }, React__default.createElement(teaComponent.Text, null, i18n.t('VIP+端口:'))), React__default.createElement(teaComponent.Col, null, React__default.createElement(teaComponent.Text, null, i18n.t('{{vip}}', {
    vip: vip
  })), vip && React__default.createElement(teaComponent.Copy, {
    text: vip
  }))), React__default.createElement(teaComponent.Row, null, React__default.createElement(teaComponent.Col, {
    span: 4
  }, React__default.createElement(teaComponent.Text, null, i18n.t('RS IP:'))), !(rsIps === null || rsIps === void 0 ? void 0 : rsIps.length) ? React__default.createElement(teaComponent.Col, {
    span: 20
  }, React__default.createElement(teaComponent.Text, null, i18n.t('未设置'))) : null), rsIps === null || rsIps === void 0 ? void 0 : rsIps.map(function (ip, i) {
    return React__default.createElement(teaComponent.Row, null, React__default.createElement(teaComponent.Col, {
      span: 4
    }), React__default.createElement(teaComponent.Col, {
      span: 20,
      key: "rowkey".concat(i)
    }, React__default.createElement(teaComponent.Text, null, ip), ip && React__default.createElement(teaComponent.Copy, {
      text: ip
    })));
  }))))), React__default.createElement(teaComponent.Col, {
    span: 8,
    style: {
      padding: 0
    }
  }, React__default.createElement(ffComponent.FormPanel, {
    title: '实例参数',
    style: {
      padding: 0
    }
  }, isBaseLoading && React__default.createElement(LoadingPanel, null), !isBaseLoading && (instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.map(function (item) {
    return React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('{{name}}', {
        name: (item === null || item === void 0 ? void 0 : item.label) + ':'
      }),
      text: true,
      key: item === null || item === void 0 ? void 0 : item.name
    }, React__default.createElement(teaComponent.Text, {
      style: {
        textAlign: 'right'
      }
    }, (item === null || item === void 0 ? void 0 : item.value) || '-'));
  }))))));
}

function InstanceMonitorPanel() {
  return React__default.createElement(React__default.Fragment, null, "InstanceMonitorPanel comp...");
}

function CreateServiceBindingDialog(props) {
  var _this = this;
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q;
  var _r = props.detail,
    showCreateResourceDialog = _r.showCreateResourceDialog,
    serviceBindingEdit = _r.serviceBindingEdit,
    namespaces = _r.namespaces,
    serviceInstanceSchema = _r.serviceInstanceSchema,
    resourceDetail = _r.resourceDetail,
    _s = props.base,
    platform = _s.platform,
    route = _s.route,
    isI18n = _s.isI18n,
    _t = props.list,
    createResourceWorkflow = _t.createResourceWorkflow,
    serviceResources = _t.serviceResources,
    services = _t.services,
    actions = props.actions,
    _u = props.mode;
  var loading = (createResourceWorkflow === null || createResourceWorkflow === void 0 ? void 0 : createResourceWorkflow.operationState) === ffRedux.OperationState.Performing;
  var failed = createResourceWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(createResourceWorkflow);
  var instanceParamsFields = (_b = (_a = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.bindingCreateParameterSchema;
  var updateFormData = (_c = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _c === void 0 ? void 0 : _c.updateServiceBinding;
  var reduceServiceBindingJson = function reduceServiceBindingJson(data) {
    var _a, _b, _c, _d, _e, _f;
    var _g = data === null || data === void 0 ? void 0 : data.formData,
      name = _g.name,
      namespace = _g.namespace;
    //拼接中间件实例参数部分属性值
    var parameters = (_c = (_b = (_a = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.bindingCreateParameterSchema) === null || _c === void 0 ? void 0 : _c.reduce(function (pre, cur) {
      var _a;
      var _b, _c, _d, _e;
      return Object.assign(pre, ((_b = data === null || data === void 0 ? void 0 : data.formData) === null || _b === void 0 ? void 0 : _b[cur === null || cur === void 0 ? void 0 : cur.name]) !== '' ? (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.name] = formatPlanSchemaSubmitData(cur, (_c = data === null || data === void 0 ? void 0 : data.formData) === null || _c === void 0 ? void 0 : _c[cur === null || cur === void 0 ? void 0 : cur.name], (_e = (_d = data === null || data === void 0 ? void 0 : data.formData) === null || _d === void 0 ? void 0 : _d['unitMap']) === null || _e === void 0 ? void 0 : _e[cur === null || cur === void 0 ? void 0 : cur.name]), _a) : {});
    }, {});
    var instancename = (route === null || route === void 0 ? void 0 : route.queries).instancename;
    var jsonData = {
      id: uuid(),
      apiVersion: 'infra.tce.io/v1',
      kind: (_d = ResourceTypeMap === null || ResourceTypeMap === void 0 ? void 0 : ResourceTypeMap[ResourceTypeEnum.ServiceBinding]) === null || _d === void 0 ? void 0 : _d.resourceKind,
      metadata: {
        labels: {
          "ssm.infra.tce.io/instance-id": Util === null || Util === void 0 ? void 0 : Util.getInstanceId(platform, (_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.instancename, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection)
        },
        name: name,
        namespace: namespace
      },
      spec: {
        serviceClass: (_f = route === null || route === void 0 ? void 0 : route.queries) === null || _f === void 0 ? void 0 : _f.servicename,
        parameters: parameters,
        instanceRef: {
          name: instancename,
          namespace: DefaultNamespace
        }
      }
    };
    return JSON.stringify(jsonData);
  };
  var _submit = function _submit() {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var instance, result, regionId, params;
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
      return tslib.__generator(this, function (_l) {
        switch (_l.label) {
          case 0:
            if (!NotSupportBindingSchemaVendors.includes((_a = services === null || services === void 0 ? void 0 : services.selection) === null || _a === void 0 ? void 0 : _a.name) && !(instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.length)) {
              teaComponent.message === null || teaComponent.message === void 0 ? void 0 : teaComponent.message.warning({
                content: i18n.t('Schema数据异常,请刷新重试')
              });
              return [2 /*return*/];
            }

            return [4 /*yield*/, (_b = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _b === void 0 ? void 0 : _b.validateAllServiceBinding()];
          case 1:
            _l.sent();
            instance = (_e = (_d = (_c = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.list) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.records) === null || _e === void 0 ? void 0 : _e.find(function (item) {
              var _a, _b;
              return ((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) === ((_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.instancename);
            });
            result = ServiceBinding === null || ServiceBinding === void 0 ? void 0 : ServiceBinding._validateAll(serviceBindingEdit, instanceParamsFields, (_f = services === null || services === void 0 ? void 0 : services.selection) === null || _f === void 0 ? void 0 : _f.name, instance);
            if (result === null || result === void 0 ? void 0 : result.valid) {
              regionId = HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
              params = {
                platform: platform,
                regionId: regionId,
                clusterId: (_g = route === null || route === void 0 ? void 0 : route.queries) === null || _g === void 0 ? void 0 : _g.clusterid,
                jsonData: reduceServiceBindingJson(serviceBindingEdit),
                resourceType: ResourceTypeEnum.ServiceBinding,
                specificOperate: CreateSpecificOperatorEnum === null || CreateSpecificOperatorEnum === void 0 ? void 0 : CreateSpecificOperatorEnum.CreateServiceBinding,
                namespace: (_h = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _h === void 0 ? void 0 : _h.namespace
              };
              (_j = actions === null || actions === void 0 ? void 0 : actions.create) === null || _j === void 0 ? void 0 : _j.createResource.start([params], regionId);
              (_k = actions === null || actions === void 0 ? void 0 : actions.create) === null || _k === void 0 ? void 0 : _k.createResource.perform();
            }
            return [2 /*return*/];
        }
      });
    });
  };

  var _cancel = function _cancel() {
    var _a;
    actions.detail.showCreateResourceDialog(false);
    (_a = actions === null || actions === void 0 ? void 0 : actions.create.createResource) === null || _a === void 0 ? void 0 : _a.reset();
  };
  var _renderButtons = function _renderButtons() {
    var buttons = [{
      handleFunc: _submit,
      text: failed ? i18n.t('重试') : i18n.t('确定'),
      type: 'primary',
      loading: loading
    }, {
      handleFunc: _cancel,
      text: i18n.t('取消'),
      loading: false
    }];
    return buttons === null || buttons === void 0 ? void 0 : buttons.map(function (item, index) {
      return React__default.createElement(teaComponent.Button, {
        loading: item === null || item === void 0 ? void 0 : item.loading,
        key: index,
        type: item === null || item === void 0 ? void 0 : item.type,
        onClick: item === null || item === void 0 ? void 0 : item.handleFunc
      }, item === null || item === void 0 ? void 0 : item.text);
    });
  };
  var isLoadingSchema = !NotSupportBindingSchemaVendors.includes((_d = services === null || services === void 0 ? void 0 : services.selection) === null || _d === void 0 ? void 0 : _d.name) ? !((_e = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _e === void 0 ? void 0 : _e.fetched) || ((_f = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _f === void 0 ? void 0 : _f.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Fetching) : false;
  var loadSchemaFailed = !NotSupportBindingSchemaVendors.includes((_g = services === null || services === void 0 ? void 0 : services.selection) === null || _g === void 0 ? void 0 : _g.name) ? ((_h = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _h === void 0 ? void 0 : _h.error) || ((_j = serviceInstanceSchema === null || serviceInstanceSchema === void 0 ? void 0 : serviceInstanceSchema.object) === null || _j === void 0 ? void 0 : _j.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Failed) : false;
  return React__default.createElement(teaComponent.Modal, {
    visible: showCreateResourceDialog,
    caption: i18n.t('新建绑定'),
    onClose: function onClose() {
      actions.detail.showCreateResourceDialog(false);
    },
    size: 'm'
  }, React__default.createElement(teaComponent.Modal.Body, null, React__default.createElement(ffComponent.FormPanel, {
    isNeedCard: false
  }, React__default.createElement(ffComponent.FormPanel.Item, {
    label: React__default.createElement(teaComponent.Text, {
      style: {
        display: 'flex',
        alignItems: 'center'
      }
    }, React__default.createElement(teaComponent.Text, null, i18n.t('名称')), React__default.createElement(teaComponent.Text, {
      className: 'text-danger tea-pt-1n'
    }, "*")),
    validator: (_k = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.validator) === null || _k === void 0 ? void 0 : _k.name,
    input: {
      value: (_l = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _l === void 0 ? void 0 : _l.name,
      placeholder: i18n.t('请输入名称'),
      onChange: function onChange(value) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.updateServiceBinding('name', value);
      }
    }
  }), React__default.createElement(ffComponent.FormPanel.Item, {
    label: React__default.createElement(teaComponent.Text, {
      style: {
        display: 'flex',
        alignItems: 'center'
      }
    }, React__default.createElement(teaComponent.Text, null, i18n.t('命名空间')), React__default.createElement(teaComponent.Text, {
      className: 'text-danger tea-pt-1n'
    }, "*")),
    validator: (_m = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.validator) === null || _m === void 0 ? void 0 : _m.namespace,
    select: {
      model: namespaces,
      action: (_o = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _o === void 0 ? void 0 : _o.namespaces,
      displayField: function displayField(record) {
        var _a;
        return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name;
      },
      valueField: function valueField(record) {
        var _a;
        return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name;
      },
      onChange: function onChange(value) {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.updateServiceBinding('namespace', value);
      }
    }
  }), isLoadingSchema && React__default.createElement(LoadingPanel, {
    text: i18n.t('Schema')
  }), !isLoadingSchema && loadSchemaFailed && React__default.createElement(RetryPanel, {
    style: {
      minWidth: 150
    },
    loadingText: i18n.t('Schema加载失败'),
    action: (_q = (_p = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _p === void 0 ? void 0 : _p.serviceInstanceSchema) === null || _q === void 0 ? void 0 : _q.fetch
  }), !isLoadingSchema && !loadSchemaFailed && (instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.map(function (item) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m;
    var values = tslib.__assign({}, serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData);
    if (item.enabledCondition && values) {
      var _o = item === null || item === void 0 ? void 0 : item.enabledCondition.split('=='),
        conditionKey = _o[0],
        conditionValue = _o[1];
      var value = values[conditionKey];
      if (String(value) !== String(conditionValue)) {
        return null;
      }
    }
    var _p = (_b = (_a = item.description) === null || _a === void 0 ? void 0 : _a.split('---')) !== null && _b !== void 0 ? _b : [],
      english = _p[0],
      chinese = _p[1];
    return !hideSchema(item) ? React__default.createElement(ffComponent.FormPanel.Item, {
      label: React__default.createElement(teaComponent.Text, {
        style: {
          display: 'flex',
          alignItems: 'center'
        }
      }, React__default.createElement(teaComponent.Text, null, i18n.t('{{name}}', {
        name: prefixForSchema(item, (_c = route === null || route === void 0 ? void 0 : route.queries) === null || _c === void 0 ? void 0 : _c.servicename) + (item === null || item === void 0 ? void 0 : item.label) + suffixUnitForSchema(item)
      })), React__default.createElement(teaComponent.Icon, {
        type: "info",
        tooltip: isI18n ? english : chinese
      }), !(item === null || item === void 0 ? void 0 : item.optional) && React__default.createElement(teaComponent.Text, {
        className: 'text-danger tea-pt-1n'
      }, "*")),
      key: item === null || item === void 0 ? void 0 : item.name,
      validator: (_d = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.validator) === null || _d === void 0 ? void 0 : _d[item === null || item === void 0 ? void 0 : item.name]
    }, getFormItemType(item) === FormItemType.Select && React__default.createElement(ffComponent.FormPanel.Select, {
      placeholder: i18n.t('{{title}}', {
        title: '请选择' + (item === null || item === void 0 ? void 0 : item.label)
      }),
      options: (_e = item === null || item === void 0 ? void 0 : item.candidates) === null || _e === void 0 ? void 0 : _e.map(function (candidate) {
        return {
          value: candidate,
          text: candidate
        };
      }),
      value: (_f = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _f === void 0 ? void 0 : _f[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        updateFormData(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.Switch && React__default.createElement(teaComponent.Switch, {
      value: (_g = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _g === void 0 ? void 0 : _g[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        updateFormData(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.Input && React__default.createElement(teaComponent.Input, {
      value: (_h = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _h === void 0 ? void 0 : _h[item === null || item === void 0 ? void 0 : item.name],
      placeholder: i18n.t('{{title}}', {
        title: '请输入' + (item === null || item === void 0 ? void 0 : item.label)
      }),
      onChange: function onChange(e) {
        updateFormData(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.Paasword && React__default.createElement(InputPassword.InputPassword, {
      value: (_j = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _j === void 0 ? void 0 : _j[item === null || item === void 0 ? void 0 : item.name],
      placeholder: i18n.t('{{title}}', {
        title: '请输入' + (item === null || item === void 0 ? void 0 : item.label)
      }),
      onChange: function onChange(e) {
        updateFormData(item === null || item === void 0 ? void 0 : item.name, e);
      }
    }), getFormItemType(item) === FormItemType.InputNumber && React__default.createElement(teaComponent.InputNumber, {
      value: (_k = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _k === void 0 ? void 0 : _k[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        updateFormData(item === null || item === void 0 ? void 0 : item.name, e);
      },
      min: SchemaInputNumOption === null || SchemaInputNumOption === void 0 ? void 0 : SchemaInputNumOption.min
    }), getFormItemType(item) === FormItemType.MapField && React__default.createElement(MapField, {
      plan: item,
      onChange: function onChange(_a) {
        var field = _a.field,
          value = _a.value;
        actions === null || actions === void 0 ? void 0 : actions.create.updateInstance(item === null || item === void 0 ? void 0 : item.name, JSON.stringify(value === null || value === void 0 ? void 0 : value.reduce(function (pre, cur) {
          var _a;
          return tslib.__assign(tslib.__assign({}, pre), (_a = {}, _a[cur === null || cur === void 0 ? void 0 : cur.key] = cur === null || cur === void 0 ? void 0 : cur.value, _a));
        }, {})));
      }
    }), showUnitOptions(item) && React__default.createElement(ffComponent.FormPanel.Select, {
      size: 's',
      className: 'tea-ml-2n',
      placeholder: i18n.t('{{title}}', {
        title: '请选择unit'
      }),
      options: getUnitOptions(item),
      value: (_m = (_l = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _l === void 0 ? void 0 : _l['unitMap']) === null || _m === void 0 ? void 0 : _m[item === null || item === void 0 ? void 0 : item.name],
      onChange: function onChange(e) {
        var _a;
        var _b;
        actions === null || actions === void 0 ? void 0 : actions.create.updateInstance('unitMap', tslib.__assign(tslib.__assign({}, (_b = serviceBindingEdit === null || serviceBindingEdit === void 0 ? void 0 : serviceBindingEdit.formData) === null || _b === void 0 ? void 0 : _b['unitMap']), (_a = {}, _a[item === null || item === void 0 ? void 0 : item.name] = e, _a)));
      }
    })) : null;
  })))), React__default.createElement(teaComponent.Modal.Footer, null, _renderButtons(), React__default.createElement(TipInfo, {
    isShow: failed,
    type: "error",
    isForm: true
  }, getWorkflowError(createResourceWorkflow))));
}

function CertificateField(props) {
  var _a = props.id,
    id = _a === void 0 ? 0 : _a,
    schema = props.schema;
  return schema ? React__default.createElement(ffComponent.FormPanel.Item, {
    text: true,
    label: schema === null || schema === void 0 ? void 0 : schema.name,
    isShow: !!(schema === null || schema === void 0 ? void 0 : schema.value)
  }, React__default.createElement("div", {
    className: "form-unit tea-mt-1n"
  }, React__default.createElement("div", {
    className: "rich-textarea hide-number",
    style: {
      width: '100%'
    }
  }, React__default.createElement("div", {
    className: "copy-btn"
  }, React__default.createElement(teaComponent.Copy, {
    text: schema === null || schema === void 0 ? void 0 : schema.value
  }, i18n.t('复制'))), React__default.createElement("a", {
    href: "javascript:void(0)",
    onClick: function onClick(e) {
      return downloadKubeconfig(schema === null || schema === void 0 ? void 0 : schema.value, "".concat(id, "-").concat(schema.name));
    },
    className: "copy-btn",
    style: {
      right: '50px'
    }
  }, i18n.t('下载')), React__default.createElement("div", {
    className: "rich-content"
  }, React__default.createElement("pre", {
    className: "rich-text",
    style: {
      whiteSpace: 'pre-wrap',
      overflow: 'auto',
      height: '300px'
    }
  }, schema === null || schema === void 0 ? void 0 : schema.value))))) : null;
}

function ResourceDetailPanel(props) {
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o;
  var _p = props.detail,
    selectDetailResource = _p.selectDetailResource,
    resourceDetail = _p.resourceDetail,
    actions = props.actions;
  if (!(selectDetailResource === null || selectDetailResource === void 0 ? void 0 : selectDetailResource.length)) {
    return React__default.createElement("noscript", null);
  }
  var _cancel = function _cancel() {
    var _a, _b, _c, _d, _e;
    (_b = (_a = props === null || props === void 0 ? void 0 : props.actions) === null || _a === void 0 ? void 0 : _a.detail) === null || _b === void 0 ? void 0 : _b.selectDetailResource([]);
    (_e = (_d = (_c = props === null || props === void 0 ? void 0 : props.actions) === null || _c === void 0 ? void 0 : _c.detail) === null || _d === void 0 ? void 0 : _d.instanceResource) === null || _e === void 0 ? void 0 : _e.selects([]);
  };
  var resource = (_b = (_a = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.resource;
  var instanceParamsFields = (_d = (_c = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.instanceSchema;
  var loading = !((_e = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _e === void 0 ? void 0 : _e.fetched) || ((_f = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _f === void 0 ? void 0 : _f.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Fetching);
  var failed = ((_g = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _g === void 0 ? void 0 : _g.error) || ((_h = resourceDetail === null || resourceDetail === void 0 ? void 0 : resourceDetail.object) === null || _h === void 0 ? void 0 : _h.fetchState) === (ffRedux.FetchState === null || ffRedux.FetchState === void 0 ? void 0 : ffRedux.FetchState.Failed);
  return React__default.createElement(teaComponent.Modal, {
    visible: true,
    caption: i18n.t('服务绑定详情'),
    onClose: _cancel,
    size: 's'
  }, React__default.createElement(teaComponent.Modal.Body, null, React__default.createElement(ffComponent.FormPanel, {
    isNeedCard: false
  }, loading && React__default.createElement(LoadingPanel, null), !loading && failed && React__default.createElement(RetryPanel, {
    style: {
      minWidth: 150,
      width: 150
    },
    loadingText: i18n.t('查询详情失败'),
    action: (_k = (_j = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _j === void 0 ? void 0 : _j.instanceDetail) === null || _k === void 0 ? void 0 : _k.fetch
  }), !loading && !failed && React__default.createElement(React__default.Fragment, null, React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('ID'),
    text: true
  }, i18n.t('{{name}}', {
    name: ((_l = resource === null || resource === void 0 ? void 0 : resource.spec) === null || _l === void 0 ? void 0 : _l.externalID) || '-'
  })), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('名称'),
    text: true
  }, i18n.t('{{name}}', {
    name: ((_m = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _m === void 0 ? void 0 : _m.name) || '-'
  })), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('命名空间'),
    text: true
  }, i18n.t('{{name}}', {
    name: ((_o = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _o === void 0 ? void 0 : _o.namespace) || '-'
  })), instanceParamsFields === null || instanceParamsFields === void 0 ? void 0 : instanceParamsFields.map(function (item) {
    var _a, _b, _c, _d;
    var value;
    try {
      value = ((item === null || item === void 0 ? void 0 : item.type) === (SchemaType === null || SchemaType === void 0 ? void 0 : SchemaType.List) ? (_b = (_a = JSON === null || JSON === void 0 ? void 0 : JSON.parse(item === null || item === void 0 ? void 0 : item.value)) !== null && _a !== void 0 ? _a : []) === null || _b === void 0 ? void 0 : _b.join(',') : item === null || item === void 0 ? void 0 : item.value) || '-';
    } catch (error) {
      value = '-';
    }
    return !((_c = ['ca_pem', 'client_pem', 'client_key_pem']) === null || _c === void 0 ? void 0 : _c.includes(item === null || item === void 0 ? void 0 : item.name)) ? React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('{{name}}', {
        name: item === null || item === void 0 ? void 0 : item.label
      }),
      key: item === null || item === void 0 ? void 0 : item.name,
      text: true
    }, i18n.t('{{value}}', {
      value: value
    })) : React__default.createElement(CertificateField, {
      schema: item,
      id: (_d = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _d === void 0 ? void 0 : _d.name
    });
  })))));
}

var Body = teaComponent.Layout.Body,
  Content = teaComponent.Layout.Content;
function ServiceBindingPanel(props) {
  var _a, _b;
  var instanceResource = props.detail.instanceResource,
    actions = props.actions;
  var columns = [{
    key: "bindingId",
    header: "服务绑定ID",
    render: function render(item) {
      var _a, _b, _c;
      return React__default.createElement("p", null, React__default.createElement(teaComponent.Text, {
        overflow: true,
        tooltip: true
      }, (_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.externalID) !== null && _b !== void 0 ? _b : '-'), React__default.createElement(teaComponent.Copy, {
        text: (_c = item === null || item === void 0 ? void 0 : item.spec) === null || _c === void 0 ? void 0 : _c.externalID
      }));
    }
  }, {
    key: "bindingName",
    header: "服务绑定名称",
    render: function render(item) {
      var _a;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, {
        overflow: true,
        tooltip: true
      }, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name));
    }
  }, {
    key: "namespace",
    header: "命名空间",
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, (_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.namespace) !== null && _b !== void 0 ? _b : '-'));
    }
  }, {
    key: "instanceId",
    header: "实例ID",
    render: function render(item) {
      var _a, _b, _c, _d;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, ((_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['clusternet.io/instanceId']) || ((_d = (_c = item === null || item === void 0 ? void 0 : item.metadata) === null || _c === void 0 ? void 0 : _c.labels) === null || _d === void 0 ? void 0 : _d['ssm.infra.tce.io/instance-id']) || '-'));
    }
  }, {
    key: "bindingTime",
    header: "绑定时间",
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, (_b = dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')) !== null && _b !== void 0 ? _b : '-'));
    }
  }, {
    key: "status",
    header: "状态",
    render: function render(item) {
      var _a, _b, _c, _d, _e, _f, _g, _h;
      var state = ((_a = item === null || item === void 0 ? void 0 : item.status) === null || _a === void 0 ? void 0 : _a.state) || '-';
      var className = ((_b = item === null || item === void 0 ? void 0 : item.status) === null || _b === void 0 ? void 0 : _b.state) === (ServiceBindingStatusNum === null || ServiceBindingStatusNum === void 0 ? void 0 : ServiceBindingStatusNum.Ready) ? 'text-success' : '';
      var isDeleting = showResourceDeleteLoading(item, (_d = (_c = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.list) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.records);
      return React__default.createElement(React__default.Fragment, null, isDeleting && React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Icon, {
        type: 'loading'
      }), i18n.t('解除绑定中')), !isDeleting && React__default.createElement(teaComponent.Text, {
        className: className + " tea-mar-1n"
      }, state), !isDeleting && !!((_f = (_e = item === null || item === void 0 ? void 0 : item.status) === null || _e === void 0 ? void 0 : _e.conditions) === null || _f === void 0 ? void 0 : _f.length) && React__default.createElement(teaComponent.Bubble, {
        content: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.List, {
          type: "number",
          style: {
            width: '100%'
          }
        }, (_h = (_g = item === null || item === void 0 ? void 0 : item.status) === null || _g === void 0 ? void 0 : _g.conditions) === null || _h === void 0 ? void 0 : _h.map(function (item) {
          return React__default.createElement(teaComponent.List.Item, null, React__default.createElement(teaComponent.Text, {
            className: 'tea-mr-2n'
          }, " ", "".concat(item === null || item === void 0 ? void 0 : item.type, " : ").concat((item === null || item === void 0 ? void 0 : item.reason) || (item === null || item === void 0 ? void 0 : item.message))), React__default.createElement(teaComponent.Icon, {
            type: (item === null || item === void 0 ? void 0 : item.status) === 'True' ? 'success' : 'error'
          }));
        })))
      }, React__default.createElement(teaComponent.Icon, {
        type: "info",
        className: "tea-mr-2n"
      })));
    }
  }, {
    key: "operate",
    header: "操作",
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
        type: "link",
        disabled: !!((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp),
        onClick: function onClick() {
          var _a;
          (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.selectDetailResource([item]);
        }
      }, i18n.t('详情')), React__default.createElement(teaComponent.Button, {
        type: "link",
        disabled: !!((_b = item === null || item === void 0 ? void 0 : item.metadata) === null || _b === void 0 ? void 0 : _b.deletionTimestamp),
        onClick: function onClick() {
          var _a;
          (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.selectDeleteResources([item]);
        }
      }, i18n.t('解除绑定')));
    }
  }];
  return React__default.createElement(teaComponent.Layout, null, React__default.createElement(teaComponent.Table.ActionPanel, {
    className: 'tea-mb-5n'
  }, React__default.createElement(teaComponent.Justify, {
    left: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
      type: 'primary',
      onClick: function onClick() {
        actions.detail.showCreateResourceDialog(true);
      }
    }, i18n.t('新建绑定')), React__default.createElement(teaComponent.Button, {
      type: 'primary',
      onClick: function onClick() {
        var _a;
        (_a = actions === null || actions === void 0 ? void 0 : actions.list) === null || _a === void 0 ? void 0 : _a.selectDeleteResources(instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.selections);
      },
      disabled: !((_a = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.selections) === null || _a === void 0 ? void 0 : _a.length),
      tooltip: !((_b = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.selections) === null || _b === void 0 ? void 0 : _b.length) ? i18n.t('请您选择需要删除的资源') : i18n.t('')
    }, i18n.t('删除'))),
    right: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
      icon: "refresh",
      onClick: function onClick() {
        var _a, _b;
        (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.instanceResource) === null || _b === void 0 ? void 0 : _b.fetch();
      }
    }))
  })), React__default.createElement(teaComponent.Card, null, React__default.createElement(ffComponent.TablePanel, {
    recordKey: function recordKey(record) {
      var _a, _b;
      return ((_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name) + ((_b = record === null || record === void 0 ? void 0 : record.metadata) === null || _b === void 0 ? void 0 : _b.namespace);
    },
    columns: columns,
    model: instanceResource,
    action: actions.detail.instanceResource,
    isNeedPagination: true,
    selectable: {
      value: instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.selections.map(function (item) {
        var _a, _b;
        return ((_a = item.metadata) === null || _a === void 0 ? void 0 : _a.name) + ((_b = item.metadata) === null || _b === void 0 ? void 0 : _b.namespace);
      }),
      onChange: function onChange(keys, context) {
        var _a, _b;
        (_b = (_a = actions.detail) === null || _a === void 0 ? void 0 : _a.instanceResource) === null || _b === void 0 ? void 0 : _b.selects(instanceResource.list.data.records.filter(function (item) {
          var _a, _b;
          return keys.includes(((_a = item.metadata) === null || _a === void 0 ? void 0 : _a.name) + ((_b = item.metadata) === null || _b === void 0 ? void 0 : _b.namespace) + '');
        }));
      },
      rowSelect: false
    },
    rowDisabled: function rowDisabled(record) {
      var _a;
      return !!((_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp);
    }
  })), React__default.createElement(CreateServiceBindingDialog, tslib.__assign({}, props)), React__default.createElement(ResourceDetailPanel, tslib.__assign({}, props)));
}

/* eslint-disable no-unused-expressions */
function getMediumOpration(medium) {
  var _a, _b, _c, _d;
  var operation = {
    backup: {}
  };
  if (medium) {
    if (((_b = (_a = medium === null || medium === void 0 ? void 0 : medium.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b[PaasMedium.StorageTypeLabel]) === PaasMedium.StorageTypeEnum.S3) {
      operation = {
        backup: {
          destination: {
            s3Objects: {
              secretRef: {
                namespace: 'ssm',
                name: (_c = medium === null || medium === void 0 ? void 0 : medium.metadata) === null || _c === void 0 ? void 0 : _c.name
              }
            }
          }
        }
      };
    } else {
      operation = {
        backup: {
          destination: {
            nfsObjects: {
              secretRef: {
                namespace: 'ssm',
                name: (_d = medium === null || medium === void 0 ? void 0 : medium.metadata) === null || _d === void 0 ? void 0 : _d.name
              }
            }
          }
        }
      };
    }
  }
  return operation;
}
function BackUpNowDialog(props) {
  var _this = this;
  var _a, _b, _c, _d, _e, _f, _g, _h;
  var _j = props.detail,
    showBackupStrategyDialog = _j.showBackupStrategyDialog,
    mediums = _j.mediums,
    _k = props.base,
    platform = _k.platform,
    route = _k.route,
    hubCluster = _k.hubCluster,
    regionId = _k.regionId,
    _l = props.list,
    backupNowWorkflow = _l.backupNowWorkflow,
    serviceResources = _l.serviceResources,
    services = _l.services,
    actions = props.actions,
    _m = props.mode;
  var clusterId = (_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.clusterid;
  React.useEffect(function () {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m;
    if (backupNowWorkflow.operationState === ffRedux.OperationState.Done) {
      (_c = (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.mediums) === null || _b === void 0 ? void 0 : _b.mediums) === null || _c === void 0 ? void 0 : _c.selectByValue((_h = (_g = (_f = (_e = (_d = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _d === void 0 ? void 0 : _d.list) === null || _e === void 0 ? void 0 : _e.data) === null || _f === void 0 ? void 0 : _f.records) === null || _g === void 0 ? void 0 : _g[0]) === null || _h === void 0 ? void 0 : _h.clusterId);
      (_k = (_j = actions === null || actions === void 0 ? void 0 : actions.create) === null || _j === void 0 ? void 0 : _j.backupNowWorkflow) === null || _k === void 0 ? void 0 : _k.reset();
      (_m = (_l = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _l === void 0 ? void 0 : _l.instanceResource) === null || _m === void 0 ? void 0 : _m.fetch();
    }
  }, [(_b = actions === null || actions === void 0 ? void 0 : actions.create) === null || _b === void 0 ? void 0 : _b.backupNowWorkflow, (_c = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _c === void 0 ? void 0 : _c.instanceResource, (_e = (_d = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _d === void 0 ? void 0 : _d.mediums) === null || _e === void 0 ? void 0 : _e.mediums, backupNowWorkflow.operationState, (_h = (_g = (_f = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _f === void 0 ? void 0 : _f.list) === null || _g === void 0 ? void 0 : _g.data) === null || _h === void 0 ? void 0 : _h.records]);
  if (backupNowWorkflow.operationState === ffRedux.OperationState.Pending) return null;
  var failed = backupNowWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(backupNowWorkflow);
  var perform = function perform() {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      return tslib.__generator(this, function (_a) {
        actions.detail.mediums.validator.validate(null, function (validateResult) {
          var _a, _b, _c, _d, _e, _f, _g;
          if (ffValidator.isValid(validateResult)) {
            var regionId_1 = HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
            var params = {
              platform: platform,
              clusterId: clusterId,
              regionId: regionId_1,
              resourceType: ResourceTypeEnum.Backup,
              specificOperate: CreateSpecificOperatorEnum === null || CreateSpecificOperatorEnum === void 0 ? void 0 : CreateSpecificOperatorEnum.BackupNow,
              jsonData: JSON.stringify({
                apiVersion: 'infra.tce.io/v1',
                kind: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup,
                metadata: {
                  name: "backup-".concat((_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.instancename, "-").concat(new Date().getTime()),
                  namespace: SystemNamespace
                },
                spec: {
                  enabled: true,
                  operation: getMediumOpration((_b = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _b === void 0 ? void 0 : _b.selection),
                  retain: {
                    days: Backup === null || Backup === void 0 ? void 0 : Backup.maxReserveDay
                  },
                  target: {
                    instanceID: Util === null || Util === void 0 ? void 0 : Util.getInstanceId(platform, (_c = route === null || route === void 0 ? void 0 : route.queries) === null || _c === void 0 ? void 0 : _c.instancename, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection),
                    serviceClass: (_d = route === null || route === void 0 ? void 0 : route.queries) === null || _d === void 0 ? void 0 : _d.servicename
                  },
                  trigger: {
                    type: BackupTypeNum.Manual
                  }
                }
              }),
              instanceName: "backup-".concat((_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.instancename, "-").concat(new Date().getTime()),
              namespace: SystemNamespace
            };
            (_f = actions === null || actions === void 0 ? void 0 : actions.create) === null || _f === void 0 ? void 0 : _f.backupNowWorkflow.start([params], regionId_1);
            (_g = actions === null || actions === void 0 ? void 0 : actions.create) === null || _g === void 0 ? void 0 : _g.backupNowWorkflow.perform();
          } else {
            bridge.tips.error(i18n.t('请选择备份介质'));
          }
        });
        return [2 /*return*/];
      });
    });
  };

  var cancel = function cancel() {
    var _a, _b, _c, _d, _e, _f, _g, _h;
    var workflow = backupNowWorkflow;
    if (workflow.operationState === ffRedux.OperationState.Done) {
      actions.create.backupNowWorkflow.reset();
    }
    if (workflow.operationState === ffRedux.OperationState.Started) {
      actions.create.backupNowWorkflow.cancel();
    }
    //重置选中第一个
    (_c = (_b = (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.mediums) === null || _b === void 0 ? void 0 : _b.mediums) === null || _c === void 0 ? void 0 : _c.selectByValue((_h = (_g = (_f = (_e = (_d = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _d === void 0 ? void 0 : _d.list) === null || _e === void 0 ? void 0 : _e.data) === null || _f === void 0 ? void 0 : _f.records) === null || _g === void 0 ? void 0 : _g[0]) === null || _h === void 0 ? void 0 : _h.clusterId);
  };
  return React__default.createElement(teaComponent.Modal, {
    visible: true,
    caption: i18n.t('立即开始备份'),
    onClose: cancel,
    size: 'm'
  }, React__default.createElement(teaComponent.Modal.Body, null, React__default.createElement(teaComponent.Alert, null, i18n.t('中间件实例将开始备份，备份不会造成业务中断，是否立即开始备份')), React__default.createElement(ffComponent.FormPanel, {
    isNeedCard: false
  }, React__default.createElement(MediumSelectPanel.Component, {
    model: mediums,
    action: actions.detail.mediums,
    platform: platform,
    clusterId: clusterId,
    regionId: regionId
  }))), React__default.createElement(teaComponent.Modal.Footer, null, React__default.createElement(teaComponent.Button, {
    type: "primary",
    className: "tea-mr-2n",
    disabled: backupNowWorkflow.operationState === ffRedux.OperationState.Performing,
    onClick: perform
  }, failed ? i18n.t('重试') : i18n.t('完成')), React__default.createElement(teaComponent.Button, {
    title: i18n.t('取消'),
    onClick: cancel
  }, i18n.t('取消')), React__default.createElement(TipInfo, {
    isShow: failed,
    type: "error",
    isForm: true
  }, getWorkflowError(backupNowWorkflow))));
}

function BackUpStrategyDialog(props) {
  var _this = this;
  var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u;
  var _v = props.detail,
    showBackupStrategyDialog = _v.showBackupStrategyDialog,
    backupStrategyEdit = _v.backupStrategyEdit,
    backupStrategy = _v.backupStrategy,
    mediums = _v.mediums,
    _w = props.base,
    platform = _w.platform,
    route = _w.route,
    hubCluster = _w.hubCluster,
    regionId = _w.regionId,
    _x = props.list,
    createResourceWorkflow = _x.createResourceWorkflow,
    serviceResources = _x.serviceResources,
    services = _x.services,
    actions = props.actions,
    _y = props.mode;
  var _z = React.useState(''),
    errorMsg = _z[0],
    setErrorMsg = _z[1];
  var clusterId = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection);
  React.useEffect(function () {
    var _a, _b, _c, _d;
    if (platform && regionId && (serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection) && showBackupStrategyDialog) {
      var clusterId_1 = Util === null || Util === void 0 ? void 0 : Util.getClusterId(platform, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection);
      var instanceId = (_b = (_a = serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection) === null || _a === void 0 ? void 0 : _a.spec) === null || _b === void 0 ? void 0 : _b.externalID;
      (_d = (_c = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _c === void 0 ? void 0 : _c.backupStrategy) === null || _d === void 0 ? void 0 : _d.applyFilter({
        regionId: regionId,
        clusterId: clusterId_1,
        resourceType: ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.Backup,
        platform: platform,
        instanceId: instanceId
      });
    }
  }, [platform, regionId, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection, showBackupStrategyDialog]);
  //初始化备份介质
  React.useEffect(function () {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v, _w, _x;
    if (((_a = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _a === void 0 ? void 0 : _a.data) && ((_c = (_b = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _b === void 0 ? void 0 : _b.list) === null || _c === void 0 ? void 0 : _c.fetched)) {
      var backupMedium = ((_l = (_k = (_j = (_h = (_g = (_f = (_e = (_d = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _d === void 0 ? void 0 : _d.data) === null || _e === void 0 ? void 0 : _e.spec) === null || _f === void 0 ? void 0 : _f.operation) === null || _g === void 0 ? void 0 : _g.backup) === null || _h === void 0 ? void 0 : _h.destination) === null || _j === void 0 ? void 0 : _j.nfsObjects) === null || _k === void 0 ? void 0 : _k.secretRef) === null || _l === void 0 ? void 0 : _l.name) || ((_u = (_t = (_s = (_r = (_q = (_p = (_o = (_m = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _m === void 0 ? void 0 : _m.data) === null || _o === void 0 ? void 0 : _o.spec) === null || _p === void 0 ? void 0 : _p.operation) === null || _q === void 0 ? void 0 : _q.backup) === null || _r === void 0 ? void 0 : _r.destination) === null || _s === void 0 ? void 0 : _s.s3Objects) === null || _t === void 0 ? void 0 : _t.secretRef) === null || _u === void 0 ? void 0 : _u.name);
      console.log('6666', backupMedium);
      (_x = (_w = (_v = actions.detail) === null || _v === void 0 ? void 0 : _v.mediums) === null || _w === void 0 ? void 0 : _w.mediums) === null || _x === void 0 ? void 0 : _x.selectByValue(backupMedium);
    }
  }, [(_b = (_a = actions.detail) === null || _a === void 0 ? void 0 : _a.mediums) === null || _b === void 0 ? void 0 : _b.mediums, (_c = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _c === void 0 ? void 0 : _c.data, (_f = (_e = (_d = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _d === void 0 ? void 0 : _d.list) === null || _e === void 0 ? void 0 : _e.data) === null || _f === void 0 ? void 0 : _f.recordCount, (_h = (_g = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _g === void 0 ? void 0 : _g.list) === null || _h === void 0 ? void 0 : _h.fetched]);
  var backupStrategyLoading = ((_j = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _j === void 0 ? void 0 : _j.loading) || !((_k = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _k === void 0 ? void 0 : _k.fetched);
  var backupStrategyFailed = (_l = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _l === void 0 ? void 0 : _l.error;
  var operateLoading = (createResourceWorkflow === null || createResourceWorkflow === void 0 ? void 0 : createResourceWorkflow.operationState) === ffRedux.OperationState.Performing;
  var failed = createResourceWorkflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(createResourceWorkflow);
  var _submit = function _submit() {
    return tslib.__awaiter(_this, void 0, void 0, function () {
      var validateFun, validateResult, result, formData, regionId_1, params;
      var _this = this;
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q;
      return tslib.__generator(this, function (_r) {
        switch (_r.label) {
          case 0:
            validateFun = function validateFun() {
              return tslib.__awaiter(_this, void 0, void 0, function () {
                return tslib.__generator(this, function (_a) {
                  return [2 /*return*/, new Promise(function (resove, reject) {
                    actions.detail.mediums.validator.validate(null, function (result) {
                      resove(ffValidator.isValid(result));
                    });
                  })];
                });
              });
            };
            return [4 /*yield*/, validateFun()];
          case 1:
            validateResult = _r.sent();
            return [4 /*yield*/, (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.validateAll()];
          case 2:
            _r.sent();
            result = Backup === null || Backup === void 0 ? void 0 : Backup._validateAll(backupStrategyEdit);
            formData = backupStrategyEdit.formData;
            if ((result === null || result === void 0 ? void 0 : result.valid) && validateResult) {
              regionId_1 = HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
              params = {
                platform: platform,
                regionId: regionId_1,
                clusterId: (_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.clusterid,
                jsonData: Backup.reduceBackupStrategyJson({
                  mode: ((_c = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _c === void 0 ? void 0 : _c.data) ? 'edit' : 'create',
                  name: ((_d = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _d === void 0 ? void 0 : _d.data) ? formData === null || formData === void 0 ? void 0 : formData.name : '',
                  enable: formData === null || formData === void 0 ? void 0 : formData.enable,
                  backupDate: formData === null || formData === void 0 ? void 0 : formData.backupDate,
                  backupTime: formData === null || formData === void 0 ? void 0 : formData.backupTime,
                  backupReserveDay: formData === null || formData === void 0 ? void 0 : formData.backupReserveDay,
                  instanceId: Util === null || Util === void 0 ? void 0 : Util.getInstanceId(platform, (_e = route === null || route === void 0 ? void 0 : route.queries) === null || _e === void 0 ? void 0 : _e.instancename, serviceResources === null || serviceResources === void 0 ? void 0 : serviceResources.selection),
                  instanceName: (_f = route === null || route === void 0 ? void 0 : route.queries) === null || _f === void 0 ? void 0 : _f.instancename,
                  serviceName: ((_g = services === null || services === void 0 ? void 0 : services.selection) === null || _g === void 0 ? void 0 : _g.name) || ((_h = route === null || route === void 0 ? void 0 : route.queries) === null || _h === void 0 ? void 0 : _h.servicename),
                  medium: (_j = mediums === null || mediums === void 0 ? void 0 : mediums.mediums) === null || _j === void 0 ? void 0 : _j.selection
                }),
                resourceType: ResourceTypeEnum.Backup,
                instanceName: "backup-".concat((_k = route === null || route === void 0 ? void 0 : route.queries) === null || _k === void 0 ? void 0 : _k.instancename, "-").concat(new Date().getTime()),
                specificOperate: CreateSpecificOperatorEnum === null || CreateSpecificOperatorEnum === void 0 ? void 0 : CreateSpecificOperatorEnum.BackupStrategy
              };
              //若实例的备份策略已存在，则更新备份策略；若实例的备份策略不存在，则新建备份策略
              if (!((_l = backupStrategy === null || backupStrategy === void 0 ? void 0 : backupStrategy.object) === null || _l === void 0 ? void 0 : _l.data)) {
                (_m = actions === null || actions === void 0 ? void 0 : actions.create) === null || _m === void 0 ? void 0 : _m.createResource.start([params], regionId_1);
                (_o = actions === null || actions === void 0 ? void 0 : actions.create) === null || _o === void 0 ? void 0 : _o.createResource.perform();
              } else {
                params = tslib.__assign(tslib.__assign({}, params), {
                  instanceName: formData === null || formData === void 0 ? void 0 : formData.name
                });
                (_p = actions === null || actions === void 0 ? void 0 : actions.create) === null || _p === void 0 ? void 0 : _p.updateResource.start([params], regionId_1);
                (_q = actions === null || actions === void 0 ? void 0 : actions.create) === null || _q === void 0 ? void 0 : _q.updateResource.perform();
              }
            } else {
              bridge.tips.error(result === null || result === void 0 ? void 0 : result.message);
            }
            return [2 /*return*/];
        }
      });
    });
  };

  var _cancel = function _cancel() {
    var _a;
    actions.detail.showBackupDialog(false);
    (_a = actions === null || actions === void 0 ? void 0 : actions.create.createResource) === null || _a === void 0 ? void 0 : _a.reset();
  };
  var renderButtons = function renderButtons() {
    var buttons = [{
      handleFunc: _submit,
      text: failed ? i18n.t('重试') : i18n.t('确定'),
      type: 'primary',
      operateLoading: operateLoading
    }, {
      handleFunc: _cancel,
      text: i18n.t('取消'),
      operateLoading: false
    }];
    return buttons === null || buttons === void 0 ? void 0 : buttons.map(function (item, index) {
      return React__default.createElement(teaComponent.Button, {
        loading: item === null || item === void 0 ? void 0 : item.operateLoading,
        key: index,
        type: item === null || item === void 0 ? void 0 : item.type,
        onClick: item === null || item === void 0 ? void 0 : item.handleFunc
      }, item === null || item === void 0 ? void 0 : item.text);
    });
  };
  return React__default.createElement(teaComponent.Modal, {
    visible: showBackupStrategyDialog,
    caption: i18n.t('备份策略管理'),
    onClose: function onClose() {
      setErrorMsg('');
      actions.detail.showBackupDialog(false);
    },
    size: 'm'
  }, React__default.createElement(teaComponent.Modal.Body, null, backupStrategyLoading && React__default.createElement(LoadingPanel, null), !backupStrategyLoading && backupStrategyFailed && React__default.createElement(RetryPanel, {
    retryText: i18n.t('查询失败'),
    action: (_m = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _m === void 0 ? void 0 : _m.backupStrategy
  }), !backupStrategyLoading && !backupStrategyFailed && React__default.createElement(ffComponent.FormPanel, {
    isNeedCard: false
  }, React__default.createElement(MediumSelectPanel.Component, {
    model: mediums,
    action: actions.detail.mediums,
    platform: platform,
    clusterId: clusterId,
    regionId: regionId
  }), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('启用备份策略'),
    message: errorMsg ? React__default.createElement(teaComponent.Text, {
      className: "text-danger"
    }, errorMsg) : null
  }, React__default.createElement(ffComponent.FormPanel.Switch, {
    value: (_o = backupStrategyEdit === null || backupStrategyEdit === void 0 ? void 0 : backupStrategyEdit.formData) === null || _o === void 0 ? void 0 : _o.enable,
    onChange: function onChange(value) {
      return tslib.__awaiter(_this, void 0, void 0, function () {
        var _a;
        return tslib.__generator(this, function (_b) {
          (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.updateBackUpStrategy('enable', value);
          return [2 /*return*/];
        });
      });
    }
  })), ((_p = backupStrategyEdit === null || backupStrategyEdit === void 0 ? void 0 : backupStrategyEdit.formData) === null || _p === void 0 ? void 0 : _p.enable) ? React__default.createElement(React__default.Fragment, null, React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('备份日期')
  }, React__default.createElement(teaComponent.Checkbox.Group, {
    value: (_q = backupStrategyEdit === null || backupStrategyEdit === void 0 ? void 0 : backupStrategyEdit.formData) === null || _q === void 0 ? void 0 : _q.backupDate,
    onChange: function onChange(value) {
      var _a;
      (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.updateBackUpStrategy('backupDate', value);
    }
  }, (_r = Backup === null || Backup === void 0 ? void 0 : Backup.weekConfig) === null || _r === void 0 ? void 0 : _r.map(function (item) {
    return React__default.createElement(teaComponent.Checkbox, {
      key: item === null || item === void 0 ? void 0 : item.value,
      name: item === null || item === void 0 ? void 0 : item.value,
      className: "tea-mb-2n",
      style: {
        borderRadius: 5
      }
    }, i18n.t(item === null || item === void 0 ? void 0 : item.text));
  }))), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('备份时间点')
  }, React__default.createElement(teaComponent.Checkbox.Group, {
    value: (_s = backupStrategyEdit === null || backupStrategyEdit === void 0 ? void 0 : backupStrategyEdit.formData) === null || _s === void 0 ? void 0 : _s.backupTime,
    onChange: function onChange(value) {
      var _a;
      (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.updateBackUpStrategy('backupTime', value);
    }
  }, (_t = Backup === null || Backup === void 0 ? void 0 : Backup.hourConfig) === null || _t === void 0 ? void 0 : _t.map(function (item) {
    return React__default.createElement(teaComponent.Checkbox, {
      key: item === null || item === void 0 ? void 0 : item.value,
      name: item === null || item === void 0 ? void 0 : item.value,
      className: "tea-mb-2n",
      style: {
        borderRadius: 5
      }
    }, i18n.t(item === null || item === void 0 ? void 0 : item.text));
  }))), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t('备份保留时间(天)')
  }, React__default.createElement(ffComponent.FormPanel.InputNumber, {
    min: Backup.minReserveDay,
    max: Backup.maxReserveDay,
    value: (_u = backupStrategyEdit === null || backupStrategyEdit === void 0 ? void 0 : backupStrategyEdit.formData) === null || _u === void 0 ? void 0 : _u.backupReserveDay,
    onChange: function onChange(value) {
      var _a;
      (_a = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _a === void 0 ? void 0 : _a.updateBackUpStrategy('backupReserveDay', value);
    }
  }), React__default.createElement(teaComponent.Text, null, i18n.t('天后自动删除')))) : null)), React__default.createElement(teaComponent.Modal.Footer, null, renderButtons(), React__default.createElement(TipInfo, {
    isShow: failed,
    type: "error",
    isForm: true
  }, getWorkflowError(createResourceWorkflow))));
}

function InstanceBackupPanel(props) {
  var _a, _b, _c, _d;
  var _e = props.detail,
    instanceResource = _e.instanceResource,
    backupResourceLoading = _e.backupResourceLoading,
    _f = props.base,
    platform = _f.platform,
    route = _f.route,
    userInfo = _f.userInfo,
    regionId = _f.regionId,
    _g = props.list,
    serviceResources = _g.serviceResources,
    services = _g.services,
    servicesInstance = _g.servicesInstance,
    actions = props.actions;
  var _h = route === null || route === void 0 ? void 0 : route.queries,
    instancename = _h.instancename,
    clusterid = _h.clusterid;

  var columns = [{
    key: 'name',
    header: i18n.t('名称'),
    render: function render(item) {
      var _a, _b, _c;
      return React__default.createElement("p", null, React__default.createElement(teaComponent.Text, {
        overflow: true,
        tooltip: true
      }, (_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name) !== null && _b !== void 0 ? _b : '-'), React__default.createElement(teaComponent.Copy, {
        text: (_c = item === null || item === void 0 ? void 0 : item.metadata) === null || _c === void 0 ? void 0 : _c.name
      }));
    }
  }, {
    key: 'status',
    header: i18n.t('状态'),
    render: function render(item) {
      var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t;
      var isFailed = ((_a = item === null || item === void 0 ? void 0 : item.status) === null || _a === void 0 ? void 0 : _a.phase) && ![BackupStatusNum.Waiting, BackupStatusNum.Success].includes((_b = item === null || item === void 0 ? void 0 : item.status) === null || _b === void 0 ? void 0 : _b.phase);
      var state = isFailed ? BackupStatusNum === null || BackupStatusNum === void 0 ? void 0 : BackupStatusNum.Failed : !((_c = item === null || item === void 0 ? void 0 : item.status) === null || _c === void 0 ? void 0 : _c.phase) ? BackupStatusNum.Waiting : (_d = item === null || item === void 0 ? void 0 : item.status) === null || _d === void 0 ? void 0 : _d.phase;
      var isSuccess = ((_e = item === null || item === void 0 ? void 0 : item.status) === null || _e === void 0 ? void 0 : _e.phase) === BackupStatusNum.Success;
      var isDeleting = showResourceDeleteLoading(item, (_g = (_f = instanceResource === null || instanceResource === void 0 ? void 0 : instanceResource.list) === null || _f === void 0 ? void 0 : _f.data) === null || _g === void 0 ? void 0 : _g.records);
      if (isDeleting) {
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Icon, {
          type: "loading"
        }), i18n.t('删除中'));
      }
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, {
        className: "".concat((_h = BackupStatusMap === null || BackupStatusMap === void 0 ? void 0 : BackupStatusMap[state]) === null || _h === void 0 ? void 0 : _h.className, " tea-mr-1n"),
        overflow: true
      }, (_j = BackupStatusMap[state]) === null || _j === void 0 ? void 0 : _j.text), isFailed && React__default.createElement(teaComponent.Bubble, {
        content: React__default.createElement(teaComponent.List, {
          type: "number",
          style: {
            width: '100%'
          }
        }, (_l = (_k = item === null || item === void 0 ? void 0 : item.status) === null || _k === void 0 ? void 0 : _k.conditions) === null || _l === void 0 ? void 0 : _l.map(function (item, index) {
          return React__default.createElement(teaComponent.List.Item, {
            key: index
          }, React__default.createElement(teaComponent.Text, {
            className: "tea-mr-2n"
          }, " ", "".concat(item === null || item === void 0 ? void 0 : item.type, " : ").concat((item === null || item === void 0 ? void 0 : item.reason) || (item === null || item === void 0 ? void 0 : item.message))), React__default.createElement(teaComponent.Icon, {
            type: (item === null || item === void 0 ? void 0 : item.status) === 'True' ? 'success' : 'error'
          }));
        }))
      }, React__default.createElement(teaComponent.Icon, {
        type: "info",
        className: "tea-mr-2n"
      })), isSuccess && React__default.createElement(teaComponent.Bubble, {
        style: {
          width: 500
        },
        content: React__default.createElement(ffComponent.FormPanel, {
          title: i18n.t('备份桶详情:'),
          isNeedCard: false
        }, (_p = (_o = (_m = item === null || item === void 0 ? void 0 : item.status) === null || _m === void 0 ? void 0 : _m.results) === null || _o === void 0 ? void 0 : _o.BACKUP_FILE_PATH) === null || _p === void 0 ? void 0 : _p.map(function (path, index) {
          return React__default.createElement(ffComponent.FormPanel.Item, {
            key: index,
            label: i18n.t('{{name}}', {
              // eslint-disable-next-line prefer-template
              name: i18n.t('备份地址{{attr0}}', {
                attr0: index + 1
              })
            }),
            text: true
          }, path, React__default.createElement(teaComponent.Copy, {
            text: path
          }));
        }), React__default.createElement(ffComponent.FormPanel.Item, {
          label: i18n.t('备份大小'),
          text: true
        }, Util === null || Util === void 0 ? void 0 : Util.getReadableFileSizeString((_r = (_q = item === null || item === void 0 ? void 0 : item.status) === null || _q === void 0 ? void 0 : _q.results) === null || _r === void 0 ? void 0 : _r.BACKUP_FILE_SIZE)), React__default.createElement(ffComponent.FormPanel.Item, {
          label: i18n.t('开始时间:'),
          text: true
        }, dateFormatter(new Date((_s = item === null || item === void 0 ? void 0 : item.status) === null || _s === void 0 ? void 0 : _s.startTime), 'YYYY-MM-DD HH:mm:ss')), React__default.createElement(ffComponent.FormPanel.Item, {
          label: i18n.t('结束时间:'),
          text: true
        }, dateFormatter(new Date((_t = item === null || item === void 0 ? void 0 : item.status) === null || _t === void 0 ? void 0 : _t.endTime), 'YYYY-MM-DD HH:mm:ss')))
      }, React__default.createElement(teaComponent.Icon, {
        type: "info",
        className: "tea-mr-2n"
      })));
    }
  }, {
    key: 'instanceId',
    header: i18n.t('实例ID'),
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, {
        overflow: true,
        tooltip: true
      }, ((_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.target) === null || _b === void 0 ? void 0 : _b.instanceID) || '-'));
    }
  }, {
    key: 'size',
    header: i18n.t('大小'),
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, Util === null || Util === void 0 ? void 0 : Util.getReadableFileSizeString((_b = (_a = item === null || item === void 0 ? void 0 : item.status) === null || _a === void 0 ? void 0 : _a.results) === null || _b === void 0 ? void 0 : _b.BACKUP_FILE_SIZE)));
    }
  }, {
    key: 'backupType',
    header: i18n.t('备份类型'),
    render: function render(item) {
      var _a, _b, _c, _d, _e, _f;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, {
        className: (_c = BackupTypeMap[(_b = (_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.trigger) !== null && _b !== void 0 ? _b : BackupTypeNum.Unknown]) === null || _c === void 0 ? void 0 : _c.className
      }, (_f = BackupTypeMap[(_e = (_d = item === null || item === void 0 ? void 0 : item.spec) === null || _d === void 0 ? void 0 : _d.trigger) !== null && _e !== void 0 ? _e : BackupTypeNum.Unknown]) === null || _f === void 0 ? void 0 : _f.text));
    }
  }, {
    key: 'createTime',
    header: i18n.t('创建时间'),
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, (_b = dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')) !== null && _b !== void 0 ? _b : '-'));
    }
  }, {
    key: 'reserveTime',
    header: i18n.t('保留时间'),
    render: function render(item) {
      var _a, _b, _c;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, ((_a = item === null || item === void 0 ? void 0 : item.spec) === null || _a === void 0 ? void 0 : _a.trigger) === (BackupTypeNum === null || BackupTypeNum === void 0 ? void 0 : BackupTypeNum.Manual) ? i18n.t('永久保留') : (_c = dateFormatter(new Date((_b = item === null || item === void 0 ? void 0 : item.spec) === null || _b === void 0 ? void 0 : _b.retainTime), 'YYYY-MM-DD HH:mm:ss')) !== null && _c !== void 0 ? _c : '-'));
    }
  }, {
    key: 'deleteTime',
    header: i18n.t('删除时间'),
    render: function render(item) {
      var _a, _b;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Text, null, (_b = dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp), 'YYYY-MM-DD HH:mm:ss')) !== null && _b !== void 0 ? _b : '-'));
    }
  }, {
    key: 'operate',
    header: i18n.t('操作'),
    render: function render(item) {
      var _a;
      return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
        type: "link",
        disabled: !!((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp),
        onClick: function onClick() {
          actions.list.selectDeleteResources([item]);
        }
      }, i18n.t('删除')));
    }
  }];
  // 是否支持备份功能
  var showBackUpOperation = (_c = (_b = (_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.supportedOperations) === null || _c === void 0 ? void 0 : _c.some(function (operation) {
    return (operation === null || operation === void 0 ? void 0 : operation.operation) === (SupportedOperationsEnum === null || SupportedOperationsEnum === void 0 ? void 0 : SupportedOperationsEnum.Backup);
  });
  return React__default.createElement(teaComponent.Layout, null, React__default.createElement(teaComponent.Table.ActionPanel, {
    className: 'tea-mb-5n'
  }, React__default.createElement(teaComponent.Justify, {
    left: React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
      type: 'primary',
      onClick: function onClick() {
        actions === null || actions === void 0 ? void 0 : actions.detail.showBackupDialog(true);
      }
    }, i18n.t('备份策略')), React__default.createElement(teaComponent.Button, {
      loading: backupResourceLoading,
      type: 'primary',
      onClick: function onClick() {
        // actions.detail?.mediums?.mediums?.select(null);
        actions.create.backupNowWorkflow.start([]);
      },
      disabled: backupResourceLoading
    }, i18n.t('立即备份'))),
    right: React__default.createElement(teaComponent.Button, {
      icon: "refresh",
      onClick: function onClick() {
        actions.detail.instanceResource.fetch();
      }
    })
  })), React__default.createElement(ffComponent.TablePanel, {
    recordKey: function recordKey(record) {
      var _a;
      return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.uid;
    },
    columns: columns,
    model: instanceResource,
    action: (_d = actions === null || actions === void 0 ? void 0 : actions.detail) === null || _d === void 0 ? void 0 : _d.instanceResource,
    isNeedPagination: true,
    rowDisabled: function rowDisabled(record) {
      var _a;
      return !!((_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.deletionTimestamp);
    }
  }), React__default.createElement(BackUpStrategyDialog, tslib.__assign({}, props)), React__default.createElement(BackUpNowDialog, tslib.__assign({}, props)));
}

function InstanceDetailContainer(props) {
  var _a, _b;
  var actions = props.actions,
    _c = props.detail,
    selectedDetailTab = _c.selectedDetailTab,
    serviceInstanceSchema = _c.serviceInstanceSchema,
    _d = props.base,
    route = _d.route,
    platform = _d.platform,
    regionId = _d.regionId,
    servicesInstance = props.list.servicesInstance;
  var urlParams = router === null || router === void 0 ? void 0 : router.resolve(route);
  React.useEffect(function () {
    var _a;
    if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === 'detail' && ((_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.fetched)) {
      actions.detail.selectDetailTab(DetailTabType.Detail);
    }
  }, [urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub, servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object]);
  var instanceDetailTabs = React.useMemo(function () {
    var _a, _b, _c;
    if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === 'detail' && ((_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.fetched) && !((_b = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _b === void 0 ? void 0 : _b.error)) {
      return detailTabs === null || detailTabs === void 0 ? void 0 : detailTabs.filter(function (item) {
        var _a, _b, _c;
        if ((item === null || item === void 0 ? void 0 : item.id) === (DetailTabType === null || DetailTabType === void 0 ? void 0 : DetailTabType.BackUp)) {
          return (_c = (_b = (_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.supportedOperations) === null || _c === void 0 ? void 0 : _c.some(function (operation) {
            return (operation === null || operation === void 0 ? void 0 : operation.operation) === (SupportedOperationsEnum === null || SupportedOperationsEnum === void 0 ? void 0 : SupportedOperationsEnum.Backup);
          });
        } else {
          return true;
        }
      });
    } else if ((_c = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _c === void 0 ? void 0 : _c.error) {
      return detailTabs;
    } else {
      return [];
    }
  }, [urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub, servicesInstance]);
  var loadingDetailTabs = ((_a = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _a === void 0 ? void 0 : _a.loading) || !((_b = servicesInstance === null || servicesInstance === void 0 ? void 0 : servicesInstance.object) === null || _b === void 0 ? void 0 : _b.fetched);
  var _renderTabContent = function _renderTabContent(tabId) {
    var content;
    if (tabId === DetailTabType.Detail) {
      content = React__default.createElement(InstanceBaseDetail, tslib.__assign({}, props));
    } else if (tabId === DetailTabType.BackUp) {
      content = React__default.createElement(InstanceBackupPanel, tslib.__assign({}, props));
    } else if (tabId === DetailTabType.Monitor) {
      content = React__default.createElement(InstanceMonitorPanel, null);
    } else if (tabId === DetailTabType.ServiceBinding) {
      content = React__default.createElement(ServiceBindingPanel, tslib.__assign({}, props));
    }
    return content;
  };
  return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Layout, {
    style: {
      boxShadow: 'none'
    }
  }, React__default.createElement(teaComponent.Layout.Header, {
    style: {
      height: '100%'
    }
  }, React__default.createElement(teaComponent.Card, {
    style: {
      display: 'flex',
      alignItems: 'center'
    },
    className: 'tea-pt-4n tea-pb-4n tea-pl-2n tea-lr-2n'
  }, loadingDetailTabs && React__default.createElement(LoadingPanel, null), !loadingDetailTabs && (instanceDetailTabs === null || instanceDetailTabs === void 0 ? void 0 : instanceDetailTabs.map(function (item, key) {
    return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
      key: item === null || item === void 0 ? void 0 : item.id,
      style: {
        borderRadius: 20,
        color: selectedDetailTab === (item === null || item === void 0 ? void 0 : item.id) ? '#fff' : '#95979B',
        marginRight: 30,
        minWidth: 90,
        height: '30px'
      },
      type: selectedDetailTab === (item === null || item === void 0 ? void 0 : item.id) ? 'primary' : 'text',
      onClick: function onClick() {
        actions.detail.selectDetailTab(item === null || item === void 0 ? void 0 : item.id);
      }
    }, item === null || item === void 0 ? void 0 : item.label));
  })))), React__default.createElement(teaComponent.Layout.Body, {
    style: {
      overflow: 'hidden',
      paddingTop: 20
    }
  }, _renderTabContent(selectedDetailTab))));
}

var routerSea$4 = seajs === null || seajs === void 0 ? void 0 : seajs.require('router');
var store = configStore();
var MiddlewareAppContainer = /** @class */function (_super) {
  tslib.__extends(MiddlewareAppContainer, _super);
  function MiddlewareAppContainer(props, context) {
    return _super.call(this, props, context) || this;
  }
  // 页面离开时，清空store
  MiddlewareAppContainer.prototype.componentWillUnmount = function () {
    store.dispatch({
      type: ResetStoreAction
    });
  };
  MiddlewareAppContainer.prototype.render = function () {
    return React.createElement(reactRedux.Provider, {
      store: store
    }, React.createElement(MiddlewareApp, tslib.__assign({}, this.props)));
  };
  return MiddlewareAppContainer;
}(React.Component);
var mapDispatchToProps = function mapDispatchToProps(dispatch) {
  return Object.assign({}, ffRedux.bindActionCreators({
    actions: allActions
  }, dispatch), {
    dispatch: dispatch
  });
};
var MiddlewareApp = /** @class */function (_super) {
  tslib.__extends(MiddlewareApp, _super);
  function MiddlewareApp(props, context) {
    return _super.call(this, props, context) || this;
  }
  MiddlewareApp.prototype.componentDidMount = function () {
    var _a, _b, _c;
    var _d = this.props,
      actions = _d.actions,
      platform = _d.platform,
      isI18n = _d.base.isI18n;
    if (window['VERSION'] === 'en' && !isI18n) {
      actions.base.toggleIsI18n(true);
    }
    if ((_a = this.props) === null || _a === void 0 ? void 0 : _a.platform) {
      (_b = actions === null || actions === void 0 ? void 0 : actions.base) === null || _b === void 0 ? void 0 : _b.fetchPlatform(platform, (_c = this === null || this === void 0 ? void 0 : this.props) === null || _c === void 0 ? void 0 : _c.regionId);
    }
  };
  MiddlewareApp.prototype.render = function () {
    var _a, _b;
    var _c = this.props.base,
      route = _c.route,
      platform = _c.platform;
    var urlParams = router === null || router === void 0 ? void 0 : router.resolve(route);
    var queries = route === null || route === void 0 ? void 0 : route.queries;
    var content;
    if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === 'detail' || ((_a = route === null || route === void 0 ? void 0 : route.queries) === null || _a === void 0 ? void 0 : _a.mode) === 'detail') {
      content = React.createElement(InstanceDetailContainer, tslib.__assign({}, this.props));
    } else if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === 'create' || ((_b = route === null || route === void 0 ? void 0 : route.queries) === null || _b === void 0 ? void 0 : _b.mode) === 'create') {
      content = React.createElement(ServiceCreate, tslib.__assign({}, this.props));
    } else {
      content = React.createElement(PaasContent, tslib.__assign({}, this.props));
    }
    return React.createElement(teaComponent.Layout, {
      style: {
        marginTop: 0
      }
    }, React.createElement(teaComponent.Layout.Header, null, React.createElement(PaasHeader, tslib.__assign({}, this.props))), React.createElement(teaComponent.Layout.Content.Body, {
      style: {
        padding: 20
      },
      full: true
    }, content, React.createElement(ResourceDeleteDialog, tslib.__assign({}, this.props))));
  };
  MiddlewareApp = tslib.__decorate([reactRedux.connect(function (state) {
    return state;
  }, mapDispatchToProps)], MiddlewareApp);
  return MiddlewareApp;
}(React.Component);

var MediumEditDialog;
(function (MediumEditDialog) {
  var _this = this;
  MediumEditDialog.ComponentName = 'MediumEditDialog';
  var ActionType;
  (function (ActionType) {
    ActionType["NodeList"] = "NodeList";
    ActionType["Workflow"] = "Workflow";
    ActionType["Validator"] = "Validator";
    ActionType["NodeUnitName"] = "NodeUnitName";
    //介质管理类型
    ActionType["Type"] = "StorageType";
    //对象存储配置
    ActionType["CosConfig"] = "CosConfig";
    ActionType["NfsConfig"] = "NfsConfig";
    ActionType["Clear"] = "Clear";
  })(ActionType = MediumEditDialog.ActionType || (MediumEditDialog.ActionType = {}));
  BComponent$1.createActionType(MediumEditDialog.ComponentName, ActionType);
  MediumEditDialog.createValidateSchema = function (_a) {
    var pageName = _a.pageName;
    var schema = {
      formKey: BComponent$1.getActionType(pageName, ActionType.Validator),
      fields: [{
        vKey: 'mediumName',
        label: i18n.t('名称'),
        rules: [{
          type: ffValidator.RuleTypeEnum.custom,
          customFunc: function customFunc(value, store) {
            var mediumName = store.mediumName;
            // eslint-disable-next-line prefer-const
            var reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
              status = ffValidator.ValidatorStatusEnum.Success,
              message = '';
            // 验证服务名称
            if (!mediumName) {
              status = ffValidator.ValidatorStatusEnum.Failed;
              message = i18n.t('名称不能为空');
            } else if (mediumName.length > 63) {
              status = ffValidator.ValidatorStatusEnum.Failed;
              message = i18n.t('名称不能超过63个字符');
            } else if (!reg.test(mediumName)) {
              status = ffValidator.ValidatorStatusEnum.Failed;
              message = i18n.t('名称格式不正确');
            }
            return {
              status: status,
              message: message
            };
          }
        }]
      }, {
        vKey: 'cosConfig.bucketNames',
        label: i18n.t('存储名称'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.S3;
        } //对象存储
      }, {
        vKey: 'cosConfig.port',
        label: i18n.t('端口'),
        rules: [{
          type: ffValidator.RuleTypeEnum.custom,
          customFunc: function customFunc(value, store) {
            var cosConfig = store.cosConfig;
            var port = cosConfig === null || cosConfig === void 0 ? void 0 : cosConfig.port;
            // eslint-disable-next-line prefer-const
            var status = ffValidator.ValidatorStatusEnum.Success,
              message = '';
            // 验证端口号
            if (port && (!Number.isInteger(+port) || +port < 1 || +port > 65535)) {
              status = ffValidator.ValidatorStatusEnum.Failed;
              message = i18n.t('检查端口范围必须在1~65535之间');
            }
            return {
              status: status,
              message: message
            };
          }
        }],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.S3;
        } //对象存储
      }, {
        vKey: 'cosConfig.domain',
        label: i18n.t('域名'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.S3;
        } //对象存储
      }, {
        vKey: 'cosConfig.secretId',
        label: i18n.t('secretId'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.S3;
        } //对象存储
      }, {
        vKey: 'cosConfig.secretKey',
        label: i18n.t('secretKey'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.S3;
        } //对象存储
      }, {
        vKey: 'nfsConfig.serverName',
        label: i18n.t('ServerName'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.NFS;
        } //文件存储
      }, {
        vKey: 'nfsConfig.server',
        label: i18n.t('Server'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.NFS;
        } //文件存储
      }, {
        vKey: 'nfsConfig.path',
        label: i18n.t('Path'),
        rules: [ffValidator.RuleTypeEnum.isRequire],
        condition: function condition(value, store) {
          return store.storageType === PaasMedium.StorageTypeEnum.NFS;
        } //文件存储
      }]
    };

    return schema;
  };
  MediumEditDialog.createActions = function (_a) {
    var _b;
    var pageName = _a.pageName,
      getRecord = _a.getRecord;
    var actions = {
      // 获取列表
      // nodeList: createFFListActions<EdgeNode, EdgeNode.Filter>({
      //   actionName: BComponent.getActionType(pageName, ActionType.NodeList),
      //   getRecord: getState => getRecord(getState).nodeList,
      //   fetcher: async query => {
      //     let response = await WebApi.edge.fetchAllEdgeNodeList(query);
      //     return response;
      //   }
      // }),
      setName: function setName(value) {
        return function (dispatch) {
          dispatch({
            type: BComponent$1.getActionType(pageName, ActionType.NodeUnitName),
            payload: value
          });
        };
      },
      setType: function setType(value) {
        return function (dispatch) {
          dispatch({
            type: BComponent$1.getActionType(pageName, ActionType.Type),
            payload: value
          });
        };
      },
      cosConfig: function cosConfig(data) {
        return function (dispach, getState) {
          var cosConfig = getRecord(getState).cosConfig;
          var newData = tslib.__assign(tslib.__assign({}, cosConfig), data);
          dispach({
            type: BComponent$1.getActionType(pageName, ActionType.CosConfig),
            payload: newData
          });
        };
      },
      nfsConfig: function nfsConfig(data) {
        return function (dispach, getState) {
          var nfsConfig = getRecord(getState).nfsConfig;
          var newData = tslib.__assign(tslib.__assign({}, nfsConfig), data);
          dispach({
            type: BComponent$1.getActionType(pageName, ActionType.NfsConfig),
            payload: newData
          });
        };
      },
      workflow: ffRedux.generateWorkflowActionCreator({
        actionType: ActionType.Workflow,
        workflowStateLocator: function workflowStateLocator(state) {
          return getRecord(function () {
            return state;
          }).workflow;
        },
        operationExecutor: function operationExecutor(targets, params, dispatch, getState) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var result, error_1;
            return tslib.__generator(this, function (_a) {
              switch (_a.label) {
                case 0:
                  _a.trys.push([0, 2,, 3]);
                  return [4 /*yield*/, medium.createCosConfigForTdcc(targets, params)];
                case 1:
                  result = _a.sent();
                  return [2 /*return*/, result];
                case 2:
                  error_1 = _a.sent();
                  return [2 /*return*/, operationResult(targets, error_1)];
                case 3:
                  return [2 /*return*/];
              }
            });
          });
        },

        after: (_b = {}, _b[ffRedux.OperationTrigger.Done] = function (dispatch, getState) {}, _b)
      }),
      initForUpdate: function initForUpdate(resource) {
        return function (dispatch, getState) {
          var _a, _b, _c, _d, _e, _f;
          dispatch(actions.setName(resource === null || resource === void 0 ? void 0 : resource.metadata.name));
          var storageType = (_b = (_a = resource === null || resource === void 0 ? void 0 : resource.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b[PaasMedium.StorageTypeLabel];
          dispatch(actions.setType(storageType));
          if (storageType === PaasMedium.StorageTypeEnum.S3) {
            var cosConfigInfo = resource === null || resource === void 0 ? void 0 : resource.data;
            var bucketNames = '';
            try {
              var decode = jsBase64.Base64.decode(cosConfigInfo === null || cosConfigInfo === void 0 ? void 0 : cosConfigInfo.bucketNames);
              bucketNames = (_d = (_c = decode === null || decode === void 0 ? void 0 : decode.replace('[', '')) === null || _c === void 0 ? void 0 : _c.replace(']', '')) === null || _d === void 0 ? void 0 : _d.replace(/\"/g, '');
            } catch (err) {}
            var host = jsBase64.Base64.decode(cosConfigInfo === null || cosConfigInfo === void 0 ? void 0 : cosConfigInfo.host);
            var cosConfig = {
              domain: host === null || host === void 0 ? void 0 : host.split(':')[0],
              port: host === null || host === void 0 ? void 0 : host.split(':')[1],
              bucketNames: bucketNames,
              secretId: jsBase64.Base64.decode(cosConfigInfo === null || cosConfigInfo === void 0 ? void 0 : cosConfigInfo.access_key),
              secretKey: jsBase64.Base64.decode(cosConfigInfo === null || cosConfigInfo === void 0 ? void 0 : cosConfigInfo.secret_key)
            };
            dispatch(actions.cosConfig(cosConfig));
          } else {
            var nfsConfigInfo = resource === null || resource === void 0 ? void 0 : resource.data;
            var mountOptions = '';
            try {
              var decode = jsBase64.Base64.decode(nfsConfigInfo === null || nfsConfigInfo === void 0 ? void 0 : nfsConfigInfo.mountOptions);
              mountOptions = (_f = (_e = decode === null || decode === void 0 ? void 0 : decode.replace('[', '')) === null || _e === void 0 ? void 0 : _e.replace(']', '')) === null || _f === void 0 ? void 0 : _f.replace(/\"/g, '');
            } catch (err) {}
            var nfsConfig = {
              mountOptions: mountOptions,
              serverName: jsBase64.Base64.decode(nfsConfigInfo === null || nfsConfigInfo === void 0 ? void 0 : nfsConfigInfo.name),
              path: jsBase64.Base64.decode(nfsConfigInfo === null || nfsConfigInfo === void 0 ? void 0 : nfsConfigInfo.path),
              server: jsBase64.Base64.decode(nfsConfigInfo === null || nfsConfigInfo === void 0 ? void 0 : nfsConfigInfo.server)
            };
            dispatch(actions.nfsConfig(nfsConfig));
          }
        };
      },
      validator: ffValidator.createValidatorActions({
        userDefinedSchema: MediumEditDialog.createValidateSchema({
          pageName: pageName
        }),
        validateStateLocator: function validateStateLocator(state) {
          return getRecord(function () {
            return state;
          });
        },
        validatorStateLocation: function validatorStateLocation(state) {
          return getRecord(function () {
            return state;
          }).validator;
        }
      }),
      clear: function clear() {
        return {
          type: BComponent$1.getActionType(pageName, ActionType.Clear)
        };
      }
    };
    return actions;
  };
  MediumEditDialog.createReducer = function (_a) {
    var pageName = _a.pageName;
    var reducer = redux.combineReducers({
      mediumName: ffRedux.reduceToPayload(BComponent$1.getActionType(pageName, ActionType.NodeUnitName), ''),
      storageType: ffRedux.reduceToPayload(BComponent$1.getActionType(pageName, ActionType.Type), PaasMedium.StorageTypeEnum.S3),
      cosConfig: ffRedux.reduceToPayload(BComponent$1.getActionType(pageName, ActionType.CosConfig), {}),
      nfsConfig: ffRedux.reduceToPayload(BComponent$1.getActionType(pageName, ActionType.NfsConfig), {}),
      workflow: ffRedux.generateWorkflowReducer({
        actionType: ActionType.Workflow
      }),
      validator: ffValidator.createValidatorReducer(MediumEditDialog.createValidateSchema({
        pageName: pageName
      }))
    });
    return function (state, action) {
      var newState = state;
      if (action.type === BComponent$1.getActionType(pageName, ActionType.Clear)) {
        newState = undefined;
      }
      return reducer(newState, action);
    };
  };
  MediumEditDialog.Component = function (props) {
    var _a = props.model,
      workflow = _a.workflow,
      mediumName = _a.mediumName,
      validator = _a.validator,
      storageType = _a.storageType,
      _b = _a.cosConfig,
      domain = _b.domain,
      bucketNames = _b.bucketNames,
      secretId = _b.secretId,
      secretKey = _b.secretKey,
      port = _b.port,
      _c = _a.nfsConfig,
      server = _c.server,
      path = _c.path,
      mountOptions = _c.mountOptions,
      serverName = _c.serverName,
      actions = props.actions,
      platform = props.platform,
      operationType = props.operationType,
      userInfo = props.userInfo,
      onSuccess = props.onSuccess,
      clusterId = props.clusterId,
      clusterName = props.clusterName,
      existNames = props.existNames,
      resourceInstance = props.resourceInstance,
      regionId = props.regionId;
    var isCreate = operationType === BComponent$1.OperationTypeEnum.Create;
    var failed = workflow.operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(workflow);
    var _d = React.useState(''),
      nameInUseTip = _d[0],
      setNameInUseTip = _d[1];
    React.useEffect(function () {
      if (workflow.operationState === ffRedux.OperationState.Done) {
        if (!failed && typeof onSuccess === 'function') {
          onSuccess();
        }
      }
    }, [failed, onSuccess, workflow.operationState]);
    var cancel = function cancel() {
      actions.clear();
      actions.workflow.reset();
    };
    var perform = function perform() {
      return tslib.__awaiter(_this, void 0, void 0, function () {
        var host, bucketNamesString, mountOPtionsString, jsonData, resource, find;
        var _a, _b, _c, _d;
        return tslib.__generator(this, function (_e) {
          host = domain;
          if (port) {
            host = "".concat(domain, ":").concat(port);
          }
          bucketNamesString = "[".concat((_a = bucketNames === null || bucketNames === void 0 ? void 0 : bucketNames.split(',')) === null || _a === void 0 ? void 0 : _a.map(function (item) {
            return "\"".concat(item, "\"");
          }), "]");
          mountOPtionsString = "[".concat((_b = mountOptions === null || mountOptions === void 0 ? void 0 : mountOptions.split(',')) === null || _b === void 0 ? void 0 : _b.map(function (item) {
            return "\"".concat(item, "\"");
          }), "]");
          jsonData = {
            apiVersion: 'v1',
            kind: 'Secret',
            metadata: {
              name: mediumName,
              namespace: 'ssm',
              labels: {
                'tdcc.cloud.tencent.com/paas-storage-medium': storageType,
                'tdcc.cloud.tencent.com/creator': (_d = (_c = userInfo === null || userInfo === void 0 ? void 0 : userInfo.object) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.name
              }
            },
            data: storageType === PaasMedium.StorageTypeEnum.S3 ? {
              host: jsBase64.Base64.encode(host),
              bucketNames: jsBase64.Base64.encode(bucketNamesString),
              access_key: jsBase64.Base64.encode(secretId),
              secret_key: jsBase64.Base64.encode(secretKey)
            } : {
              name: jsBase64.Base64.encode(serverName),
              server: jsBase64.Base64.encode(server),
              path: jsBase64.Base64.encode(path),
              mountOptions: jsBase64.Base64.encode(mountOPtionsString)
            }
          };
          resource = {
            id: uuid(),
            mode: isCreate ? 'create' : 'update',
            namespace: SystemNamespace,
            clusterId: clusterId,
            jsonData: JSON.stringify(jsonData),
            isNavToEvent: false,
            base64encode: true
          };
          if (isCreate) {
            find = existNames === null || existNames === void 0 ? void 0 : existNames.find(function (item) {
              return item === mediumName;
            });
            if (find) {
              setNameInUseTip(i18n.t('已存在相同名字的备份介质'));
              return [2 /*return*/];
            }
          }

          actions.validator.validate(null, function (validateResult) {
            if (ffValidator.isValid(validateResult)) {
              actions.workflow.start([resource], {
                regionId: regionId,
                clusterId: clusterId,
                platform: platform
              });
              actions.workflow.perform();
            }
          });
          return [2 /*return*/];
        });
      });
    };

    if (workflow.operationState === ffRedux.OperationState.Pending) return null;
    return React__default.createElement(teaComponent.Modal, {
      visible: true,
      caption: isCreate ? i18n.t('添加自定义介质') : i18n.t('更新自定义介质'),
      onClose: cancel,
      size: 800,
      disableEscape: true
    }, React__default.createElement(teaComponent.Modal.Body, null, React__default.createElement(ffComponent.FormPanel, {
      isNeedCard: false,
      className: "tea-mb-2n"
    }, !isCreate ? React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('名称')
    }, ' ', React__default.createElement(ffComponent.FormPanel.Text, null, mediumName)) : React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('名称'),
      message: i18n.t('最长60个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      formvalidator: validator,
      vactions: actions.validator,
      vkey: "mediumName",
      value: mediumName,
      onChange: function onChange(v) {
        setNameInUseTip('');
        actions.setName(v);
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('集群')
    }, React__default.createElement(ffComponent.FormPanel.Text, null, "".concat(clusterId, "(").concat(clusterName, ")"))), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('存储类型')
    }, React__default.createElement(ffComponent.FormPanel.Select, {
      options: PaasMedium.storageTypeOptions,
      onChange: actions.setType,
      vkey: "type",
      value: storageType,
      size: "l"
    })), storageType === PaasMedium.StorageTypeEnum.S3 && React__default.createElement(React__default.Fragment, null, React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('域名')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "cosConfig.domain",
      value: domain,
      onChange: function onChange(v) {
        return actions.cosConfig({
          domain: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('端口')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "cosConfig.port",
      value: port,
      onChange: function onChange(v) {
        return actions.cosConfig({
          port: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('存储桶名称'),
      message: i18n.t('支持填写多个存储桶；用,分割')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "cosConfig.bucketNames",
      value: bucketNames,
      // className="tea-mr-2n"
      onChange: function onChange(v) {
        return actions.cosConfig({
          bucketNames: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('secretId'),
      message: React__default.createElement(ffComponent.FormPanel.HelpText, {
        parent: 'div'
      }, React__default.createElement(i18n.Trans, null, "\u8BBF\u95EECOS\u7684\u5BC6\u94A5Key,\u5BF9\u5E94\u5BC6\u94A5\u7BA1\u7406\u7684SecretId\uFF0C \u70B9\u51FB", React__default.createElement(teaComponent.ExternalLink, {
        href: '/cam/capi'
      }, "\u8FD9\u91CC"), "\u67E5\u770B\u76F8\u5173\u4FE1\u606F"))
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vkey: "cosConfig.secretId",
      type: "password",
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      value: secretId,
      onChange: function onChange(v) {
        return actions.cosConfig({
          secretId: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('secretKey'),
      message: React__default.createElement(ffComponent.FormPanel.HelpText, {
        parent: 'div'
      }, React__default.createElement(i18n.Trans, null, "\u8BBF\u95EECOS\u7684\u5BC6\u94A5Key,\u5BF9\u5E94\u5BC6\u94A5\u7BA1\u7406\u7684SecretKey\uFF0C \u70B9\u51FB", React__default.createElement(teaComponent.ExternalLink, {
        href: '/cam/capi'
      }, "\u8FD9\u91CC"), "\u67E5\u770B\u76F8\u5173\u4FE1\u606F"))
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vkey: "cosConfig.secretKey",
      type: "password",
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      value: secretKey,
      onChange: function onChange(v) {
        return actions.cosConfig({
          secretKey: v
        });
      },
      size: "l"
    }))), storageType === PaasMedium.StorageTypeEnum.NFS && React__default.createElement(React__default.Fragment, null, React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('Server Name')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "nfsConfig.serverName",
      value: serverName,
      onChange: function onChange(v) {
        return actions.nfsConfig({
          serverName: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('Server')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "nfsConfig.server",
      value: server,
      onChange: function onChange(v) {
        return actions.nfsConfig({
          server: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('Path')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "nfsConfig.path",
      value: path,
      onChange: function onChange(v) {
        return actions.nfsConfig({
          path: v
        });
      },
      size: "l"
    })), React__default.createElement(ffComponent.FormPanel.Item, {
      label: i18n.t('MountOptions'),
      message: i18n.t('支持填写多个；用,分割')
    }, React__default.createElement(ffComponent.FormPanel.Input, {
      vactions: actions.validator,
      formvalidator: validator,
      errorTipsStyle: 'Bubble',
      vkey: "nfsConfig.mountOptions",
      value: mountOptions,
      className: "tea-mr-2n",
      onChange: function onChange(v) {
        return actions.nfsConfig({
          mountOptions: v
        });
      },
      size: "l"
    }))))), React__default.createElement(teaComponent.Modal.Footer, null, React__default.createElement(teaComponent.Button, {
      type: "primary",
      className: "tea-mr-2n",
      disabled: workflow.operationState === ffRedux.OperationState.Performing,
      onClick: perform
    }, failed ? i18n.t('重试') : i18n.t('完成')), React__default.createElement(teaComponent.Button, {
      title: i18n.t('取消'),
      onClick: cancel
    }, i18n.t('取消')), React__default.createElement(TipInfo, {
      isShow: failed,
      type: "error",
      isForm: true
    }, getWorkflowError(workflow)), React__default.createElement(TipInfo, {
      isShow: !!nameInUseTip,
      type: "error",
      isForm: true
    }, nameInUseTip)));
  };
})(MediumEditDialog || (MediumEditDialog = {}));

var MediumDeleteDialog;
(function (MediumDeleteDialog) {
  var _this = this;
  MediumDeleteDialog.ComponentName = "MediumDeleteDialog";
  var ActionTypes;
  (function (ActionTypes) {
    ActionTypes["DELETE_RESOURCE_WORKFLOW"] = "DELETE_RESOURCE_WORKFLOW";
    ActionTypes["CLEAR_STATE"] = "CLEAR_STATE";
    ActionTypes["SET_VISIBLE"] = "SET_VISIBLE";
  })(ActionTypes = MediumDeleteDialog.ActionTypes || (MediumDeleteDialog.ActionTypes = {}));
  // 将ActionType加上ns的隔离
  BComponent$1.createActionType(MediumDeleteDialog.ComponentName, ActionTypes);
  MediumDeleteDialog.createActions = function (_a) {
    var pageName = _a.pageName,
      getRecord = _a.getRecord;
    var actions = {
      destory: function destory() {
        return {
          type: BComponent$1.getActionType(pageName, ActionTypes.CLEAR_STATE)
        };
      },
      setVisible: function setVisible(visible) {
        return {
          type: BComponent$1.getActionType(pageName, ActionTypes.SET_VISIBLE),
          payload: visible
        };
      },
      deleteResourceWorkflow: ffRedux.generateWorkflowActionCreator({
        actionType: BComponent$1.getActionType(pageName, ActionTypes.DELETE_RESOURCE_WORKFLOW),
        // operationExecutor: WebApi.k8sResource.deleteMutliResourceIns,
        operationExecutor: function operationExecutor(targets, options) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var result;
            return tslib.__generator(this, function (_a) {
              switch (_a.label) {
                case 0:
                  return [4 /*yield*/, medium.deleteResource(targets === null || targets === void 0 ? void 0 : targets[0], options)];
                case 1:
                  result = _a.sent();
                  return [2 /*return*/, result];
              }
            });
          });
        },
        workflowStateLocator: function workflowStateLocator(state) {
          return getRecord(function () {
            return state;
          }).deleteResourceWorkflow;
        },
        /** DONT USE THIS*/
        after: {}
      })
    };
    return actions;
  };
  MediumDeleteDialog.createReducer = function (_a) {
    var pageName = _a.pageName;
    var reducer = redux.combineReducers({
      visible: ffRedux.reduceToPayload(BComponent$1.getActionType(pageName, ActionTypes.SET_VISIBLE), false),
      deleteResourceWorkflow: ffRedux.generateWorkflowReducer({
        actionType: BComponent$1.getActionType(pageName, ActionTypes.DELETE_RESOURCE_WORKFLOW)
      })
    });
    var finalReducer = function finalReducer(state, action) {
      var newState = state;
      // 销毁组件
      if (action.type === BComponent$1.getActionType(pageName, ActionTypes.CLEAR_STATE)) {
        newState = undefined;
      }
      return reducer(newState, action);
    };
    return finalReducer;
  };
  MediumDeleteDialog.Component = function (props) {
    var model = props.model,
      _a = props.destoryOnClose,
      destoryOnClose = _a === void 0 ? true : _a,
      regionId = props.regionId,
      clusterId = props.clusterId,
      platform = props.platform,
      resourceIns = props.resourceIns,
      resourceInfo = props.resourceInfo,
      onFail = props.onFail,
      onCancelFromProps = props.onCancel,
      onCloseFromProps = props.onClose,
      onSuccess = props.onSuccess,
      _b = props.headTitle,
      headTitle = _b === void 0 ? i18n.t("资源") : _b,
      deleteTips = props.deleteTips,
      action = props.action,
      _c = props.caption,
      caption = _c === void 0 ? i18n.t("删除资源") : _c,
      _d = props.width,
      width = _d === void 0 ? 495 : _d,
      _e = props.disabledConfirm,
      disabledConfirm = _e === void 0 ? false : _e,
      body = props.body,
      _f = props.showIcon,
      showIcon = _f === void 0 ? false : _f;
    var visible = model.visible,
      deleteResourceWorkflow = model.deleteResourceWorkflow;
    React.useEffect(function () {
      return function () {
        if (destoryOnClose) {
          action.destory();
        }
      };
    }, [destoryOnClose, action.destory, action]);
    var onClose = function onClose() {
      if (typeof onCloseFromProps === "function") {
        onCloseFromProps();
      } else {
        action.setVisible(false);
        action.deleteResourceWorkflow.reset();
      }
    };
    var onCancel = function onCancel() {
      if (typeof onCancelFromProps === "function") {
        onCancelFromProps();
      } else {
        action.setVisible(false);
        action.deleteResourceWorkflow.reset();
      }
    };
    var operationState = deleteResourceWorkflow.operationState;
    var failed = operationState === ffRedux.OperationState.Done && !ffRedux.isSuccessWorkflow(deleteResourceWorkflow);
    React.useEffect(function () {
      if (operationState === ffRedux.OperationState.Done) {
        if (failed && typeof onFail === "function") {
          onFail();
        } else if (!failed && typeof onSuccess === "function") {
          onSuccess();
        }
      }
    }, [operationState, failed, onFail, onSuccess]);
    var finalBody = body ? body : React__default.createElement("div", {
      style: {
        fontSize: "14px",
        lineHeight: "20px"
      }
    }, React__default.createElement("p", {
      style: {
        wordWrap: "break-word"
      }
    }, showIcon && React__default.createElement(teaComponent.Icon, {
      type: "warning",
      size: "l",
      style: {
        marginRight: "8px"
      }
    }), React__default.createElement("strong", null, i18n.t("您确定要删除{{headTitle}}：{{resourceIns}}吗？", {
      headTitle: headTitle,
      resourceIns: resourceIns
    }))), deleteTips && React__default.createElement("div", {
      style: {
        marginLeft: showIcon ? "40px" : ""
      }
    }, deleteTips));
    return React__default.createElement(ModalMain.Modal, {
      caption: caption,
      size: width,
      onClose: onClose,
      visible: visible,
      disableEscape: true
    }, React__default.createElement(ModalMain.Modal.Body, null, finalBody, React__default.createElement(TipInfo, {
      isShow: failed,
      type: "error",
      isForm: true
    }, getWorkflowError(deleteResourceWorkflow))), React__default.createElement(ModalMain.Modal.Footer, null, React__default.createElement(teaComponent.Button, {
      loading: deleteResourceWorkflow.operationState === ffRedux.OperationState.Performing,
      type: "primary",
      disabled: disabledConfirm || deleteResourceWorkflow.operationState === ffRedux.OperationState.Performing /** inputDeleteTarget可能是空數組？ */,
      onClick: function onClick() {
        action.deleteResourceWorkflow.start([{
          clusterId: clusterId,
          resourceIns: resourceIns,
          platform: platform,
          regionId: +regionId,
          resourceInfo: resourceInfo
        }], +regionId);
        action.deleteResourceWorkflow.perform();
      }
    }, failed ? React__default.createElement("span", null, i18n.t("重试")) : React__default.createElement("span", null, i18n.t("确定"))), React__default.createElement(teaComponent.Button, {
      onClick: onCancel
    }, i18n.t("取消"))));
  };
})(MediumDeleteDialog || (MediumDeleteDialog = {}));

var SubEnum$1;
(function (SubEnum) {
  SubEnum["List"] = "list";
  SubEnum["Detail"] = "detail";
  SubEnum["Create"] = "create";
  //编辑实例
  SubEnum["Edit"] = "edit";
})(SubEnum$1 || (SubEnum$1 = {}));
var TabEnum$1;
(function (TabEnum) {
  TabEnum["Info"] = "info";
  TabEnum["Instance"] = "instance";
  TabEnum["Yaml"] = "yaml";
})(TabEnum$1 || (TabEnum$1 = {}));
var getRouterPath = function getRouterPath(pathname) {
  var _a;
  var path;
  if (pathname === null || pathname === void 0 ? void 0 : pathname.includes((_a = PlatformType.TKESTACK) === null || _a === void 0 ? void 0 : _a.toLowerCase())) {
    path = "/tkestack/medium";
  } else {
    path = "/tdcc/medium";
  }
  return path;
};
/**
 * @param sub 当前的模式，create | update | detail
 */
var router$1 = new Router("".concat(getRouterPath(location === null || location === void 0 ? void 0 : location.pathname), "(/:sub)(/:tab)"), {
  sub: SubEnum$1.List,
  tab: ""
});

var MediumTablePanel;
(function (MediumTablePanel) {
  var _this = this;
  MediumTablePanel.ComponentName = 'MediumTablePanel';
  MediumTablePanel.ActionType = Object.assign({}, BComponent$1.BaseActionType, {
    List: 'List',
    ExternalCluster: 'FETCH_EXTERNAL_CLUSTERS',
    Subscription: 'Subscription',
    AddonList: 'AddonList',
    ConfigMap: 'ConfigMap',
    HealthCheckWorkflow: 'HealthCheckWorkflow'
  });
  // 将ActionType加上ns的隔离
  BComponent$1.createActionType(MediumTablePanel.ComponentName, MediumTablePanel.ActionType);
  MediumTablePanel.createActions = function (_a) {
    var pageName = _a.pageName,
      _getRecord = _a.getRecord;
    var actions = {
      list: ffRedux.createFFListActions({
        actionName: BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.List),
        fetcher: function fetcher(query) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var k8sQueryObj, result;
            return tslib.__generator(this, function (_a) {
              switch (_a.label) {
                case 0:
                  k8sQueryObj = {
                    labelSelector: {
                      'tdcc.cloud.tencent.com/paas-storage-medium': ['s3', 'nfs']
                    }
                  };
                  if (!(query === null || query === void 0 ? void 0 : query.filter)) return [3 /*break*/, 2];
                  return [4 /*yield*/, medium.fetchResourceList(tslib.__assign(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter), {
                    k8sQueryObj: k8sQueryObj
                  }))];
                case 1:
                  result = _a.sent();
                  _a.label = 2;
                case 2:
                  return [2 /*return*/, result];
              }
            });
          });
        },
        getRecord: function getRecord(getState) {
          return _getRecord(getState).list;
        },
        onFinish: function onFinish(record, dispatch) {
          dispatch(actions.list.clearPolling());
        }
      }),
      getClusterAdminRole: GetRbacAdminDialog.createActions({
        pageName: BComponent$1.getActionType(pageName, MediumTablePanel.ComponentName),
        getRecord: function getRecord(getState) {
          return _getRecord(getState).getClusterAdminRole;
        }
      }),
      externalClusters: ffRedux.createFFListActions({
        actionName: BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.ExternalCluster),
        fetcher: function fetcher(query, getState) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var _a, regionId, platform, result;
            var _b, _c;
            return tslib.__generator(this, function (_d) {
              switch (_d.label) {
                case 0:
                  _a = query.filter, regionId = _a.regionId, platform = _a.platform;
                  if (!(platform && regionId)) return [3 /*break*/, 2];
                  return [4 /*yield*/, fetchExternalClusters({
                    platform: platform,
                    regionId: regionId
                  })];
                case 1:
                  result = _d.sent();
                  // 过滤状态为运行状态的注册集群
                  result.records = (_b = result === null || result === void 0 ? void 0 : result.records) === null || _b === void 0 ? void 0 : _b.filter(function (item) {
                    var _a;
                    return (item === null || item === void 0 ? void 0 : item.status) === ((_a = ExternalCluster === null || ExternalCluster === void 0 ? void 0 : ExternalCluster.StatusEnum) === null || _a === void 0 ? void 0 : _a.Running);
                  });
                  result.recordCount = (_c = result === null || result === void 0 ? void 0 : result.records) === null || _c === void 0 ? void 0 : _c.length;
                  _d.label = 2;
                case 2:
                  return [2 /*return*/, result];
              }
            });
          });
        },
        getRecord: function getRecord(getState) {
          return _getRecord(getState).externalClusters;
        },
        selectFirst: false
      }),
      create: MediumEditDialog.createActions({
        pageName: BComponent$1.getActionType(pageName, MediumTablePanel.ComponentName),
        getRecord: function getRecord(getState) {
          return _getRecord(getState).create;
        }
      }),
      "delete": MediumDeleteDialog.createActions({
        pageName: BComponent$1.getActionType(pageName, MediumTablePanel.ComponentName),
        getRecord: function getRecord(getState) {
          return _getRecord(getState)["delete"];
        }
      }),
      clear: function clear() {
        return function (dispatch, getState) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            return tslib.__generator(this, function (_a) {
              dispatch({
                type: BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.Clear)
              });
              return [2 /*return*/];
            });
          });
        };
      }
    };

    return actions;
  };
  MediumTablePanel.createValidateSchema = function (_a) {
    var pageName = _a.pageName;
    var schema = {
      formKey: BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.Validator),
      fields: []
    };
    return schema;
  };
  MediumTablePanel.createReducer = function (_a) {
    var pageName = _a.pageName;
    var TempReducer = redux.combineReducers({
      list: ffRedux.createFFListReducer({
        actionName: BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.List),
        displayField: function displayField(r) {
          return r.metadata.name;
        },
        valueField: function valueField(r) {
          return r.metadata.name;
        }
      }),
      getClusterAdminRole: GetRbacAdminDialog.createReducer({
        pageName: BComponent$1.getActionType(pageName, MediumTablePanel.ComponentName)
      }),
      externalClusters: ffRedux.createFFListReducer({
        actionName: BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.ExternalCluster),
        displayField: 'clusterName',
        valueField: 'clusterId'
      }),
      create: MediumEditDialog.createReducer({
        pageName: BComponent$1.getActionType(pageName, MediumTablePanel.ComponentName)
      }),
      "delete": MediumDeleteDialog.createReducer({
        pageName: BComponent$1.getActionType(pageName, MediumTablePanel.ComponentName)
      })
    });
    var Reducer = function Reducer(state, action) {
      var newState = state;
      // 销毁页面
      if (action.type === BComponent$1.getActionType(pageName, MediumTablePanel.ActionType.Clear)) {
        newState = undefined;
      }
      return TempReducer(newState, action);
    };
    return Reducer;
  };
  MediumTablePanel.Component = React__default.memo(function (props) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t;
    var action = props.action,
      _u = props.model,
      list = _u.list,
      externalClusters = _u.externalClusters,
      create = _u.create,
      getClusterAdminRole = _u.getClusterAdminRole,
      model = props.model,
      _v = props.filter,
      regionId = _v.regionId,
      platform = _v.platform,
      route = _v.route,
      k8sVersion = _v.k8sVersion,
      userInfo = _v.userInfo;
    var _w = React.useState(null),
      operateType = _w[0],
      setOperateType = _w[1];
    var clusterId = ((_a = route.queries) !== null && _a !== void 0 ? _a : {}).clusterId;
    React.useEffect(function () {
      var _a, _b, _c, _d, _e, _f;
      if ((_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _a === void 0 ? void 0 : _a.fetched) {
        (_b = action.externalClusters) === null || _b === void 0 ? void 0 : _b.selectByValue(clusterId || ((_f = (_e = (_d = (_c = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _c === void 0 ? void 0 : _c.data) === null || _d === void 0 ? void 0 : _d.records) === null || _e === void 0 ? void 0 : _e[0]) === null || _f === void 0 ? void 0 : _f.clusterId));
      }
    }, [action.externalClusters, clusterId, (_c = (_b = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _b === void 0 ? void 0 : _b.data) === null || _c === void 0 ? void 0 : _c.records, (_d = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.list) === null || _d === void 0 ? void 0 : _d.fetched]);
    React.useEffect(function () {
      var _a, _b;
      var filter = {
        platform: platform,
        regionId: regionId,
        clusterId: (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId,
        namespace: 'ssm'
      };
      if (((_b = externalClusters.list) === null || _b === void 0 ? void 0 : _b.fetched) && filter.clusterId && BComponent$1.isNeedFetch(list, filter)) {
        action.list.applyFilter(filter);
      }
    }, [action.list, list, externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection, platform, regionId, (_e = externalClusters.list) === null || _e === void 0 ? void 0 : _e.fetched]);
    React.useEffect(function () {
      var filter = {
        platform: platform,
        clusterIds: [],
        regionId: regionId
      };
      if (BComponent$1.isNeedFetch(externalClusters, filter)) {
        action.externalClusters.applyFilter({
          platform: platform,
          clusterIds: [],
          regionId: regionId
        });
      }
    }, [externalClusters, action.externalClusters]);
    //无权限提示
    var errorText = React__default.useMemo(function () {
      var _a, _b, _c, _d;
      var errorText;
      try {
        var listRbacErrorInfo = JSON.parse((_c = (_b = (_a = list === null || list === void 0 ? void 0 : list.list) === null || _a === void 0 ? void 0 : _a.error) === null || _b === void 0 ? void 0 : _b.message) !== null && _c !== void 0 ? _c : null);
        errorText = (listRbacErrorInfo === null || listRbacErrorInfo === void 0 ? void 0 : listRbacErrorInfo.code) === 403 ? React__default.createElement(i18n.Trans, null, React__default.createElement(teaComponent.Text, {
          verticalAlign: "middle"
        }, "\u6743\u9650\u4E0D\u8DB3\uFF0C\u8BF7\u8054\u7CFB\u96C6\u7FA4\u7BA1\u7406\u5458\u6DFB\u52A0\u6743\u9650\uFF1B\u82E5\u60A8\u672C\u8EAB\u662F\u96C6\u7FA4\u7BA1\u7406\u5458\uFF0C\u53EF\u76F4\u63A5"), React__default.createElement(teaComponent.Button, {
          type: "link",
          onClick: function onClick() {
            // 弹出获取集群admin角色的按钮
            action.getClusterAdminRole.getClusterAdminRole.start([]);
          }
        }, "\u83B7\u53D6\u96C6\u7FA4admin\u89D2\u8272"), React__default.createElement(teaComponent.Text, {
          verticalAlign: "middle"
        }, React__default.createElement(i18n.Slot, {
          content: "(".concat((_d = listRbacErrorInfo === null || listRbacErrorInfo === void 0 ? void 0 : listRbacErrorInfo.message) !== null && _d !== void 0 ? _d : '-', ")")
        }))) : undefined;
      } catch (error) {}
      return errorText;
    }, [action.getClusterAdminRole.getClusterAdminRole, (_g = (_f = list === null || list === void 0 ? void 0 : list.list) === null || _f === void 0 ? void 0 : _f.error) === null || _g === void 0 ? void 0 : _g.message]);
    var navigateDetail = function navigateDetail(resource) {
      var _a;
      router$1.navigate({
        sub: 'detail'
      }, {
        resourceIns: resource.metadata.name,
        clusterId: (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId
      });
    };
    React.useEffect(function () {
      return function () {
        action.clear();
      };
    }, []);
    var columns = [{
      key: 'name',
      header: i18n.t('名称'),
      render: function render(item) {
        var _a;
        return React__default.createElement(teaComponent.Button, {
          type: "link",
          onClick: function onClick() {
            navigateDetail(item);
          }
        }, (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name);
      }
    }, {
      key: 'type',
      header: i18n.t('类型'),
      render: function render(item) {
        var _a, _b;
        return React__default.createElement("p", null, PaasMedium.StorageTypeMap[(_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['tdcc.cloud.tencent.com/paas-storage-medium']]);
      }
    },
    // {
    //   key: "clusterId",
    //   // header: t('集群ID(集群名称)'),
    //   header: t("集群ID"),
    //   render: (item) => (
    //     <Text>
    //       {t("{{clusterId}}", {
    //         clusterId: Util.getClusterId(platform, item, route),
    //       })}
    //       {/* {t('{{clusterName}}', { clusterName: `(${Util.getClusterName(platform,item,route)})`})} */}
    //     </Text>
    //   ),
    // },
    {
      key: 'user',
      header: i18n.t('创建人'),
      render: function render(item) {
        var _a, _b, _c;
        return React__default.createElement("p", null, (_c = (_b = (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.labels) === null || _b === void 0 ? void 0 : _b['tdcc.cloud.tencent.com/creator']) !== null && _c !== void 0 ? _c : '-');
      }
    }, {
      key: 'creationTimestamp',
      header: i18n.t('时间戳'),
      render: function render(item) {
        var _a;
        return React__default.createElement("p", null, dateFormatter(new Date((_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.creationTimestamp), 'YYYY-MM-DD HH:mm:ss') || '-');
      }
    }, {
      key: 'operate',
      header: i18n.t('操作'),
      render: function render(item) {
        return React__default.createElement(React__default.Fragment, null, React__default.createElement(teaComponent.Button, {
          type: "link",
          onClick: function onClick() {
            action.list.select(item);
            setOperateType(BComponent$1.OperationTypeEnum.Modify);
            action.create.initForUpdate(item);
            action.create.workflow.start([]);
          }
        }, i18n.t('更新配置')), React__default.createElement(teaComponent.Button, {
          key: 'delete',
          type: "link",
          onClick: function onClick() {
            action.list.select(item);
            action["delete"].setVisible(true);
          }
        }, i18n.t('删除')));
      }
    }];
    return React__default.createElement(teaComponent.Layout, null, React__default.createElement(teaComponent.Layout.Body, null, React__default.createElement(teaComponent.Layout.Content, null, React__default.createElement(teaComponent.Layout.Content.Header, {
      title: i18n.t('备份介质')
    }), React__default.createElement(teaComponent.Layout.Content.Body, {
      full: true
    }, React__default.createElement(teaComponent.Table.ActionPanel, null, React__default.createElement(teaComponent.Justify, {
      left: React__default.createElement(teaComponent.Button, {
        type: "primary",
        // disabled={!isCanCreateExternal}
        onClick: function onClick() {
          setOperateType(BComponent$1.OperationTypeEnum.Create);
          action.create.workflow.start([]);
        }
      }, i18n.t('创建')),
      right: React__default.createElement(React__default.Fragment, null, React__default.createElement("div", {
        style: {
          display: 'inline-block',
          fontSize: '12px',
          verticalAlign: 'middle'
        }
      }, React__default.createElement(teaComponent.Text, {
        theme: "label",
        verticalAlign: "middle"
      }, i18n.t('注册集群')), React__default.createElement(ffComponent.FormPanel.Select, {
        model: externalClusters,
        onChange: function onChange(value) {
          action.externalClusters.selectByValue(value);
          router$1.navigate({}, {
            clusterId: value
          });
        },
        action: action.externalClusters,
        size: "m",
        errorTipsStyle: "Bubble"
      })), React__default.createElement("div", {
        style: {
          width: 350,
          display: 'inline-block'
        }
      }, React__default.createElement(ffComponent.TablePanelTagSearchBox, {
        disabled: !((_j = (_h = model.list) === null || _h === void 0 ? void 0 : _h.list) === null || _j === void 0 ? void 0 : _j.fetched),
        onChange: function onChange(searchFilter) {
          var _a;
          action.list.changeSearchFilter(searchFilter);
          action.list.applyFilter({
            platform: platform,
            regionId: regionId,
            clusterId: (_a = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _a === void 0 ? void 0 : _a.clusterId,
            namespace: SystemNamespace,
            searchFilter: searchFilter
          });
        },
        attributes: [{
          type: 'input',
          key: 'resourceName',
          name: i18n.t('名称')
        }],
        model: model.list,
        tips: i18n.t('名称只能搜索一个关键字')
      })), React__default.createElement(teaComponent.Button, {
        icon: "refresh",
        title: i18n.t('刷新'),
        onClick: function onClick() {
          action.list.fetch();
        }
      }))
    })), React__default.createElement(ffComponent.TablePanel, {
      operationsWidth: 300,
      recordKey: function recordKey(record) {
        var _a;
        return (_a = record === null || record === void 0 ? void 0 : record.metadata) === null || _a === void 0 ? void 0 : _a.name;
      },
      columns: columns,
      model: list,
      action: action.list,
      isNeedPagination: true,
      errorText: errorText
    }), React__default.createElement(MediumDeleteDialog.Component, {
      action: action["delete"],
      model: model["delete"],
      regionId: regionId,
      clusterId: (_k = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _k === void 0 ? void 0 : _k.clusterId,
      resourceIns: (_m = (_l = list === null || list === void 0 ? void 0 : list.selection) === null || _l === void 0 ? void 0 : _l.metadata) === null || _m === void 0 ? void 0 : _m.name,
      resourceInfo: list === null || list === void 0 ? void 0 : list.selection,
      platform: platform,
      headTitle: i18n.t('存储介质'),
      onSuccess: function onSuccess() {
        action["delete"].setVisible(false);
        action.list.startPolling({
          delayTime: 10000
        });
        action.list.fetch();
        action["delete"].deleteResourceWorkflow.reset();
        action.list.resetPaging();
      }
    }), React__default.createElement(MediumEditDialog.Component, {
      actions: action.create,
      model: create,
      regionId: regionId,
      clusterId: (_o = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _o === void 0 ? void 0 : _o.clusterId,
      clusterName: (_p = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _p === void 0 ? void 0 : _p.clusterName,
      k8sVersion: k8sVersion,
      operationType: operateType,
      userInfo: userInfo,
      platform: platform,
      existNames: (_s = (_r = (_q = list === null || list === void 0 ? void 0 : list.list) === null || _q === void 0 ? void 0 : _q.data) === null || _r === void 0 ? void 0 : _r.records) === null || _s === void 0 ? void 0 : _s.map(function (item) {
        var _a;
        return (_a = item === null || item === void 0 ? void 0 : item.metadata) === null || _a === void 0 ? void 0 : _a.name;
      }),
      resourceInstance: operateType === BComponent$1.OperationTypeEnum.Modify ? list === null || list === void 0 ? void 0 : list.selection : undefined,
      onSuccess: function onSuccess() {
        action.list.fetch();
        action.create.clear();
        action.list.resetPaging();
        action.create.workflow.reset();
      }
    }), React__default.createElement(GetRbacAdminDialog.Component, {
      model: getClusterAdminRole,
      action: action.getClusterAdminRole,
      filter: {
        platform: platform,
        regionId: regionId,
        clusterId: (_t = externalClusters === null || externalClusters === void 0 ? void 0 : externalClusters.selection) === null || _t === void 0 ? void 0 : _t.clusterId
      },
      onSuccess: function onSuccess() {
        var _a;
        (_a = action.list) === null || _a === void 0 ? void 0 : _a.fetch();
      }
    })))));
  });
})(MediumTablePanel || (MediumTablePanel = {}));

/**
 *
 * @param filter
 * @returns
 */
function isFilterReady(filter) {
  if (_typeof(filter) !== "object") {
    return false;
  }
  var keys = Object.keys(filter);
  var isOk = true;
  for (var _i = 0, keys_1 = keys; _i < keys_1.length; _i++) {
    var key = keys_1[_i];
    //保证filter内的key都是可预设的值
    isOk = isOk && filter[key] !== undefined && filter[key] !== null;
    if (!isOk) {
      break;
    }
  }
  return isOk;
}

var MediumInfoPanel = React__default.memo(function (props) {
  var _a;
  var resourceInstance = props.model.resourceInstance;
  var _b = React__default.useMemo(function () {
      var resourceIns = resourceInstance.object.data;
      var isNeedLoading = resourceInstance.object.fetched !== true || resourceInstance.object.fetchState === ffRedux.FetchState.Fetching || resourceIns === null;
      return {
        resourceIns: resourceIns,
        isNeedLoading: isNeedLoading
      };
    }, [resourceInstance.object]),
    resourceIns = _b.resourceIns,
    isNeedLoading = _b.isNeedLoading;
  return isNeedLoading ? React__default.createElement(teaComponent.Card, null, React__default.createElement(teaComponent.Card.Body, null, React__default.createElement(teaComponent.Icon, {
    type: "loading"
  }))) : React__default.createElement(ffComponent.FormPanel, {
    title: i18n.t("基本信息")
  }, React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t("名称"),
    text: true
  }, (_a = resourceIns.metadata) === null || _a === void 0 ? void 0 : _a.name), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t("运行的节点列表"),
    text: true
  }, React__default.createElement("div", null, "test")), React__default.createElement(ffComponent.FormPanel.Item, {
    label: i18n.t("期望的节点列表"),
    text: true
  }, React__default.createElement("div", null, "test")));
});

var MediumDetailPanel;
(function (MediumDetailPanel) {
  var _this = this;
  MediumDetailPanel.ComponentName = "MediumDetailPanel";
  MediumDetailPanel.ActionType = Object.assign({}, BComponent$1.BaseActionType, {
    ResourceInstance: "ResourceInstance"
  });
  // 将ActionType加上ns的隔离
  BComponent$1.createActionType(MediumDetailPanel.ComponentName, MediumDetailPanel.ActionType);
  MediumDetailPanel.createActions = function (_a) {
    var pageName = _a.pageName,
      _getRecord = _a.getRecord;
    var actions = {
      resourceInstance: ffRedux.createFFObjectActions({
        actionName: BComponent$1.getActionType(pageName, MediumDetailPanel.ActionType.ResourceInstance),
        fetcher: function fetcher(query) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            var result;
            return tslib.__generator(this, function (_a) {
              switch (_a.label) {
                case 0:
                  return [4 /*yield*/, medium.fetchResourceDetail(tslib.__assign({}, query === null || query === void 0 ? void 0 : query.filter))];
                case 1:
                  result = _a.sent();
                  return [2 /*return*/, result];
              }
            });
          });
        },
        getRecord: function getRecord(getState) {
          return _getRecord(getState).resourceInstance;
        }
      }),
      clear: function clear() {
        return function (dispatch, getState) {
          return tslib.__awaiter(_this, void 0, void 0, function () {
            return tslib.__generator(this, function (_a) {
              dispatch({
                type: BComponent$1.getActionType(pageName, MediumDetailPanel.ActionType.Clear)
              });
              return [2 /*return*/];
            });
          });
        };
      }
    };

    return actions;
  };
  MediumDetailPanel.createValidateSchema = function (_a) {
    var pageName = _a.pageName;
    var schema = {
      formKey: BComponent$1.getActionType(pageName, MediumDetailPanel.ActionType.Validator),
      fields: []
    };
    return schema;
  };
  MediumDetailPanel.createReducer = function (_a) {
    var pageName = _a.pageName;
    var TempReducer = redux.combineReducers({
      resourceInstance: ffRedux.createFFObjectReducer({
        actionName: BComponent$1.getActionType(pageName, MediumDetailPanel.ActionType.ResourceInstance)
      })
    });
    var Reducer = function Reducer(state, action) {
      var newState = state;
      // 销毁页面
      if (action.type === BComponent$1.getActionType(pageName, MediumDetailPanel.ActionType.Clear)) {
        newState = undefined;
      }
      return TempReducer(newState, action);
    };
    return Reducer;
  };
  MediumDetailPanel.Component = React__default.memo(function (props) {
    // const { navigate, urlParams } = useRoute();
    var action = props.action,
      resourceType = props.resourceType,
      resourceInstance = props.model.resourceInstance,
      filter = props.filter,
      _a = props.filter,
      regionId = _a.regionId,
      clusterId = _a.clusterId,
      k8sVersion = _a.k8sVersion,
      resourceIns = _a.resourceIns,
      platform = _a.platform,
      route = _a.route;
    React__default.useEffect(function () {
      return function () {
        action.clear();
      };
    }, [action]);
    //获取当前资源实例
    React.useEffect(function () {
      var filter = {
        regionId: regionId,
        clusterId: clusterId,
        platform: platform,
        instanceName: resourceIns,
        // specificName: resourceIns,
        k8sVersion: k8sVersion
      };
      if (isFilterReady(filter) && BComponent$1.isNeedFetch(resourceInstance, filter)) {
        action.resourceInstance.applyFilter(filter);
      }
    }, [action, resourceInstance, k8sVersion, regionId, clusterId, resourceIns, platform]);
    return React__default.createElement(teaComponent$1.Layout, null, React__default.createElement(teaComponent$1.Layout.Body, null, React__default.createElement(teaComponent$1.Layout.Content, null, React__default.createElement(teaComponent$1.Layout.Content.Header, {
      showBackButton: true,
      onBackButtonClick: function onBackButtonClick() {
        var clusterId = route.queries.clusterId;
        var urlParams = router$1 === null || router$1 === void 0 ? void 0 : router$1.resolve(route);
        router$1.navigate(tslib.__assign(tslib.__assign({}, urlParams), {
          sub: "list"
        }), {
          clusterId: clusterId
        });
      },
      title: i18n.t("备份介质")
    }), React__default.createElement(teaComponent$1.Layout.Content.Body, {
      full: true
    }, React__default.createElement(MediumInfoPanel, tslib.__assign({}, props))))));
  });
})(MediumDetailPanel || (MediumDetailPanel = {}));

var MediumPage;
(function (MediumPage) {
  MediumPage.List = MediumTablePanel;
  MediumPage.ComponentName = "MediumPage";
  MediumPage.ActionType = Object.assign({}, BComponent$1.BaseActionType, {});
  // 将ActionType加上ns的隔离
  BComponent$1.createActionType(MediumPage.ComponentName, MediumPage.ActionType);
  MediumPage.createActions = function (_a) {
    var pageName = _a.pageName,
      _getRecord = _a.getRecord;
    var actions = {
      list: MediumTablePanel.createActions({
        pageName: pageName,
        getRecord: function getRecord(getState) {
          return _getRecord(getState).list;
        }
      }),
      detail: MediumDetailPanel.createActions({
        pageName: pageName,
        getRecord: function getRecord(getState) {
          return _getRecord(getState).detail;
        }
      })
    };
    return actions;
  };
  MediumPage.createReducer = function (_a) {
    var pageName = _a.pageName;
    var TempReducer = redux.combineReducers({
      list: MediumTablePanel.createReducer({
        pageName: pageName
      }),
      detail: MediumDetailPanel.createReducer({
        pageName: pageName
      })
    });
    var Reducer = function Reducer(state, action) {
      var newState = state;
      return TempReducer(newState, action);
    };
    return Reducer;
  };
  MediumPage.Component = React__default.memo(function (props) {
    var _a;
    var action = props.action,
      route = props.filter.route,
      filter = props.filter,
      _b = props.model,
      list = _b.list,
      detail = _b.detail;
    var urlParams = router$1 === null || router$1 === void 0 ? void 0 : router$1.resolve(route);
    var _c = (_a = route.queries) !== null && _a !== void 0 ? _a : {},
      resourceIns = _c.resourceIns,
      clusterId = _c.clusterId;
    var content = null;
    if ((urlParams === null || urlParams === void 0 ? void 0 : urlParams.sub) === "detail") {
      content = React__default.createElement(MediumDetailPanel.Component, {
        action: action.detail,
        model: detail,
        filter: tslib.__assign(tslib.__assign({}, filter), {
          resourceIns: resourceIns,
          clusterId: clusterId
        })
      });
    } else {
      content = React__default.createElement(MediumTablePanel.Component, {
        action: action.list,
        model: list,
        filter: filter
      });
    }
    return React__default.createElement(React__default.Fragment, null, content);
  });
})(MediumPage || (MediumPage = {}));

function prefixPageId$1(obj, pageId) {
  Object.keys(obj).forEach(function (key) {
    obj[key] = "".concat(pageId, "_").concat(obj[key]);
  });
}
var Base$1;
(function (Base) {
  Base["IsI18n"] = "IsI18n";
  Base["FETCH_PLATFORM"] = "FETCH_PLATFORM";
  Base["FETCH_REGION"] = "FETCH_REGION";
  Base["HubCluster"] = "HubCluster";
  Base["ClusterVersion"] = "ClusterVersion";
  Base["SELECT_TAB"] = "SELECT_TAB";
  Base["Clear"] = "Clear";
  Base["FETCH_UserInfo"] = "FETCH_UserInfo";
  Base["UPDATE_ROUTE"] = "UPDATE_ROUTE";
  Base["GetClusterAdminRoleFlow"] = "GetClusterAdminRoleFlow";
})(Base$1 || (Base$1 = {}));
var Create$1;
(function (Create) {
  Create["SERVICE_INSTANCE_EDIT"] = "SERVICE_INSTANCE_EDIT";
  Create["SERVICE_INSTANCE_EDIT_VALIDATOR"] = "SERVICE_INSTANCE_EDIT_VALIDATOR";
  Create["CREATE_SERVICE_RESOURCE"] = "CREATE_SERVICE_RESOURCE";
  Create["UPDATE_SERVICE_RESOURCE"] = "UPDATE_SERVICE_RESOURCE";
  Create["SERVICE_PLAN_EDIT"] = "SERVICE_PLAN_EDIT";
  Create["CREATE_SERVICE_INSTANCE"] = "CREATE_SERVICE_INSTANCE";
})(Create$1 || (Create$1 = {}));
var List$1;
(function (List) {
  List["FETCH_SERVICES"] = "FETCH_SERVICES";
  List["FETCH_CREATE_RESOURCE_SCHEMAS"] = "FETCH_CREATE_RESOURCE_SCHEMAS";
  List["FETCH_SERVICE_RESOURCE"] = "FETCH_SERVICE_RESOURCE";
  List["FETCH_SERVICE_PLANS"] = "FETCH_SERVICE_PLANS";
  List["FETCH_Service_Plan_Map"] = "FETCH_Service_Plan_Map";
  List["FETCH_SERVICE_RESOURCE_LIST"] = "FETCH_SERVICE_RESOURCE_LIST";
  List["SELECT_DELETE_RESOURCE"] = "SELECT_DELETE_RESOURCE";
  List["DELETE_SERVICE_RESOURCE"] = "DELETE_SERVICE_RESOURCE";
  List["SHOW_INSTANCE_TABLE_DIALOG"] = "SHOW_INSTANCE_TABLE_DIALOG";
  List["SHOW_CREATE_RESOURCE_DIALOG"] = "SHOW_CREATE_RESOURCE_DIALOG";
  List["FETCH_EXTERNAL_CLUSTERS"] = "FETCH_EXTERNAL_CLUSTERS";
  List["Clear"] = "Clear";
})(List$1 || (List$1 = {}));
var Detail$1;
(function (Detail) {
  Detail["RESOURCE_DETAIL"] = "RESOURCE_DETAIL";
  Detail["Clear"] = "Clear";
  Detail["Select_Detail_Tab"] = "Select_Detail_Tab";
  Detail["FETCH_INSTANCE_RESOURCE"] = "FETCH_INSTANCE_RESOURCE";
  Detail["BACKUP_RESOURCE"] = "BACKUP_RESOURCE";
  Detail["CHECK_COS"] = "CHECK_COS";
  Detail["BACKUP_RESOURCE_LOADING"] = "BACKUP_RESOURCE_LOADING";
  Detail["SHOW_BACKUP_STRATEGY_DIALOG"] = "SHOW_BACKUP_STRATEGY_DIALOG";
  Detail["BACKUP_STRATEGY_EDIT"] = "BACKUP_STRATEGY_EDIT";
  Detail["OPEN_CONSOLE_WORKFLOW"] = "OPEN_CONSOLE_WORKFLOW";
  Detail["SHOW_CREATE_RESOURCE_DIALOG"] = "SHOW_CREATE_RESOURCE_DIALOG";
  Detail["FETCH_NAMESPACES"] = "FETCH_NAMESPACES";
  Detail["SERVICE_BINDING_EDIT"] = "SERVICE_BINDING_EDIT";
  Detail["FETCH_SERVICE_INSTANCE_SCHEMA"] = "FETCH_SERVICE_INSTANCE_SCHEMA";
  Detail["SELECT_DETAIL_RESOURCE"] = "SELECT_DETAIL_RESOURCE";
  Detail["FETCH_BACKUP_STRATEGY"] = "FETCH_BACKUP_STRATEGY";
})(Detail$1 || (Detail$1 = {}));
prefixPageId$1(Base$1, "Base");
prefixPageId$1(Create$1, "Create");
prefixPageId$1(List$1, "List");
prefixPageId$1(Detail$1, "Detail");
var ActionType$1 = {
  Base: Base$1,
  Create: Create$1,
  List: List$1,
  Detail: Detail$1
};

var routerSea$5 = seajs.require('router');
var baseActions$1 = {
  /** 国际版 */
  toggleIsI18n: function toggleIsI18n(isI18n) {
    return {
      type: ActionType$1.Base.IsI18n,
      payload: isI18n
    };
  },
  fetchRegion: function fetchRegion(regionId) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, route, platform;
        return tslib.__generator(this, function (_b) {
          dispatch({
            type: ActionType$1.Base.FETCH_REGION,
            payload: regionId
          });
          _a = getState().base, route = _a.route, platform = _a.platform;
          // tdcc首先拉取hub集群,然后拉取已开通的服务列表;tkeStack同时拉取目标集群和拉取已开通的服务列表
          if (platform === (PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TDCC)) {
            dispatch(baseActions$1.hubCluster.applyFilter({
              regionId: regionId,
              platform: platform
            }));
          } else if (platform === (PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TKESTACK)) ;
          return [2 /*return*/];
        });
      });
    };
  },

  fetchPlatform: function fetchPlatform(platform, region) {
    return function (dispatch, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var route;
        return tslib.__generator(this, function (_a) {
          dispatch({
            type: ActionType$1.Base.FETCH_PLATFORM,
            payload: platform
          });
          window['platform'] = platform;
          route = getState().base.route;
          return [2 /*return*/];
        });
      });
    };
  },

  hubCluster: ffRedux.createFFObjectActions({
    actionName: ActionType$1.Base.HubCluster,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var _a, platform, route, response, data;
        var _b;
        return tslib.__generator(this, function (_c) {
          switch (_c.label) {
            case 0:
              _a = getState().base, platform = _a.platform, route = _a.route;
              return [4 /*yield*/, fetchHubCluster(query === null || query === void 0 ? void 0 : query.filter)];
            case 1:
              response = _c.sent();
              data = (_b = response === null || response === void 0 ? void 0 : response.records) === null || _b === void 0 ? void 0 : _b[0];
              if (data && !(data === null || data === void 0 ? void 0 : data.regionId)) {
                data.regionId = HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion;
              }
              return [2 /*return*/, data];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().base.hubCluster;
    },
    onFinish: function onFinish(record, dispatch, getState) {
      var _a, _b;
      var _c = getState().base,
        route = _c.route,
        platform = _c.platform,
        regionId = _c.regionId;
      dispatch({
        type: ActionType$1.Base.ClusterVersion,
        payload: '1.18.4'
      });
      if (!record.data) {
        if (platform === PlatformType.TDCC) {
          routerSea$5.navigate('/tdcc/paasoverview/startup');
        }
      } else {
        // 更新region
        dispatch({
          payload: (_a = record === null || record === void 0 ? void 0 : record.data) === null || _a === void 0 ? void 0 : _a.regionId,
          type: (_b = ActionType$1 === null || ActionType$1 === void 0 ? void 0 : ActionType$1.Base) === null || _b === void 0 ? void 0 : _b.FETCH_REGION
        });
        // dispatch(
        //   listActions.services.applyFilter({
        //     platform,
        //     clusterId: record.data?.clusterId,
        //     regionId: record?.data?.regionId,
        //   })
        // );
        // dispatch(
        //   listActions?.externalClusters.applyFilter({
        //     platform,
        //     clusterIds: [],
        //     regionId: record?.data?.regionId,
        //   })
        // );
      }
    }
  }),

  userInfo: ffRedux.createFFObjectActions({
    actionName: ActionType$1.Base.FETCH_UserInfo,
    fetcher: function fetcher(query, getState) {
      return tslib.__awaiter(void 0, void 0, void 0, function () {
        var response;
        var _a;
        return tslib.__generator(this, function (_b) {
          switch (_b.label) {
            case 0:
              if (!(((_a = query === null || query === void 0 ? void 0 : query.filter) === null || _a === void 0 ? void 0 : _a.platform) === PlatformType.TDCC)) return [3 /*break*/, 1];
              response = {
                name: Util === null || Util === void 0 ? void 0 : Util.getUserName()
              };
              return [3 /*break*/, 3];
            case 1:
              return [4 /*yield*/, fetchUserInfo(query === null || query === void 0 ? void 0 : query.filter)];
            case 2:
              response = _b.sent();
              _b.label = 3;
            case 3:
              return [2 /*return*/, response !== null && response !== void 0 ? response : {}];
          }
        });
      });
    },
    getRecord: function getRecord(getState) {
      return getState().base.userInfo;
    }
  }),
  // 获取集群Admin权限
  getClusterAdminRole: GetRbacAdminDialog.createActions({
    pageName: ActionType$1.Base.GetClusterAdminRoleFlow,
    getRecord: function getRecord(getState) {
      return getState().base.getClusterAdminRole;
    }
  })
};

var createActions$1 = {};

var detailActions$1 = {};

var listActions$1 = {};

var allActions$1 = {
  base: baseActions$1,
  create: createActions$1,
  list: listActions$1,
  detail: detailActions$1,
  mediumPage: MediumPage.createActions({
    pageName: "TdccMedium",
    getRecord: function getRecord(getState) {
      return getState().mediumPage;
    }
  })
};

/* eslint-disable no-undef */
/**
 * 重置redux store，用于离开页面时清空状态
 */
var ResetStoreAction$1 = "ResetStore";
/**
 * 生成可重置的reducer，用于rootReducer简单包装
 * @return 可重置的reducer，当接收到 ResetStoreAction 时重置之
 */
// eslint-disable-next-line no-unused-vars
var generateResetableReducer$1 = function generateResetableReducer(rootReducer) {
  return function (state, action) {
    var newState = state;
    // 销毁页面
    if (action.type === ResetStoreAction$1) {
      newState = undefined;
    }
    return rootReducer(newState, action);
  };
};

var TempReducer$3 = redux.combineReducers({
  route: router.getReducer(),
  isI18n: ffRedux.reduceToPayload(ActionType.Base.IsI18n, false),
  platform: ffRedux.reduceToPayload(ActionType.Base.FETCH_PLATFORM, PlatformType === null || PlatformType === void 0 ? void 0 : PlatformType.TDCC),
  regionId: ffRedux.reduceToPayload(ActionType.Base.FETCH_REGION, HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion),
  hubCluster: ffRedux.createFFObjectReducer({
    actionName: ActionType.Base.HubCluster
  }),
  selectedTab: ffRedux.reduceToPayload(ActionType.Base.SELECT_TAB, ResourceTypeEnum === null || ResourceTypeEnum === void 0 ? void 0 : ResourceTypeEnum.ServiceResource),
  clusterVersion: ffRedux.reduceToPayload(ActionType.Base.ClusterVersion, ""),
  userInfo: ffRedux.createFFObjectReducer({
    actionName: ActionType.Base.FETCH_UserInfo
  }),
  getClusterAdminRole: GetRbacAdminDialog.createReducer({
    pageName: ActionType.Base.GetClusterAdminRoleFlow
  })
});
var baseReducer$1 = function baseReducer(state, action) {
  var newState = state;
  // 销毁页面
  if (action.type === ActionType.Base.Clear) {
    newState = undefined;
  }
  return TempReducer$3(newState, action);
};

var RootReducer$1 = redux.combineReducers({
  base: baseReducer$1,
  mediumPage: MediumPage.createReducer({
    pageName: "TdccMedium"
  })
});

var createStore$1 = process.env.NODE_ENV === 'development' ? redux.applyMiddleware(thunk, reduxLogger.createLogger({
  collapsed: true,
  diff: true
}))(redux.createStore) : redux.applyMiddleware(thunk)(redux.createStore);

function configStore$1() {
  var store = createStore$1(generateResetableReducer$1(RootReducer$1));
  // hot reloading
  // if (typeof module !== 'undefined' && (module as any).hot) {
  //   (module as any).hot.accept('../reducers/RootReducer', () => {
  //     store.replaceReducer(generateResetableReducer(require('../reducers/RootReducer').RootReducer));
  //   });
  // }
  return store;
}

var store$1 = configStore$1();
var MediumAppContainer = /** @class */function (_super) {
  tslib.__extends(MediumAppContainer, _super);
  function MediumAppContainer(props, context) {
    return _super.call(this, props, context) || this;
  }
  // 页面离开时，清空store
  MediumAppContainer.prototype.componentWillUnmount = function () {
    store$1.dispatch({
      type: ResetStoreAction$1
    });
  };
  MediumAppContainer.prototype.render = function () {
    return React.createElement(reactRedux.Provider, {
      store: store$1
    }, React.createElement(MediumApp, tslib.__assign({}, this.props)));
  };
  return MediumAppContainer;
}(React.Component);
var mapDispatchToProps$1 = function mapDispatchToProps(dispatch) {
  return Object.assign({}, ffRedux.bindActionCreators({
    actions: allActions$1
  }, dispatch), {
    dispatch: dispatch
  });
};
var MediumApp = /** @class */function (_super) {
  tslib.__extends(MediumApp, _super);
  function MediumApp(props, context) {
    return _super.call(this, props, context) || this;
  }
  MediumApp.prototype.componentDidMount = function () {
    var _a, _b, _c, _d, _e, _f, _g;
    var _h = this.props,
      actions = _h.actions,
      platform = _h.platform,
      isI18n = _h.base.isI18n,
      _j = _h.regionId,
      regionId = _j === void 0 ? HubCluster === null || HubCluster === void 0 ? void 0 : HubCluster.DefaultRegion : _j;
    if (window['VERSION'] === 'en' && !isI18n) {
      actions.base.toggleIsI18n(true);
    }
    if ((_a = this.props) === null || _a === void 0 ? void 0 : _a.platform) {
      (_b = actions === null || actions === void 0 ? void 0 : actions.base) === null || _b === void 0 ? void 0 : _b.fetchPlatform(platform, (_c = this === null || this === void 0 ? void 0 : this.props) === null || _c === void 0 ? void 0 : _c.regionId);
      (_e = (_d = actions === null || actions === void 0 ? void 0 : actions.base) === null || _d === void 0 ? void 0 : _d.hubCluster) === null || _e === void 0 ? void 0 : _e.applyFilter({
        regionId: regionId,
        platform: platform
      });
      (_g = (_f = actions.base) === null || _f === void 0 ? void 0 : _f.userInfo) === null || _g === void 0 ? void 0 : _g.applyFilter({
        platform: platform
      });
    }
  };
  MediumApp.prototype.render = function () {
    var _a = this.props,
      _b = _a.base,
      route = _b.route,
      platform = _b.platform,
      regionId = _b.regionId,
      clusterVersion = _b.clusterVersion,
      userInfo = _b.userInfo,
      actions = _a.actions,
      mediumPage = _a.mediumPage;
    // eslint-disable-next-line no-undef
    var content = React.createElement(MediumPage.Component, {
      filter: {
        platform: platform,
        route: route,
        regionId: regionId,
        k8sVersion: clusterVersion,
        userInfo: userInfo
      },
      action: actions.mediumPage,
      model: mediumPage
    });
    return React.createElement(React.Fragment, null, content);
  };
  MediumApp = tslib.__decorate([reactRedux.connect(function (state) {
    return state;
  }, mapDispatchToProps$1), router$1.serve()], MediumApp);
  return MediumApp;
}(React.Component);

exports.MediumApp = MediumApp;
exports.MediumAppContainer = MediumAppContainer;
exports.MiddlewareApp = MiddlewareApp;
exports.MiddlewareAppContainer = MiddlewareAppContainer;
exports.store = store;
