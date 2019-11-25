import { Identifiable } from '@tencent/qcloud-lib';
import { Validation } from '../../common/models';

export interface NamespaceCreation extends Identifiable {
  /**命名空间名称 */
  name?: string;

  v_name?: Validation;

  /**命名空间备注*/
  description?: string;

  v_description?: Validation;

  /**集群 */
  clusterId?: string;

  /**地域 */
  regionId?: number;
}
