import { Identifiable } from '@tencent/ff-redux';

import { ResourceInfo } from './';

export interface CreateResource extends Identifiable {
  /** resourceInfo */
  resourceInfo?: ResourceInfo;

  /** 用户当前选择的命名空间 */
  namespace?: string;

  isSpetialNamespace?: boolean;

  /** yaml的数据 */
  yamlData?: string;

  /** 模式 create | update */
  mode?: string;

  /** 具体的resource资源的名称，如某个具体的 deployment的具体实例，update的时候使用 */
  resourceIns?: string;

  /** 当前的clusterId */
  clusterId?: string;

  /** 当前的logAgentName */
  logAgentName?: string;

  /** yamlJsonData 更新pod的数量、更新镜像等，都通过jsonData直接传过去 */
  jsonData?: string;

  /** 使用merge的方式，merge有几种方式的merge，k8s自己实现的以及JSON官方的 */
  isStrategic?: boolean;

  /** 使用merge的方式，merge有几种方式的merge，k8s自己实现的以及JSON官方的 */
  mergeType?: string;

  meshId?: string;

  /** 集群版本 */
  clusterVersion?: string;
}

export const MergeType = {
  Merge: ' application/merge-patch+json',
  Json: 'application/json-patch+json',
  StrategicMerge: 'application/strategic-merge-patch+json'
};
