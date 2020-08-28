import { Identifiable } from '@tencent/ff-redux';

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name: string;

  displayName: string;

  clusterDisplayName?: string;

  clusterName?: string;

  namespace?: string;
}

export interface NamespaceFilter {
  /** clusterId */
  clusterId?: string;

  /** regionId */
  regionId?: number;

  /**是否选择默认值 */

  default?: boolean;
}
