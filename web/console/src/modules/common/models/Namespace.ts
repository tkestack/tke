import { Identifiable } from '@tencent/ff-redux';

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name?: string;

  /**命名空间 */
  namespace?: string;

  /**描述 */
  description?: string;

  /**状态 */
  status?: string;

  /**创建时间 */
  createdAt?: string;

  metadata?;
}

export interface NamespaceFilter {
  /**集群Id */
  clusterId?: string;

  /**地域Id */
  regionId?: number;
}

export interface NamespaceOperator {
  /**集群Id */
  clusterId?: string;

  /**地域Id */
  regionId?: number;
}
