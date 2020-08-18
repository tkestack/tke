import { Identifiable } from '@tencent/ff-redux';

import { Validation } from '../../common/models';

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name: string;

  displayName: string;

  //业务侧使用
  clusterVersion?: string;

  clusterId?: string;

  clusterDisplayName?: string;
}

/** 可视化创建的namespace的相关配置 */
export interface NamespaceEdit {
  /** name */
  name?: string;
  v_name?: Validation;

  /** 描述 */
  description?: string;
  v_description?: Validation;
}

/** 创建Namespace的时候，提交的jasonSchema */
export interface NamespaceEditJSONYaml {
  /** 资源的类型 */
  kind: string;

  /** api的版本 */
  apiVersion: string;

  /** metadata */
  metadata: NamespaceMetadata;

  /** spec */
  spec?: {};

  /** status */
  status?: {};
}

/** metadata的配置，非全部配置项 */
interface NamespaceMetadata {
  /** 插件能力 */
  annotations?: {
    [props: string]: string;
  };

  /** namespace的名称 */
  name: string;
}
