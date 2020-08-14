import { Identifiable } from '@tencent/ff-redux';
import { Cluster } from './Cluster';

export interface Namespace extends Identifiable {
  /** 命名空间名称 */
  name?: string;

  /**命名空间 */
  namespace?: string;

  /** 用在业务侧的命名空间全名 */
  namespaceValue?: string;

  /**描述 */
  description?: string;

  /**状态 */
  status?: string;

  /**创建时间 */
  createdAt?: string;

  metadata?;

  cluster?: Cluster;
}

export interface NamespaceFilter {
  /**业务 */
  projectName?: string;

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
