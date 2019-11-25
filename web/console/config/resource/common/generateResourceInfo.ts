import {
  RequestType,
  ResourceInfo,
  DisplayField,
  ActionField,
  DetailField,
  isEmpty
} from '../../../src/modules/common';
import { apiVersion, ApiVersionKeyName } from './apiVersion';

interface GenerateResourceInfo {
  /** k8s的版本 */
  k8sVersion: string;

  /** 资源的名称 */
  resourceName: ApiVersionKeyName;

  /** requestType */
  requestType: RequestType;

  /** 是否与ns有关 */
  isRelevantToNamespace?: boolean;

  /** displayField */
  displayField?: DisplayField;

  /** actionField */
  actionField?: ActionField;

  /** detailField */
  detailField?: DetailField;
}

export const generateResourceInfo = (options: GenerateResourceInfo): ResourceInfo => {
  let { k8sVersion, resourceName, isRelevantToNamespace = false, requestType, ...restOptions } = options;
  // apiVersion的配置
  const apiKind = apiVersion[k8sVersion][resourceName];
  // TKEStack当中，有自己的版本控制
  let serverVersionConfig: any;
  let watchModule = apiKind.watchModule;
  if (watchModule) {
    serverVersionConfig = window['modules'] && window['modules'][watchModule] ? window['modules'][watchModule] : {};
  }
  let config: ResourceInfo = {
    headTitle: apiKind.headTitle,
    basicEntry: apiKind.basicEntry,
    group: isEmpty(serverVersionConfig) || apiKind.group === '' ? apiKind.group : serverVersionConfig.groupName,
    version: isEmpty(serverVersionConfig) ? apiKind.version : serverVersionConfig.version,
    namespaces: isRelevantToNamespace ? 'namespaces' : '',
    requestType
  };

  if (!isEmpty(restOptions)) {
    config = Object.assign({}, config, restOptions);
  }
  return config;
};
